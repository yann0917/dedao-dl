package services

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

type HotTab struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	SceneId    int       `json:"scene_id"`
	Status     int       `json:"status"`
	SortNum    int       `json:"sort_num"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	SceneName  string    `json:"scene_name"`
	List       []HTab    `json:"list"`
}

type HTab struct {
	Id             int         `json:"id"`
	Title          string      `json:"title"`
	SubTitle       string      `json:"subTitle"`
	SearchKey      string      `json:"searchKey"`
	Url            string      `json:"url"`
	PinnedDeadLine interface{} `json:"pinnedDeadLine"`
	SceneId        int         `json:"sceneId"`
	Uv             int         `json:"uv"`
	Click          int         `json:"click"`
	Ctr            int         `json:"ctr"`
	TrackInfo      string      `json:"trackInfo"`
	TurnType       int         `json:"turnType"`
	HotType        int         `json:"hotType"`
	LogId          int         `json:"log_id"`
	LogType        string      `json:"log_type"`
}

type Recommend struct {
	PtypeSceneName string `json:"ptype_scene_name"`
	Type           int    `json:"type"`
	List           []RTab `json:"list"`
}

type RTab struct {
	Id         int    `json:"id"`
	Idstr      string `json:"idstr"`
	Name       string `json:"name"`
	Rank       int    `json:"rank"`
	Status     int    `json:"status"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
	FinishTime int    `json:"finish_time"`
	Type       int    `json:"type"`
	IsHot      int    `json:"is_hot"`
	IsTop      int    `json:"is_top"`
	TrackInfo  string `json:"track_info"`
	SearchKey  string `json:"search_key"`
}

type SearchTot struct {
	HotTabList   []HotTab    `json:"hot_tab_list"`
	RecommendMap []Recommend `json:"recommend_map"`
}

type Navigation struct {
	Enid         string  `json:"enid"`
	Id           int     `json:"id"`
	RelationId   int     `json:"relation_id"`
	RelationName string  `json:"relation_name"`
	Type         int     `json:"type"`
	NavType      int     `json:"nav_type"`
	Name         string  `json:"name"`
	Icon         string  `json:"icon"`
	EnglishName  string  `json:"english_name"`
	LabelList    []Label `json:"label_list"`
}

type Label struct {
	Enid string `json:"enid"`
	Name string `json:"name"`
}

type SunflowerLabelList struct {
	List []Navigation `json:"list"`
}
type ProductSimple struct {
	ProductType     int      `json:"product_type"`
	ProductEnid     string   `json:"product_enid"`
	Title           string   `json:"title"`
	Intro           string   `json:"intro"`
	Introduction    string   `json:"introduction"`
	IndexImage      string   `json:"index_image"`
	Score           string   `json:"score"`
	UserScoreCount  int      `json:"user_score_count"`
	HorizontalImage string   `json:"horizontal_image"`
	LearnUserCount  int      `json:"learn_user_count"`
	AuthorList      []string `json:"author_list"`
	TrackInfo       string   `json:"trackinfo"`
	LogType         string   `json:"log_type"`
}

type SunflowerContent struct {
	ProductList    []ProductSimple `json:"product_list"`
	CurrentEnid    string          `json:"current_enid"`
	NavigationList []Navigation    `json:"navigation_list"`
	PageId         int             `json:"page_id"`
	PageSize       int             `json:"page_size"`
	IsMore         int             `json:"is_more"` // 0,1
	RequestId      string          `json:"request_id"`
}

type SunflowerResource struct {
	Id          int     `json:"id"`
	Enid        string  `json:"enid"`
	Name        string  `json:"name"`
	Intro       string  `json:"intro"`
	Logo        string  `json:"logo"`
	ProductType int     `json:"product_type"`
	ProductId   int     `json:"product_id"`
	Score       float64 `json:"score"`
	ClassType   int     `json:"class_type"`
	Status      int     `json:"status"`
}

type SunflowerResourceList struct {
	List []SunflowerResource `json:"list"`
}

type AlgoFilterParam struct {
	ClassfcName  string        `json:"classfc_name"` //  default "心理学"
	LabelId      string        `json:"label_id"`
	NavType      int           `json:"nav_type"`
	NavigationId string        `json:"navigation_id"`
	Page         int           `json:"page"`          // default 0
	PageSize     int           `json:"page_size"`     // default 18
	ProductTypes string        `json:"product_types"` // default "66"
	RequestId    string        `json:"request_id"`
	SortStrategy string        `json:"sort_strategy"` // default "HOT" "NEW"
	TagsIds      []interface{} `json:"tags_ids"`
}

