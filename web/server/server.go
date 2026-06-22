package server

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	webapi "github.com/yann0917/dedao-dl/web/api"
)

type Options struct {
	Host        string
	Port        int
	OpenBrowser bool
}

func Start(opts Options) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	apiGroup := router.Group("/api")
	webapi.RegisterRoutes(apiGroup)

	if err := registerStaticRoutes(router); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	if opts.OpenBrowser {
		go func() {
			time.Sleep(400 * time.Millisecond)
			_ = openBrowser("http://" + addr)
		}()
	}

	return router.Run(addr)
}

func openBrowser(target string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}

	return cmd.Start()
}

func serveFSFile(c *gin.Context, filePath string) {
	data, err := fs.ReadFile(distFS, filePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	contentType := mime.TypeByExtension(path.Ext(filePath))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	c.Data(http.StatusOK, contentType, data)
}
