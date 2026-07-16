package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"syscall"
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
	router.Use(gin.Recovery())

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

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case err := <-serverErrCh:
		return err
	case <-signalCtx.Done():
		fmt.Println("收到退出信号，正在优雅关闭 Web 服务...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown web server failed: %w", err)
		}

		if err := <-serverErrCh; err != nil {
			return err
		}

		fmt.Println("Web 服务已安全关闭")
		return nil
	}
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
