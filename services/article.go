package services

import (
	"strconv"
	"time"
)

// ArticleDetail article content
// GET query params token,sign,appid
type ArticleDetail struct {
	Article Article `json:"article"`
	Content string  `json:"content"`
}

// Content article Content
type Content struct {
	Title        string        `json:"title"`
	TemplateType string        `json:"templateType"`
	Content      []interface{} `json:"content"`
	Plaintext    string        `json:"plaintext"`
	CoverImage   string        `json:"coverImage"`
	CSSLink      string        `json:"cssLink"`
	Author       struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"author"`
	Message struct {
		Title     string `json:"title"`
		WriteText string `json:"writeText"`
		Tips      string `json:"tips"`
	} `json:"message"`
}

// Article metadata
type Article struct {
	ID          int64  `json:"Id"`
	AppID       int    `json:"AppId"`
	Version     int    `json:"Version"`
	CreateTime  int    `json:"CreateTime"`
	UpdateTime  int    `json:"UpdateTime"`
	PublishTime int    `json:"PublishTime"`
	Status      int    `json:"Status"`
	IDStr       string `json:"IdStr"`
	AppIDStr    string `json:"AppIdStr"`
}

// ArticleIntro article introduce
type ArticleIntro struct {
	ArticleBase
	MediaBaseInfo []MediaBaseInfo `json:"media_base_info"`
	Aduio         Audio           `json:"audio"`
}

// ArticleBase article base info
type ArticleBase struct {
	ID             int      `json:"id"`
	IDStr          string   `json:"id_str"`
	Enid           string   `json:"enid"`
	ClassEnid      string   `json:"class_enid"`
	OriginID       int      `json:"origin_id"`
	OriginIDStr    string   `json:"origin_id_str"`
	ProductType    int      `json:"product_type"`
	ProductID      int      `json:"product_id"`
	ProductIDStr   string   `json:"product_id_str"`
	ClassID        int      `json:"class_id"`
	ClassIDStr     string   `json:"class_id_str"`
	ChapterID      int      `json:"chapter_id"`
	ChapterIDStr   string   `json:"chapter_id_str"`
	Title          string   `json:"title"`
	Logo           string   `json:"logo"`
	URL            string   `json:"url"`
	Summary        string   `json:"summary"`
	Mold           int      `json:"mold"`
	PushContent    string   `json:"push_content"`
	PublishTime    int      `json:"publish_time"`
	PushTime       int      `json:"push_time"`
	PushStatus     int      `json:"push_status"`
	ShareTitle     string   `json:"share_title"`
	ShareContent   string   `json:"share_content"`
	ShareSwitch    int      `json:"share_switch"`
	DdArticleID    int64    `json:"dd_article_id"`
	DdArticleIDStr string   `json:"dd_article_id_str"`
	DdArticleToken string   `json:"dd_article_token"`
	Status         int      `json:"status"`
	CreateTime     int      `json:"create_time"`
	UpdateTime     int      `json:"update_time"`
	CurLearnCount  int      `json:"cur_learn_count"`
	IsFreeTry      bool     `json:"is_free_try"`
	IsUserFreeTry  bool     `json:"is_user_free_try"`
	OrderNum       int      `json:"order_num"`
	IsLike         bool     `json:"is_like"`
	ShareURL       string   `json:"share_url"`
	TrialShareURL  string   `json:"trial_share_url"`
	IsRead         bool     `json:"is_read"`
	LogID          string   `json:"log_id"`
	LogType        string   `json:"log_type"`
	RecommendTitle string   `json:"recommend_title"`
	AudioAliasIds  []string `json:"audio_alias_ids"`
	IsBuy          bool     `json:"is_buy"`
	DdMediaID      int      `json:"dd_media_id"`
	DdMediaIDStr   string   `json:"dd_media_id_str"`
	VideoStatus    int      `json:"video_status"`
}

// ArticleList clasee article list
type ArticleList struct {
	List    []ArticleIntro `json:"article_list"`
	ClassID int            `json:"class_id"`
	Ptype   int            `json:"ptype"`
	PID     int            `json:"pid"`
	Reverse bool           `json:"reverse"`
	MaxID   int            `json:"max_id"`
}

// ArticlePoint article point note
type ArticlePoint struct {
	HasArticlePoint int    `json:"has_article_point"`
	Content         string `json:"content"`
}

// ArticleInfo article info
type ArticleInfo struct {
	ClassID           int          `json:"class_id"`
	ClassEnid         string       `json:"class_enid"`
	Ptype             int          `json:"ptype"`
	Pid               int          `json:"pid"`
	ArticleID         int          `json:"article_id"`
	OriginArticleID   int          `json:"origin_article_id"`
	DdArticleID       int64        `json:"dd_article_id"`
	DdArticleToken    string       `json:"dd_article_token"`
	LikeNum           int          `json:"like_num"`
	IsLike            bool         `json:"is_like"`
	IsBuy             int          `json:"is_buy"`
	ShareSwitch       int          `json:"share_switch"`
	IsFreeTry         bool         `json:"is_free_try"`
	IsUserFreeTry     bool         `json:"is_user_free_try"`
	PrevArticleID     int          `json:"prev_article_id"`
	PrevArticleEnid   string       `json:"prev_article_enid"`
	NextArticleID     int          `json:"next_article_id"`
	NextArticleEnid   string       `json:"next_article_enid"`
	ArticleInfo       ArticleIntro `json:"article_info"`
	ClassInfo         ClassInfo    `json:"class_info"`
	Extra             string       `json:"extra"`
	TrialReadCount    int          `json:"trial_read_count"`
	TrialMaxReadCount int          `json:"trial_max_read_count"`
	ShareImage        string       `json:"share_image"`
	ShareURL          string       `json:"share_url"`
	ArticleTitle      string       `json:"article_title"`
	ClassTitle        string       `json:"class_title"`
	PaymentAudioID    int          `json:"payment_audio_id"`
	ResourceID        int          `json:"resource_id"`
	ResourceType      int          `json:"resource_type"`
}

// ArticleList get class article list
func (s *Service) ArticleList(id string, maxID int) (list *ArticleList, err error) {
	cacheFile := "articleList:" + id + ":" + strconv.Itoa(maxID)
	list = new(ArticleList)
	list.MaxID = maxID
	x, ok := getCache(list, cacheFile)
	if ok {
		list = x.(*ArticleList)
		return
	}
	body, err := s.reqArticleList(id, maxID)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	setCache(list, cacheFile)
	return
}

// ArticleInfo get article info
func (s *Service) ArticleInfo(enid string) (info *ArticleInfo, err error) {
	cacheFile := "articleInfo:" + enid
	x, ok := info.getCache(cacheFile)
	if ok {
		info = x.(*ArticleInfo)
		return
	}
	body, err := s.reqArticleInfo(enid)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	info.setCache(cacheFile)
	return
}

// ArticleDetail get article detail
func (s *Service) ArticleDetail(token, id, appID string) (detail *ArticleDetail, err error) {
	cacheFile := "articleDetail:" + id
	x, ok := detail.getCache(cacheFile)
	if ok {
		detail = x.(*ArticleDetail)
		return
	}
	body, err := s.reqArticleDetail(token, appID)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	detail.setCache(cacheFile)
	return
}

// ArticlePoint get article point
func (s *Service) ArticlePoint(id string, pType int) (detail *ArticleDetail, err error) {
	body, err := s.reqArticlePoint(id, pType)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

func (c *ArticleList) getCacheKey() string {
	return "articleList:" + strconv.Itoa(c.MaxID)
}

func (c *ArticleList) getCache(fileName string) (interface{}, bool) {
	err := LoadCacheFile(fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *ArticleList) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := SaveCacheFile(fileName)
	return err
}

func (c *ArticleDetail) getCacheKey() string {
	return "articleDetail"
}

func (c *ArticleDetail) getCache(fileName string) (interface{}, bool) {
	err := LoadCacheFile(fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *ArticleDetail) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := SaveCacheFile(fileName)
	return err
}

func (c *ArticleInfo) getCacheKey() string {
	return "articleInfo"
}

func (c *ArticleInfo) getCache(fileName string) (interface{}, bool) {
	err := LoadCacheFile(fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *ArticleInfo) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := SaveCacheFile(fileName)
	return err
}
