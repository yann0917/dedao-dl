package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runMergeCmd(cmd *exec.Cmd, paths []string, mergeFilePath string) error {
	return runMergeCmdContext(context.Background(), cmd, paths, mergeFilePath)
}

func runMergeCmdContext(ctx context.Context, cmd *exec.Cmd, paths []string, mergeFilePath string) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return fmt.Errorf("%w\n%s", ctxErr, stderr.String())
		}
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}

	if mergeFilePath != "" {
		os.Remove(mergeFilePath) // nolint
	}
	// remove parts
	for _, path := range paths {
		os.Remove(path) // nolint
	}
	return nil
}

// MergeAudio merge audio
func MergeAudio(paths []string, mergedFilePath string) error {
	return MergeAudioContext(context.Background(), paths, mergedFilePath)
}

// MergeAudioContext merge audio with a deadline/cancellation.
// ffmpeg 读 m3u8 时可能长时间卡在某个分片，调用方应传入带超时的 context。
func MergeAudioContext(ctx context.Context, paths []string, mergedFilePath string) error {
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", mergedFilePath)
	return runMergeCmdContext(ctx, exec.CommandContext(ctx, "ffmpeg", cmds...), paths, "")
}

// MergeAudioAndVideo merge audio and video
func MergeAudioAndVideo(paths []string, mergedFilePath string) error {
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", "-c:a", "copy", mergedFilePath)
	return runMergeCmd(exec.Command("ffmpeg", cmds...), paths, "")
}

// MergeToMP4 merges video parts to an MP4 file.
func MergeToMP4(paths []string, mergedFilePath string, filename string) error {
	mergeFilePath := filename + ".txt" // merge list file should be in the current directory
	// write ffmpeg input file list
	mergeFile, _ := os.Create(mergeFilePath)
	for _, path := range paths {
		_, _ = mergeFile.Write([]byte(fmt.Sprintf("file '%s'\n", path))) // nolint
	}
	err := mergeFile.Close() // nolint
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"ffmpeg", "-y", "-f", "concat", "-safe", "-1",
		"-i", mergeFilePath, "-c", "copy", "-bsf:a", "aac_adtstoasc", mergedFilePath,
	)
	return runMergeCmd(cmd, paths, mergeFilePath)
}

// ProbeDuration returns the media duration in seconds using ffprobe.
func ProbeDuration(filePath string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	out, err := exec.CommandContext(
		ctx,
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	).Output()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return 0, fmt.Errorf("ffprobe timed out: %w", ctxErr)
		}
		return 0, err
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "N/A" {
		return 0, fmt.Errorf("ffprobe returned no duration")
	}
	duration, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ffprobe duration %q: %w", raw, err)
	}
	return duration, nil
}
