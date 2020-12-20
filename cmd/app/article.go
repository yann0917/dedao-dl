package app

import (
	"fmt"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// ArticleList 已购课程文章列表
func ArticleList(id int) {
	courseDetail, err := getService().CourseDetail("bauhinia", id)
	if err != nil {
		return
	}
	enid := courseDetail.Enid
	list, err := getService().ArticleList(enid)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "课程名称", "更新时间", "是否阅读"})
	table.SetAutoWrapText(false)

	for i, p := range list.List {
		isRead := "❌"
		if p.IsRead {
			isRead = "✔"
		}

		table.Append([]string{strconv.Itoa(i),
			p.IDStr, p.Title,
			utils.Unix2String(int64(p.UpdateTime)),
			isRead,
		})
	}
	table.Render()
	return
}

// ArticleInfo article info
//
// get article token, audio token, media security token etc
func articleInfo(id, aid int) (info *services.ArticleInfo, err error) {
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

	// get article enid
	var articlleIntro services.ArticleIntro
	for _, p := range list.List {
		if p.ClassID == id && p.ID == aid {
			articlleIntro = p
		}
	}

	info, err = getService().ArticleInfo(articlleIntro.Enid)
	if err != nil {
		return
	}
	return
}

// ArticleDetail article detail
func ArticleDetail(id, aid int) {
	info, err := articleInfo(id, aid)
	if err != nil {
		return
	}
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err := getService().ArticleDetail(token, info.ClassInfo.Enid, appid)

	out := os.Stdout
	table := tablewriter.NewWriter(out)

	var content services.Content
	jsoniter.UnmarshalFromString(detail.Content, &content)
	fmt.Fprint(out, content.Plaintext)
	fmt.Fprintln(out)
	table.Render()
	return
}
