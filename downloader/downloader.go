package downloader

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/yann0917/dedao-dl/utils"
)

// Options defines options used in downloading.
type Options struct {
	InfoOnly       bool
	Stream         string
	Refer          string
	OutputPath     string
	OutputName     string
	FileNameLength int
	Caption        bool

	MultiThread  bool
	ThreadNumber int
	RetryTimes   int
	ChunkSizeMB  int
}

// Downloader is the default downloader.
type Downloader struct {
	bar    *pb.ProgressBar
	option Options
}

func progressBar(size int, prefix string) *pb.ProgressBar {
	bar := pb.New(size).
		Set(pb.Bytes, true).
		SetRefreshRate(time.Millisecond * 10).
		SetMaxWidth(1000)

	return bar
}

// New returns a new Downloader implementation.
func New(option Options) *Downloader {
	downloader := &Downloader{
		option: option,
	}
	return downloader
}

//PrintToPDF print to pdf
func PrintToPDF(v Datum, cookies map[string]string, path string) error {

	name := utils.FileName(v.Title, "pdf")

	filename := filepath.Join(path, name)
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", name)

	_, exist, err := utils.FileSize(filename)

	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败")
		return err
	}

	if exist {
		fmt.Printf("\033[33;1m%s\033[0m\n", "已存在")
		return nil
	}

	err = utils.ColumnPrintToPDF(v.ID, filename, cookies)

	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败")
		return err
	}

	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")

	return nil
}
