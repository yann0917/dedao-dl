package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var (
	channelID int
)

func init() {
	rootCmd.AddCommand(channelCmd)
	channelCmd.AddCommand(channelInfoCmd)
	channelCmd.AddCommand(channelHomepageCmd)
	channelCmd.AddCommand(channelVipCmd)

	channelInfoCmd.Flags().IntVarP(&channelID, "id", "i", 0, "频道ID")
	channelVipCmd.Flags().IntVarP(&channelID, "id", "i", 0, "频道ID")
	channelHomepageCmd.Flags().IntVarP(&channelID, "id", "i", 0, "频道ID")
}

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "学习圈相关操作",
}

var channelInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "获取学习圈信息",
	Example: "dedao-dl channel info --id 1000",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if channelID <= 0 {
			return fmt.Errorf("请使用 --id 指定频道ID")
		}
		info, err := app.ChannelInfo(channelID)
		if err != nil {
			return err
		}
		return printJSON(info)
	},
}

var channelHomepageCmd = &cobra.Command{
	Use:     "homepage",
	Short:   "获取学习圈首页分类",
	Example: "dedao-dl channel homepage --id 1000",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if channelID <= 0 {
			return fmt.Errorf("请使用 --id 指定频道ID")
		}
		cats, err := app.ChannelHomepage(channelID)
		if err != nil {
			return err
		}
		return printJSON(cats)
	},
}

var channelVipCmd = &cobra.Command{
	Use:     "vip",
	Short:   "获取学习圈VIP/权限信息",
	Example: "dedao-dl channel vip --id 1000",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if channelID <= 0 {
			return fmt.Errorf("请使用 --id 指定频道ID")
		}
		info, err := app.ChannelVipInfo(channelID)
		if err != nil {
			return err
		}
		return printJSON(info)
	},
}
