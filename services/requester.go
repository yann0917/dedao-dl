package services

import (
	"fmt"
	"io"
	"strings"
)

// reqGetLoginAccessToken 扫码请求token
func (s *Service) reqGetLoginAccessToken() (string, error) {
	// request index get csrf-token
	index, err := s.client.R().Get("")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return "", err
	}

	setCookie := index.GetHeaderValues("Set-Cookie")
	cookies := strings.Split(strings.Join(setCookie, "; "), "; ")
	csrfToken := ""
	for _, v := range cookies {
		item := strings.Split(v, "=")
		if len(item) > 1 && item[0] == "csrfToken" {
			csrfToken = item[1]
			break
		}
	}

	resp, err := s.client.R().
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Xi-Csrf-Token", csrfToken).
		SetHeader("Xi-DT", "web").
		Post("/loginapi/getAccessToken")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return "", err
	}
	accessToken, err := resp.ToString()
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return "", err
	}
	return accessToken, err
}

// reqGetQrcode 扫码登录二维码
// token: X-Oauth-Access-Token from /loginapi/getAccessToken
func (s *Service) reqGetQrcode(token string) (qr *QrCodeResp, err error) {
	_, err = s.client.R().
		SetHeader("X-Oauth-Access-Token", token).
		SetResult(&qr).
		Get("/oauth/api/embedded/qrcode")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return
	}
	return
}

// reqCheckLogin 轮询扫码登录结果
// token: X-Oauth-Access-Token from /loginapi/getAccessToken
// qrCode: qrCodeString from /oauth/api/embedded/qrcode
func (s *Service) reqCheckLogin(token, qrCode string) (check *CheckLoginResp, cookie string, err error) {
	resp, err := s.client.R().
		SetHeader("X-Oauth-Access-Token", token).
		SetBody(map[string]interface{}{
			"keepLogin": true,
			"pname":     "igetoauthpc",
			"qrCode":    qrCode,
			"scene":     "registerlogin",
		}).
		SetResult(&check).
		Post("/oauth/api/embedded/qrcode/check_login")
	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return
	}
	cookies := resp.GetHeaderValues("Set-Cookie")
	cookie = strings.Join(cookies, "; ")
	return
}

// reqUser 请求token
func (s *Service) reqToken() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Get("/ddph/v2/token/create")
	return handleHTTPResponse(resp, err)
}

// reqUser 请求用户信息
func (s *Service) reqUser() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Get("/api/pc/user/info")

	return handleHTTPResponse(resp, err)
}

// reqCourseType 请求首页课程分类列表
func (s *Service) reqCourseType() (io.ReadCloser, error) {
	resp, err := s.client.R().Post("/api/hades/v1/index/detail")
	return handleHTTPResponse(resp, err)
}

// reqCourseList 请求课程列表
func (s *Service) reqCourseList(category, order string, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.R().SetBodyJsonMarshal(map[string]interface{}{
		"category":        category,
		"order":           order,
		"filter_complete": 0,
		"page":            page,
		"page_size":       limit,
	}).Post("/api/hades/v1/product/list")
	return handleHTTPResponse(resp, err)
}

// reqCourseInfo 请求课程介绍
func (s *Service) reqCourseInfo(ID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"detail_id": ID,
			"is_login":  1,
		}).
		Post("/pc/bauhinia/pc/class/info")
	return handleHTTPResponse(resp, err)
}

// reqArticleList 请求文章列表
// chapterID = "" 获取所有的文章列表，否则只获取该章节的文章列表
func (s *Service) reqArticleList(ID, chapterID string, maxID int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"chapter_id":      chapterID,
			"count":           30,
			"detail_id":       ID,
			"include_edge":    false,
			"is_unlearn":      false,
			"max_id":          maxID,
			"max_order_num":   0,
			"reverse":         false,
			"since_id":        0,
			"since_order_num": 0,
			"unlearn_switch":  false,
		}).Post("/api/pc/bauhinia/pc/class/purchase/article_list")
	return handleHTTPResponse(resp, err)
}

// reqArticleCommentList 请求文章热门留言列表
// enId 文章 ID
// sort like-最热 create-最新
func (s *Service) reqArticleCommentList(enId, sort string, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"detail_enid":  enId,
			"note_type":    2,
			"only_replied": false,
			"page":         page,
			"page_count":   limit,
			"sort_by":      sort,
			"source_type":  65,
		}).Post("/pc/ledgers/notes/article_comment_list")
	return handleHTTPResponse(resp, err)
}

