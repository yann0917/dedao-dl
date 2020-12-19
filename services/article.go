package services

import (
	"fmt"
	"time"
)

// ArticleDetail article content
// GET query params token,sign,appid
type ArticleDetail struct {
	Article Article `json:"article"`
	Content string  `json:"content"`
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
	ArticleInfo
	MediaBaseInfo []MediaBaseInfo `json:"media_base_info"`
	Aduio         Audio           `json:"audio"`
}

// ArticleInfo article info
type ArticleInfo struct {
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
}

// ArticlePoint article point note
type ArticlePoint struct {
	HasArticlePoint int    `json:"has_article_point"`
	Content         string `json:"content"`
}

// ArticleList get class article list
func (s *Service) ArticleList(id string) (list *ArticleList, err error) {
	cacheFile := "articleList:" + id
	x, ok := list.getCache(cacheFile)
	if ok {
		list = x.(*ArticleList)
		return
	}
	body, err := s.reqArticleList(id)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	list.setCache(cacheFile)
	return
}

// ArticleDetail get article detail
func (s *Service) ArticleDetail(token, sign, appID string) (detail *ArticleDetail, err error) {
	body, err := s.reqArticleDetail(token, sign, appID)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

// ArticlePoint get article point
func (s *Service) ArticlePoint(token, sign, appID string) (detail *ArticleDetail, err error) {
	body, err := s.reqArticleDetail(token, sign, appID)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	return
}

func (c *ArticleList) getCacheKey() string {
	return "articleList"
}

func (c *ArticleList) getCache(fileName string) (interface{}, bool) {
	err := Cache.LoadFile(cacheDir + fileName)
	fmt.Println(err)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *ArticleList) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := Cache.SaveFile(cacheDir + fileName)
	return err
}
