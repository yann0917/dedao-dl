package services

import (
	"io"
)

// reqUser 请求token
func (s *Service) reqToken() (io.ReadCloser, error) {
	resp, err := s.client.Request("GET", "/ddph/v2/token/create")
	return handleHTTPResponse(resp, err)
}

// reqUser 请求用户信息
func (s *Service) reqUser() (io.ReadCloser, error) {
	resp, err := s.client.Request("GET", "/api/pc/user/info")
	return handleHTTPResponse(resp, err)
}

// reqCourseType 请求首页课程分类列表
func (s *Service) reqCourseType() (io.ReadCloser, error) {
	resp, err := s.client.Request("POST", "/api/hades/v1/index/detail")
	return handleHTTPResponse(resp, err)
}

// reqCourseList 请求课程列表
func (s *Service) reqCourseList(category, order string, page, limit int) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"category":        category,
		"order":           order,
		"filter_complete": 0,
		"page":            page,
		"page_size":       limit,
	}).Request("POST", "/api/hades/v1/product/list")
	return handleHTTPResponse(resp, err)
}

// reqCourseInfo 请求课程介绍
func (s *Service) reqCourseInfo(ID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"detail_id": ID,
		"is_login":  1,
	}).Request("POST", "/pc/bauhinia/pc/class/info")
	return handleHTTPResponse(resp, err)
}

// reqArticleList 请求文章列表
func (s *Service) reqArticleList(ID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"chapter_id":      "",
		"count":           30,
		"detail_id":       ID,
		"include_edge":    false,
		"is_unlearn":      false,
		"max_id":          0,
		"max_order_num":   0,
		"reverse":         false,
		"since_id":        0,
		"since_order_num": 0,
		"unlearn_switch":  false,
	}).Request("POST", "/api/pc/bauhinia/pc/class/purchase/article_list")
	return handleHTTPResponse(resp, err)
}

// reqArticleInfo 请求文章 token
func (s *Service) reqArticleInfo(ID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"detail_id": ID,
	}).Request("POST", "/pc/bauhinia/pc/article/info")
	return handleHTTPResponse(resp, err)
}

// reqArticleDetail 请求文章详情
func (s *Service) reqArticleDetail(token, appID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"token": token,
		// "sign":  sign,
		"appid": appID,
	}).Request("GET", "/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqArticlePoint 请求文章重点
func (s *Service) reqArticlePoint(enid string, pType int) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"article_id_hazy": enid,
		"product_type":    pType,
	}).Request("GET", "/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqAudioByAlias 请求音频详情
func (s *Service) reqAudioByAlias(ids string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"ids":            ids,
		"get_extra_data": 1,
	}).Request("POST", "/pc/bauhinia/v1/audio/mutiget_by_alias")
	return handleHTTPResponse(resp, err)
}

// reqEbookDetail 请求电子书详情
func (s *Service) reqEbookDetail(enid string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"id": enid,
	}).Request("GET", "/pc/ebook2/v1/pc/detail")
	return handleHTTPResponse(resp, err)
}

// reqEbookReadToken 请求电子书阅读 token
func (s *Service) reqEbookReadToken(enid string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"id": enid,
	}).Request("POST", "/api/pc/ebook2/v1/pc/read/token")
	return handleHTTPResponse(resp, err)
}

// reqEbookInfo 请求电子书 info
func (s *Service) reqEbookInfo(token string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"token": token,
	}).Request("GET", "/ebk_web/v1/get_book_info")
	return handleHTTPResponse(resp, err)
}
