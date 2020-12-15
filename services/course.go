package services

// Course course metadata
type Course struct {
	AssetsType     int         `json:"assets_type"`
	AudioDetail    interface{} `json:"audio_detail"`
	Author         string      `json:"author"`
	ClassFinished  bool        `json:"class_finished"`
	ClassID        int         `json:"class_id"`
	ClassType      int         `json:"class_type"`
	CourseNum      int         `json:"course_num"`
	CreateTime     int         `json:"create_time"`
	DdExtURL       string      `json:"dd_ext_url"`
	DdURL          string      `json:"dd_url"`
	DrmToken       string      `json:"drm_token"`
	Duration       int         `json:"duration"`
	Enid           string      `json:"enid"`
	ExtInfo        interface{} `json:"ext_info"`
	HasExtra       bool        `json:"has_extra"`
	HasPlayAuth    bool        `json:"has_play_auth"`
	Icon           string      `json:"icon"`
	ID             int         `json:"id"`
	Intro          string      `json:"intro"`
	IsCollected    bool        `json:"is_collected"`
	IsFinished     int         `json:"is_finished"`
	IsNew          int         `json:"is_new"`
	IsTop          int         `json:"is_top"`
	LastActionTime int         `json:"last_action_time"`
	LastRead       string      `json:"last_read"`
	LogID          string      `json:"log_id"`
	LogType        string      `json:"log_type"`
	Price          string      `json:"price"`
	ProductIntro   string      `json:"product_intro"`
	ProductPrice   int         `json:"product_price"`
	Progress       int         `json:"progress"`
	PublishNum     int         `json:"publish_num"`
	Size           string      `json:"size"`
	Status         int         `json:"status"`
	Title          string      `json:"title"`
	Type           int         `json:"type"`
	WendaExtInfo   struct {
		AnswerID int `json:"answer_id"`
	} `json:"wenda_ext_info"`
}

// CourseList product list
type CourseList struct {
	List   []Course `json:"list"`
	ISMore int      `json:"is_more"`
}

// CourseList get course list
func (s *Service) CourseList(category, order string, page, limit int) (list *CourseList, err error) {
	body, err := s.reqCourseList(category, order, page, limit)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}
