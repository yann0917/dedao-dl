package app

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// ArticleList 已购课程文章列表
func ArticleList(id int) (list *services.ArticleList, err error) {
	courseDetail, err := getService().CourseDetail("bauhinia", id)
	if err != nil {
		return
	}
	enid := courseDetail.Enid
	list, err = getService().ArticleList(enid)
	if err != nil {
		return
	}

	return

}

// ArticleInfo article info
//
// get article token, audio token, media security token etc
func ArticleInfo(id, aid int) (info *services.ArticleInfo, err error) {
	courseDetail, err := getService().CourseDetail("bauhinia", id)
	if err != nil {
		return
	}

	// get course enid
	enid := courseDetail.Enid
	list, err := getService().ArticleList(enid)
	if err != nil {
		return
	}

	aids := []int{}

	// get article enid
	var articlleIntro services.ArticleIntro
	for _, p := range list.List {
		aids = append(aids, p.ID)
		if p.ClassID == id && p.ID == aid {
			articlleIntro = p
		}
	}

	if !utils.Contains(aids, aid) {
		err = errors.New("找不到该文章 ID，请检查输入是否正确")
		return
	}

	info, err = getService().ArticleInfo(articlleIntro.Enid)
	if err != nil {
		return
	}
	return
}

// ArticleDetail article detail
func ArticleDetail(id, aid int) (err error) {
	info, err := ArticleInfo(id, aid)
	if err != nil {
		return
	}
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err := getService().ArticleDetail(token, info.ClassInfo.Enid, appid)
	if err != nil {
		return
	}
	out := os.Stdout
	table := tablewriter.NewWriter(out)

	var content services.Content
	jsoniter.UnmarshalFromString(detail.Content, &content)
	fmt.Fprint(out, content.Plaintext)
	fmt.Fprintln(out)
	table.Render()
	return
}
