package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

type homePortalData struct {
	HomeData           services.HomeData               `json:"homeData"`
	FreeResources      *services.SunflowerResourceList `json:"freeResources,omitempty"`
	FreeResourcesError string                          `json:"freeResourcesError,omitempty"`
	EbookLabels        *services.SunflowerLabelList    `json:"ebookLabels,omitempty"`
	EbookLabelsError   string                          `json:"ebookLabelsError,omitempty"`
	EbookContent       *services.SunflowerContent      `json:"ebookContent,omitempty"`
	EbookContentError  string                          `json:"ebookContentError,omitempty"`
	CourseLabels       *services.SunflowerLabelList    `json:"courseLabels,omitempty"`
	CourseLabelsError  string                          `json:"courseLabelsError,omitempty"`
	CourseContent      *services.SunflowerContent      `json:"courseContent,omitempty"`
	CourseContentError string                          `json:"courseContentError,omitempty"`
}

func registerHomeRoutes(group *gin.RouterGroup) {
	home := group.Group("/home")
	home.GET("/portal", getHomePortal)
	home.GET("/label-content", getHomeLabelContent)
}

func getHomePortal(c *gin.Context) {
	service := config.Instance.ActiveUserService()

	state, err := service.GetHomeInitialState()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	data := homePortalData{
		HomeData: state.HomeData,
	}

	freeResources, err := service.SunflowerResourceList()
	if err != nil {
		data.FreeResourcesError = err.Error()
	} else {
		data.FreeResources = freeResources
	}

	data.EbookLabels, data.EbookContent, data.EbookLabelsError, data.EbookContentError = loadHomeLabelSection(service, 2, 10)
	data.CourseLabels, data.CourseContent, data.CourseLabelsError, data.CourseContentError = loadHomeLabelSection(service, 4, 4)

	ok(c, data)
}

func getHomeLabelContent(c *gin.Context) {
	nType := readQueryInt(c, "type", 0)
	if nType != 2 && nType != 4 {
		fail(c, http.StatusBadRequest, "type 仅支持 2(电子书) 或 4(课程)")
		return
	}

	page := readQueryInt(c, "page", 0)
	pageSize := readQueryInt(c, "pageSize", defaultHomePageSize(nType))
	enid := c.Query("enid")

	data, err := config.Instance.ActiveUserService().SunflowerLabelContent(enid, nType, page, pageSize)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func loadHomeLabelSection(service *services.Service, nType int, pageSize int) (
	labels *services.SunflowerLabelList,
	content *services.SunflowerContent,
	labelsErr string,
	contentErr string,
) {
	labels, err := service.SunflowerLabelList(nType)
	if err != nil {
		labelsErr = err.Error()
		return
	}

	currentEnid := ""
	if len(labels.List) > 0 {
		currentEnid = labels.List[0].Enid
	}

	content, err = service.SunflowerLabelContent(currentEnid, nType, 0, pageSize)
	if err != nil {
		contentErr = err.Error()
	}

	return
}

func defaultHomePageSize(nType int) int {
	if nType == 2 {
		return 10
	}
	return 4
}
