package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/config"
)

var ebookNotesCmd = &cobra.Command{
	Use:     "notes",
	Short:   "获取电子书笔记列表",
	Long:    `使用 dedao-dl ebook notes 获取指定电子书的笔记列表`,
	Example: "dedao-dl ebook notes -i 12456",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if bookID == 0 {
			return fmt.Errorf("请提供电子书的ID")
		}
		return getEbookNotesList(bookID)
	},
}

func init() {
	ebookCmd.AddCommand(ebookNotesCmd)
	ebookNotesCmd.PersistentFlags().IntVarP(&bookID, "id", "i", 0, "电子书ID")
}

func getEbookNotesList(id int) error {
	service := config.Instance.ActiveUserService()
	if service == nil {
		return fmt.Errorf("服务未初始化")
	}

	detail, err := app.EbookDetail(id)
	if err != nil {
		return fmt.Errorf("获取电子书详情失败: %v", err)
	}

	enid := detail.Enid
	// 获取电子书笔记列表
	noteList, err := service.EbookNoteList(enid)
	if err != nil {
		return fmt.Errorf("获取电子书笔记列表失败: %v", err)
	}

	if len(noteList.List) == 0 {
		fmt.Println("该电子书暂无笔记")
		return nil
	}

	// 使用公共函数排序笔记列表
	sortedNotes := app.SortNotesByChapter(noteList.List, detail.CatalogList)

	// 使用公共函数创建章节映射器
	getChapterName := app.CreateChapterMapper(detail.CatalogList)

	// 准备表格数据
	var tableData [][]string
	for _, note := range sortedNotes {
		content := note.NoteLine
		if content == "" {
			content = "无内容"
		}

		// 使用公共函数获取章节信息
		chapter := app.GetChapterWithFallback(note, getChapterName, detail.CatalogList)

		// 使用公共函数格式化数据
		createTime := app.FormatNoteTime(note.CreateTime)
		status := app.GetNoteStatus(note)
		interaction := app.FormatNoteInteraction(note)

		tableData = append(tableData, []string{
			chapter,
			content,
			createTime,
			status,
			interaction,
		})
	}

	// 使用表格展示笔记列表，支持章节合并
	out := os.Stdout
	table := tablewriter.NewTable(out,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{
				Merging: tw.CellMerging{Mode: tw.MergeHierarchical},
				Formatting: tw.CellFormatting{
					AutoWrap:  tw.WrapBreak, // 自动换行
					Alignment: tw.AlignLeft, // 左对齐
				},
				ColMaxWidths: tw.CellWidth{Global: 80}, // 设置全局列宽
			},
		}),
	)
	table.Header([]string{
		"章节", "笔记内容", "创建时间", "状态", "互动数据",
	})
	table.Bulk(tableData)

	fmt.Printf("电子书《%s》笔记列表 (共 %d 条笔记):\n\n", detail.OperatingTitle, len(sortedNotes))
	table.Render()

	// // 显示详细信息
	// fmt.Println("\n详细信息:")
	// for i, note := range sortedNotes {
	// 	fmt.Printf("\n--- 笔记 %d ---\n", i+1)

	// 	// 使用公共函数获取章节信息并转换为章节名称
	// 	chapter := app.GetChapterWithFallback(note, getChapterName, detail.CatalogList)
	// 	if chapter != "" {
	// 		fmt.Printf("章节: %s\n", chapter)
	// 	}

	// 	fmt.Printf("笔记内容: %s\n", note.NoteLine)
	// 	if note.Note != "" {
	// 		fmt.Printf("笔记标题: %s\n", note.NoteTitle)
	// 	}
	// 	if note.Extra.BookStartPos > 0 && note.Extra.BookOffset > 0 {
	// 		fmt.Printf("位置: %d-%d\n", note.Extra.BookStartPos, note.Extra.BookOffset)
	// 	}
	// 	fmt.Printf("创建时间: %s\n", app.FormatNoteTime(note.CreateTime))
	// 	fmt.Printf("分享链接: %s\n", note.ShareURL)
	// 	if note.Tips != "" {
	// 		fmt.Printf("提示: %s\n", note.Tips)
	// 	}
	// }

	return nil
}
