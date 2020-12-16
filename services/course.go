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

// CourseIntro course introduce
type CourseIntro struct {
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CourseList product list
type CourseList struct {
	List   []Course `json:"list"`
	ISMore int      `json:"is_more"`
}

// CourseInfo product intro info
type CourseInfo struct {
	ClassInfo    ClassInfo     `json:"class_info"`
	Items        []CourseIntro `json:"items"`
	ArticleIntro ArticleIntro  `json:"intro_article"`
	ChapterList  []Chapter     `json:"chapter_list"`
}

type ClassInfo struct {
	LogID                string `json:"log_id"`
	LogType              string `json:"log_type"`
	ID                   int    `json:"id"`
	IDStr                string `json:"id_str"`
	Enid                 string `json:"enid"`
	ProductID            int    `json:"product_id"`
	ProductType          int    `json:"product_type"`
	HasChapter           int    `json:"has_chapter"`
	Name                 string `json:"name"`
	Intro                string `json:"intro"`
	PhaseNum             int    `json:"phase_num"`
	LearnUserCount       int    `json:"learn_user_count"`
	CurrentArticleCount  int    `json:"current_article_count"`
	Highlight            string `json:"highlight"`
	Price                int    `json:"price"`
	IndexImg             string `json:"index_img"`
	IndexImgApplet       string `json:"index_img_applet"`
	Logo                 string `json:"logo"`
	LogoIphonex          string `json:"logo_iphonex"`
	SquareImg            string `json:"square_img"`
	OutlineImg           string `json:"outline_img"`
	PlayerImg            string `json:"player_img"`
	ShareTitle           string `json:"share_title"`
	ShareSummary         string `json:"share_summary"`
	Status               int    `json:"status"`
	OrderNum             int    `json:"order_num"`
	ShzfURL              string `json:"shzf_url"`
	ShzfURLQr            string `json:"shzf_url_qr"`
	PriceDesc            string `json:"price_desc"`
	ArticleTime          int    `json:"article_time"`
	ArticleTitle         string `json:"article_title"`
	IsSubscribe          int    `json:"is_subscribe"`
	LecturerUID          int    `json:"lecturer_uid"`
	LecturerName         string `json:"lecturer_name"`
	LecturerTitle        string `json:"lecturer_title"`
	LecturerIntro        string `json:"lecturer_intro"`
	LecturerNameAndTitle string `json:"lecturer_name_and_title"`
	LecturerAvatar       string `json:"lecturer_avatar"`
	IsFinished           int    `json:"is_finished"`
	UpdateTime           int    `json:"update_time"`
	ShareURL             string `json:"share_url"`
	DefaultSortReverse   bool   `json:"default_sort_reverse"`
	PresaleURL           string `json:"presale_url"`
	WithoutAudio         bool   `json:"without_audio"`
	ViewType             int    `json:"view_type"`
	H5URLName            string `json:"h5_url_name"`
	PackageManagerSwitch bool   `json:"package_manager_switch"`
	LectureAvatorSpecial string `json:"lecture_avator_special"`
	MiniShareImg         string `json:"mini_share_img"`
	EstimatedShelfTime   int    `json:"estimated_shelf_time"`
	EstimatedDownTime    int    `json:"estimated_down_time"`
	CornerImg            string `json:"corner_img"`
	WithoutGiving        bool   `json:"without_giving"`
	DdURL                string `json:"dd_url"`
	PublishTime          int    `json:"publish_time"`
	DdMediaID            string `json:"dd_media_id"`
	VideoCover           string `json:"video_cover"`
	IsInVip              bool   `json:"is_in_vip"`
	IsVip                bool   `json:"is_vip"`
	Collection           struct {
		IsCollected     bool `json:"is_collected"`
		CollectionCount int  `json:"collection_count"`
	} `json:"collection"`
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

// CourseInfo get course info
func (s *Service) CourseInfo(ID string) (info *CourseInfo, err error) {
	body, err := s.reqCourseInfo(ID)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}
