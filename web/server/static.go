package server

import (
	"io/fs"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/webui"
)

var distFS fs.FS

func registerStaticRoutes(router *gin.Engine) error {
	subFS, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return err
	}
	distFS = subFS

	router.NoRoute(func(c *gin.Context) {
		filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		if _, err := fs.Stat(distFS, filePath); err == nil {
			serveFSFile(c, filePath)
			return
		}

		serveFSFile(c, "index.html")
	})

	return nil
}
