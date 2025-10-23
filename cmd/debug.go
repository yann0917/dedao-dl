package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
)

var (
	debugGroupID  int
	debugCategory string
	debugVerbose  bool
)

var debugCmd = &cobra.Command{
	Use:     "debug",
	Short:   "调试工具：查看API原始响应",
	Long:    `使用 dedao-dl debug 查看API原始响应，用于调试分组功能`,
	Example: "dedao-dl debug --category ebook --group-id 0",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		return debugAPI()
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)

	debugCmd.Flags().StringVar(&debugCategory, "category", "ebook", "类别 (ebook/course/odob/ace)")
	debugCmd.Flags().IntVar(&debugGroupID, "group-id", 0, "分组ID (0=顶层列表)")
	debugCmd.Flags().BoolVarP(&debugVerbose, "verbose", "v", false, "显示详细信息")
}

func debugAPI() error {
	fmt.Printf("=== 调试模式：category=%s, group_id=%d ===\n\n", debugCategory, debugGroupID)

	var list *services.CourseList
	var err error

	// Get the course list or group items
	if debugGroupID > 0 {
		fmt.Printf("正在获取分组 %d 的内容...\n", debugGroupID)
		list, err = app.GetGroupItems(debugCategory, debugGroupID)
	} else {
		fmt.Printf("正在获取顶层列表...\n")
		list, err = app.CourseList(debugCategory)
	}

	if err != nil {
		return fmt.Errorf("获取列表失败: %v", err)
	}

	// Save raw JSON to file
	filename := fmt.Sprintf("debug_api_response_%s_group%d.json", debugCategory, debugGroupID)
	jsonBytes, _ := json.MarshalIndent(list, "", "  ")
	if err := os.WriteFile(filename, jsonBytes, 0644); err != nil {
		fmt.Printf("警告: 无法保存文件 %s: %v\n", filename, err)
	} else {
		fmt.Printf("✓ 原始响应已保存到: %s\n\n", filename)
	}

	// Analyze the response
	fmt.Printf("📊 总览\n")
	fmt.Printf("─────────────────────────────────────\n")
	fmt.Printf("总数量: %d\n", len(list.List))
	fmt.Printf("API返回的total: %d\n", list.Total)
	fmt.Printf("has_single_book: %v\n\n", list.HasSingleBook)

	// Count groups and items
	var (
		groupCount     = 0
		inGroupCount   = 0
		standaloneCount = 0
		groups         []string
	)

	for _, item := range list.List {
		if item.IsGroup {
			groupCount++
			groupInfo := fmt.Sprintf("ID=%d, GroupID=%d, Type=%d, Title=%s",
				item.ID, item.GroupID, item.GroupType, item.Title)
			groups = append(groups, groupInfo)
		} else if item.GroupID > 0 {
			inGroupCount++
		} else {
			standaloneCount++
		}
	}

	fmt.Printf("📁 分组统计\n")
	fmt.Printf("─────────────────────────────────────\n")
	fmt.Printf("分组数量: %d\n", groupCount)
	fmt.Printf("独立项目: %d (不在任何分组中)\n", standaloneCount)
	fmt.Printf("分组内项目: %d (group_id > 0)\n\n", inGroupCount)

	if groupCount > 0 {
		fmt.Printf("🔍 发现的分组详情\n")
		fmt.Printf("─────────────────────────────────────\n")
		for i, item := range list.List {
			if !item.IsGroup {
				continue
			}

			fmt.Printf("\n[分组 #%d]\n", i+1)
			fmt.Printf("  ID: %d\n", item.ID)
			fmt.Printf("  GroupID: %d (分组自身的group_id)\n", item.GroupID)
			fmt.Printf("  GroupType: %d\n", item.GroupType)
			fmt.Printf("  IsSelfBuildGroup: %v\n", item.IsSelfBuildGroup)
			fmt.Printf("  标题: %s\n", item.Title)
			fmt.Printf("  作者: %s\n", item.Author)
			fmt.Printf("  LabelID: %d\n", item.LabelID)

			// Check for possible count fields
			fmt.Printf("\n  可能的数量字段:\n")
			fmt.Printf("    CourseNum: %d\n", item.CourseNum)
			fmt.Printf("    PublishNum: %d\n", item.PublishNum)

			if debugVerbose {
				fmt.Printf("\n  完整结构: %+v\n", item)
			}

			// Suggest how to fetch this group's contents
			fmt.Printf("\n  💡 获取此分组内容的命令:\n")
			fmt.Printf("     dedao-dl debug --category %s --group-id %d\n", debugCategory, item.ID)
		}
	}

	// Show items in groups (if any)
	if inGroupCount > 0 {
		fmt.Printf("\n📦 属于分组的项目\n")
		fmt.Printf("─────────────────────────────────────\n")
		for i, item := range list.List {
			if item.GroupID > 0 && !item.IsGroup {
				fmt.Printf("[项目 #%d] ID=%d, GroupID=%d, Title=%s\n",
					i+1, item.ID, item.GroupID, item.Title)
			}
		}
	}

	// Provide next steps
	fmt.Printf("\n\n📝 调查建议\n")
	fmt.Printf("═════════════════════════════════════\n")
	if groupCount > 0 && debugGroupID == 0 {
		fmt.Printf("1. 在网页端点击进入一个分组\n")
		fmt.Printf("2. 在浏览器DevTools中复制新的API请求\n")
		fmt.Printf("3. 或者使用上面建议的命令测试获取分组内容\n")
		fmt.Printf("   (注意: 当前代码可能还不支持group_id参数)\n")
	} else if debugGroupID > 0 {
		fmt.Printf("当前正在查看 group_id=%d 的内容\n", debugGroupID)
		fmt.Printf("检查返回的项目是否都有 group_id=%d\n", debugGroupID)
	}

	fmt.Printf("\n✓ 调试完成\n")
	return nil
}
