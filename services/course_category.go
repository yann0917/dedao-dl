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

type NavbarChild struct {
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Filter    string `json:"filter"`
	ShowCount bool   `json:"show_count"`
}

type NavbarItem struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Category    string        `json:"category"`
	ChannelType string        `json:"channel_type"`
	ItemType    string        `json:"item_type"`
	Children    []NavbarChild `json:"children"`
}

type NavbarData struct {
	List []NavbarItem `json:"list"`
}

// GetNavbar 获取已购导航与筛选配置
func (s *Service) GetNavbar() (data *NavbarData, err error) {
	body, err := s.reqNavbar()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &data); err != nil {
		return
	}
	return
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
