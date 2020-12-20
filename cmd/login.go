package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var cookie string

// Login login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录得到 pc 端 https://www.dedao.cn",
	Long:  `use dedao-dl login to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		app.LoginByCookie(cookie)
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		app.Who()
	},
}

var whoCmd = &cobra.Command{
	Use:     "who",
	Short:   "查看当前登录的用户",
	Long:    `use dedao-dl who to get current login user info`,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		app.Who()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoCmd)

	defaultCookie := app.GetCookie()
	cookie = defaultCookie
	loginCmd.Flags().StringVarP(&cookie, "cookie", "c", defaultCookie, "cookie from www.dedao.cn")
}
