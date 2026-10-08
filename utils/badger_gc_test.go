package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// vlogGCStat 统计目录下 .vlog 文件的个数与总字节数
func vlogGCStat(t *testing.T, dir string) (int, int64) {
	t.Helper()
	n, size := 0, int64(0)
	files, err := filepath.Glob(filepath.Join(dir, "*.vlog"))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if st, err := os.Stat(f); err == nil {
			n++
			size += st.Size()
		}
	}
	return n, size
}

// 模拟连续多次短命 dle：打开库 → 写 3 个 1.5MB、TTL 24h 的章节 → 按下载成功路径删掉本书前缀 → 关闭。
// 章节内容超过 badger 的 ValueThreshold（1MB）会写进 .vlog，删除只是墓碑；
// 修复前每轮留下一个再不回收的 .vlog（5 轮共 23MB），修复后稳态只保留当前轮的活跃文件。
func TestValueLogGCOnClose(t *testing.T) {
	dir := t.TempDir()
	page := strings.Repeat("x", 1500*1024)
	for run := 1; run <= 5; run++ {
		db, err := NewBadgerDB(dir)
		if err != nil {
			t.Fatal(err)
		}
		for ch := 0; ch < 3; ch++ {
			key := fmt.Sprintf("ebook:page:book%d:ch%d", run, ch)
			if err := db.SetWithTTL(key, []string{page}, 24*time.Hour); err != nil {
				t.Fatal(err)
			}
		}
		if err := db.DeleteWithPrefix(fmt.Sprintf("ebook:page:book%d:", run)); err != nil {
			t.Fatal(err)
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}

		n, size := vlogGCStat(t, dir)
		t.Logf("第 %d 次运行后：.vlog %d 个，共 %.1f MB", run, n, float64(size)/1e6)

		// 每轮关闭应回收上一轮的垃圾；活跃文件（本轮写入）要等下一轮，属正常滞后
		if n > 2 {
			t.Fatalf("第 %d 轮后仍有 %d 个 .vlog，历史轮次的垃圾没有被回收", run, n)
		}
		if size > 8*1024*1024 {
			t.Fatalf("第 %d 轮后 .vlog 共 %.1f MB，超过稳态上限（约一轮的写入量）", run, float64(size)/1e6)
		}
	}
}

// 库里存在多轮历史垃圾时，一次关闭应把能收的都收掉
func TestValueLogGCCollectsBacklog(t *testing.T) {
	dir := t.TempDir()
	page := strings.Repeat("x", 1500*1024)

	// 制造 3 轮历史垃圾（写入后删除，模拟修复前的累积）
	for run := 1; run <= 3; run++ {
		db, err := NewBadgerDB(dir)
		if err != nil {
			t.Fatal(err)
		}
		for ch := 0; ch < 3; ch++ {
			key := fmt.Sprintf("ebook:page:legacy%d:ch%d", run, ch)
			if err := db.SetWithTTL(key, []string{page}, 24*time.Hour); err != nil {
				t.Fatal(err)
			}
		}
		if err := db.DeleteWithPrefix(fmt.Sprintf("ebook:page:legacy%d:", run)); err != nil {
			t.Fatal(err)
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}

	// 再开一轮，不做任何写入，直接关闭：应回收全部历史垃圾（无活跃文件）
	db, err := NewBadgerDB(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	n, size := vlogGCStat(t, dir)
	t.Logf("清理历史垃圾后：.vlog %d 个，共 %.2f MB", n, float64(size)/1e6)
	if n > 1 || size > 1024*1024 {
		t.Fatalf("历史垃圾未回收干净：.vlog %d 个，共 %.2f MB", n, float64(size)/1e6)
	}
}
