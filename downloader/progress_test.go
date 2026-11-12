package downloader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestProgressStorePersistsRecords(t *testing.T) {
	dir := t.TempDir()
	store, err := NewProgressStore(dir, "course-enid")
	if err != nil {
		t.Fatalf("创建进度清单失败: %v", err)
	}

	file := filepath.Join(dir, "MP3", "001.第一讲.mp3")
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	doneItem := Datum{ID: 101, Enid: "article-enid", Title: "第一讲"}
	if err := store.MarkDone(progressKey(doneItem), doneItem, file, 1, 60); err != nil {
		t.Fatalf("记录完成失败: %v", err)
	}

	reloaded, err := NewProgressStore(dir, "course-enid")
	if err != nil {
		t.Fatalf("重新加载失败: %v", err)
	}
	record, ok := reloaded.Get("101")
	if !ok {
		t.Fatal("未找到 articleID=101 的记录")
	}
	if record.Status != progressStatusDone {
		t.Fatalf("status = %q, want %q", record.Status, progressStatusDone)
	}
	if record.File != filepath.Join("MP3", "001.第一讲.mp3") {
		t.Fatalf("file = %q", record.File)
	}
	if record.Size != 1 || record.Duration != 60 {
		t.Fatalf("size/duration = %d/%d", record.Size, record.Duration)
	}

	failedItem := Datum{ID: 102, Title: "第二讲"}
	if err := reloaded.MarkFailed(progressKey(failedItem), failedItem, filepath.Join(dir, "MP3", "002.mp3"), errors.New("boom")); err != nil {
		t.Fatalf("记录失败失败: %v", err)
	}

	again, err := NewProgressStore(dir, "course-enid")
	if err != nil {
		t.Fatalf("再次加载失败: %v", err)
	}
	record, ok = again.Get("102")
	if !ok || record.Status != progressStatusFailed {
		t.Fatalf("失败记录不正确: %#v", record)
	}
	if !strings.Contains(record.Error, "boom") {
		t.Fatalf("error = %q", record.Error)
	}

	done, failed := again.Count()
	if done != 1 || failed != 1 {
		t.Fatalf("count = %d/%d, want 1/1", done, failed)
	}
	if _, err := os.Stat(again.Path()); err != nil {
		t.Fatalf("进度清单不存在: %v", err)
	}
}

// 多个 store 同时写同一份清单时，不能互相覆盖记录。
func TestProgressStoreConcurrentWriters(t *testing.T) {
	dir := t.TempDir()
	const total = 20

	stores := make([]*ProgressStore, total)
	for i := 0; i < total; i++ {
		store, err := NewProgressStore(dir, "course-enid")
		if err != nil {
			t.Fatalf("创建 store 失败: %v", err)
		}
		stores[i] = store
	}

	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			item := Datum{ID: i + 1, Title: fmt.Sprintf("lesson-%d", i+1)}
			file := filepath.Join(dir, "MP3", fmt.Sprintf("%03d.mp3", i+1))
			if err := stores[i].MarkDone(progressKey(item), item, file, 1, 60); err != nil {
				t.Errorf("写入第 %d 条失败: %v", i+1, err)
			}
		}(i)
	}
	wg.Wait()

	reloaded, err := NewProgressStore(dir, "course-enid")
	if err != nil {
		t.Fatalf("重新加载失败: %v", err)
	}
	done, failed := reloaded.Count()
	if done != total || failed != 0 {
		t.Fatalf("并发写入后 count = %d/%d, want %d/0", done, failed, total)
	}
}
