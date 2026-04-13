package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var (
	searchQuery string
	searchType  int
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "搜索建议",
	Long:    `使用 dedao-dl search 获取搜索建议结果`,
	Args:    cobra.NoArgs,
	Example: `dedao-dl search --query "基层中国的运行逻辑" --type 0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(searchQuery) == "" {
			return errors.New("query 不能为空，请使用 --query 传入搜索关键词")
		}
		return searchSuggest()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().StringVarP(&searchQuery, "query", "q", "", "搜索关键词")
	searchCmd.PersistentFlags().IntVarP(&searchType, "type", "t", 0, "搜索类型，默认 0")
	_ = searchCmd.MarkPersistentFlagRequired("query")
}

func searchSuggest() (err error) {
	resp, err := app.SearchSuggest(searchQuery, searchType)
	if err != nil {
		return
	}

	table := tablewriter.NewTable(os.Stdout, tablewriter.WithConfig(tablewriter.Config{
		Row: tw.CellConfig{
			Formatting: tw.CellFormatting{
				AutoWrap:  tw.WrapBreak,
				Alignment: tw.AlignLeft,
			},
			ColMaxWidths: tw.CellWidth{Global: 45},
		},
	}))
	table.Header([]string{"#", "TabType", "分类", "ID", "ENID", "标题", "作者", "简介"})

	row := 0
	for _, group := range resp.List {
		for _, item := range group.List {
			category := item.Tname
			if category == "" {
				category = group.TrackName
			}
			table.Append([]string{
				strconv.Itoa(row),
				strconv.Itoa(group.TabType),
				category,
				strconv.Itoa(item.ID),
				shortenText(strings.TrimSpace(item.Extra.Enid), 12),
				cleanSearchTitle(item.Title),
				item.Author,
				strings.TrimSpace(item.Content),
			})
			row++
		}
	}

	if row == 0 {
		fmt.Println("未找到搜索结果")
		return nil
	}

	table.Render()
	return
}

func cleanSearchTitle(s string) string {
	replacer := strings.NewReplacer("<hl>", "", "</hl>", "")
	return strings.TrimSpace(replacer.Replace(s))
}

func shortenText(s string, keep int) string {
	runes := []rune(strings.TrimSpace(s))
	if keep <= 0 || len(runes) <= keep {
		return string(runes)
	}
	return string(runes[:keep]) + "..."
}
