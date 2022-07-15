package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/pkg/errors"
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
			return
		}
	})
	return
}

func LoginByQr() error {
	token, err := getService().LoginAccessToken()
	if err != nil {
		return err
	}
	code, err := getService().GetQrcode(token)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	fmt.Println("同时支持「得到App」和「微信」扫码")
	for {
		select {
		case <-ticker.C:
			check, cookie, err := getService().CheckLogin(token, code.Data.QrCodeString)
			if err != nil {
				return err
			}
			if check.Data.Status == 1 {
				LoginByCookie(cookie)
				fmt.Println("扫码成功")
				return nil
			} else if check.Data.Status == 2 {
				err = errors.New("登录失败，二维码已过期")
				return err
			}
		case <-ticker.C:
		//10分钟后二维码失效
		case <-time.After(600 * time.Second):
			err = errors.New("登录失败，二维码已过期")
			return err
		}
	}
}
