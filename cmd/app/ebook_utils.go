package app

import (
	"fmt"
	"sort"
	"time"

	"github.com/yann0917/dedao-dl/services"
)

// SortNotesByChapter 按章节排序笔记列表
func SortNotesByChapter(notes []services.EbookNote, catalogList []services.Catelog) []services.EbookNote {
	sortedNotes := make([]services.EbookNote, len(notes))
	copy(sortedNotes, notes)

	// 创建章节顺序映射
	chapterOrderMap := createChapterOrderMap(catalogList)
	getNoteOrder := createNoteOrderFunction(chapterOrderMap)

	// 使用 Go 的 sort.Slice 进行排序
	sort.Slice(sortedNotes, func(i, j int) bool {
		orderI := getNoteOrder(sortedNotes[i])
		orderJ := getNoteOrder(sortedNotes[j])

		if orderI != orderJ {
			return orderI < orderJ
		}
		// 章节顺序相同，按创建时间排序
		return sortedNotes[i].CreateTime < sortedNotes[j].CreateTime
	})

	return sortedNotes
}

// CreateChapterMapper 创建章节映射函数
func CreateChapterMapper(catalogList []services.Catelog) func(string) string {
	chapterMap := make(map[string]string)

	for _, catalog := range catalogList {
		if catalog.Href != "" {
			// 将 href 作为键，章节文本作为值
			// 例如: href="#chapter_4_4" -> text="第四章 章节标题"
			chapterMap[catalog.Href] = catalog.Text
		}
	}

	// 返回章节名称获取函数
	return func(sectionID string) string {
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
}

// GetChapterWithFallback 获取笔记的章节名称，支持多种后备匹配方式
func GetChapterWithFallback(note services.EbookNote, getChapterName func(string) string, catalogList []services.Catelog) string {
	// 优先使用 Extra.BookSection 并转换为章节名称
	chapter := ""
	if note.Extra.BookSection != "" {
		chapter = getChapterName(note.Extra.BookSection)
	}

	if chapter == "" {
		// 尝试通过其他方式匹配章节
		if note.NoteTitle != "" {
			// 检查是否能在目录中找到匹配的章节
			for _, catalog := range catalogList {
				if catalog.Text != "" && catalog.Text == note.NoteTitle {
					chapter = catalog.Text
					break
				}
			}
			if chapter == "" {
				// 模糊匹配：检查是否包含目录中的章节名称
				for _, catalog := range catalogList {
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
			chapter = note.Extra.Title
			// 检查是否能在目录中找到匹配的章节
			for _, catalog := range catalogList {
				if catalog.Text != "" && catalog.Text == chapter {
					break
				}
			}
		} else {
			chapter = "未知章节"
		}
	}

	return chapter
}

// FormatNoteTime 格式化笔记创建时间
func FormatNoteTime(createTime int64) string {
	if createTime > 0 {
		return time.Unix(createTime, 0).Format("2006-01-02 15:04:05")
	}
	return "未知"
}

// GetNoteStatus 获取笔记状态描述
func GetNoteStatus(note services.EbookNote) string {
	if note.Tips != "" {
		return "不公开"
	}
	return "正常"
}

// FormatNoteInteraction 格式化笔记互动数据
func FormatNoteInteraction(note services.EbookNote) string {
	return fmt.Sprintf("👍%d 💬%d", note.NotesCount.LikeCount, note.NotesCount.CommentCount)
}

// createChapterOrderMap 创建章节顺序映射的辅助函数
func createChapterOrderMap(catalogList []services.Catelog) map[string]int {
	chapterOrderMap := make(map[string]int)
	for order, catalog := range catalogList {
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
	return chapterOrderMap
}

// createNoteOrderFunction 创建笔记排序函数的辅助函数
func createNoteOrderFunction(chapterOrderMap map[string]int) func(services.EbookNote) int {
	return func(note services.EbookNote) int {
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
}
