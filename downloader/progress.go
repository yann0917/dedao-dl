package downloader

import (
	"encoding/json"
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
	path string

	mu   sync.Mutex
	data progressFile
}

// NewProgressStore 创建或加载进度清单。name 建议用课程/听书的 enid。
func NewProgressStore(courseDir, name string) (*ProgressStore, error) {
	name = utils.FileName(strings.TrimSpace(name), "")
	if name == "" {
		name = "progress"
	}

	store := &ProgressStore{
		dir:  courseDir,
		path: filepath.Join(courseDir, "progress", name+".json"),
		data: progressFile{
			Enid:  name,
			Items: map[string]ProgressRecord{},
		},
	}

	if err := os.MkdirAll(filepath.Dir(store.path), 0755); err != nil {
		return nil, fmt.Errorf("创建进度目录失败: %w", err)
	}

	raw, err := os.ReadFile(store.path)
	if err == nil {
		var loaded progressFile
		if unmarshalErr := json.Unmarshal(raw, &loaded); unmarshalErr == nil {
			if loaded.Items == nil {
				loaded.Items = map[string]ProgressRecord{}
			}
			store.data = loaded
			if store.data.Enid == "" {
				store.data.Enid = name
			}
		} else {
			fmt.Printf("警告：进度文件 %s 解析失败，将重新记录：%v\n", store.path, unmarshalErr)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取进度文件失败: %w", err)
	}

	return store, nil
}

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
	p.mu.Lock()
	defer p.mu.Unlock()

	p.data.Items[key] = record
	p.data.UpdatedAt = record.UpdatedAt
	return p.saveLocked()
}

func (p *ProgressStore) saveLocked() error {
	raw, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化进度失败: %w", err)
	}

	tmp := p.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0644); err != nil {
		return fmt.Errorf("写入进度临时文件失败: %w", err)
	}
	// Windows 下 Rename 不能覆盖已存在文件
	if err := os.Remove(p.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("替换进度文件失败: %w", err)
	}
	if err := os.Rename(tmp, p.path); err != nil {
		return fmt.Errorf("更新进度文件失败: %w", err)
	}
	return nil
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
