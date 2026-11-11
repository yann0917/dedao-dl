package downloader

import (
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/yann0917/dedao-dl/utils"
)

func TestAudioDurationTolerance(t *testing.T) {
	tests := []struct {
		expected int
		want     float64
	}{
		{0, 5},
		{60, 5},
		{300, 6},
		{614, 12.28},
		{3600, 72},
	}
	for _, tt := range tests {
		got := audioDurationTolerance(tt.expected)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("audioDurationTolerance(%d) = %v, want %v", tt.expected, got, tt.want)
		}
	}
}

func TestIsAudioComplete(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skipf("本机无 ffmpeg，跳过: %v", err)
	}
	if !audioProbeAvailable {
		t.Skip("本机无 ffprobe，跳过")
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "sample.mp3")
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "lavfi",
		"-i", "sine=frequency=440:duration=2",
		"-q:a", "9",
		file,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("生成测试音频失败: %v\n%s", err, out)
	}

	ok, err := isAudioComplete(file, 2)
	if err != nil || !ok {
		t.Fatalf("2 秒文件对 2 秒期望应为完整: ok=%v err=%v", ok, err)
	}

	actual, err := utils.ProbeDuration(file)
	if err != nil {
		t.Fatalf("读取测试音频时长失败: %v", err)
	}
	inside := int(math.Floor(actual + audioDurationTolerance(int(actual))))
	outside := inside + 1

	ok, err = isAudioComplete(file, inside)
	if err != nil || !ok {
		t.Fatalf("实际 %.2fs 对 %d 秒期望应仍在容差内: ok=%v err=%v", actual, inside, ok, err)
	}

	ok, err = isAudioComplete(file, outside)
	if err != nil || ok {
		t.Fatalf("实际 %.2fs 对 %d 秒期望应超出容差: ok=%v err=%v", actual, outside, ok, err)
	}

	ok, err = isAudioComplete(file, 60)
	if err != nil || ok {
		t.Fatalf("2 秒文件对 60 秒期望应为不完整: ok=%v err=%v", ok, err)
	}

	ok, err = isAudioComplete(filepath.Join(dir, "missing.mp3"), 2)
	if err != nil || ok {
		t.Fatalf("不存在的文件应为不完整: ok=%v err=%v", ok, err)
	}
}

func TestNewPartFileUnique(t *testing.T) {
	dir := t.TempDir()

	first, err := newPartFile(dir)
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(first)

	second, err := newPartFile(dir)
	if err != nil {
		t.Fatalf("创建第二个临时文件失败: %v", err)
	}
	defer os.Remove(second)

	if first == second {
		t.Fatalf("两次下载使用了同一个临时文件: %s", first)
	}
	if filepath.Ext(first) != ".mp3" {
		t.Fatalf("临时文件扩展名不是 .mp3: %s", first)
	}
}
