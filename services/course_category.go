package services

const (
	// CateCourse 课程
	CateCourse = "bauhinia"
	// CateAudioBook 听书
	CateAudioBook = "odob"
	// CateEbook 电子书
	CateEbook = "ebook"
	// CateAce 锦囊
	CateAce = "compass"
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

	body, err := s.reqCourseType()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}

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
