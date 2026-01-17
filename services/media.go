package services

// MediaBaseInfo media info
type MediaBaseInfo struct {
	MediaType     int    `json:"media_type"` // 1-audio,2-video
	SourceID      string `json:"source_id"`
	SecurityToken string `json:"security_token"`
	SourceURL     string `json:"source_url"`
}

// Audio audio 每天听本书音频详情
type Audio struct {
	TopicEncodeID      string  `json:"topic_encode_id"`       // 主题编码ID
	AliasID            string  `json:"alias_id"`              // 别名ID
	AudioListIcon      string  `json:"audio_list_icon"`       // 音频列表图标
	AudioType          int     `json:"audio_type"`            // 音频类型
	BoredCount         int     `json:"bored_count"`           // 无聊计数
	ClassArticleID     string  `json:"class_article_id"`      // 课程文章ID
	ClassCourseID      string  `json:"class_course_id"`       // 课程ID
	ClassID            int     `json:"class_id"`              // 课程ID
	Collection         int     `json:"collection"`            // 收藏
	Count              int     `json:"count"`                 // 计数
	DdArticleID        string  `json:"dd_article_id"`         // 得到文章ID
	DrmVersion         int     `json:"drm_version"`           // DRM版本
	Duration           int     `json:"duration"`              // 时长
	Etag               string  `json:"etag"`                  // Etag
	Extrainfo          string  `json:"extrainfo"`             // 额外信息
	Icon               string  `json:"icon"`                  // 图标
	IndexImg           string  `json:"index_img"`             // 索引图片
	ListenFinished     bool    `json:"listen_finished"`       // 是否听完
	ListenProgress     float64 `json:"listen_progress"`       // 听的进度
	LogID              string  `json:"log_id"`                // 日志ID
	LogInterface       string  `json:"log_interface"`         // 日志接口
	LogType            string  `json:"log_type"`              // 日志类型
	Price              int     `json:"price"`                 // 价格
	Reader             string  `json:"reader"`                // 读者
	ReaderName         string  `json:"reader_name"`           // 读者名称
	Schedule           int     `json:"schedule"`              // 计划
	ShareSummary       string  `json:"share_summary"`         // 分享摘要
	ShareTitle         string  `json:"share_title"`           // 分享标题
	Size               int     `json:"size"`                  // 大小
	SourceIcon         string  `json:"source_icon"`           // 来源图标
	SourceID           int     `json:"source_id"`             // 来源ID
	SourceName         string  `json:"source_name"`           // 来源名称
	SourceType         int     `json:"source_type"`           // 来源类型
	Summary            string  `json:"summary"`               // 摘要
	Title              string  `json:"title"`                 // 标题
	TopicID            int     `json:"topic_id"`              // 主题ID
	Trackinfo          string  `json:"trackinfo"`             // 跟踪信息
	UsedDRM            int     `json:"used_drm"`              // 使用的DRM
	MP3PlayURL         string  `json:"mp3_play_url"`          // MP3播放URL
	OdobGroupEnid      string  `json:"odob_group_enid"`       // 每天听本书分组加密ID
	Slogan             string  `json:"slogan"`                // 标语
	Token              string  `json:"token"`                 // 令牌
	PlayerImg          string  `json:"player_img"`            // 播放器图片
	IsPlayLater        bool    `json:"is_play_later"`         // 是否稍后播放
	HasPlayAuth        bool    `json:"has_play_auth"`         // 是否有播放权限
	IsVIP              bool    `json:"is_vip"`                // 是否VIP
	Category           int     `json:"category"`              // 分类
	IsFree             int     `json:"is_free"`               // 是否免费
	TrialDuration      int     `json:"trial_duration"`        // 试听时长
	UpdateTips         string  `json:"update_tips"`           // 更新提示
	TrialListenTips    string  `json:"trial_listen_tips"`     // 试听提示
	TrialListenEndTips string  `json:"trial_listen_end_tips"` // 试听结束提示
	TrialListenEndURL  string  `json:"trial_listen_end_url"`  // 试听结束URL
	PlayCount          int     `json:"play_count"`            // 播放计数
	PlayCountTips      string  `json:"play_count_tips"`       // 播放计数提示
	BookShelfStatus    int     `json:"book_shelf_status"`     // 书架状态
	ShareURL           string  `json:"share_url"`             // 分享URL
	PlayDDURL          string  `json:"play_dd_url"`           // 播放得到URL
	IsSubscribed       bool    `json:"is_subscribed"`         // 是否已订阅
	PackagePID         int     `json:"package_pid"`           // 包PID
	PackagePType       int     `json:"package_ptype"`         // 包类型
	PackageTitle       string  `json:"package_title"`         // 包标题
	Superscript        string  `json:"superscript"`           // 上标
	PodcastIcon        string  `json:"podcast_icon"`          // 播客图标
	IsNewest           int     `json:"is_newest"`             // 是否最新
	IsLatestLearning   int     `json:"is_latest_learning"`    // 是否最新学习

	M3u8Token string `json:"m3u8_token"`
}

type Video struct {
	Token            string  `json:"token"`
	TokenVersion     int     `json:"token_version"`
	CoverImg         string  `json:"cover_img"`
	DdMediaID        int64   `json:"dd_media_id"`
	DdMediaIDStr     string  `json:"dd_media_id_str"`
	M3u8Token        string  `json:"m3u8_token"`
	Duration         int     `json:"duration"`
	Bitrate480       string  `json:"bitrate_480"`
	Bitrate480Size   int     `json:"bitrate_480_size"`
	Bitrate480Audio  string  `json:"bitrate_480_audio"`
	Bitrate720       string  `json:"bitrate_720"`
	Bitrate720Size   int     `json:"bitrate_720_size"`
	Bitrate720Audio  string  `json:"bitrate_720_audio"`
	Bitrate1080      string  `json:"bitrate_1080"`
	Bitrate1080Size  int     `json:"bitrate_1080_size"`
	Bitrate1080Audio string  `json:"bitrate_1080_audio"`
	IsDrm            bool    `json:"is_drm"`
	ListenProgress   float64 `json:"listen_progress"`
	ListenFinished   bool    `json:"listen_finished"`
	LogID            string  `json:"log_id"`
	LogType          string  `json:"log_type"`
	Caption          string  `json:"caption"`
	VttCaption       string  `json:"vtt_caption"`
	VideoAudio       string  `json:"video_audio"`
	MediaFiles       any     `json:"media_files"`
	SourceName       string  `json:"source_name"`
	SourceType       int     `json:"source_type"`
	SourceID         int     `json:"source_id"`
}

// AudioList audio basic info list
type AudioList struct {
	List []Audio `json:"list"`
}

type AudioDetail struct {
	Detail Audio `json:"audio_detail"`
}

// AudioByAlias get article audio info
func (s *Service) AudioByAlias(ID string) (list *AudioList, err error) {
	body, err := s.reqAudioByAlias(ID)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &list); err != nil {
		return
	}
	return
}

// AudioDetailAlias get odob audio info
func (s *Service) AudioDetailAlias(ID string) (detail *Audio, err error) {
	adetail := AudioDetail{}
	body, err := s.reqOdobAudioDetail(ID)
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &adetail); err != nil {
		return
	}
	detail = &adetail.Detail
	return
}
