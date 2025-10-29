package services

import (
	"math"

	"errors"
)

// ReplierInfo Replier Info
type ReplierInfo struct {
	ReplierUID         int    `json:"replier_uid"`
	ReplierName        string `json:"replier_name"`
	ReplierImg         string `json:"replier_img"`
	ReplierIntro       string `json:"replier_intro"`
	ReplierVStatus     bool   `json:"replier_v_status"`
	ReplierVStateValue int    `json:"replier_v_state_value"`
	ReplierTitle       string `json:"replier_title"`
}

// CourseIntro course introduce
type CourseIntro struct {
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CourseInfo product intro info
type CourseInfo struct {
	ClassInfo              ClassInfo     `json:"class_info"`
	Items                  []CourseIntro `json:"items"`
	ArticleIntro           ArticleIntro  `json:"intro_article"`
	ChapterList            []Chapter     `json:"chapter_list"`
	FlatArticleList        []ArticleBase `json:"flat_article_list"`
	UserType               string        `json:"user_type"`
	HasMoreFlatArticleList bool          `json:"has_more_flat_article_list"`
	IsShowGrading          bool          `json:"is_show_grading"`
}

// ClassInfo class info
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
	CornerImgVertical    string `json:"corner_img_vertical"`
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
	FormalArticleCount    int    `json:"formal_article_count"`
	VideoClass            int    `json:"video_class"`
	VideoClassIntro       string `json:"video_class_intro"`
	ActivityIcon          string `json:"activity_icon"`
	ActivityTitle         string `json:"activity_title"`
	ActivityURL           string `json:"activity_url"`
	RealityExtraCount     int    `json:"reality_extra_count"`
	RealityFormalCount    int    `json:"reality_formal_count"`
	IntroArticleIds       []int  `json:"intro_article_ids"`
	IsPreferential        int    `json:"is_preferential"`
	IsCountDown           int    `json:"is_count_down"`
	PreferentialStartTime int    `json:"preferential_start_time"`
	PreferentialEndTime   int    `json:"preferential_end_time"`
	EarlyBirdPrice        int    `json:"early_bird_price"`
	TrialCount            int    `json:"trial_count"`
}

// CourseDetail get course list
func (s *Service) CourseDetail(category string, id int) (detail *Course, err error) {
	list, err := s.CourseListAll(category, "study")
	if err != nil {
		return
	}

	for _, v := range list.List {
		switch category {
		case CateCourse:
			if v.ClassID == id {
				detail = &v
				return
			}
		case CateEbook, CateAudioBook:
			if v.ID == id {
				detail = &v
				return
			}
		default:
			err = errors.New("please make sure to enter the correct course ID")
			return
		}
	}
	err = errors.New("you have not purchased the course, cannot get course information")
	return
}

// CourseInfo get course info
func (s *Service) CourseInfo(enid string) (info *CourseInfo, err error) {

	body, err := s.reqCourseInfo(enid)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

// HasAudio include audio
func (c *CourseInfo) HasAudio() bool {
	return !c.ClassInfo.WithoutAudio
}

// IsSubscribe Is Subscribe
func (c *CourseInfo) IsSubscribe() bool {
	return c.ClassInfo.IsSubscribe == 1
}

// HasAudio include audio
func (c *Course) HasAudio() bool {
	return c.AudioDetail.LogType == "audio"
}

func EnlightenClub() (detail Course) {

	return Course{
		Enid:           "5L9DznlwYyOVdwasGdKmbWABv0Zk4a",
		ID:             252,
		Type:           0,
		ClassType:      0,
		ClassID:        252,
		HasExtra:       false,
		ClassFinished:  false,
		Title:          "罗辑思维·启发俱乐部",
		Intro:          "罗振宇，又称\"罗胖\"，得到App和罗辑思维创始人。",
		Author:         "罗振宇·得到App创始人",
		Icon:           "https://piccdn3.umiwi.com/img/202004/05/202004050004416065909398.jpeg",
		CreateTime:     1472925194,
		LastRead:       "",
		Progress:       0,
		Duration:       0,
		CourseNum:      0,
		PublishNum:     1076,
		LogID:          "252",
		LogType:        "free_column",
		IsTop:          0,
		LastActionTime: 0,
		IsNew:          0,
		IsFinished:     0,
		Size:           "",
		DdURL:          "",
		AssetsType:     0,
		DrmToken:       "",
		AudioDetail:    Audio{},
		ProductPrice:   0,
		Price:          "0.00",
		ProductIntro:   "",
		HasPlayAuth:    false,
		ExtInfo:        nil,
		Status:         0,
		DdExtURL:       "",
		IsCollected:    false,
		WendaExtInfo: struct {
			AnswerID int `json:"answer_id"`
		}{},
		ClassExtReview: ClassExtReview{},
		PlanStatus:     0,
	}
}

// CourseListV2 获取V2版本的课程列表
func (s *Service) CourseList(category, order string, page, limit int) (response *CourseList, err error) {
	body, err := s.reqCourseList(category, order, page, limit)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &response); err != nil {
		return
	}
	return
}

