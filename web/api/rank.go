package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
)

func registerRankRoutes(group *gin.RouterGroup) {
	rank := group.Group("/rank")
	rank.GET("/base-info", getRankBaseInfo)
	rank.GET("/list", getRankList)
}

func getRankBaseInfo(c *gin.Context) {
	data, err := config.Instance.ActiveUserService().RankSquareBaseInfo()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func getRankList(c *gin.Context) {
	rankType := readQueryInt(c, "rankType", 0)
	if rankType <= 0 {
		fail(c, http.StatusBadRequest, "rankType 参数无效")
		return
	}

	data, err := config.Instance.ActiveUserService().RankListData(rankType)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}
