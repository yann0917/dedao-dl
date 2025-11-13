package app

import "github.com/yann0917/dedao-dl/services"

// ChannelInfo 获取学习圈频道信息
func ChannelInfo(channelID int) (*services.ChannelInfo, error) {
	s := getService()
	return s.ChannelInfo(channelID)
}

// ChannelHomepage 获取学习圈频道首页分类
func ChannelHomepage(channelID int) ([]services.ChannelHomepageCategory, error) {
	s := getService()
	return s.ChannelHomepage(channelID)
}

// ChannelVipInfo 获取学习圈VIP/权限信息
func ChannelVipInfo(channelID int) (*services.ChannelVipInfo, error) {
	s := getService()
	return s.ChannelVipInfo(channelID)
}
