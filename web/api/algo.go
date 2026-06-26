package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

func registerAlgoRoutes(group *gin.RouterGroup) {
	algo := group.Group("/algo")
	algo.GET("/filter", getAlgoFilter)
	algo.GET("/products", getAlgoProducts)
}

func getAlgoFilter(c *gin.Context) {
	data, err := config.Instance.ActiveUserService().AlgoFilter(readAlgoFilterParam(c))
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func getAlgoProducts(c *gin.Context) {
	data, err := config.Instance.ActiveUserService().AlgoProduct(readAlgoFilterParam(c))
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, data)
}

func readAlgoFilterParam(c *gin.Context) services.AlgoFilterParam {
	return services.AlgoFilterParam{
		ClassfcName:  c.DefaultQuery("classfcName", "全部"),
		LabelId:      c.Query("labelId"),
		NavType:      readQueryInt(c, "navType", 0),
		NavigationId: c.Query("navigationId"),
		Page:         readQueryInt(c, "page", 0),
		PageSize:     readQueryInt(c, "pageSize", 18),
		ProductTypes: c.DefaultQuery("productTypes", "66"),
		RequestId:    c.Query("requestId"),
		SortStrategy: c.DefaultQuery("sortStrategy", "HOT"),
		TagsIds:      []interface{}{},
	}
}
