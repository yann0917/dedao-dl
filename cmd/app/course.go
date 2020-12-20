package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/yann0917/dedao-dl/utils"
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
func CourseList(category string) {
	limit, _ := getService().CourseCount(category)
	list, err := getService().CourseList(category, "study", 1, limit)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "课程名称", "作者", "购买日期", "价格", "学习进度"})
	table.SetAutoWrapText(false)

	for i, p := range list.List {

		table.Append([]string{strconv.Itoa(i),
			strconv.Itoa(p.ID), p.Title, p.Author,
			utils.Unix2String(int64(p.CreateTime)),
			p.Price,
			strconv.Itoa(p.Progress),
		})
	}
	table.Render()
	return
}

// CourseInfo 已购课程列表
func CourseInfo(id int) {
	courseDetail, err := getService().CourseDetail("bauhinia", id)
	if err != nil {
		return
	}
	enID := courseDetail.Enid
	info, err := getService().CourseInfo(enID)
	if err != nil {
		return
	}
	out := os.Stdout
	table := tablewriter.NewWriter(out)

	fmt.Fprint(out, "专栏名称："+info.ClassInfo.Name+"\n")
	fmt.Fprint(out, "专栏作者："+info.ClassInfo.LecturerNameAndTitle+"\n")
	fmt.Fprint(out, "更新进度："+strconv.Itoa(info.ClassInfo.CurrentArticleCount)+
		"/"+strconv.Itoa(info.ClassInfo.PhaseNum)+"\n")
	fmt.Fprint(out, "课程亮点："+info.ClassInfo.Highlight+"\n")
	fmt.Fprintln(out)

	table.SetHeader([]string{"#", "ID", "章节", "更新时间", "是否更新完成"})
	table.SetAutoWrapText(false)

	for i, p := range info.ChapterList {
		isFinished := "❌"
		if p.IsFinished == 1 {
			isFinished = "✔"
		}
		table.Append([]string{strconv.Itoa(i),
			p.IDStr, p.Name,
			utils.Unix2String(int64(p.UpdateTime)),
			isFinished,
		})
	}
	table.Render()
	return
}
