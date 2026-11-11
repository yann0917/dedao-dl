package downloader

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/yann0917/dedao-dl/request"
	"github.com/yann0917/dedao-dl/utils"
)

// Download download data.
// 返回的 bool 表示该文件已存在且校验通过、本次被跳过。
func Download(v Datum, stream, path string) (bool, error) {
	// 按大到小排序
	v.genSortedStreams()

	title := utils.FileName(v.Title, "")
	if stream == "" {
		if len(v.sortedStreams) == 0 {
			return false, fmt.Errorf("没有可下载的流")
		}
		stream = v.sortedStreams[0].name
	}
	data, ok := v.Streams[stream]
	if !ok {
		return false, fmt.Errorf("指定要下载的类型不存在：%s", stream)
	}

	// 判断下载连接是否存在
	if len(data.URLs) == 0 {
		return false, nil
	}

	filePreName := filepath.Join(path, title)
	fileName, err := utils.FilePath(filePreName, "mp3", false)
	if err != nil {
		return false, err
	}

	if v.Type == "audio" {
		return downloadAudio(v, fileName)
	}

	_, mergedFileExists, err := utils.FileSize(fileName)
	if err != nil {
		return false, err
	}

	// After the merge, the file size has changed, so we do not check whether the size matches
	if mergedFileExists {
		return false, nil
	}

	chunkSizeMB := 1

	if len(data.URLs) == 1 {
		if err := Save(data.URLs[0], filePreName, chunkSizeMB); err != nil {
			return false, err
		}
		return false, nil
	}

	wgp := utils.NewWaitGroupPool(10)

	errs := make([]error, 0)
	lock := sync.Mutex{}
	parts := make([]string, len(data.URLs))

	for index, url := range data.URLs {
		if len(errs) > 0 {
			break
		}

		partFileName := fmt.Sprintf("%s[%d]", filePreName, index)
		partFilePath, err := utils.FilePath(partFileName, url.Ext, false)
		if err != nil {
			return false, err
		}
		parts[index] = partFilePath

		wgp.Add()
		go func(url URL, fileName string) {
			defer wgp.Done()
			err := Save(url, fileName, chunkSizeMB)
			if err != nil {
				lock.Lock()
				errs = append(errs, err)
				lock.Unlock()
			}
		}(url, partFileName)
	}

	wgp.Wait()

	if len(errs) > 0 {
		return false, errs[0]
	}

	switch v.Type {
	case "audio":
		err = utils.MergeAudio(parts, fileName)
	case "video":
		err = utils.MergeAudioAndVideo(parts, fileName)
	}

	if v.Type != "audio" && v.Type != "video" {
		return false, nil
	}

	return false, err
}

var (
	audioDownloadTimeout  = 15 * time.Minute
	audioDownloadAttempts = 2
	audioProbeAvailable   = func() bool {
		_, err := exec.LookPath("ffprobe")
		return err == nil
	}()
)

func downloadAudio(v Datum, fileName string) (bool, error) {
	if v.Force {
		if _, exists, err := utils.FileSize(fileName); err != nil {
			return false, err
		} else if exists {
			fmt.Printf("%s: 强制重新下载\n", fileName)
			if err := os.Remove(fileName); err != nil {
				return false, err
			}
		}
	} else {
		complete, err := isAudioComplete(fileName, v.Duration)
		if err != nil {
			return false, err
		}
		if complete {
			fmt.Printf("%s: 已存在且时长校验通过，跳过\n", fileName)
			markAudioDone(v, fileName)
			return true, nil
		}

		if _, exists, err := utils.FileSize(fileName); err != nil {
			return false, err
		} else if exists {
			fmt.Printf("%s: 文件不完整，重新下载\n", fileName)
			if err := os.Remove(fileName); err != nil {
				return false, err
			}
		}
	}

	// ffmpeg 直接写最终文件时，进程中断会留下看起来完整的残缺文件；
	// 先写 .part.mp3，成功后再原子改名。
	partFile := fileName + ".part.mp3"
	var lastErr error

	for attempt := 1; attempt <= audioDownloadAttempts; attempt++ {
		if err := os.Remove(partFile); err != nil && !os.IsNotExist(err) {
			return false, fmt.Errorf("清理临时文件失败: %w", err)
		}
		if attempt > 1 {
			fmt.Printf("重试第 %d/%d 次：%s\n", attempt, audioDownloadAttempts, fileName)
		}

		err := downloadAudioOnce(v.M3U8URL, partFile, v.Duration)
		if err == nil {
			if renameErr := os.Rename(partFile, fileName); renameErr != nil {
				markAudioFailed(v, fileName, renameErr)
				return false, fmt.Errorf("重命名临时文件失败: %w", renameErr)
			}
			markAudioDone(v, fileName)
			return false, nil
		}

		lastErr = err
		_ = os.Remove(partFile)
		fmt.Printf("下载失败（第 %d/%d 次）：%s，%v\n", attempt, audioDownloadAttempts, fileName, err)
		if attempt < audioDownloadAttempts {
			time.Sleep(2 * time.Second)
		}
	}

	markAudioFailed(v, fileName, lastErr)
	return false, lastErr
}

func markAudioDone(v Datum, fileName string) {
	if v.Progress == nil {
		return
	}

	size, _, err := utils.FileSize(fileName)
	if err != nil {
		fmt.Printf("警告：读取文件大小失败，进度未记录：%v\n", err)
		return
	}

	duration := 0
	if audioProbeAvailable {
		if actual, probeErr := utils.ProbeDuration(fileName); probeErr == nil {
			duration = int(actual + 0.5)
		}
	}

	if err := v.Progress.MarkDone(progressKey(v), v, fileName, int64(size), duration); err != nil {
		fmt.Printf("警告：写入进度失败：%v\n", err)
	}
}

