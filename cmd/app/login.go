package app

import (
	"strings"

	"github.com/go-rod/rod"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// LoginByCookie login by cookie
func LoginByCookie(cookie string) (err error) {
	var u config.Dedao
	err = services.ParseCookies(cookie, &u.CookieOptions)
	if err != nil {
		return
	}
	// save config
	u.CookieStr = cookie
	config.Instance.SetUser(&u)
	config.Instance.Save()
	return
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
