package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

var loginSessions = newQRSessionStore()

type authMeData struct {
	LoggedIn bool           `json:"loggedIn"`
	User     *services.User `json:"user,omitempty"`
}

type authAccountSummary struct {
	UIDHazy string `json:"uidHazy"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Active  bool   `json:"active"`
}

func registerAuthRoutes(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	auth.POST("/qrcode", createQRCodeSession)
	auth.GET("/qrcode/:sessionID/status", pollQRCodeStatus)
	auth.GET("/me", getCurrentUser)
	auth.GET("/accounts", getAccounts)
	auth.POST("/switch", switchAccount)
	auth.POST("/logout", logoutCurrentUser)
}

func createQRCodeSession(c *gin.Context) {
	service := services.NewService(&services.CookieOptions{})

	token, err := service.LoginAccessToken()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	qrCode, err := service.GetQrcode(token)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	if qrCode == nil {
		fail(c, http.StatusBadGateway, "二维码生成失败")
		return
	}

	session := loginSessions.Create(token, qrCode.Data.QrCode, qrCode.Data.QrCodeString)
	ok(c, gin.H{
		"sessionId":    session.ID,
		"qrCode":       qrCode.Data.QrCode,
		"qrCodeString": qrCode.Data.QrCodeString,
		"expiresAt":    session.ExpiresAt.Unix(),
	})
}

func pollQRCodeStatus(c *gin.Context) {
	sessionID := c.Param("sessionID")
	session, exists := loginSessions.Get(sessionID)
	if !exists {
		fail(c, http.StatusNotFound, "登录会话不存在或已过期")
		return
	}

	if time.Now().After(session.ExpiresAt) {
		loginSessions.Delete(sessionID)
		ok(c, gin.H{"status": 2, "expiresAt": session.ExpiresAt.Unix()})
		return
	}

	service := services.NewService(&services.CookieOptions{})
	check, cookie, err := service.CheckLogin(session.Token, session.QRCodeString)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	if check == nil {
		ok(c, gin.H{"status": 0, "expiresAt": session.ExpiresAt.Unix()})
		return
	}

	payload := gin.H{
		"status":    check.Data.Status,
		"expiresAt": session.ExpiresAt.Unix(),
	}

	switch check.Data.Status {
	case 1:
		user, err := loginByCookie(cookie)
		if err != nil {
			fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		payload["user"] = user
		loginSessions.Delete(sessionID)
	case 2:
		loginSessions.Delete(sessionID)
	}

	ok(c, payload)
}

func getCurrentUser(c *gin.Context) {
	if !config.Instance.HasActiveUser() {
		ok(c, authMeData{LoggedIn: false})
		return
	}

	user, err := config.Instance.ActiveUserService().User()
	if err != nil {
		ok(c, authMeData{LoggedIn: false})
		return
	}

	ok(c, authMeData{LoggedIn: true, User: user})
}

func logoutCurrentUser(c *gin.Context) {
	if err := config.Instance.LogoutActiveUser(); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, authMeData{LoggedIn: false})
}

func getAccounts(c *gin.Context) {
	ok(c, buildAccountSummaries())
}

func switchAccount(c *gin.Context) {
	var req struct {
		UIDHazy string `json:"uidHazy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.UIDHazy) == "" {
		fail(c, http.StatusBadRequest, "缺少 uidHazy 参数")
		return
	}

	if err := config.Instance.SwitchUser(&config.User{UIDHazy: strings.TrimSpace(req.UIDHazy)}); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := config.Instance.ActiveUserService().User()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, authMeData{LoggedIn: true, User: user})
}

func loginByCookie(cookie string) (*services.User, error) {
	var userConfig config.Dedao
	cookie = strings.TrimSpace(cookie)
	if err := services.ParseCookies(cookie, &userConfig.CookieOptions); err != nil {
		return nil, err
	}

	if cookie != "" {
		userConfig.CookieStr = cookie
	}

	if _, err := config.Instance.SetUser(&userConfig); err != nil {
		return nil, err
	}
	if err := config.Instance.Save(); err != nil {
		return nil, err
	}

	return config.Instance.ActiveUserService().User()
}

func buildAccountSummaries() []authAccountSummary {
	summaries := make([]authAccountSummary, 0, len(config.Instance.Users))
	for _, user := range config.Instance.Users {
		summaries = append(summaries, authAccountSummary{
			UIDHazy: user.UIDHazy,
			Name:    user.Name,
			Avatar:  user.Avatar,
			Active:  config.Instance.ActiveUID == user.UIDHazy,
		})
	}
	return summaries
}
