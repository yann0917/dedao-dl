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
		"category":  category,
		"order":     order,
		"page":      page,
		"page_size": limit,
	}).Request("POST", "/api/hades/v1/product/list")
	return handleHTTPResponse(resp, err)
}

// reqCourseType 请求首页课程分类列表
func (s *Service) reqArticleDetail() (io.ReadCloser, error) {
	resp, err := s.client.SetData(map[string]string{
		"token": "KWn/CP3W2txbAhtO3K0USs1+F4U4kG+D0Y4siJEj5ScUK/80XS61f8byPlarazHnAYdCfIc5uitXQqvXl/M+r77asns=",
		"sign":  "b23a426b357d1b83",
		"appid": "1632426125495894021",
	}).Request("GET", "/pc/ddarticle/v1/article/get/v2")
	return handleHTTPResponse(resp, err)
}
