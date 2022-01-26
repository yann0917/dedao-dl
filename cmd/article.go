package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func articleList(id int) (err error) {
	list, err := app.ArticleList(id, "")
	if err != nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "课程名称", "更新时间", "是否阅读"})
	table.SetAutoWrapText(false)

	for i, p := range list.List {
		isRead := "❌"
		if p.IsRead {
			isRead = "✔"
		}

		table.Append([]string{strconv.Itoa(i),
			p.IDStr, p.Title,
			utils.Unix2String(int64(p.UpdateTime)),
			isRead,
		})
	}
	table.Render()
	return
}

func articleDetail(id, aid int) (err error) {
	detail, err := app.ArticleDetail(id, aid)

	if err != nil {
		return
	}
	out := os.Stdout
	table := tablewriter.NewWriter(out)

	var content []services.Content
	jsoniter.UnmarshalFromString(detail.Content, &content)
	fmt.Fprint(out, contentsToMarkdown(content))
	fmt.Fprintln(out)
	table.Render()
	return
}

func contentsToMarkdown(contents []services.Content) (res string) {
	for _, content := range contents {
		switch content.Type {
		case "audio":
			title := strings.TrimRight(content.Title, ".mp3")
			res += getMdHeader(1) + title + "\r\n\r\n"

		case "header":
			res += getMdHeader(content.Level) + content.Text + "\r\n\r\n"
		case "paragraph":
			for _, item := range content.Contents {
				subContent := strings.Trim(item.Text.Content, " ")
				switch item.Type {
				case "text":
					if item.Text.Bold {
						res += " **" + subContent + "** "
					} else if item.Text.Highlight {
						res += " *" + subContent + "* "
					} else {
						res += subContent
					}
				}
			}
			res = strings.Trim(res, " ")
			res += "\r\n\r\n"
		case "elite": // 划重点
			res += getMdHeader(2) + "划重点\r\n\r\n" + content.Text + "\r\n\r\n"

		case "image":
			res += "![" + content.URL + "](" + content.URL + ")" + "\r\n\r\n"
		case "label-group":
			res += getMdHeader(2) + "`" + content.Text + "`" + "\r\n\r\n"
		}
	}

	res += "---\r\n"
	return
}

func getMdHeader(level int) string {
	switch level {
	case 1:
		return "# "
	case 2:
		return "## "
	case 3:
		return "### "
	case 4:
		return "#### "
	case 5:
		return "##### "
	case 6:
		return "###### "
	default:
		return ""
	}
}

func DownloadMarkdown(cType string, id int, filePath string) error {
	switch cType {
	case app.CateCourse:
		list, err := app.ArticleList(id, "")
		if err != nil {
			return err
		}
		for _, v := range list.List {
			detail, err := app.ArticleDetail(id, v.ID)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			var content []services.Content
			err = jsoniter.UnmarshalFromString(detail.Content, &content)
			if err != nil {
				return err
			}
			res := contentsToMarkdown(content)
			f, err := os.OpenFile(filePath+"/"+v.Title+".md", os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening file")
				return err
			}
			if err = f.Close(); err != nil {
				if err != nil {
					return err
				}
			}
			f.WriteString(res)
		}
	case app.CateAudioBook:

	}

	return nil

}