// CourseGroupList fetches a single page of items within a specific group.
// 获取分组内的课程列表（单页）
func (s *Service) CourseGroupList(category, order string, groupID, page, limit int) (response *CourseList, err error) {
	body, err := s.reqCourseGroupList(category, order, groupID, page, limit)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &response); err != nil {
		return
	}
	return
}

// CourseGroupListAll fetches all items within a specific group across all pages.
// It handles pagination automatically and aggregates results.
// 获取分组内的所有课程列表（自动处理分页）
func (s *Service) CourseGroupListAll(category, order string, groupID int) (data *CourseList, err error) {
	resp, err := s.CourseGroupList(category, order, groupID, 1, 18)
	if err != nil {
		return
	}

	if resp.Total == 0 {
		data = resp
		return
	}

	total := resp.Total
	limit := 18
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// 已经获取第一页数据
	var allCourses []Course
	allCourses = append(allCourses, resp.List...)

	// 获取剩余页面数据
	for page := 2; page <= totalPages; page++ {
		pageResp, err := s.CourseGroupList(category, order, groupID, page, limit)
		if err != nil {
			return data, err
		}
		allCourses = append(allCourses, pageResp.List...)
	}

	// 构建完整结果
	data = &CourseList{
		List:          allCourses,
		Total:         total,
		IsMore:        0, // 已获取全部，没有更多
		HasSingleBook: resp.HasSingleBook,
	}

	return
}

// CourseListV2All 获取V2版本的所有课程列表
func (s *Service) CourseListAll(category, order string) (data *CourseList, err error) {
	resp, err := s.CourseList(category, order, 1, 18)
	if err != nil {
		return
	}

	if resp.Total == 0 {
		data = resp
		return
	}

	total := resp.Total
	limit := 18
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// 已经获取第一页数据
	var allCourses []Course
	allCourses = append(allCourses, resp.List...)

	// 获取剩余页面数据
	for page := 2; page <= totalPages; page++ {
		pageResp, err := s.CourseList(category, order, page, limit)
		if err != nil {
			return data, err
		}
		allCourses = append(allCourses, pageResp.List...)
	}
	// 启发俱乐部
	if category == CateCourse {
		allCourses = append(allCourses, EnlightenClub())
	}

	// 构建完整结果
	data = &CourseList{
		List:          allCourses,
		Total:         total,
		IsMore:        0, // 已获取全部，没有更多
		HasSingleBook: resp.HasSingleBook,
	}

	return
}

// OutsideDetail 获取名家讲书课程详情
func (s *Service) OutsideDetail(enid string) (detail *OutsideDetail, err error) {
	body, err := s.reqOutsideDetail(enid)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

// TopicPkgOdobDetails 获取名家讲书每天听本书音频集合详情
func (s *Service) TopicPkgOdobDetails(enid string) (detail *TopicPkgOdobDetails, err error) {
	body, err := s.reqTopicPkgOdobDetails(enid)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}
