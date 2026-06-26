package api

import (
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

const articleAppID = "1632426125495894021"

type articleDetailData struct {
	Info     *services.ArticleInfo   `json:"info,omitempty"`
	Detail   *services.ArticleDetail `json:"detail,omitempty"`
	Markdown string                  `json:"markdown"`
}

func registerArticleRoutes(group *gin.RouterGroup) {
	article := group.Group("/article")
	article.GET("/detail", getArticleDetail)
}

func getArticleDetail(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	aType := readQueryInt(c, "aType", 1)
	if aType != 1 && aType != 2 {
		fail(c, http.StatusBadRequest, "aType 仅支持 1(课程) 或 2(听书)")
		return
	}

	service := config.Instance.ActiveUserService()

	info, err := service.ArticleInfo(enid, aType)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	detail, err := service.ArticleDetail(info.DdArticleToken, articleAppID)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	var contents []services.Content
	if err := jsoniter.UnmarshalFromString(detail.Content, &contents); err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	ok(c, articleDetailData{
		Info:     info,
		Detail:   detail,
		Markdown: contentsToMarkdown(contents),
	})
}

func contentsToMarkdown(contents []services.Content) (res string) {
	for _, content := range contents {
		switch content.Type {
		case "audio":
			title := strings.TrimRight(content.Title, ".mp3")
			res += getMDHeader(1) + title + "\r\n\r\n"
		case "header":
			content.Text = strings.Trim(content.Text, " ")
			if len(content.Text) > 0 {
				res += getMDHeader(content.Level) + content.Text + "\r\n\r\n"
			}
		case "blockquote":
			texts := strings.Split(content.Text, "\n")
			for _, text := range texts {
				res += "> " + text + "\r\n"
				res += "> \r\n"
			}
			res = strings.TrimRight(res, "> \r\n")
			res += "\r\n\r\n"
		case "paragraph":
			resP, err := paragraphToMarkdown(content.Contents)
			if err != nil {
				return
			}
			res += resP
		case "list":
			resL, err := listToMarkdown(content.Contents)
			if err != nil {
				return
			}
			res += resL
		case "elite":
			res += getMDHeader(2) + "划重点\r\n\r\n" + content.Text + "\r\n\r\n"
		case "image":
			res += "![" + content.URL + "](" + content.URL
			if content.Legend != "" {
				res += " \"" + content.Legend + "\""
			}
			res += ")" + "\r\n\r\n"
		case "label-group":
			res += getMDHeader(2) + "`" + content.Text + "`" + "\r\n\r\n"
		}
	}

	res += "---\r\n"
	return
}

func paragraphToMarkdown(content interface{}) (res string, err error) {
	tmpJSON, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}

	cont := services.Contents{}
	if err = jsoniter.Unmarshal(tmpJSON, &cont); err != nil {
		return
	}

	for _, item := range cont {
		subContent := strings.Trim(item.Text.Content, " ")
		switch item.Type {
		case "text":
			if item.Text.Bold {
				res += " **" + subContent + "** "
			} else if item.Text.Highlight {
				res += " *" + subContent + "* "
			} else {
				res += subContent
			}
		}
	}

	res = strings.Trim(res, " ")
	res = strings.Trim(res, "\r\n")
	res += "\r\n\r\n"
	return
}

func listToMarkdown(content interface{}) (res string, err error) {
	tmpJSON, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}

	var cont []services.Contents
	if err = jsoniter.Unmarshal(tmpJSON, &cont); err != nil {
		return
	}

	for _, group := range cont {
		for _, item := range group {
			subContent := strings.Trim(item.Text.Content, " ")
			switch item.Type {
			case "text":
				if item.Text.Bold {
					res += "* **" + subContent + "** "
				} else if item.Text.Highlight {
					res += "* *" + subContent + "* "
				} else {
					res += "* " + subContent
				}
			}
		}
		res += "\r\n\r\n"
	}

	return
}

func getMDHeader(level int) string {
	headers := map[int]string{
		1: "# ",
		2: "## ",
		3: "### ",
		4: "#### ",
		5: "##### ",
		6: "###### ",
	}

	if value, ok := headers[level]; ok {
		return value
	}

	return ""
}
