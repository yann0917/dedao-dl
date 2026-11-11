package downloader

import (
	"math"
	"os/exec"
	"path/filepath"
	"testing"
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

	ok, err = isAudioComplete(file, 60)
	if err != nil || ok {
		t.Fatalf("2 秒文件对 60 秒期望应为不完整: ok=%v err=%v", ok, err)
	}

	ok, err = isAudioComplete(filepath.Join(dir, "missing.mp3"), 2)
	if err != nil || ok {
		t.Fatalf("不存在的文件应为不完整: ok=%v err=%v", ok, err)
	}
}
