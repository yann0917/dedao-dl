package services

import (
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Course course metadata
type Course struct {
	Enid           string        `json:"enid"`
	ID             int           `json:"id"`
	Type           int           `json:"type"`
	ClassType      int           `json:"class_type"`
	ClassID        int           `json:"class_id"`
	HasExtra       bool          `json:"has_extra"`
	ClassFinished  bool          `json:"class_finished"`
	Title          string        `json:"title"`
	Intro          string        `json:"intro"`
	Author         string        `json:"author"`
	Icon           string        `json:"icon"`
	CreateTime     int           `json:"create_time"`
	LastRead       string        `json:"last_read"`
	Progress       int           `json:"progress"`
	Duration       int           `json:"duration"`
	CourseNum      int           `json:"course_num"`
	PublishNum     int           `json:"publish_num"`
	LogID          string        `json:"log_id"`
	LogType        string        `json:"log_type"`
	IsTop          int           `json:"is_top"`
	LastActionTime int           `json:"last_action_time"`
	IsNew          int           `json:"is_new"`
	IsFinished     int           `json:"is_finished"`
	Size           string        `json:"size"`
	DdURL          string        `json:"dd_url"`
	AssetsType     int           `json:"assets_type"`
	DrmToken       string        `json:"drm_token"`
	AudioDetail    Audio         `json:"audio_detail"`
	ProductPrice   int           `json:"product_price"`
	Price          string        `json:"price"`
	ProductIntro   string        `json:"product_intro"`
	HasPlayAuth    bool          `json:"has_play_auth"`
	ExtInfo        []ReplierInfo `json:"ext_info"`
	Status         int           `json:"status"`
	DdExtURL       string        `json:"dd_ext_url"`
	IsCollected    bool          `json:"is_collected"`
	WendaExtInfo   struct {
		AnswerID int `json:"answer_id"`
	} `json:"wenda_ext_info"`
}

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

// CourseList product list
type CourseList struct {
	List   []Course `json:"list"`
	ISMore int      `json:"is_more"`
	Page   int      `json:"page"`
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

// FlatArticleList flat
type FlatArticleList struct {
	ID             int           `json:"id"`
	IDStr          string        `json:"id_str"`
	Enid           string        `json:"enid"`
	ClassEnid      string        `json:"class_enid"`
	OriginID       int           `json:"origin_id"`
	OriginIDStr    string        `json:"origin_id_str"`
	ProductType    int           `json:"product_type"`
	ProductID      int           `json:"product_id"`
	ProductIDStr   string        `json:"product_id_str"`
	ClassID        int           `json:"class_id"`
	ClassIDStr     string        `json:"class_id_str"`
	ChapterID      int           `json:"chapter_id"`
	ChapterIDStr   string        `json:"chapter_id_str"`
	Title          string        `json:"title"`
	Logo           string        `json:"logo"`
	URL            string        `json:"url"`
	Summary        string        `json:"summary"`
	Mold           int           `json:"mold"`
	PushContent    string        `json:"push_content"`
	PublishTime    int           `json:"publish_time"`
	PushTime       int           `json:"push_time"`
	PushStatus     int           `json:"push_status"`
	ShareTitle     string        `json:"share_title"`
	ShareContent   string        `json:"share_content"`
	ShareSwitch    int           `json:"share_switch"`
	DdArticleID    int64         `json:"dd_article_id"`
	DdArticleIDStr string        `json:"dd_article_id_str"`
	DdArticleToken string        `json:"dd_article_token"`
	Status         int           `json:"status"`
	CreateTime     int           `json:"create_time"`
	UpdateTime     int           `json:"update_time"`
	CurLearnCount  int           `json:"cur_learn_count"`
	IsFreeTry      bool          `json:"is_free_try"`
	IsUserFreeTry  bool          `json:"is_user_free_try"`
	OrderNum       int           `json:"order_num"`
	IsLike         bool          `json:"is_like"`
	ShareURL       string        `json:"share_url"`
	TrialShareURL  string        `json:"trial_share_url"`
	IsRead         bool          `json:"is_read"`
	LogID          string        `json:"log_id"`
	LogType        string        `json:"log_type"`
	RecommendTitle string        `json:"recommend_title"`
	AudioAliasIds  []interface{} `json:"audio_alias_ids"`
	IsBuy          bool          `json:"is_buy"`
	DdMediaID      int           `json:"dd_media_id"`
	DdMediaIDStr   string        `json:"dd_media_id_str"`
	VideoStatus    int           `json:"video_status"`
	DdLiveID       int           `json:"dd_live_id"`
	NotJoinPlan    int           `json:"not_join_plan"`
}

// CourseList get course list by page
func (s *Service) CourseList(category, order string, page, limit int) (list *CourseList, err error) {
	cacheFile := "courseList:" + category + ":" + strconv.Itoa(page)
	list = new(CourseList)
	list.Page = page
	x, ok := list.getCache(cacheFile)
	if ok {
		list = x.(*CourseList)
		return
	}
	body, err := s.reqCourseList(category, order, page, limit)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	list.setCache(cacheFile)
	return
}

// CourseListAll get all course list
func (s *Service) CourseListAll(category, order string) (list *CourseList, err error) {
	count, err := s.CourseCount(category)
	if err != nil {
		return
	}
	limit := 18.0
	page := int(math.Ceil(float64(count) / limit))
	var lists []Course
	for i := 1; i <= page; i++ {
		list, err = s.CourseList(category, order, i, int(limit))
		if err != nil {
			return
		}
		lists = append(lists, list.List...)
	}
	// 启发俱乐部
	if category == CateCourse {
		lists = append(lists, EnlightenClub())
	}
	if page == 0 {
		list = new(CourseList)
	}
	list.List = lists
	return
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
		case CateAudioBook:
			if v.ID == id {
				detail = &v
				return
			}
		}
	}
	err = errors.New("You have not purchased the course, cannot get course information")
	return
}

// CourseInfo get course info
func (s *Service) CourseInfo(enid string) (info *CourseInfo, err error) {
	cacheFile := "courseInfo:" + enid
	x, ok := info.getCache(cacheFile)
	if ok {
		info = x.(*CourseInfo)
		return
	}
	body, err := s.reqCourseInfo(enid)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	info.setCache(cacheFile)
	return
}

// HasAudio include audio
func (c *CourseInfo) HasAudio() bool {
	return c.ClassInfo.WithoutAudio == false
}

// IsSubscribe Is Subscribe
func (c *CourseInfo) IsSubscribe() bool {
	return c.ClassInfo.IsSubscribe == 1
}

// HasAudio include audio
func (c *Course) HasAudio() bool {
	return c.AudioDetail.LogType == "audio"
}

func (c *CourseList) getCacheKey() string {
	return "courseList:" + strconv.Itoa(c.Page)
}

func (c *CourseList) getCache(fileName string) (interface{}, bool) {
	err := LoadCacheFile(fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *CourseList) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := SaveCacheFile(fileName)
	return err
}

func (c *CourseInfo) getCacheKey() string {
	return "courseInfo"
}

func (c *CourseInfo) getCache(fileName string) (interface{}, bool) {
	err := LoadCacheFile(fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *CourseInfo) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := SaveCacheFile(fileName)
	return err
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
		Intro:          "罗振宇，又称“罗胖”，得到App和罗辑思维创始人。",
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
	}
}
