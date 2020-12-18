package app

import (
	"os"
	"strings"

	"github.com/go-rod/rod"
	"github.com/olekukonko/tablewriter"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

//Who get current login user
func Who() {
	activeUser := config.Instance.ActiveUser()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"UID", "姓名", "头像"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(false)
	table.Append([]string{activeUser.UIDHazy, activeUser.Name, activeUser.Avatar})
	table.Render()
}

// LoginByCookie login by cookie
func LoginByCookie(cookie string) {
	var u config.Dedao
	services.ParseCookies(cookie, &u.CookieOptions)
	// save config
	u.CookieStr = cookie
	config.Instance.SetUser(&u)
	config.Instance.Save()
}

// GetCookie get cookie string
func GetCookie() (cookie string) {
	_ = rod.Try(func() {
		cookie = utils.Get(config.BaseURL)
		if !strings.Contains(cookie, "ISID=") {
		}
	})
	return
}
