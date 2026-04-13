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
	Long:    `使用 dedao-dl vip-odob 获取 pc/odob/v2/vipuser/vip_card_info 返回结果。`,
	Example: "dedao-dl vip-odob\n" + "dedao-dl --json vip-odob",
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
		table.Append([]string{"result", "接口返回为空，请使用 --json 排查"})
		table.Render()
		return nil
	}

	table.Append([]string{"user.uid", strconv.Itoa(info.User.UID)})
	table.Append([]string{"user.nickname", info.User.Nickname})
	table.Append([]string{"user.is_vip", strconv.FormatBool(info.User.IsVip)})
	table.Append([]string{"user.is_expire", strconv.FormatBool(info.User.IsExpire)})
	table.Append([]string{"user.end_time", strconv.FormatInt(info.User.EndTime, 10)})
	table.Append([]string{"cards.count", strconv.Itoa(len(info.Card))})

	if len(info.Card) > 0 {
		card := info.Card[0]
		table.Append([]string{"card[0].name", card.Name})
		table.Append([]string{"card[0].price", card.Price})
		table.Append([]string{"card[0].price_desc", card.PriceDesc})
		table.Append([]string{"card[0].is_subscribed", strconv.Itoa(card.IsSubscribed)})
	}

	table.Render()
	return nil
}
