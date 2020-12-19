package services

import (
	"time"
)

const (
	// CateCourse 课程
	CateCourse = "bauhinia"
	// CateAudioBook 听书
	CateAudioBook = "odob"
	// CateEbook 电子书
	CateEbook = "ebook"
	// CateLecture 讲座
	CateLecture = "navigator"
	// CateAce 锦囊
	CateAce = "compass"
	// CateOther 其他
	CateOther = "other"
	// CatAll 全部
	CatAll = "all"
)

// CourseCategory course category metadata
type CourseCategory struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Category string `json:"category"`
}

// CourseCategoryList course type list
type CourseCategoryList struct {
	Data struct {
		List      []CourseCategory `json:"list"`
		IsShowURL bool             `json:"is_show_url"`
		PCURL     string           `json:"pc_url"`
	} `json:"data"`
}

// CourseType get course type list
func (s *Service) CourseType() (list *CourseCategoryList, err error) {
	cacheFile := "courseTypeList"

	x, ok := list.getCache(cacheFile)
	if ok {
		list = x.(*CourseCategoryList)
		return
	}
	body, err := s.reqCourseType()
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &list); err != nil {
		return
	}

	list.setCache(cacheFile)
	return
}

// CourseCount 获取课程数量 by 分类
func (s *Service) CourseCount(category string) (count int, err error) {
	list, err := s.CourseType()
	if err != nil {
		return
	}
	for _, v := range list.Data.List {
		if v.Category == category {
			count = v.Count
			return
		}
	}
	return
}

func (c *CourseCategoryList) getCacheKey() string {
	return "courseType"
}

func (c *CourseCategoryList) getCache(fileName string) (interface{}, bool) {
	err := Cache.LoadFile(cacheDir + fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *CourseCategoryList) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := Cache.SaveFile(cacheDir + fileName)
	return err
}
