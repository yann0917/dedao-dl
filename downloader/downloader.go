package downloader

import (
	"time"

	"github.com/cheggaaa/pb/v3"
)

func progressBar(size int, prefix string) *pb.ProgressBar {
	bar := pb.New(size).
		Set(pb.Bytes, true).
		SetRefreshRate(time.Millisecond * 10).
		SetMaxWidth(1000)

	return bar
}
