package app

import "github.com/yann0917/dedao-dl/services"

// User 获取当前登录用户信息
func User() (*services.User, error) {
	return getService().User()
}

// EbookVIPInfo 获取电子书 VIP 信息
func EbookVIPInfo() (*services.EbookVIPInfo, error) {
	return getService().EbookVIPInfo()
}

// OdobVIPInfo 获取每天听本书 VIP 信息
func OdobVIPInfo() (*services.OdobVipUser, error) {
	return getService().OdobVIPInfo()
}
