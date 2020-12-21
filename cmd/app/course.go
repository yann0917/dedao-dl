package app

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/yann0917/dedao-dl/services"
)

// CourseType 课程分类
func CourseType() (err error) {
	list, err := getService().CourseType()
	if err != nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "名称", "统计", "分类标签"})
	table.SetAutoWrapText(false)

	for i, p := range list.Data.List {

		table.Append([]string{strconv.Itoa(i), p.Name, strconv.Itoa(p.Count), p.Category})
	}
	table.Render()
	return
}

// CourseList 已购课程列表
func CourseList(category string) (list *services.CourseList, err error) {
	limit, _ := getService().CourseCount(category)
	list, err = getService().CourseList(category, "study", 1, limit)
	if err != nil {
		return
	}
	return

}

// CourseInfo 已购课程列表
func CourseInfo(id int) (info *services.CourseInfo, err error) {
	courseDetail, err := getService().CourseDetail("bauhinia", id)
	if err != nil {
		return
	}
	enID := courseDetail.Enid
	info, err = getService().CourseInfo(enID)
	if err != nil {
		return
	}
	return
}
