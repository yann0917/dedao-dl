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
	Long:    `获取电子书 VIP 信息。默认展示：用户ID(uid)、昵称(nickname)、是否会员(is_vip)、是否过期(is_expire)、月度配额(month_count)、周度配额(week_count)、总配额(total_count)、已用配额(finished_count)、价格说明(price_desc)、节省金额(save_price)、错误提示(err_tips)；使用 --json 可查看完整字段。`,
	Example: "dedao-dl vip-ebook -h\n" + "dedao-dl vip-ebook\n" + "dedao-dl --json vip-ebook",
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
		table.Append([]string{"结果", "接口返回为空，请使用 --json 排查"})
		table.Render()
		return nil
	}

	table.Append([]string{"用户ID(uid)", strconv.Itoa(info.UID)})
	table.Append([]string{"昵称(nickname)", info.Nickname})
	table.Append([]string{"是否会员(is_vip)", strconv.FormatBool(info.IsVip)})
	table.Append([]string{"是否过期(is_expire)", strconv.FormatBool(info.IsExpire)})
	table.Append([]string{"月度配额(month_count)", strconv.Itoa(info.MonthCount)})
	table.Append([]string{"周度配额(week_count)", strconv.Itoa(info.WeekCount)})
	table.Append([]string{"总配额(total_count)", strconv.Itoa(info.TotalCount)})
	table.Append([]string{"已用配额(finished_count)", strconv.Itoa(info.FinishedCount)})
	table.Append([]string{"价格说明(price_desc)", info.PriceDesc})
	table.Append([]string{"节省金额(save_price)", info.SavePrice})
	table.Append([]string{"错误提示(err_tips)", info.ErrTips})
	table.Render()
	return nil
}
