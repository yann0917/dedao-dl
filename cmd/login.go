package cmd

import (
	"strings"

	"github.com/go-rod/rod"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// Login login
type Login struct {
	phone    string
	password string
	services.CookieOptions
}

// IsByCookie cookie login
func (l *Login) IsByCookie() bool {
	return l.GAT != "" && l.ISID != "" && l.SID != ""
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "use `dedao-dl login` to login https://www.dedao.cn",
	Long:  `use dedao-dl login to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
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
}

// Dedao dedao client
var Dedao *services.Service

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoCmd)

	defaultCookie := getCookie()
	var u config.Dedao
	Dedao = u.New()
	services.ParseCookies(defaultCookie, &u.CookieOptions)
	u.CookieStr = defaultCookie
	config.Instance.SetUser(&u)
	config.Instance.Save()
	loginCmd.Flags().StringVarP(&u.CookieStr, "cookie", "c", defaultCookie, "cookie from www.dedao.cn")
}

func getCookie() (cookie string) {
	_ = rod.Try(func() {
		cookie = utils.Get("https://www.dedao.cn")
		if !strings.Contains(cookie, "ISID=") {
		}
	})
	return
}
