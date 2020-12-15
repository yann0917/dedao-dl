package services

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

// CourseCourseCategoryList course type list
type CourseCourseCategoryList struct {
	Data struct {
		List      []CourseCategory `json:"list"`
		IsShowURL bool             `json:"is_show_url"`
		PCURL     string           `json:"pc_url"`
	} `json:"data"`
}

// CourseType get course type list
func (s *Service) CourseType() (list *CourseCourseCategoryList, err error) {
	body, err := s.reqCourseType()
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}
