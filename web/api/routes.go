package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup) {
	registerAuthRoutes(group)

	authed := group.Group("")
	authed.Use(requireAuth())
	registerHomeRoutes(authed)
	registerAlgoRoutes(authed)
	registerDownloadRoutes(authed)
	registerEbookRoutes(authed)
	registerAudioRoutes(authed)
	registerArticleRoutes(authed)
	registerUserRoutes(authed)
	registerCourseRoutes(authed)
	registerSearchRoutes(authed)
	registerRankRoutes(authed)
}
