package services

// MediaBaseInfo media info
type MediaBaseInfo struct {
	MediaType     int    `json:"media_type"` // 1-audio,2-video
	SourceID      string `json:"source_id"`
	SecurityToken string `json:"security_token"`
}

// Audio audio
type Audio struct {
	AliasID        string  `json:"alias_id"`
	Extrainfo      string  `json:"extrainfo"`
	ClassID        int     `json:"class_id"`
	Title          string  `json:"title"`
	ShareTitle     string  `json:"share_title"`
	Mp3PlayURL     string  `json:"mp3_play_url"`
	Duration       int     `json:"duration"`
	Schedule       int     `json:"schedule"`
	TopicID        int     `json:"topic_id"`
	Summary        string  `json:"summary"`
	Price          int     `json:"price"`
	Icon           string  `json:"icon"`
	Size           int     `json:"size"`
	Etag           string  `json:"etag"`
	Token          string  `json:"token"`
	ShareSummary   string  `json:"share_summary"`
	Collection     int     `json:"collection"`
	Count          int     `json:"count"`
	BoredCount     int     `json:"bored_count"`
	AudioType      int     `json:"audio_type"`
	DrmVersion     int     `json:"drm_version"`
	SourceID       int     `json:"source_id"`
	SourceType     int     `json:"source_type"`
	SourceIcon     string  `json:"source_icon"`
	SourceName     string  `json:"source_name"`
	ListenProgress float64 `json:"listen_progress"`
	ListenFinished bool    `json:"listen_finished"`
	DdArticleID    string  `json:"dd_article_id"`
	AudioListIcon  string  `json:"audio_list_icon"`
	ClassCourseID  string  `json:"class_course_id"`
	ClassArticleID string  `json:"class_article_id"`
	LogID          string  `json:"log_id"`
	LogType        string  `json:"log_type"`
	LogInterface   string  `json:"log_interface"`
	Trackinfo      string  `json:"trackinfo"`
	UsedDrm        int     `json:"used_drm"`
	IndexImg       string  `json:"index_img"`
	Reader         string  `json:"reader"`
	ReaderName     string  `json:"reader_name"`
}

type Video struct {
	Token            string  `json:"token"`
	TokenVersion     int     `json:"token_version"`
	CoverImg         string  `json:"cover_img"`
	DdMediaID        int64   `json:"dd_media_id"`
	DdMediaIDStr     string  `json:"dd_media_id_str"`
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
}

// AudioList audio baseic info list
type AudioList struct {
	List []Audio `json:"list"`
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