type Option struct {
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	SubOptions []Option `json:"sub_options,omitempty"`
}

type Strategy struct {
	Title      string   `json:"title"`
	IsMultiple bool     `json:"is_multiple"`
	IsHide     bool     `json:"is_hide"`
	Options    []Option `json:"options"`
}

type AlgoFilter struct {
	Title            string   `json:"title"`
	ProductTypes     Strategy `json:"product_types"`
	SortStrategy     Strategy `json:"sort_strategy"`
	ProgressStrategy Strategy `json:"progress_strategy"`
	BuyStrategy      Strategy `json:"buy_strategy"`
	Navigations      Strategy `json:"navigations"`
	Tags             Strategy `json:"tags"`
}

type AlgoFilterResp struct {
	Filter  AlgoFilter      `json:"filter"`
	Total   int             `json:"total"`
	Request AlgoFilterParam `json:"request"`
}

type CostIntro struct {
	Price string `json:"price"`
}

type ArticleSimpleInfo struct {
	ArticleId            int    `json:"article_id"`
	OriginArticleId      int    `json:"origin_article_id"`
	ArticleAudioDuration int    `json:"article_audio_duration"`
	ArticleTitle         string `json:"article_title"`
	ArticleIntro         string `json:"article_intro"`
	Image                string `json:"image"`
}

type AlgoProduct struct {
	ItemType              int               `json:"item_type"`
	Id                    int               `json:"id"`
	ProductType           int               `json:"product_type"`
	ProductId             int               `json:"product_id"`
	Name                  string            `json:"name"`
	Intro                 string            `json:"intro"`
	PhaseNum              int               `json:"phase_num"`
	LearnUserCount        int               `json:"learn_user_count"`
	Price                 int               `json:"price"`
	CurrentPrice          int               `json:"current_price"`
	IndexImg              string            `json:"index_img"`
	HorizontalImage       string            `json:"horizontal_image"`
	LecturersImg          string            `json:"lecturers_img"`
	PriceDesc             string            `json:"price_desc"`
	IsBuy                 int               `json:"is_buy"`
	IsVip                 bool              `json:"is_vip"`
	Trackinfo             string            `json:"trackinfo"`
	LogId                 string            `json:"log_id"`
	LogType               string            `json:"log_type"`
	CornerImg             string            `json:"corner_img"`
	LecturerNameAndTitle  string            `json:"lecturer_name_and_title"`
	LecturersVStatus      int               `json:"lecturers_v_status"`
	CornerImgVertical     string            `json:"corner_img_vertical"`
	LecturerName          string            `json:"lecturer_name"`
	LecturerTitle         string            `json:"lecturer_title"`
	AuthorList            []string          `json:"author_list"`
	IdOut                 string            `json:"id_out"`
	IsLimitFree           bool              `json:"is_limit_free"`
	HasPlayAuth           bool              `json:"has_play_auth"`
	ButtonType            int               `json:"button_type"`
	AudioId               string            `json:"audio_id"`
	AliasId               string            `json:"alias_id"`
	InBookrack            bool              `json:"in_bookrack"`
	IsShowNewbook         bool              `json:"is_show_newbook"`
	IsSubscribe           bool              `json:"is_subscribe"`
	OnlineTime            string            `json:"online_time"`
	ArticleInfo           ArticleSimpleInfo `json:"article_info"`
	NeedLogin             int               `json:"need_login"`
	IsOnBookshelf         bool              `json:"is_on_bookshelf"`
	Duration              int               `json:"duration"`
	CollectionNum         int               `json:"collection_num"`
	DdUrl                 string            `json:"dd_url"`
	LearningDaysDesc      string            `json:"learning_days_desc"`
	IsPreferential        int               `json:"is_preferential"`
	IsCountDown           int               `json:"is_count_down"`
	PreferentialStartTime int               `json:"preferential_start_time"`
	PreferentialEndTime   int               `json:"preferential_end_time"`
	EarlyBirdPrice        int               `json:"early_bird_price"`
	TimeNow               int               `json:"time_now"`
	HotLearnDesc          string            `json:"hot_learn_desc"`
	Score                 string            `json:"score"`
	Introduction          string            `json:"introduction"`
	UserScoreCount        int               `json:"user_score_count"`
	Progress              int               `json:"progress"`
	Metrics               string            `json:"metrics"`
	CostIntro             CostIntro         `json:"cost_intro"`
	ConsumedNum           int               `json:"consumed_num"`
	FreeMaximum           int               `json:"free_maximum"`
	BSellingChannelGroup  int               `json:"b_selling_channel_group"`
	TrackInfo             string            `json:"track_info"`
	SortValues            []interface{}     `json:"sort_values"`
}

