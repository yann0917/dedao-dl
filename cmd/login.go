package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

// Cookie cookie from https://www.dedao.cn
var Cookie string

// Login login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录得到 pc 端 https://www.dedao.cn",
	Long:  `use dedao-dl login to login https://www.dedao.cn`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Cookie == "" {
			defaultCookie := app.GetCookie()
			if defaultCookie == "" {
				return errors.New("自动获取 cookie 失败，请使用参数设置 cookie")
			}
			Cookie = defaultCookie
		}
		err := app.LoginByCookie(Cookie)
		return err
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		who()
	},
}

var whoCmd = &cobra.Command{
	Use:     "who",
	Short:   "查看当前登录的用户",
	Long:    `use dedao-dl who to get current login user info`,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		who()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoCmd)

	loginCmd.Flags().StringVarP(&Cookie, "cookie", "c", "", "cookie from https://www.dedao.cn")
}

// LoginedCookies cookie sting to map for chromedp print pdf
func LoginedCookies() (cookies map[string]string) {
	services.ParseCookies(Cookie, &cookies)
	return
}

// current login user
func who() {
	activeUser := config.Instance.ActiveUser()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"UID", "姓名", "头像"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(false)
	table.Append([]string{activeUser.UIDHazy, activeUser.Name, activeUser.Avatar})
	table.Render()
}
