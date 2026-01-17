package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var downloadType, courseMerge, courseComment, courseOrder = 1, false, false, false

var downloadCmd = &cobra.Command{
	Use:   "dl",
	Short: "下载已购买课程，并转换成 PDF & 音频",
	Long: `使用 dedao-dl dl 下载已购买课程, 并转换成 PDF & 音频 & markdown
-t 指定下载格式, 1:mp3, 2:PDF文档, 3:markdown文档, 默认 mp3
-m 是否合并课程文稿(仅支持markdown), 默认不合并
-c 是否下载课程热门留言(仅支持markdown), 默认不下载`,
	Example: "dedao-dl dl 123 -t 1 -m",
	Args:    cobra.RangeArgs(1, 2),
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		enid := ""
		if err != nil {
			id = 0
			enid = args[0]
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
			EnID:         enid,
			AID:          aid,
			IsMerge:      courseMerge,
			IsComment:    courseComment,
			IsOrder:      courseOrder,
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
-t 指定下载格式, 1:html, 2:PDF文档, 3:epub, 4:markdown笔记, 默认 html`,
	Example: "dedao-dl dle 123 -t 1",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) > 1 {
			return errors.New("参数错误")
		}

		id, err := strconv.Atoi(args[0])
		var d app.DeDaoDownloader

		if downloadType == 4 {
			// 笔记下载格式
			if err != nil {
				// args[0] is not an integer, treat as EnID
				d = &app.EBookNotesDownload{
					DownloadType: downloadType,
					EnID:         args[0],
				}
			} else {
				// args[0] is an integer ID
				d = &app.EBookNotesDownload{
					DownloadType: downloadType,
					ID:           id,
				}
			}
		} else {
			if err != nil {
				d = &app.EBookDownload{
					DownloadType: downloadType,
					EnID:         args[0],
				}
			} else {
				d = &app.EBookDownload{
					DownloadType: downloadType,
					ID:           id,
				}
			}
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
	downloadCmd.PersistentFlags().BoolVarP(&courseMerge, "merge", "m", false, "是否合并课程章节")
	downloadCmd.PersistentFlags().BoolVarP(&courseComment, "comment", "c", false, "是否下载课程热门留言, 仅针对 markdown 文档")
	downloadCmd.PersistentFlags().BoolVarP(&courseOrder, "order", "o", false, "是否按顺序展示, 如果为true, 则文件名前缀会加上序号, 如 00x.")

	dlOdobCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:mp3, 2:PDF文档, 3:markdown文档")
	dlEbookCmd.PersistentFlags().IntVarP(&downloadType, "downloadType", "t", 1, "下载格式, 1:html, 2:PDF文档, 3:epub, 4:markdown笔记")
}
