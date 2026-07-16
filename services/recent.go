package services

import "fmt"

// RecentRequest 最近学习接口请求参数
type RecentRequest struct {
	FilterProductType bool   `json:"filter_product_type"`
	MaxID             int64  `json:"max_id"`
	PageSize          int    `json:"page_size"`
	ProductType       string `json:"product_type"`
	UID               *int64 `json:"uid"`
	UIDHazy           string `json:"uid_hazy"`
}

// RecentResponse 最近学习接口响应
type RecentResponse struct {
	List        []RecentItem `json:"list"`
	Timestamp   int64        `json:"timestamp"`
	CurrentTime int          `json:"current_time"`
	HasMore     bool         `json:"has_more"`
}

// RecentItem 最近学习列表项
type RecentItem struct {
	ProductID       int             `json:"product_id"`
	ProductType     int             `json:"product_type"`
	ProductIDHazy   string          `json:"product_id_hazy"`
	OriginType      int             `json:"origin_type,omitempty"`
	OriginID        int             `json:"origin_id,omitempty"`
	TypeName        string          `json:"type_name"`
	Title           string          `json:"title"`
	Author          string          `json:"author"`
	Status          int             `json:"status"`
	HasEnded        int             `json:"has_ended"`
	CurrentCount    int             `json:"current_count"`
	IndexImg        string          `json:"index_img"`
	SquareImg       string          `json:"square_img"`
	Logo            string          `json:"logo"`
	OnlineTime      string          `json:"online_time"`
	OnlineTimeStamp int             `json:"online_time_stamp,omitempty"`
	Timestamp       int64           `json:"timestamp"`
	LastInfo        string          `json:"last_info"`
	Resource        *RecentResource `json:"resource,omitempty"`
	ProgressIntro   RecentProgress  `json:"progress_intro"`
}

// RecentResource 最近学习资源
type RecentResource struct {
	AudioID      string `json:"audio_id"`
	Duration     int    `json:"duration,omitempty"`
	ResourceID   int    `json:"resource_id"`
	ResourceType int    `json:"resource_type"`
	Title        string `json:"title,omitempty"`
	DdURL        string `json:"dd_url"`
	Audio        *Audio `json:"audio,omitempty"`
}

// RecentProgress 最近学习进度
type RecentProgress struct {
	Intro       string `json:"intro"`
	Progress    int    `json:"progress"`
	MaxProgress int    `json:"max_progress"`
	IsFinish    int    `json:"is_finish"`
	Uint        string `json:"uint"`
}

// Recent 获取最近学习情况
func (s *Service) Recent(param RecentRequest) (resp *RecentResponse, err error) {
	if param.UIDHazy == "" {
		err = fmt.Errorf("uid_hazy is required")
		return
	}
	if param.PageSize <= 0 {
		param.PageSize = 20
	}

	body, err := s.reqRecent(param)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &resp); err != nil {
		return
	}
	return
}
