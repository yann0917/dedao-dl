package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

func registerCourseRoutes(group *gin.RouterGroup) {
	course := group.Group("/course")
	course.GET("/categories", listCourseCategories)
	course.GET("/navbar", getCourseNavbar)
	course.GET("/list", listCourses)
	course.GET("/info", getCourseInfo)
	course.GET("/articles", listCourseArticles)
}

func listCourseCategories(c *gin.Context) {
	data, err := config.Instance.ActiveUserService().CourseType()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}

func getCourseNavbar(c *gin.Context) {
	data, err := config.Instance.ActiveUserService().GetNavbar()
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}

func listCourses(c *gin.Context) {
	category := c.DefaultQuery("category", services.CateCourse)
	order := c.DefaultQuery("order", "study")
	filter := c.DefaultQuery("filter", services.CatAll)
	page := readQueryInt(c, "page", 1)
	limit := readQueryInt(c, "limit", 18)
	groupID := readQueryInt(c, "groupId", 0)

	var (
		data *services.CourseList
		err  error
	)
	if groupID > 0 {
		data, err = config.Instance.ActiveUserService().CourseGroupListWithFilter(category, order, filter, groupID, page, limit)
	} else {
		data, err = config.Instance.ActiveUserService().CourseListWithFilter(category, order, filter, page, limit)
	}
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}

func getCourseInfo(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	data, err := config.Instance.ActiveUserService().CourseInfo(enid)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}

func listCourseArticles(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	count := readQueryInt(c, "count", 30)
	maxID := readQueryInt(c, "maxId", 0)
	reverse := c.DefaultQuery("reverse", "0") == "1"
	data, err := config.Instance.ActiveUserService().ArticleListWithOptions(enid, "", count, maxID, reverse)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}
	ok(c, data)
}

func readQueryInt(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
