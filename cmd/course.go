package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/utils"
)

var (
	cType     string
	classID   int
	articleID int
	bookID    int
	compassID int
)

var courseTypeCmd = &cobra.Command{
	Use:     "cat",
	Short:   "获取课程分类",
	Long:    `使用 dedao-dl cat 获取课程分类`,
	Example: "dedao-dl cat",
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		courseType()
	},
}

var courseCmd = &cobra.Command{
	Use:     "course",
	Short:   "获取我购买过课程",
	Long:    `使用 dedao-dl course 获取我购买过的课程`,
	Example: "dedao-dl course",
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("classID", classID)
		if classID > 0 {
			courseInfo(classID)
			return
		}
		courseList("bauhinia")
	},
}

var compassCmd = &cobra.Command{
	Use:     "ace",
	Short:   "获取我的锦囊",
	Long:    `使用 dedao-dl ace 获取我的锦囊`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ace",
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if compassID > 0 {
			return
		}
		courseList("compass")
	},
}

var odobCmd = &cobra.Command{
	Use:     "odob",
	Short:   "获取我的听书书架",
	Long:    `使用 dedao-dl odob 获取我的听书书架`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl odob",
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if compassID > 0 {
			return
		}
		courseList("odob")
	},
}

func init() {
	rootCmd.AddCommand(courseTypeCmd)
	rootCmd.AddCommand(courseCmd)
	rootCmd.AddCommand(compassCmd)
	rootCmd.AddCommand(odobCmd)

	courseCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程 ID，获取课程信息")
	compassCmd.PersistentFlags().IntVarP(&compassID, "id", "i", 0, "锦囊 ID")
	// rootCmd.PersistentFlags().StringVarP(&cType, "type", "t", "bauhinia", "课程类型(all, bauhinia, ebook, compass")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func courseType() {
	list, err := app.CourseType()
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
func courseInfo(id int) {
	info, err := app.CourseInfo(id)
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

	if len(info.ChapterList) > 0 {
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
	} else if len(info.FlatArticleList) > 0 {
		isFinished := "❌"
		if info.ClassInfo.IsFinished == 1 {
			isFinished = "✔"
		}
		for i, p := range info.FlatArticleList {
			table.Append([]string{strconv.Itoa(i),
				p.IDStr, p.Title,
				utils.Unix2String(int64(p.UpdateTime)),
				isFinished,
			})
		}
	}
	table.Render()
}

func courseList(category string) {
	list, err := app.CourseList(category)
	if err != nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "课程名称", "作者", "购买日期", "价格", "是否完结", "学习进度"})
	table.SetAutoWrapText(false)

	for i, p := range list.List {
		classFinished := "❌"
		if p.ClassFinished {
			classFinished = "✔"
		}
		table.Append([]string{strconv.Itoa(i),
			strconv.Itoa(p.ID), p.Title, p.Author,
			utils.Unix2String(int64(p.CreateTime)),
			p.Price,
			classFinished,
			strconv.Itoa(p.Progress),
		})
	}
	table.Render()
	return
}
