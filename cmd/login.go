package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var cookie string

// Login login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "use `dedao-dl login` to login https://www.dedao.cn",
	Long:  `use dedao-dl login to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		app.LoginByCookie(cookie)
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		app.Who()
	},
}

var whoCmd = &cobra.Command{
	Use:   "who",
	Short: "use `dedao-dl who` to get current login user info",
	Long:  `use dedao-dl who to get current login user info`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Who()
	},
	PreRunE: AuthFunc,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoCmd)

	defaultCookie := app.GetCookie()
	cookie = defaultCookie
	loginCmd.Flags().StringVarP(&cookie, "cookie", "c", defaultCookie, "cookie from www.dedao.cn")
}
