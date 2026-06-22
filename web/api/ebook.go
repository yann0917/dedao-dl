package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

type ebookDetailData struct {
	Detail     *services.EbookDetail           `json:"detail,omitempty"`
	Notes      *services.EbookNoteListResponse `json:"notes,omitempty"`
	NotesError string                          `json:"notesError,omitempty"`
}

func registerEbookRoutes(group *gin.RouterGroup) {
	ebook := group.Group("/ebook")
	ebook.GET("/detail", getEbookDetail)
}

func getEbookDetail(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	service := config.Instance.ActiveUserService()

	detail, err := service.EbookDetail(enid)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	data := ebookDetailData{
		Detail: detail,
	}

	notes, err := service.EbookNoteList(enid)
	if err != nil {
		data.NotesError = err.Error()
	} else {
		data.Notes = notes
	}

	ok(c, data)
}
