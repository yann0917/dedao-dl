package cmd

import (
	"os"
	"strconv"

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
	Long:  `使用 dedao-dl login to login https://www.dedao.cn`,
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
	Long:    `使用 dedao-dl who 当前登录的用户信息`,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		who()
	},
}

var usersCmd = &cobra.Command{
	Use:     "users",
	Short:   "查看登录过的用户列表",
	Long:    `使用 dedao-dl users 查看登录过的用户列表`,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		users()
	},
}

var suCmd = &cobra.Command{
	Use:     "su",
	Short:   "切换登录过的账号",
	Long:    `使用 dedao-dl su 切换当前登录的账号`,
	Args:    cobra.ExactArgs(1),
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("please input UID")
		}
		uid := args[0]
		err := switchAccount(uid)
		return err
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		who()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(suCmd)
	loginCmd.Flags().StringVarP(&Cookie, "cookie", "c", "", "cookie from https://www.dedao.cn")
}

// LoginedCookies cookie sting to map for chromedp print pdf
func LoginedCookies() (cookies map[string]string) {
	Cookie := config.Instance.ActiveUser().CookieStr
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

// users get login user list
func users() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "UID", "姓名", "头像"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(false)
	for i, user := range config.Instance.Users {
		table.Append([]string{strconv.Itoa(i), user.UIDHazy, user.Name, user.Avatar})
	}
	table.Render()

}

func switchAccount(uid string) (err error) {
	if config.Instance.LoginUserCount() == 0 {
		err = errors.New("cannot found account's")
		return
	}
	err = config.Instance.SwitchUser(&config.User{UIDHazy: uid})

	if err != nil {
		return err
	}
	return
}
