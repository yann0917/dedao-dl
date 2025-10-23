package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

var (
	classID      int
	articleID    int
	bookID       int
	compassID    int
	topicID      string
	courseGroupID int
	odobGroupID  int
	aceGroupID   int
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
	Example: "dedao-dl course\ndedao-dl course --group-id 12345",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if classID > 0 {
			return courseInfo(classID)
		}
		if courseGroupID > 0 {
			return groupList(app.CateCourse, courseGroupID)
		}
		return courseListFlat(app.CateCourse)
	},
}

var compassCmd = &cobra.Command{
	Use:     "ace",
	Short:   "获取我的锦囊",
	Long:    `使用 dedao-dl ace 获取我的锦囊`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ace\ndedao-dl ace --group-id 12345",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if compassID > 0 {
			return nil
		}
		if aceGroupID > 0 {
			return groupList(app.CateAce, aceGroupID)
		}
		return courseListFlat(app.CateAce)
	},
}

var odobCmd = &cobra.Command{
	Use:     "odob",
	Short:   "获取我的听书书架",
	Long:    `使用 dedao-dl odob 获取我的听书书架`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl odob\ndedao-dl odob --group-id 12345",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if compassID > 0 {
			return nil
		}
		if odobGroupID > 0 {
			return groupList(app.CateAudioBook, odobGroupID)
		}
		return courseListFlat(app.CateAudioBook)
	},
}

func init() {
	rootCmd.AddCommand(courseTypeCmd)
	rootCmd.AddCommand(courseCmd)
	rootCmd.AddCommand(compassCmd)
	rootCmd.AddCommand(odobCmd)

	courseCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程 ID，获取课程信息")
	courseCmd.PersistentFlags().IntVar(&courseGroupID, "group-id", 0, "分组ID，显示指定分组内的课程")

	compassCmd.PersistentFlags().IntVarP(&compassID, "id", "i", 0, "锦囊 ID")
	compassCmd.PersistentFlags().IntVar(&aceGroupID, "group-id", 0, "分组ID，显示指定分组内的锦囊")

	odobCmd.PersistentFlags().IntVar(&odobGroupID, "group-id", 0, "分组ID，显示指定分组内的听书")
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
	table.Header([]string{"#", "ID", "课程名称", "作者", "购买日期", "价格", "学习进度", "备注"})

	for i, p := range list.List {
		classID, remark := "", ""
		switch category {
		case app.CateAce:
			fallthrough
		case app.CateAudioBook:
			if p.Type == 1013 {
				remark = "名家讲书"
			}
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
			remark,
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

// courseListFlat 获取课程列表（扁平化显示，展开所有分组）
func courseListFlat(category string) (err error) {
	list, err := app.CourseList(category)
	if err != nil {
		return err
	}

	// Expand all groups
	allItems, _, expandErr := expandGroups(list.List, category)
	if expandErr != nil && len(allItems) == 0 {
		return expandErr
	}

	// Render using helper
	return renderCourseTable(allItems, category, renderOptions{})
}

// groupList 显示指定分组内的课程列表
func groupList(category string, groupID int) (err error) {
	list, err := app.GetGroupItems(category, groupID)
	if err != nil {
		return err
	}

	// Render using helper with header
	return renderCourseTable(list.List, category, renderOptions{
		header: fmt.Sprintf("📁 分组内容 (Group ID: %d)", groupID),
	})
}

// Helper functions for course list rendering

// renderOptions configures table rendering
type renderOptions struct {
	header string
}

// getItemIdentifier extracts the ID and remark for a course item based on category.
// It handles different content types (ebook, course, audiobook, ace) appropriately.
// 根据类别提取课程项目的ID和备注信息
func getItemIdentifier(item services.Course, category string) (id string, remark string) {
	switch category {
	case app.CateAce:
		fallthrough
	case app.CateAudioBook:
		if item.Type == 1013 {
			remark = "名家讲书"
		}
		fallthrough
	case app.CateEbook:
		id = strconv.Itoa(item.ID)
	case app.CateCourse:
		id = strconv.Itoa(item.ClassID)
	}
	return
}

// renderCourseTable renders a list of courses in table format with statistics.
// It displays items in a formatted table and calculates reading progress statistics.
// 以表格格式渲染课程列表，并显示统计信息
func renderCourseTable(items []services.Course, category string, options renderOptions) error {
	total, reading, done, unread := len(items), 0, 0, 0

	out := os.Stdout

	// Print header if provided
	if options.header != "" {
		fmt.Fprintf(out, "\n%s\n", options.header)
	}

	// Create table
	table := tablewriter.NewTable(out, tablewriter.WithConfig(tablewriter.Config{
		Row: tw.CellConfig{
			Formatting: tw.CellFormatting{
				AutoWrap:  tw.WrapBreak,
				Alignment: tw.AlignLeft,
			},
			ColMaxWidths: tw.CellWidth{Global: 64},
		},
	}))
	table.Header([]string{"#", "ID", "课程名称", "作者", "购买日期", "价格", "学习进度", "备注"})

	// Render rows
	for i, p := range items {
		classID, remark := getItemIdentifier(p, category)
		table.Append([]string{
			strconv.Itoa(i),
			classID,
			p.Title,
			p.Author,
			utils.Unix2String(int64(p.CreateTime)),
			p.Price,
			strconv.Itoa(p.Progress) + "%",
			remark,
		})

		// Track statistics
		if p.Progress == 0 {
			unread++
		} else if p.Progress == 100 {
			done++
		} else {
			reading++
		}
	}

	// Print statistics
	fmt.Fprintf(out, "\n共 %d 本书, 在读: %d, 读完: %d, 未读: %d\n",
		total, reading, done, unread)

	table.Render()
	return nil
}

// expandGroups expands all groups in a list and returns flattened items.
// It processes groups sequentially and provides progress feedback for large operations.
// 展开列表中的所有分组，返回扁平化的项目列表
func expandGroups(items []services.Course, category string) ([]services.Course, int, error) {
	var allItems []services.Course
	var groupCount int
	var failedGroups []string

	for _, item := range items {
		if item.IsGroup {
			groupCount++

			groupItems, err := app.GetGroupItems(category, item.ID)
			if err != nil {
				failedGroups = append(failedGroups, fmt.Sprintf("%s (ID: %d)", item.Title, item.ID))
				fmt.Fprintf(os.Stderr, "⚠️  无法获取分组 %s (ID: %d): %v\n", item.Title, item.ID, err)
				continue
			}

			// Handle empty groups explicitly
			if len(groupItems.List) == 0 {
				fmt.Fprintf(os.Stderr, "ℹ️  分组 %s (ID: %d) 为空\n", item.Title, item.ID)
			}

			allItems = append(allItems, groupItems.List...)
		} else {
			allItems = append(allItems, item)
		}
	}

	// Return partial results with aggregated error
	if len(failedGroups) > 0 {
		return allItems, groupCount, fmt.Errorf("无法获取 %d 个分组: %s", len(failedGroups), strings.Join(failedGroups, ", "))
	}

	return allItems, groupCount, nil
}
