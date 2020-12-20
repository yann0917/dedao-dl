package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
)

var Cookie string

// Login login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录得到 pc 端 https://www.dedao.cn",
	Long:  `use dedao-dl login to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		app.LoginByCookie(Cookie)
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
	Cookie = defaultCookie
	loginCmd.Flags().StringVarP(&Cookie, "cookie", "c", defaultCookie, "cookie from www.dedao.cn")
}

// LoginedCookies cookie to map for chromedp print pdf
func LoginedCookies() (cookies map[string]string) {
	services.ParseCookies(Cookie, &cookies)
	return
}
