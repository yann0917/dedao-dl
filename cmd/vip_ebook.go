package cmd

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
)

func init() {
	rootCmd.AddCommand(vipEbookCmd)
}

var vipEbookCmd = &cobra.Command{
	Use:     "vip-ebook",
	Short:   "获取电子书 VIP 信息",
	Long:    `使用 dedao-dl vip-ebook 获取 /api/pc/ebook2/v1/vip/info 返回结果。`,
	Example: "dedao-dl vip-ebook\n" + "dedao-dl --json vip-ebook",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := app.EbookVIPInfo()
		if err != nil {
			return err
		}
		if outputJSON {
			return printJSON(info)
		}
		return renderEbookVIPTable(info)
	},
}

func renderEbookVIPTable(info *services.EbookVIPInfo) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"字段", "值"})
	if info == nil {
		table.Append([]string{"result", "接口返回为空，请使用 --json 排查"})
		table.Render()
		return nil
	}

	table.Append([]string{"uid", strconv.Itoa(info.UID)})
	table.Append([]string{"nickname", info.Nickname})
	table.Append([]string{"is_vip", strconv.FormatBool(info.IsVip)})
	table.Append([]string{"is_expire", strconv.FormatBool(info.IsExpire)})
	table.Append([]string{"month_count", strconv.Itoa(info.MonthCount)})
	table.Append([]string{"week_count", strconv.Itoa(info.WeekCount)})
	table.Append([]string{"total_count", strconv.Itoa(info.TotalCount)})
	table.Append([]string{"finished_count", strconv.Itoa(info.FinishedCount)})
	table.Append([]string{"price_desc", info.PriceDesc})
	table.Append([]string{"save_price", info.SavePrice})
	table.Append([]string{"err_tips", info.ErrTips})
	table.Render()
	return nil
}
