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
)

// CourseType course type metadata
type CourseType struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Category string `json:"category"`
}

// CourseTypeList course type list
type CourseTypeList struct {
	List      []CourseType `json:"list"`
	IsShowURL bool         `json:"is_show_url"`
	PCURL     string       `json:"pc_url"`
}
