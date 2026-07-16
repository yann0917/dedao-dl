package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
)

func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Instance.HasActiveUser() {
			fail(c, http.StatusUnauthorized, "请先扫码登录")
			c.Abort()
			return
		}
		c.Next()
	}
}
