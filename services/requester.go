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

// reqCourseType 请求文章详情
func (s *Service) reqArticleDetail(token, sign, appID string) (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"token": token,
		"sign":  sign,
		"appid": appID,
	}).Request("GET", "/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}
