package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
)

func registerChannelRoutes(group *gin.RouterGroup) {
	channel := group.Group("/channel")
	channel.GET("/info", getChannelInfo)
	channel.GET("/homepage", getChannelHomepage)
	channel.GET("/vip", getChannelVipInfo)
	channel.GET("/topic-detail", getChannelTopicDetail)
}

func getChannelInfo(c *gin.Context) {
	channelID := readQueryInt(c, "channelId", 0)
	if channelID <= 0 {
		fail(c, http.StatusBadRequest, "channelId 参数无效")
		return
	}

	data, err := config.Instance.ActiveUserService().ChannelInfo(channelID)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func getChannelHomepage(c *gin.Context) {
	channelID := readQueryInt(c, "channelId", 0)
	if channelID <= 0 {
		fail(c, http.StatusBadRequest, "channelId 参数无效")
		return
	}

	data, err := config.Instance.ActiveUserService().ChannelHomepage(channelID)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func getChannelVipInfo(c *gin.Context) {
	channelID := readQueryInt(c, "channelId", 0)
	if channelID <= 0 {
		fail(c, http.StatusBadRequest, "channelId 参数无效")
		return
	}

	data, err := config.Instance.ActiveUserService().ChannelVipInfo(channelID)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func getChannelTopicDetail(c *gin.Context) {
	topicID := readQueryInt(c, "topicId", 0)
	if topicID <= 0 {
		fail(c, http.StatusBadRequest, "topicId 参数无效")
		return
	}

	data, err := config.Instance.ActiveUserService().ChannelTopicDetail(topicID)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}
