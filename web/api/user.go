package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

func registerUserRoutes(group *gin.RouterGroup) {
	user := group.Group("/user")
	user.GET("/info", getUserInfo)
	user.GET("/center", getUserCenter)
}

func getUserInfo(c *gin.Context) {
	info, err := config.Instance.ActiveUserService().User()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, info)
}

type userCenterData struct {
	User          *services.User         `json:"user,omitempty"`
	EbookVIP      *services.EbookVIPInfo `json:"ebookVip,omitempty"`
	EbookVIPError string                 `json:"ebookVipError,omitempty"`
	OdobVIP       *services.OdobVipUser  `json:"odobVip,omitempty"`
	OdobVIPError  string                 `json:"odobVipError,omitempty"`
	Accounts      []authAccountSummary   `json:"accounts"`
}

func getUserCenter(c *gin.Context) {
	service := config.Instance.ActiveUserService()
	data := userCenterData{
		Accounts: buildAccountSummaries(),
	}

	user, err := service.User()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	data.User = user

	ebookVIP, err := service.EbookVIPInfo()
	if err != nil {
		data.EbookVIPError = err.Error()
	} else {
		data.EbookVIP = ebookVIP
	}

	odobVIP, err := service.OdobVIPInfo()
	if err != nil {
		data.OdobVIPError = err.Error()
	} else {
		data.OdobVIP = odobVIP
	}

	ok(c, data)
}
