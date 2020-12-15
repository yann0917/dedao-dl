package services

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
