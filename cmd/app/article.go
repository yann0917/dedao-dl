package app

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// ArticleList 已购课程章节列表
func ArticleList(id string) {
	list, err := getService().ArticleList(id)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "章节名称"})
	table.SetAutoWrapText(false)

	for i, p := range list.List {

		table.Append([]string{strconv.Itoa(i),
			strconv.Itoa(p.ID), p.Title,
		})
	}
	table.Render()
	return
}
