package app

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// ArticleList 已购课程文章列表
func ArticleList(id int, chapterID string) (list *services.ArticleList, err error) {
	courseDetail, err := getService().CourseDetail(CateCourse, id)
	if err != nil {
		return
	}
	enid := courseDetail.Enid
	count := courseDetail.PublishNum
	page := int(math.Ceil(float64(count) / 30.0))
	maxID := 0
	var lists []services.ArticleIntro
	for i := 0; i < page; i++ {
		list, err = getService().ArticleList(enid, chapterID, maxID)
		if err != nil {
			return
		}
		maxID = list.List[len(list.List)-1].ID
		lists = append(lists, list.List...)
	}
	list.List = lists
	return

}

// ArticleInfo article info
//
// get article token, audio token, media security token etc
func ArticleInfo(id, aid int) (info *services.ArticleInfo, aEnid string, err error) {
	list, err := ArticleList(id, "")
	if err != nil {
		return
	}

	var aids []int

	// get article enid
	for _, p := range list.List {
		aids = append(aids, p.ID)
		if p.ClassID == id && p.ID == aid {
			aEnid = p.Enid
		}
	}
	if !utils.Contains(aids, aid) {
		err = errors.New("找不到该文章 ID，请检查输入是否正确")
		return
	}

	info, err = getService().ArticleInfo(aEnid)
	if err != nil {
		return
	}
	return
}

// ArticleDetail article detail
func ArticleDetail(id, aid int) (detail *services.ArticleDetail, err error) {
	info, aEnid, err := ArticleInfo(id, aid)
	if err != nil {
		return
	}
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, aEnid, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}
