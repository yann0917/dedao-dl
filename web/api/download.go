package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dedaoapp "github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/utils"
)

type courseDownloadRequest struct {
	ID           int    `json:"id"`
	EnID         string `json:"enid"`
	ArticleID    int    `json:"articleId"`
	DownloadType int    `json:"downloadType"`
	Title        string `json:"title"`
	IsMerge      bool   `json:"isMerge"`
	IsComment    bool   `json:"isComment"`
	IsOrder      bool   `json:"isOrder"`
}

type audioDownloadRequest struct {
	ID           int    `json:"id"`
	EnID         string `json:"enid"`
	DownloadType int    `json:"downloadType"`
	Title        string `json:"title"`
}

type ebookDownloadRequest struct {
	ID           int    `json:"id"`
	EnID         string `json:"enid"`
	DownloadType int    `json:"downloadType"`
	Title        string `json:"title"`
}

type downloadStartResponse struct {
	SessionID    string `json:"sessionId"`
	Target       string `json:"target"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	DownloadType int    `json:"downloadType"`
	OutputDir    string `json:"outputDir"`
	StreamURL    string `json:"streamUrl"`
}

var downloadSessions = newDownloadSessionStore()

func init() {
	downloadSessions.StartCleanupLoop(10*time.Minute, 30*time.Minute)
}

func registerDownloadRoutes(group *gin.RouterGroup) {
	download := group.Group("/download")
	download.POST("/course", downloadCourse)
	download.POST("/audio", downloadAudio)
	download.POST("/ebook", downloadEbook)
	download.GET("/stream/:sessionID", streamDownloadSession)
}

func downloadCourse(c *gin.Context) {
	var req courseDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "课程下载参数无效")
		return
	}

	req.EnID = strings.TrimSpace(req.EnID)
	if req.ID <= 0 && req.EnID == "" {
		fail(c, http.StatusBadRequest, "课程下载缺少 id 或 enid")
		return
	}
	if !isValidDownloadType(req.DownloadType) {
		fail(c, http.StatusBadRequest, "课程下载类型无效")
		return
	}

	title := resolveDownloadTitle(req.Title, "课程下载")
	outputDir, err := prepareDownloadOutputDir()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	session := downloadSessions.Create("course", title, outputDir)
	downloader := dedaoapp.CourseDownload{
		ID:           req.ID,
		EnID:         req.EnID,
		AID:          req.ArticleID,
		DownloadType: req.DownloadType,
		Title:        title,
		IsMerge:      req.IsMerge,
		IsComment:    req.IsComment,
		IsOrder:      req.IsOrder,
		ProgressCB: func(total, current, pct int, currentName string) {
			session.publish(downloadEvent{
				Type:        "progress",
				Progress:    pct,
				Current:     current,
				Total:       total,
				CurrentName: currentName,
			})
		},
		LogCB: func(level, message string) {
			session.publish(downloadEvent{
				Type:    "log",
				Level:   level,
				Message: message,
			})
		},
	}
	startDownloadAsync(session, outputDir, &downloader)
	ok(c, buildDownloadStartResponse(session, req.DownloadType))
}

func downloadAudio(c *gin.Context) {
	var req audioDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "听书下载参数无效")
		return
	}

	req.EnID = strings.TrimSpace(req.EnID)
	if req.ID <= 0 && req.EnID == "" {
		fail(c, http.StatusBadRequest, "听书下载缺少 id 或 enid")
		return
	}
	if !isValidDownloadType(req.DownloadType) {
		fail(c, http.StatusBadRequest, "听书下载类型无效")
		return
	}

	title := resolveDownloadTitle(req.Title, "听书下载")
	outputDir, err := prepareDownloadOutputDir()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	session := downloadSessions.Create("audio", title, outputDir)
	downloader := dedaoapp.OdobDownload{
		ID:           req.ID,
		EnID:         req.EnID,
		DownloadType: req.DownloadType,
		Title:        title,
		ProgressCB: func(total, current, pct int, currentName string) {
			session.publish(downloadEvent{
				Type:        "progress",
				Progress:    pct,
				Current:     current,
				Total:       total,
				CurrentName: currentName,
			})
		},
		LogCB: func(level, message string) {
			session.publish(downloadEvent{
				Type:    "log",
				Level:   level,
				Message: message,
			})
		},
	}
	startDownloadAsync(session, outputDir, &downloader)
	ok(c, buildDownloadStartResponse(session, req.DownloadType))
}

func downloadEbook(c *gin.Context) {
	var req ebookDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "电子书下载参数无效")
		return
	}

	req.EnID = strings.TrimSpace(req.EnID)
	if req.ID <= 0 && req.EnID == "" {
		fail(c, http.StatusBadRequest, "电子书下载缺少 id 或 enid")
		return
	}
	if !isValidDownloadType(req.DownloadType) {
		fail(c, http.StatusBadRequest, "电子书下载类型无效")
		return
	}

	title := resolveDownloadTitle(req.Title, "电子书下载")
	outputDir, err := prepareDownloadOutputDir()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	session := downloadSessions.Create("ebook", title, outputDir)
	downloader := dedaoapp.EBookDownload{
		ID:           req.ID,
		EnID:         req.EnID,
		DownloadType: req.DownloadType,
		Title:        title,
		ProgressCB: func(total, current, pct int, currentName string) {
			session.publish(downloadEvent{
				Type:        "progress",
				Progress:    pct,
				Current:     current,
				Total:       total,
				CurrentName: currentName,
			})
		},
		LogCB: func(level, message string) {
			session.publish(downloadEvent{
				Type:    "log",
				Level:   level,
				Message: message,
			})
		},
	}
	startDownloadAsync(session, outputDir, &downloader)
	ok(c, buildDownloadStartResponse(session, req.DownloadType))
}

func streamDownloadSession(c *gin.Context) {
	sessionID := strings.TrimSpace(c.Param("sessionID"))
	if sessionID == "" {
		fail(c, http.StatusBadRequest, "缺少 sessionID")
		return
	}

	session, events, history, ok := downloadSessions.Subscribe(sessionID)
	if !ok {
		fail(c, http.StatusNotFound, "下载会话不存在或已过期")
		return
	}
	defer downloadSessions.Unsubscribe(sessionID, events)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()
	pendingEvents := make([]downloadEvent, 0, len(history)+1)
	pendingEvents = append(pendingEvents, session.Snapshot())
	pendingEvents = append(pendingEvents, history...)

	c.Stream(func(w io.Writer) bool {
		if len(pendingEvents) > 0 {
			event := pendingEvents[0]
			pendingEvents = pendingEvents[1:]
			c.SSEvent("", event)
			return true
		}

		select {
		case <-c.Request.Context().Done():
			return false
		case event, ok := <-events:
			if !ok {
				return false
			}
			c.SSEvent("", event)
			return true
		case <-heartbeat.C:
			c.SSEvent("", downloadEvent{
				Type:      "heartbeat",
				SessionID: sessionID,
				Timestamp: time.Now().UnixMilli(),
			})
			return true
		}
	})
}

func isValidDownloadType(downloadType int) bool {
	return downloadType >= 1 && downloadType <= 3
}

func prepareDownloadOutputDir() (string, error) {
	// Keep app and utils output roots aligned so all download branches land in one place.
	outputDir, err := filepath.Abs(dedaoapp.OutputDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	dedaoapp.OutputDir = outputDir
	utils.OutputDir = outputDir
	return outputDir, nil
}

func buildDownloadStartResponse(session *downloadSession, downloadType int) downloadStartResponse {
	snapshot := session.Snapshot()
	return downloadStartResponse{
		SessionID:    session.id,
		Target:       snapshot.Target,
		Title:        snapshot.Title,
		Status:       snapshot.Status,
		DownloadType: downloadType,
		OutputDir:    snapshot.OutputDir,
		StreamURL:    "/api/download/stream/" + session.id,
	}
}

func resolveDownloadTitle(rawTitle string, fallback string) string {
	title := strings.TrimSpace(rawTitle)
	if title == "" {
		return fallback
	}
	return title
}

func startDownloadAsync(session *downloadSession, outputDir string, downloader dedaoapp.DeDaoDownloader) {
	go func() {
		session.publish(downloadEvent{
			Type:      "log",
			Level:     "info",
			Message:   fmt.Sprintf("开始执行下载，输出目录：%s", outputDir),
			OutputDir: outputDir,
		})
		if err := downloader.Download(); err != nil {
			session.publish(downloadEvent{
				Type:      "error",
				Level:     "error",
				Message:   err.Error(),
				OutputDir: outputDir,
			})
			return
		}
		session.publish(downloadEvent{
			Type:        "progress",
			Progress:    100,
			Current:     1,
			Total:       1,
			CurrentName: "下载完成",
			OutputDir:   outputDir,
		})
		session.publish(downloadEvent{
			Type:      "done",
			Level:     "info",
			Message:   "下载完成",
			OutputDir: outputDir,
		})
	}()
}
