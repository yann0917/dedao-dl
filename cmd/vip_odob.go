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
	rootCmd.AddCommand(vipOdobCmd)
}

var vipOdobCmd = &cobra.Command{
	Use:     "vip-odob",
	Short:   "获取每天听本书 VIP 信息",
	Long:    `获取每天听本书 VIP 信息。默认展示：用户ID(user.uid)、昵称(user.nickname)、是否会员(user.is_vip)、是否过期(user.is_expire)、会员到期时间(user.end_time)、卡片数量(cards.count)，以及首张卡片的名称(card[0].name)、价格(card[0].price)、价格说明(card[0].price_desc)、订阅状态(card[0].is_subscribed)；使用 --json 可查看完整字段。`,
	Example: "dedao-dl vip-odob -h\n" + "dedao-dl vip-odob\n" + "dedao-dl --json vip-odob",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := app.OdobVIPInfo()
		if err != nil {
			return err
		}
		if outputJSON {
			return printJSON(info)
		}
		return renderOdobVIPTable(info)
	},
}

func renderOdobVIPTable(info *services.OdobVipUser) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"字段", "值"})
	if info == nil {
		table.Append([]string{"结果", "接口返回为空，请使用 --json 排查"})
		table.Render()
		return nil
	}

	table.Append([]string{"用户ID(user.uid)", strconv.Itoa(info.User.UID)})
	table.Append([]string{"昵称(user.nickname)", info.User.Nickname})
	table.Append([]string{"是否会员(user.is_vip)", strconv.FormatBool(info.User.IsVip)})
	table.Append([]string{"是否过期(user.is_expire)", strconv.FormatBool(info.User.IsExpire)})
	table.Append([]string{"会员到期时间(user.end_time)", strconv.FormatInt(info.User.EndTime, 10)})
	table.Append([]string{"卡片数量(cards.count)", strconv.Itoa(len(info.Card))})

	if len(info.Card) > 0 {
		card := info.Card[0]
		table.Append([]string{"首张卡片名称(card[0].name)", card.Name})
		table.Append([]string{"首张卡片价格(card[0].price)", card.Price})
		table.Append([]string{"首张卡片价格说明(card[0].price_desc)", card.PriceDesc})
		table.Append([]string{"首张卡片订阅状态(card[0].is_subscribed)", strconv.Itoa(card.IsSubscribed)})
	}

	table.Render()
	return nil
}