// reqArticleInfo 请求文章 token
// id article id or odob audioAliasID
func (s *Service) reqArticleInfo(ID string, aType int) (io.ReadCloser, error) {
	param := make(map[string]string)
	switch aType {
	case 1:
		param["detail_id"] = ID
	case 2:
		param["audio_alias_id"] = ID
	}
	resp, err := s.client.R().
		SetBodyJsonMarshal(param).Post("/pc/bauhinia/pc/article/info")
	return handleHTTPResponse(resp, err)
}

// reqArticleDetail 请求文章详情
func (s *Service) reqArticleDetail(token, appID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"token":  token,
			"appid":  appID,
			"is_new": "1",
		}).
		Get("/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqArticlePoint 请求文章重点
func (s *Service) reqArticlePoint(enid string, pType string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"article_id_hazy": enid,
			"product_type":    pType,
		}).Get("/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqAudioByAlias 请求音频详情
func (s *Service) reqAudioByAlias(ids string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"ids":            ids,
			"get_extra_data": 1,
		}).
		Post("/pc/bauhinia/v1/audio/mutiget_by_alias")
	return handleHTTPResponse(resp, err)
}

// reqEbookDetail 请求电子书详情
func (s *Service) reqEbookDetail(enid string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParam("id", enid).
		Get("/pc/ebook2/v1/pc/detail")

	return handleHTTPResponse(resp, err)
}

// reqEbookReadToken 请求电子书阅读 token
func (s *Service) reqEbookReadToken(enid string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]string{
			"id": enid,
		}).
		Post("/api/pc/ebook2/v1/pc/read/token")
	return handleHTTPResponse(resp, err)
}

// reqEbookInfo 请求电子书 info
func (s *Service) reqEbookInfo(token string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetQueryParam("token", token).
		Get("/ebk_web/v1/get_book_info")
	return handleHTTPResponse(resp, err)
}

// reqEbookPages 获取页面详情
func (s *Service) reqEbookPages(chapterID, token string, index, count, offset int) (io.ReadCloser, error) {
	// A4纸的尺寸为210mm*297mm(宽*高)，这是国际化组织的ISO216定义的标准尺寸。
	// 设定的分辨率是72像素/英寸时，A4纸的尺寸的图像的像素是595×842，在ps中的大小为1.43M。
	// 设定的分辨率是150像素/英寸时，A4纸的尺寸的图像的像素是1240×1754，在ps中的大小为6.22M。
	// 设定的分辨率是300像素/英寸时，A4纸的尺寸的图像的像素是2479×3508，在ps中的大小为24.9M。
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"chapter_id":  chapterID,
			"count":       count,
			"index":       index,
			"offset":      offset,
			"orientation": 1,
			"config": map[string]interface{}{
				"density":         1,
				"direction":       0,
				"font_name":       "pingfang",
				"font_scale":      1,
				"font_size":       20,
				"height":          90000,
				"line_height":     "2em",
				"margin_bottom":   30,
				"margin_left":     30,
				"margin_right":    30,
				"margin_top":      0,
				"paragraph_space": "1em",
				"platform":        1,
				"width":           40000,
			},
			"token": token,
		}).
		Post("/ebk_web/v1/get_pages")
	return handleHTTPResponse(resp, err)
}

// reqEbookInfo 请求电子书vip info
func (s *Service) reqEbookVIPInfo() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Post("/api/pc/ebook2/v1/vip/info")
	return handleHTTPResponse(resp, err)
}

// reqOdobVIPInfo 请求每天听本书书 vip info
func (s *Service) reqOdobVIPInfo() (io.ReadCloser, error) {
	resp, err := s.client.R().
		Post("pc/odob/v2/vipuser/vip_card_info")

	return handleHTTPResponse(resp, err)
}

// reqTopicAll 请求推荐话题列表
func (s *Service) reqTopicAll(page, limit int, all bool) (io.ReadCloser, error) {
	r := s.client.R()
	if !all {
		r = r.SetBodyJsonMarshal(map[string]int{
			"page_id": page,
			"limit":   limit,
		})
	}
	resp, err := r.Post("/pc/ledgers/topic/all")
	return handleHTTPResponse(resp, err)
}

// reqTopicAll 请求话题详情
func (s *Service) reqTopicDetail(topicID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"incr_view_count": true,
			"topic_id_hazy":   topicID,
		}).Post("/pc/ledgers/topic/detail")
	return handleHTTPResponse(resp, err)
}

// reqTopicNotesList 请求话题笔记列表
func (s *Service) reqTopicNotesList(topicID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetBodyJsonMarshal(map[string]interface{}{
			"count":         40,
			"is_elected":    true,
			"page_id":       0,
			"version":       2,
			"topic_id_hazy": topicID,
		}).Post("/pc/ledgers/topic/notes/list")
	return handleHTTPResponse(resp, err)
}
