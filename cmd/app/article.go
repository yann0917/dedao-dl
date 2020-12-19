package app

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/yann0917/dedao-dl/utils"
)

// ArticleList 已购课程章节列表
func ArticleList(category string, id int) {
	courseDetail, err := getService().CourseDetail(category, id)
	if err != nil {
		return
	}
	enID := courseDetail.Enid
	list, err := getService().ArticleList(enID)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "章节名称", "更新时间", "是否阅读"})
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
