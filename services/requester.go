package services

import (
	"io"
)

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

// reqCourseType 请求课程列表
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

// reqCourseType 请求课程介绍
func (s *Service) reqCourseInfo(ID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"detail_id": ID,
		"is_login":  1,
	}).Request("POST", "/pc/bauhinia/pc/class/info")
	return handleHTTPResponse(resp, err)
}

// reqCourseType 请求文章列表
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

// reqCourseType 请求文章详情
func (s *Service) reqArticleDetail(token, sign, appID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"token": token,
		"sign":  sign,
		"appid": appID,
	}).Request("GET", "/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}

// reqCourseType 请求音频详情
func (s *Service) reqAudioByAlias(ids string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]interface{}{
		"ids":            ids,
		"get_extra_data": 1,
	}).Request("POST", "/pc/bauhinia/v1/audio/mutiget_by_alias")
	return handleHTTPResponse(resp, err)
}
