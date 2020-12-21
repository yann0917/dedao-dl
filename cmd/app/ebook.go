package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// EbookDetail 电子书详情
func EbookDetail(id int) {
	courseDetail, err := getService().CourseDetail("ebook", id)
	if err != nil {
		return
	}
	enID := courseDetail.Enid
	detail, err := getService().EbookDetail(enID)
	if err != nil {
		return
	}

	out := os.Stdout
	table := tablewriter.NewWriter(out)
	fmt.Fprint(out, "书名："+detail.OperatingTitle+"\n")
	fmt.Fprint(out, "单价："+detail.Price+"\n")
	fmt.Fprint(out, "作者："+detail.BookAuthor+"\n")
	fmt.Fprint(out, "类型："+detail.ClassifyName+"\n")
	fmt.Fprint(out, "专家推荐指数："+detail.ProductScore+"\n")
	fmt.Fprint(out, "豆瓣评分："+detail.DoubanScore+"\n")
	fmt.Fprint(out, "发行日期："+detail.PublishTime+"\n")
	fmt.Fprint(out, "出版社："+detail.Press.Name+"\n")
	fmt.Fprintln(out)

	table.SetHeader([]string{"#", "ID", "章节名称"})
	table.SetAutoWrapText(false)
	// table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	// table.SetCenterSeparator("|")
	for i, p := range detail.CatalogList {
		table.Append([]string{strconv.Itoa(i), strconv.Itoa(p.PlayOrder),
			p.Text,
		})
	}
	table.Render()
	return
}
