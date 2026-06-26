package api

import "github.com/gin-gonic/gin"

type envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(200, envelope{Code: 0, Msg: "", Data: data})
}

func fail(c *gin.Context, status int, msg string) {
	c.JSON(status, envelope{Code: 1, Msg: msg, Data: nil})
}
