package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

type audioGroupData struct {
	Outside    *services.OutsideDetail       `json:"outside,omitempty"`
	Group      *services.TopicPkgOdobDetails `json:"group,omitempty"`
	GroupError string                        `json:"groupError,omitempty"`
}

func registerAudioRoutes(group *gin.RouterGroup) {
	audio := group.Group("/audio")
	audio.GET("/detail", getAudioDetail)
	audio.GET("/group", getAudioGroupDetail)
}

func getAudioDetail(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	data, err := config.Instance.ActiveUserService().AudioDetailAlias(enid)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func getAudioGroupDetail(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	service := config.Instance.ActiveUserService()

	outside, err := service.OutsideDetail(enid)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	data := audioGroupData{
		Outside: outside,
	}

	group, err := service.TopicPkgOdobDetails(enid)
	if err != nil {
		data.GroupError = err.Error()
	} else {
		data.Group = group
	}

	ok(c, data)
}