func markAudioFailed(v Datum, fileName string, err error) {
	if v.Progress == nil {
		return
	}
	if markErr := v.Progress.MarkFailed(progressKey(v), v, fileName, err); markErr != nil {
		fmt.Printf("警告：写入失败进度失败：%v\n", markErr)
	}
}

func downloadAudioOnce(m3u8, partFile string, expectedDuration int) error {
	ctx, cancel := context.WithTimeout(context.Background(), audioDownloadTimeout)
	defer cancel()

	if err := utils.MergeAudioContext(ctx, []string{m3u8}, partFile); err != nil {
		return err
	}

	if expectedDuration > 0 && audioProbeAvailable {
		actual, err := utils.ProbeDuration(partFile)
		if err != nil {
			// ffmpeg 已成功；这里只做额外校验，ffprobe 读不出来不当作失败
			fmt.Printf("警告：无法校验音频时长（%v），按下载成功处理\n", err)
			return nil
		}
		if actual+audioDurationTolerance(expectedDuration) < float64(expectedDuration) {
			return fmt.Errorf("音频时长异常：期望 %ds，实际 %.0fs", expectedDuration, actual)
		}
	}
	return nil
}

// isAudioComplete 判断已有音频是否完整。有时长信息时用 ffprobe 校验，
// 避免把中断留下的残缺文件当成已完成。
func isAudioComplete(fileName string, expectedDuration int) (bool, error) {
	size, exists, err := utils.FileSize(fileName)
	if err != nil {
		return false, err
	}
	if !exists || size == 0 {
		return false, nil
	}
	if expectedDuration <= 0 || !audioProbeAvailable {
		return true, nil
	}

	actual, err := utils.ProbeDuration(fileName)
	if err != nil {
		// 文件存在但读不出时长，按损坏处理，重新下载
		return false, nil
	}
	return actual+audioDurationTolerance(expectedDuration) >= float64(expectedDuration), nil
}

func audioDurationTolerance(expected int) float64 {
	tolerance := float64(expected) * 0.02
	if tolerance < 5 {
		tolerance = 5
	}
	return tolerance
}

// Save url file
func Save(urlData URL, fileName string, chunkSizeMB int) error {
	if urlData.Size == 0 {
		size, err := request.Size(urlData.URL)
		if err != nil {
			return err
		}
		urlData.Size = size
	}

	var err error
	filePath, err := utils.FilePath(fileName, urlData.Ext, false)
	if err != nil {
		return err
	}
	fileSize, exists, err := utils.FileSize(filePath)
	if err != nil {
		return err
	}
	// Skip segment file
	// TODO: Live video URLs will not return the size
	if exists && fileSize == urlData.Size {
		return nil
	}
	tempFilePath := filePath + ".download"
	tempFileSize, _, err := utils.FileSize(tempFilePath)

	if err != nil {
		return err
	}
	headers := map[string]string{}
	var (
		file      *os.File
		fileError error
	)
	if tempFileSize > 0 {
		// range start from 0, 0-1023 means the first 1024 bytes of the file
		headers["Range"] = fmt.Sprintf("bytes=%d-", tempFileSize)
		file, fileError = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, fileError = os.Create(tempFilePath)
	}
	if fileError != nil {
		return fileError
	}

	// close and rename temp file at the end of this function
	defer func() {
		// must close the file before rename, or it will cause
		// `The process cannot access the file because it is being used by another process.` error.
		file.Close()
		if err == nil {
			os.Rename(tempFilePath, filePath)
		}
	}()

	if chunkSizeMB > 0 {
		var start, end, chunkSize int
		chunkSize = chunkSizeMB * 1024 * 1024
		remainingSize := urlData.Size
		if tempFileSize > 0 {
			start = tempFileSize
			remainingSize -= tempFileSize
		}
		chunk := remainingSize / chunkSize
		if remainingSize%chunkSize != 0 {
			chunk++
		}
		var i = 1
		for ; i <= chunk; i++ {
			end = start + chunkSize - 1
			headers["Range"] = fmt.Sprintf("bytes=%d-%d", start, end)
			temp := start
			for i := 0; ; i++ {
				written, err := writeFile(urlData.URL, file, headers)
				if err == nil {
					break
				} else if i+1 >= 3 {
					return err
				}
				temp += written
				headers["Range"] = fmt.Sprintf("bytes=%d-%d", temp, end)
				time.Sleep(1 * time.Second)
			}
			start = end + 1
		}
	} else {
		temp := tempFileSize
		for i := 0; ; i++ {
			written, err := writeFile(urlData.URL, file, headers)
			if err == nil {
				break
			} else if i+1 >= 3 {
				return err
			}
			temp += written
			headers["Range"] = fmt.Sprintf("bytes=%d-", temp)
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func writeFile(url string, file *os.File, headers map[string]string) (int, error) {
	res, err := request.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	writer := io.MultiWriter(file)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(writer, res)
	if copyErr != nil && copyErr != io.EOF {
		return int(written), fmt.Errorf("file copy error: %s", copyErr)
	}
	return int(written), nil
}

// PrintToPDF print to pdf
//func PrintToPDF(v Datum, cookies map[string]string, path string) error {
//
//	name := utils.FileName(v.Title, "pdf")
//
//	filename := filepath.Join(path, name)
//	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)
//
//	_, exist, err := utils.FileSize(filename)
//
//	if err != nil {
//		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
//		return err
//	}
//
//	if exist {
//		fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
//		return nil
//	}
//
//	err = utils.ColumnPrintToPDF(v.Enid, filename, cookies)
//
//	if err != nil {
//		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
//		return err
//	}
//
//	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
//
//	return nil
//}
