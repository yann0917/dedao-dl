package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yann0917/dedao-dl/utils"
)

const (
	progressStatusDone   = "done"
	progressStatusFailed = "failed"
)

var (
	errProgressNotExist  = errors.New("progress manifest not found")
	errProgressMalformed = errors.New("progress manifest malformed")
)

// ProgressRecord 单讲下载进度。File 是相对课程目录的路径，例如 MP3/001.xxx.mp3。
type ProgressRecord struct {
	ID        int    `json:"id"`
	Enid      string `json:"enid,omitempty"`
	Title     string `json:"title"`
	File      string `json:"file"`
	Size      int64  `json:"size"`
	Duration  int    `json:"duration,omitempty"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	UpdatedAt string `json:"updatedAt"`
}

type progressFile struct {
	Enid      string                    `json:"enid"`
	UpdatedAt string                    `json:"updatedAt"`
	Items     map[string]ProgressRecord `json:"items"`
}

// ProgressStore 持久化一个课程/听书的逐讲下载状态。
// 清单放在 <课程目录>/progress/<enid>.json，与 output 一起移动，
// 不依赖 .cache（clean cache 不会丢进度）。
type ProgressStore struct {
	dir  string
	name string
	path string

	// fileMu 按清单路径共享：同一进程里可能同时存在多个 store（例如两个 Web 下载会话），
	// 写入前必须串行并合并磁盘上的最新记录，避免旧 map 覆盖别人的更新。
	fileMu *sync.Mutex

	mu   sync.Mutex
	data progressFile
}

var progressFileLocks sync.Map // map[string]*sync.Mutex

func lockForProgressFile(path string) *sync.Mutex {
	value, _ := progressFileLocks.LoadOrStore(path, &sync.Mutex{})
	return value.(*sync.Mutex)
}

// NewProgressStore 创建或加载进度清单。name 建议用课程/听书的 enid。
func NewProgressStore(courseDir, name string) (*ProgressStore, error) {
	name = utils.FileName(strings.TrimSpace(name), "")
	if name == "" {
		name = "progress"
	}

	path := filepath.Join(courseDir, "progress", name+".json")
	fileMu := lockForProgressFile(path)
	fileMu.Lock()
	defer fileMu.Unlock()

	store := &ProgressStore{
		dir:    courseDir,
		name:   name,
		path:   path,
		fileMu: fileMu,
		data: progressFile{
			Enid:  name,
			Items: map[string]ProgressRecord{},
		},
	}

	if err := os.MkdirAll(filepath.Dir(store.path), 0755); err != nil {
		return nil, fmt.Errorf("创建进度目录失败: %w", err)
	}

	loaded, err := store.loadFromDisk()
	switch {
	case err == nil:
		store.data = loaded
		if store.data.Enid == "" {
			store.data.Enid = name
		}
	case errors.Is(err, errProgressNotExist):
		// 首次运行，保持空清单
	case errors.Is(err, errProgressMalformed):
		fmt.Printf("警告：进度文件 %s 解析失败，将重新记录\n", store.path)
	default:
		return nil, fmt.Errorf("读取进度文件失败: %w", err)
	}

	return store, nil
}

// Path 返回进度清单的路径。
func (p *ProgressStore) Path() string {
	return p.path
}

// Count 返回已记录完成/失败的数量。
func (p *ProgressStore) Count() (done, failed int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, item := range p.data.Items {
		switch item.Status {
		case progressStatusDone:
			done++
		case progressStatusFailed:
			failed++
		}
	}
	return done, failed
}

// Get 返回某个 articleID 的进度记录。
func (p *ProgressStore) Get(key string) (ProgressRecord, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, ok := p.data.Items[key]
	return item, ok
}

// MarkDone 记录一讲已下载完成。
func (p *ProgressStore) MarkDone(key string, v Datum, fileName string, size int64, duration int) error {
	return p.mark(key, ProgressRecord{
		ID:        v.ID,
		Enid:      v.Enid,
		Title:     v.Title,
		File:      p.rel(fileName),
		Size:      size,
		Duration:  duration,
		Status:    progressStatusDone,
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
}

// MarkFailed 记录一讲下载失败。
func (p *ProgressStore) MarkFailed(key string, v Datum, fileName string, err error) error {
	message := ""
	if err != nil {
		message = err.Error()
		if len(message) > 500 {
			message = message[:500]
		}
	}
	return p.mark(key, ProgressRecord{
		ID:        v.ID,
		Enid:      v.Enid,
		Title:     v.Title,
		File:      p.rel(fileName),
		Status:    progressStatusFailed,
		Error:     message,
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
}

func (p *ProgressStore) mark(key string, record ProgressRecord) error {
	p.fileMu.Lock()
	defer p.fileMu.Unlock()

	// 合并磁盘上的最新记录，避免并发 store 用旧 map 覆盖。
	latest, err := p.loadFromDisk()
	switch {
	case err == nil:
		if latest.Enid == "" {
			latest.Enid = p.name
		}
		p.mu.Lock()
		p.data = latest
		p.mu.Unlock()
	case errors.Is(err, errProgressNotExist), errors.Is(err, errProgressMalformed):
		// 首次运行或清单损坏：保留内存数据，本次写入重建清单。
	default:
		// 读失败时不要用旧内存覆盖磁盘上可能更新的记录。
		return fmt.Errorf("读取进度文件失败: %w", err)
	}

	p.mu.Lock()
	p.data.Items[key] = record
	p.data.UpdatedAt = record.UpdatedAt
	err = p.saveLocked()
	p.mu.Unlock()
	return err
}

func (p *ProgressStore) loadFromDisk() (progressFile, error) {
	raw, err := readFileRetry(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			return progressFile{}, errProgressNotExist
		}
		return progressFile{}, err
	}
	var loaded progressFile
	if err := json.Unmarshal(raw, &loaded); err != nil {
		return progressFile{}, errProgressMalformed
	}
	if loaded.Items == nil {
		loaded.Items = map[string]ProgressRecord{}
	}
	return loaded, nil
}

func (p *ProgressStore) saveLocked() error {
	raw, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化进度失败: %w", err)
	}

	// 每次写唯一临时文件，避免不同进程共用同一个 .tmp；os.Rename 会替换旧清单。
	f, err := os.CreateTemp(filepath.Dir(p.path), filepath.Base(p.path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("创建进度临时文件失败: %w", err)
	}
	tmp := f.Name()
	if _, err := f.Write(raw); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return fmt.Errorf("写入进度临时文件失败: %w", err)
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("关闭进度临时文件失败: %w", err)
	}
	if err := renameWithRetry(tmp, p.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("更新进度文件失败: %w", err)
	}
	return nil
}

// Windows 上其他进程/goroutine 短暂持有清单时，读取或替换可能返回共享冲突；
// 清单很小，短暂重试即可。
func readFileRetry(path string) ([]byte, error) {
	var err error
	for i := 0; i < 5; i++ {
		var raw []byte
		raw, err = os.ReadFile(path)
		if err == nil || os.IsNotExist(err) {
			return raw, err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, err
}

func renameWithRetry(oldPath, newPath string) error {
	var err error
	for i := 0; i < 5; i++ {
		err = os.Rename(oldPath, newPath)
		if err == nil {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return err
}

func (p *ProgressStore) rel(fileName string) string {
	rel, err := filepath.Rel(p.dir, fileName)
	if err != nil {
		return fileName
	}
	return rel
}

func progressKey(v Datum) string {
	if v.ID > 0 {
		return strconv.Itoa(v.ID)
	}
	if v.Enid != "" {
		return v.Enid
	}
	return v.Title
}
