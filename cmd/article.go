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
	Use:   "article",
	Short: "获取文章详情",
	Long: `使用 dedao-dl article 获取文章详情

支持通过课程 ID/课程 EnID 列出文章列表，也支持通过文章 ID/文章 EnID 查看文章正文（Markdown 输出）。

示例：
  dedao-dl article --i 12345
  dedao-dl article --c xxxxx
  dedao-dl article --i 12345 --a 67890
  dedao-dl article --c xxxxx --e yyyyy
  dedao-dl article --e yyyyy`,
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if articleEnID != "" {
			return articleDetailByEnID(articleEnID)
		}

		if articleID > 0 {
			if classID > 0 {
				return articleDetail(classID, articleID)
			}
			if classEnID != "" {
				return articleDetailByClassEnIDAndArticleID(classEnID, articleID)
			}
			return fmt.Errorf("必须提供课程 ID(--id) 或课程 enid(--classEnID)，才能通过文章 ID(--aid) 获取详情")
		}

		if classID > 0 {
			return articleList(classID)
		}

		if classEnID != "" {
			return articleListByEnID(classEnID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
	articleCmd.PersistentFlags().StringVarP(&classEnID, "classEnID", "c", "", "课程enid")
	articleCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程id")
	articleCmd.PersistentFlags().IntVarP(&articleID, "aid", "a", 0, "文章id")
	articleCmd.PersistentFlags().StringVarP(&articleEnID, "articleEnID", "e", "", "文章enid")
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

func articleListByEnID(enid string) (err error) {
	list, err := app.ArticleListAllByEnId(enid, "")
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

func articleDetailByEnID(articleEnid string) (err error) {
	detail, err := app.ArticleDetailByEnId(articleEnid)
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

func articleDetailByClassEnIDAndArticleID(classEnid string, aid int) (err error) {
	list, err := app.ArticleListAllByEnId(classEnid, "")
	if err != nil {
		return
	}
	articleEnid := ""
	for _, p := range list.List {
		if p.ID == aid {
			articleEnid = p.Enid
			break
		}
	}
	if articleEnid == "" {
		return fmt.Errorf("在课程 enid=%s 下找不到文章 ID=%d", classEnid, aid)
	}
	return articleDetailByEnID(articleEnid)
}
