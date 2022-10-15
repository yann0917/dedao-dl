package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type GetDownload struct {
	OnEachStart func(t *DownloadTask)
	OnEachStop  func(t *DownloadTask)
	OnEachSkip  func(t *DownloadTask)
	Header      http.Header
	Client      http.Client
}

type DownloadTask struct {
	Link string
	Path string
	Err  error
}

type DownloadTasks struct {
	tasks []*DownloadTask
}

func Default() (g GetDownload) {
	g.Header = make(http.Header)
	g.Header.Set("user-agent", UserAgent)
	return g
}

var one = Default()

func Download(dl *DownloadTask, timeout time.Duration) (err error) {
	return one.Download(dl, timeout)
}
func DownloadWithContext(ctx context.Context, dl *DownloadTask) (err error) {
	return one.DownloadWithContext(ctx, dl)
}
func Batch(tasks *DownloadTasks, concurrent int, eachTimeout time.Duration) *DownloadTasks {
	return one.Batch(tasks, concurrent, eachTimeout)
}

func (g *GetDownload) Download(task *DownloadTask, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	return g.DownloadWithContext(ctx, task)
}
func (g *GetDownload) DownloadWithContext(ctx context.Context, task *DownloadTask) (err error) {
	if g.shouldSkip(ctx, task) {
		if g.OnEachSkip != nil {
			g.OnEachSkip(task)
		}
		return
	}
	if g.OnEachStart != nil {
		g.OnEachStart(task)
	}
	defer func() {
		task.Err = err
		if g.OnEachStop != nil {
			g.OnEachStop(task)
		}
	}()

	f, err := os.OpenFile(task.Path, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return
	}
	defer f.Close()

	req, err := http.NewRequest(http.MethodGet, task.Link, nil)
	if err != nil {
		return
	}
	if s, e := f.Stat(); e == nil {
		if s.Size() > 0 {
			req.Header.Set("range", fmt.Sprintf("bytes=%d-", s.Size()))
		}
	}
	for k := range g.Header {
		req.Header[k] = g.Header[k]
	}

	rsp, err := g.Client.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	defer func() {
		_, _ = io.Copy(io.Discard, rsp.Body)
		_ = rsp.Body.Close()
	}()

	switch rsp.StatusCode {
	case http.StatusPartialContent:
		_, _ = f.Seek(0, io.SeekEnd)
	case http.StatusOK, http.StatusRequestedRangeNotSatisfiable:
		_ = f.Truncate(0)
	default:
		return fmt.Errorf("invalid status code %d(%s)", rsp.StatusCode, rsp.Status)
	}

	_, err = io.Copy(f, rsp.Body)
	if err != nil {
		return fmt.Errorf("copy error: %s", err)
	}

	mt, e := http.ParseTime(rsp.Header.Get("last-modified"))
	if e == nil {
		_ = os.Chtimes(task.Path, mt, mt)
	}
	ok, e := os.Create(task.Path + ".ok")
	if e == nil {
		_ = ok.Close()
	}

	return
}
func (g *GetDownload) Batch(tasks *DownloadTasks, concurrent int, eachTimeout time.Duration) *DownloadTasks {
	var sema = semaphore.NewWeighted(int64(concurrent))
	var grp errgroup.Group

	tasks.ForEach(func(t *DownloadTask) {
		_ = sema.Acquire(context.TODO(), 1)
		grp.Go(func() (err error) {
			defer sema.Release(1)
			t.Err = g.Download(t, eachTimeout)
			return
		})
	})

	_ = grp.Wait()

	return tasks
}
func (g *GetDownload) shouldSkip(ctx context.Context, task *DownloadTask) (skip bool) {
	// check .ok file exist
	fd, err := os.Open(task.Path + ".ok")
	if err == nil {
		_ = fd.Close()
		return true
	}

	// check target file size
	local, err := os.Stat(task.Path)
	if err != nil {
		return false
	}

	switch local.Size() {
	case 0:
		return false
	default:
		req, err := http.NewRequest(http.MethodHead, task.Link, nil)
		if err == nil {
			req.Header = g.Header
			rsp, err := g.Client.Do(req.WithContext(ctx))
			if err == nil {
				_ = rsp.Body.Close()
				return rsp.ContentLength == local.Size()
			}
		}
		return false
	}
}

func NewDownloadTask(link, path string) *DownloadTask {
	return &DownloadTask{
		Link: link,
		Path: path,
	}
}

func (d *DownloadTasks) Add(link, path string) {
	for _, t := range d.tasks {
		if t.Link == link && t.Path == path {
			return
		}
	}
	d.tasks = append(d.tasks, NewDownloadTask(link, path))
}

func (d *DownloadTasks) ForEach(f func(t *DownloadTask)) {
	for _, t := range d.tasks {
		f(t)
	}
}

func NewDownloadTasks() *DownloadTasks {
	return &DownloadTasks{}
}
