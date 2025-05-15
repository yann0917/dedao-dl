package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/utils"
)

var (
	classID   int
	articleID int
	bookID    int
	compassID int
	topicID   string
)

var courseTypeCmd = &cobra.Command{
	Use:     "cat",
	Short:   "获取课程分类",
	Long:    `使用 dedao-dl cat 获取课程分类`,
	Example: "dedao-dl cat",
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		return courseType()
	},
}

var courseCmd = &cobra.Command{
	Use:     "course",
	Short:   "获取我购买过课程",
	Long:    `使用 dedao-dl course 获取我购买过的课程`,
	Example: "dedao-dl course",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if classID > 0 {
			return courseInfo(classID)

		}
		return courseList(app.CateCourse)
	},
}

var compassCmd = &cobra.Command{
	Use:     "ace",
	Short:   "获取我的锦囊",
	Long:    `使用 dedao-dl ace 获取我的锦囊`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ace",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if compassID > 0 {
			return nil
		}
		return courseList(app.CateAce)
	},
}

var odobCmd = &cobra.Command{
	Use:     "odob",
	Short:   "获取我的听书书架",
	Long:    `使用 dedao-dl odob 获取我的听书书架`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl odob",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if compassID > 0 {
			return nil
		}
		return courseList(app.CateAudioBook)
	},
}

func init() {
	rootCmd.AddCommand(courseTypeCmd)
	rootCmd.AddCommand(courseCmd)
	rootCmd.AddCommand(compassCmd)
	rootCmd.AddCommand(odobCmd)

	courseCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程 ID，获取课程信息")
	compassCmd.PersistentFlags().IntVarP(&compassID, "id", "i", 0, "锦囊 ID")
}

func courseType() (err error) {
	list, err := app.CourseType()
	if err != nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"#", "名称", "统计", "分类标签"})

	for i, p := range list.Data.List {

		table.Append([]string{strconv.Itoa(i), p.Name, strconv.Itoa(p.Count), p.Category})
	}
	table.Render()
	return
}
func courseInfo(id int) (err error) {
	info, err := app.CourseInfo(id)
	if err != nil {
		return
	}

	out := os.Stdout
	table := tablewriter.NewWriter(out)

	fmt.Fprint(out, "专栏名称："+info.ClassInfo.Name+"\n")
	fmt.Fprint(out, "专栏作者："+info.ClassInfo.LecturerNameAndTitle+"\n")
	if info.ClassInfo.PhaseNum == 0 {
		fmt.Fprint(out, "共"+strconv.Itoa(info.ClassInfo.CurrentArticleCount)+"讲\n")
	} else {
		fmt.Fprint(out, "更新进度："+strconv.Itoa(info.ClassInfo.CurrentArticleCount)+
			"/"+strconv.Itoa(info.ClassInfo.PhaseNum)+"\n")
	}
	fmt.Fprint(out, "课程亮点："+info.ClassInfo.Highlight+"\n")
	fmt.Fprintln(out)

	table.Header([]string{"#", "ID", "章节", "讲数", "更新时间", "是否更新完成"})

	if len(info.ChapterList) > 0 {
		for i, p := range info.ChapterList {
			isFinished := "❌"
			if p.IsFinished == 1 {
				isFinished = "✔"
			}
			table.Append([]string{strconv.Itoa(i),
				p.IDStr, p.Name, strconv.Itoa(p.PhaseNum),
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
				p.IDStr, "-", p.Title,
				utils.Unix2String(int64(p.UpdateTime)),
				isFinished,
			})
		}
		if info.HasMoreFlatArticleList {
			fmt.Fprint(out, "⚠️  更多文章请使用 article -i 查看文章列表...\n")
		}
	}
	table.Render()
	return
}

func courseList(category string) (err error) {
	list, err := app.CourseList(category)
	if err != nil {
		return
	}
	total, reading, done, unread := len(list.List), 0, 0, 0

	out := os.Stdout
	table := tablewriter.NewTable(out, tablewriter.WithConfig(tablewriter.Config{
		Row: tw.CellConfig{
			Formatting: tw.CellFormatting{
				AutoWrap:  tw.WrapBreak, // Break words to fit
				Alignment: tw.AlignLeft, // Left-align rows
			},
			ColMaxWidths: tw.CellWidth{Global: 64},
		},
	}))
	table.Header([]string{"#", "ID", "课程名称", "作者", "购买日期", "价格", "学习进度"})

	for i, p := range list.List {
		classID := ""
		switch category {
		case app.CateAce:
			fallthrough
		case app.CateAudioBook:
			fallthrough
		case app.CateEbook:
			classID = strconv.Itoa(p.ID)
		case app.CateCourse:
			classID = strconv.Itoa(p.ClassID)
		}
		table.Append([]string{strconv.Itoa(i),
			classID, p.Title, p.Author,
			utils.Unix2String(int64(p.CreateTime)),
			p.Price,
			strconv.Itoa(p.Progress) + "%",
		})
		if p.Progress == 0 {
			unread++
		} else if p.Progress == 100 {
			done++
		} else {
			reading++
		}
	}

	fmt.Fprintf(out, "\n共 %d 本书, 在读: %d, 读完: %d, 未读: %d\n", total, reading, done, unread)

	table.Render()
	return
}
