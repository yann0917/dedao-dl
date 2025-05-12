package cmd

import (
	"fmt"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

var articleCmd = &cobra.Command{
	Use:     "article",
	Short:   "获取文章详情",
	Long:    `使用 dedao-dl article 获取文章详情`,
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if classID > 0 && articleID == 0 {
			if err := articleList(classID); err != nil {
				return err
			}
		}

		if articleID > 0 {
			err := articleDetail(classID, articleID)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
	articleCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程id")
	articleCmd.PersistentFlags().IntVarP(&articleID, "aid", "a", 0, "文章id")
}

func articleList(id int) (err error) {
	list, err := app.ArticleList(id, "")
	if err != nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"#", "ID", "课程名称", "更新时间", "音频进度", "是否阅读"})

	for i, p := range list.List {
		isRead := "❌"
		if p.IsRead {
			isRead = "✔"
		}
		listenProgress := "0"
		if p.Audio != nil {
			listenProgress = strconv.FormatFloat(p.Audio.ListenProgress, 'g', 5, 32)
		}
		table.Append([]string{strconv.Itoa(i),
			p.IDStr, p.Title,
			utils.Unix2String(int64(p.UpdateTime)),
			listenProgress,
			isRead,
		})
	}
	table.Render()
	return
}

func articleDetail(id, aid int) (err error) {
	detail, _, err := app.ArticleDetail(id, aid)

	if err != nil {
		return
	}
	out := os.Stdout
	table := tablewriter.NewWriter(out)

	var content []services.Content
	err = jsoniter.UnmarshalFromString(detail.Content, &content)
	_, _ = fmt.Fprint(out, app.ContentsToMarkdown(content))
	_, _ = fmt.Fprintln(out)
	table.Render()
	return
}