type AlgoProductResp struct {
	ProductList []AlgoProduct `json:"product_list"`
	RequestId   string        `json:"request_id"`
	Total       int           `json:"total"`
	IsMore      int           `json:"is_more"`
}

type HomeModule struct {
	Description string `json:"description"`
	Ext1        string `json:"ext1"`
	Ext2        string `json:"ext2"`
	Ext3        string `json:"ext3"`
	Ext4        string `json:"ext4"`
	Ext5        string `json:"ext5"`
	Id          int    `json:"id"`
	IsShow      int    `json:"isShow"`
	Name        string `json:"name"`
	Sort        int    `json:"sort"`
	Title       string `json:"title"`
	Type        string `json:"type"`
}

type HomeCategory struct {
	EnglishName  string  `json:"englishName"`
	Enid         string  `json:"enid"`
	Icon         string  `json:"icon"`
	Id           int     `json:"id"`
	LabelList    []Label `json:"labelList"`
	Name         string  `json:"name"`
	NavType      int     `json:"navType"`
	RelationId   int     `json:"relationId"`
	RelationName string  `json:"relationName"`
	Type         int     `json:"type"`
}

type Banner struct {
	BeginTime int    `json:"beginTime"`
	EndTime   int    `json:"endTime"`
	Id        int    `json:"id"`
	Img       string `json:"img"`
	IsDelete  int    `json:"isDelete"`
	ModuleId  int    `json:"moduleId"`
	Sort      int    `json:"sort"`
	SourceId  string `json:"sourceId"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	UrlType   int    `json:"urlType"`
}

type HomeData struct {
	ModuleList   []HomeModule   `json:"moduleList"`
	CategoryList []HomeCategory `json:"categoryList"`
	Banner       []Banner       `json:"banner"`
}
type HomeInitState struct {
	IsLogin  bool     `json:"isLogin"`
	HomeData HomeData `json:"homeData"`
	Uid      string   `json:"uid"`
}

var (
	CsrfToken = ""
	SetCookie []string
)

func (s *Service) GetHomeInitialState() (state HomeInitState, err error) {
	resp, err := s.client.R().Get("")
	if err != nil {
		return
	}
	if resp.IsSuccess() {
		SetCookie = resp.Header().Values("Set-Cookie")
		cookies := strings.Split(strings.Join(SetCookie, "; "), "; ")
		for _, v := range cookies {
			item := strings.Split(v, "=")
			if len(item) > 1 && item[0] == "csrfToken" {
				CsrfToken = item[1]
			}
		}
		// 匹配 <script> window.__INITIAL_STATE__=
		src := string(resp.Body())
		reg := regexp.MustCompile(`<script>[\S\s]+?</script>`)
		list := reg.FindAllString(src, -1)
		for _, item := range list {
			if strings.Contains(item, "window.__INITIAL_STATE__=") {
				reg = regexp.MustCompile(`{[\S\s]+[^}{]*?}`)
				result := reg.FindAllString(item, -1)
				if len(result) > 0 {
					err = json.Unmarshal([]byte(result[0]), &state)
				}
				break
			}
		}
	}
	return
}

// SearchHot 搜索框热门搜索
func (s *Service) SearchHot() (list *SearchTot, err error) {
	body, err := s.reqSearchHot()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

// SunflowerLabelList 首页导航标签列表
// nType 2-好看又好查的电子书, 4-精选课程
func (s *Service) SunflowerLabelList(nType int) (list *SunflowerLabelList, err error) {
	body, err := s.reqSunflowerLabelList(nType)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

// SunflowerLabelContent 首页导航标签内容列表
// nType 2-好看又好查的电子书, 4-精选课程
// page default=0 pageSize default=4
func (s *Service) SunflowerLabelContent(enID string, nType, page, pageSize int) (list *SunflowerContent, err error) {
	body, err := s.reqSunflowerLabelContent(enID, nType, page, pageSize)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

// SunflowerResourceList 首页免费专区
func (s *Service) SunflowerResourceList() (list *SunflowerResourceList, err error) {
	body, err := s.reqSunflowerResourceList()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

func (s *Service) AlgoFilter(param AlgoFilterParam) (resp *AlgoFilterResp, err error) {
	body, err := s.reqAlgoFilter(param)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}

func (s *Service) AlgoProduct(param AlgoFilterParam) (resp *AlgoProductResp, err error) {
	body, err := s.reqAlgoProduct(param)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}
