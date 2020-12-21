package app

import (
	"strings"

	"github.com/go-rod/rod"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

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
