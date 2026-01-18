package app

import (
	"fmt"
	"time"

	"errors"

	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
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
	_, err = config.Instance.SetUser(&u)
	if err != nil {
		return
	}
	err = config.Instance.Save()
	return
}

func LoginByQr() error {
	service := getService()
	if service == nil {
		return errors.New("服务初始化失败")
	}

	token, err := service.LoginAccessToken()
	if err != nil {
		return err
	}
	// fmt.Printf("token:%#v\n", token)
	code, err := service.GetQrcode(token)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	fmt.Println("同时支持「得到App」和「微信」扫码")
	for {
		select {
		case <-ticker.C:
			check, cookie, err := service.CheckLogin(token, code.Data.QrCodeString)
			if err != nil {
				return err
			}
			if check.Data.Status == 1 {
				err = LoginByCookie(cookie)
				if err != nil {
					return err
				}
				fmt.Println("扫码成功")
				return nil
			} else if check.Data.Status == 2 {

				err = errors.New("登录失败，二维码已过期")
				return err
			}
		case <-ticker.C:
		// 10分钟后二维码失效
		case <-time.After(600 * time.Second):
			err = errors.New("登录失败，二维码已过期")
			return err
		}
	}
}

func SwitchAccount(uid string) (err error) {
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
