package app

import (
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

const (
	// CateCourse 课程
	CateCourse = "bauhinia"
	// CateAudioBook 听书
	CateAudioBook = "odob"
	// CateEbook 电子书
	CateEbook = "ebook"
	// CateAce 锦囊
	CateAce = "compass"
	// CatAll 全部
	CatAll = "all"
)

func getService() *services.Service {
	return config.Instance.ActiveUserService()
}

// LoginedCookies cookie sting to map for chromedp print pdf
func LoginedCookies() (cookies map[string]string) {
	Cookie := config.Instance.ActiveUser().CookieStr
	services.ParseCookies(Cookie, &cookies)
	return
}
