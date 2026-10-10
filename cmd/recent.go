package cmd

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
)

// htmlTagRe 上游 title/last_info/intro 内嵌 <hl>、<h1> 等高亮标签，表格展示前剥除
var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func stripHTML(s string) string {
	if s == "" {
		return s
	}
	return strings.Join(strings.Fields(htmlTagRe.ReplaceAllString(s, " ")), " ")
}

var (
	recentUIDHazy           string
	recentMaxID             int64
	recentPageSize          int
	recentProductType       string
	recentFilterProductType bool
)

func init() {
	rootCmd.AddCommand(recentCmd)
	recentCmd.Flags().StringVar(&recentUIDHazy, "uid-hazy", "", "用户 uid_hazy（默认自动读取当前登录用户）")
	recentCmd.Flags().Int64Var(&recentMaxID, "max-id", 0, "分页游标，默认 0")
	recentCmd.Flags().IntVar(&recentPageSize, "page-size", 20, "每页数量，默认 20")
	recentCmd.Flags().StringVar(&recentProductType, "product-type", "", "产品类型过滤（默认不过滤）")
	recentCmd.Flags().BoolVar(&recentFilterProductType, "filter-product-type", true, "是否按 product_type 过滤")
}

var recentCmd = &cobra.Command{
	Use:     "recent",
	Short:   "查询用户最近学习情况",
	Long:    `查询用户最近学习情况。默认自动使用当前登录用户的 uid_hazy；也可通过 --uid-hazy 指定。`,
	Example: "dedao-dl recent -h\n" + "dedao-dl recent\n" + "dedao-dl recent --page-size 50 --product-type 66\n" + "dedao-dl --json recent",
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := app.Recent(
			recentMaxID,
			recentPageSize,
			recentProductType,
			recentUIDHazy,
			recentFilterProductType,
		)
		if err != nil {
			return err
		}

		if outputJSON {
			return printJSON(resp)
		}
		return renderRecentTable(resp)
	},
}

func renderRecentTable(resp *services.RecentResponse) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"#", "标题", "作者", "类型", "进度", "最近学习"})
	if resp == nil || len(resp.List) == 0 {
		table.Append([]string{"-", "暂无数据", "-", "-", "-", "-"})
		table.Render()
		return nil
	}

	for i, item := range resp.List {
		// 上游 recent 接口字段语义按类型不同：课程类的 progress 恒为 0，
		// 真实学习进度百分比放在 max_progress（与课程列表接口交叉验证过 76/50/100）；
		// 其余类型 progress/max_progress 为真实比例（如听书 18/19）。
		pi := item.ProgressIntro
		progress := ""
		switch {
		case item.TypeName == "课程" && pi.Progress == 0 && pi.MaxProgress > 0:
			progress = strconv.Itoa(pi.MaxProgress) + "%"
		case pi.MaxProgress > 0:
			percent := (pi.Progress*100 + pi.MaxProgress/2) / pi.MaxProgress
			progress = strconv.Itoa(percent) + "%"
		case pi.Intro != "":
			progress = stripHTML(pi.Intro)
		case pi.Progress > 0:
			progress = strconv.Itoa(pi.Progress)
		}
		table.Append([]string{
			strconv.Itoa(i + 1),
			stripHTML(item.Title),
			item.Author,
			item.TypeName,
			progress,
			stripHTML(item.LastInfo),
		})
	}
	table.Render()
	return nil
}
