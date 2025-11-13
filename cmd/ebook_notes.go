package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
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

	// 创建章节映射，用于根据章节ID匹配章节名称
	chapterMap := make(map[string]string)

	for _, catalog := range detail.CatalogList {
		if catalog.Href != "" {
			// 将 href 作为键，章节文本作为值
			// 例如: href="#chapter_4_4" -> text="第四章 章节标题"
			chapterMap[catalog.Href] = catalog.Text
		}
	}

	// 辅助函数：根据章节ID获取章节名称
	getChapterName := func(sectionID string) string {
		if sectionID == "" {
			return ""
		}

		// 遍历章节映射，进行前缀匹配
		for href, name := range chapterMap {
			// 检查章节ID是否是href的前缀
			// 例如: sectionID="Chapter_3_2", href="Chapter_3_2#sigil_toc_id_3"
			if len(href) >= len(sectionID) && href[:len(sectionID)] == sectionID {
				return name
			}
		}

		// 反向匹配：检查href是否是章节ID的前缀
		for href, name := range chapterMap {
			if len(sectionID) >= len(href) && sectionID[:len(href)] == href {
				return name
			}
		}

		// 包含匹配：检查章节ID是否包含在href中或href包含在章节ID中
		for href, name := range chapterMap {
			if len(sectionID) > 0 && len(href) > 0 {
				// 处理带#的情况
				sectionIDClean := sectionID
				if sectionIDClean[0] == '#' {
					sectionIDClean = sectionIDClean[1:]
				}

				hrefClean := href
				if hrefClean[0] == '#' {
					hrefClean = hrefClean[1:]
				}

				// 完整的前缀匹配，不限制字符数
				if len(sectionIDClean) <= len(hrefClean) && hrefClean[:len(sectionIDClean)] == sectionIDClean {
					return name
				}
				if len(hrefClean) <= len(sectionIDClean) && sectionIDClean[:len(hrefClean)] == hrefClean {
					return name
				}
			}
		}

		// 如果都找不到，返回原始ID
		return sectionID
	}

	// 按章节排序笔记列表 - 按照电子书目录的顺序
	sortedNotes := make([]services.EbookNote, len(noteList.List))
	copy(sortedNotes, noteList.List)

	// 创建章节顺序映射：章节ID/名称 -> 目录顺序
	chapterOrderMap := make(map[string]int)
	for order, catalog := range detail.CatalogList {
		if catalog.Href != "" {
			// 记录完整的href的顺序
			chapterOrderMap[catalog.Href] = order
			// 也记录去掉#前缀的顺序
			if len(catalog.Href) > 1 && catalog.Href[0] == '#' {
				chapterOrderMap[catalog.Href[1:]] = order
			}
			// 记录章节文本的顺序
			if catalog.Text != "" {
				chapterOrderMap[catalog.Text] = order
			}
		}
	}

	// 获取笔记在目录中的顺序
	getNoteOrder := func(note services.EbookNote) int {
		if note.Extra.BookSection != "" {
			// 直接查找章节ID
			if order, exists := chapterOrderMap[note.Extra.BookSection]; exists {
				return order
			}
			// 查找带#的版本
			if order, exists := chapterOrderMap["#"+note.Extra.BookSection]; exists {
				return order
			}
			// 模糊匹配：检查是否是某个href的前缀
			for href, order := range chapterOrderMap {
				if len(href) >= len(note.Extra.BookSection) && href[:len(note.Extra.BookSection)] == note.Extra.BookSection {
					return order
				}
			}
		}
		// 如果找不到章节，检查笔记标题
		if note.NoteTitle != "" {
			if order, exists := chapterOrderMap[note.NoteTitle]; exists {
				return order
			}
		}
		// 如果都找不到，返回一个很大的数，让它在最后
		return 999999
	}

	// 按照章节顺序和创建时间排序
	for i := 0; i < len(sortedNotes)-1; i++ {
		for j := i + 1; j < len(sortedNotes); j++ {
			orderI := getNoteOrder(sortedNotes[i])
			orderJ := getNoteOrder(sortedNotes[j])

			// 主要按章节顺序排序
			if orderI != orderJ {
				if orderI > orderJ {
					sortedNotes[i], sortedNotes[j] = sortedNotes[j], sortedNotes[i]
				}
			} else {
				// 章节顺序相同，按创建时间排序
				if sortedNotes[i].CreateTime > sortedNotes[j].CreateTime {
					sortedNotes[i], sortedNotes[j] = sortedNotes[j], sortedNotes[i]
				}
			}
		}
	}

	// 准备表格数据
	var tableData [][]string
	for _, note := range sortedNotes {
		content := note.NoteLine
		if content == "" {
			content = "无内容"
		}

		// 获取章节信息，优先使用 Extra.BookSection 并转换为章节名称
		chapter := ""
		if note.Extra.BookSection != "" {
			chapter = getChapterName(note.Extra.BookSection)
		}

		if chapter == "" {
			// 尝试通过其他方式匹配章节
			if note.NoteTitle != "" {
				// 检查是否能在目录中找到匹配的章节
				if foundChapter, exists := chapterMap[note.NoteTitle]; exists {
					chapter = foundChapter
				} else {
					// 模糊匹配：检查是否包含目录中的章节名称
					for _, catalog := range detail.CatalogList {
						if catalog.Text != "" && (len(catalog.Text) < 50) { // 避免匹配过长的标题
							// 简单的包含匹配，可根据需要优化
							if len(note.NoteTitle) > len(catalog.Text) &&
								note.NoteTitle[:len(catalog.Text)] == catalog.Text {
								chapter = catalog.Text
								break
							}
						}
					}
					if chapter == "" {
						chapter = note.NoteTitle
					}
				}
			} else if note.Extra.Title != "" {
				// 尝试使用 Extra.Title
				if foundChapter, exists := chapterMap[note.Extra.Title]; exists {
					chapter = foundChapter
				} else {
					chapter = note.Extra.Title
				}
			} else {
				chapter = "未知章节"
			}
		}

		// 格式化创建时间
		createTime := "未知"
		if note.CreateTime > 0 {
			createTime = time.Unix(note.CreateTime, 0).Format("2006-01-02 15:04:05")
		}

		// 获取状态信息
		status := "正常"
		if note.Tips != "" {
			status = "不公开"
		}

		// 格式化互动数据
		interaction := fmt.Sprintf("👍%d 💬%d", note.NotesCount.LikeCount, note.NotesCount.CommentCount)

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

	// 显示详细信息
	fmt.Println("\n详细信息:")
	for i, note := range sortedNotes {
		fmt.Printf("\n--- 笔记 %d ---\n", i+1)

		// 获取章节信息并转换为章节名称
		chapter := ""
		if note.Extra.BookSection != "" {
			chapter = getChapterName(note.Extra.BookSection)
		}
		if chapter != "" {
			fmt.Printf("章节: %s\n", chapter)
		}

		fmt.Printf("笔记内容: %s\n", note.NoteLine)
		if note.Note != "" {
			fmt.Printf("笔记标题: %s\n", note.NoteTitle)
		}
		if note.Extra.BookStartPos > 0 && note.Extra.BookOffset > 0 {
			fmt.Printf("位置: %d-%d\n", note.Extra.BookStartPos, note.Extra.BookOffset)
		}
		fmt.Printf("创建时间: %s\n", time.Unix(note.CreateTime, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("分享链接: %s\n", note.ShareURL)
		if note.Tips != "" {
			fmt.Printf("提示: %s\n", note.Tips)
		}
	}

	return nil
}
