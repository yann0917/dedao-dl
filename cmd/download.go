package cmd

import (
	"strconv"

	"errors"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var downloadType int
var downloadCmd = &cobra.Command{
	Use:   "dl",
	Short: "下载已购买课程，并转换成 PDF & 音频",
	Long: `使用 dedao-dl dl 下载已购买课程, 并转换成 PDF & 音频 & markdown
-t 指定下载格式, 1:mp3, 2:PDF文档, 3:markdown文档, 默认 mp3`,
	Example: "dedao-dl dl 123 -t 1",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("课程ID错误")
		}
		aid := 0
		if len(args) > 1 {
			aid, err = strconv.Atoi(args[1])
			if err != nil {
				return errors.New("文章ID错误")
			}
		}

		d := &app.CourseDownload{
			DownloadType: downloadType,
			ID:           id,
			AID:          aid,
		}
		err = app.Download(d)

		return err
	},
}

var dlOdobCmd = &cobra.Command{
	Use:   "dlo",
	Short: "下载每天听本书音频 & 文稿",
	Long: `使用 dedao-dl dlo 下载每天听本书音频, 并转换成 PDF & 音频 & markdown
-t 指定下载格式, 1:mp3, 2:PDF文档, 3:markdown文档, 默认 mp3`,
	Example: "dedao-dl dlo 123 -t 1",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("听书ID错误")
		}
		if len(args) > 1 {
			return errors.New("参数错误")
		}
		d := &app.OdobDownload{
			DownloadType: downloadType,
			ID:           id,
		}
		err = app.Download(d)
		return err
	},
}

var dlEbookCmd = &cobra.Command{
	Use:   "dle",
	Short: "下载电子书",
	Long: `使用 dedao-dl dle 下载电子书
-t 指定下载格式, 1:html, 2:PDF文档, 3:epub, 默认 html`,
	Example: "dedao-dl dle 123 -t 1",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("电子书ID错误")
		}
		if len(args) > 1 {
			return errors.New("参数错误")
		}
		d := &app.EBookDownload{
			DownloadType: downloadType,
			ID:           id,
		}
		err = app.Download(d)
		return err
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(dlOdobCmd)
	rootCmd.AddCommand(dlEbookCmd)
	downloadCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:mp3, 2:PDF文档, 3:markdown文档")
	dlOdobCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:mp3, 2:PDF文档, 3:markdown文档")
	dlEbookCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:html, 2:PDF文档, 3:epub")
}
