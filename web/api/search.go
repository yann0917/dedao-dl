package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
)

func registerSearchRoutes(group *gin.RouterGroup) {
	search := group.Group("/search")
	search.GET("/hot", getSearchHot)
	search.GET("/suggest", getSearchSuggest)
}

func getSearchHot(c *gin.Context) {
	data, err := config.Instance.ActiveUserService().SearchHot()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}

func getSearchSuggest(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		ok(c, gin.H{"list": []interface{}{}})
		return
	}

	searchType := readQueryInt(c, "searchType", 10)
	data, err := config.Instance.ActiveUserService().SearchSuggest(query, searchType)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}
