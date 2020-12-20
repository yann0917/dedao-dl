package app

import (
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
	list, err := getService().EbookDetail(enID)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "章节名称"})
	table.SetAutoWrapText(false)

	for i, p := range list.CatalogList {

		table.Append([]string{strconv.Itoa(i), strconv.Itoa(p.PlayOrder),
			p.Text,
		})
	}
	table.Render()
	return
}
