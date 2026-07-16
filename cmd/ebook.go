package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var (
	ebookGroupID int
)

var ebookCmd = &cobra.Command{
	Use:     "ebook",
	Short:   "获取我的电子书架",
	Long:    `使用 dedao-dl ebook 获取我的电子书架`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ebook\ndedao-dl ebook --group-id 12345",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if bookID > 0 {
			return ebookDetail(bookID)
		}
		query, err := buildListQuery(app.CateEbook, listOrder, listPage, listLimit)
		if err != nil {
			return err
		}
		if ebookGroupID > 0 {
			return groupList(app.CateEbook, ebookGroupID, query)
		}
		return courseListByQuery(app.CateEbook, query)
	},
}

func init() {
	rootCmd.AddCommand(ebookCmd)

	ebookCmd.PersistentFlags().IntVarP(&bookID, "id", "i", 0, "电子书ID")
	ebookCmd.PersistentFlags().IntVarP(&ebookGroupID, "group-id", "g", 0, "分组ID，显示指定分组内的电子书")
	ebookCmd.PersistentFlags().IntVarP(&listPage, "page", "p", 0, "页码（与 --limit 一起使用）")
	ebookCmd.PersistentFlags().IntVarP(&listLimit, "limit", "l", 0, "每页数量（与 --page 一起使用）")
	ebookCmd.PersistentFlags().StringVar(&listOrder, "order", "study", "排序方式：study（默认）")
}

func ebookDetail(id int) (err error) {
	detail, err := app.EbookDetail(id)
	if err != nil {
		return
	}
	if outputJSON {
		return printJSON(detail)
	}

	out := os.Stdout
	table := tablewriter.NewWriter(out)
	_, _ = fmt.Fprint(out, "书名："+detail.OperatingTitle+"\n")
	_, _ = fmt.Fprint(out, "单价："+detail.Price+"\n")
	_, _ = fmt.Fprint(out, "作者："+detail.BookAuthor+"\n")
	_, _ = fmt.Fprint(out, "类型："+detail.ClassifyName+"\n")
	_, _ = fmt.Fprint(out, "专家推荐指数："+detail.ProductScore+"\n")
	_, _ = fmt.Fprint(out, "豆瓣评分："+detail.DoubanScore+"\n")
	_, _ = fmt.Fprint(out, "发行日期："+detail.PublishTime+"\n")
	_, _ = fmt.Fprint(out, "出版社："+detail.Press.Name+"\n")
	_, _ = fmt.Fprintln(out)

	table.Header([]string{"#", "ID", "章节名称"})
	for i, p := range detail.CatalogList {
		table.Append([]string{strconv.Itoa(i), strconv.Itoa(p.PlayOrder),
			p.Text,
		})
	}
	table.Render()
	return
}
