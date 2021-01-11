package services

// MediaBaseInfo media info
type MediaBaseInfo struct {
	MediaType     int    `json:"media_type"`
	SourceID      string `json:"source_id"`
	SecurityToken string `json:"security_token"`
}

// Audio audio
type Audio struct {
	AliasID        string `json:"alias_id"`
	Extrainfo      string `json:"extrainfo"`
	ClassID        int    `json:"class_id"`
	Title          string `json:"title"`
	ShareTitle     string `json:"share_title"`
	Mp3PlayURL     string `json:"mp3_play_url"`
	Duration       int    `json:"duration"`
	Schedule       int    `json:"schedule"`
	TopicID        int    `json:"topic_id"`
	Summary        string `json:"summary"`
	Price          int    `json:"price"`
	Icon           string `json:"icon"`
	Size           int    `json:"size"`
	Etag           string `json:"etag"`
	Token          string `json:"token"`
	ShareSummary   string `json:"share_summary"`
	Collection     int    `json:"collection"`
	Count          int    `json:"count"`
	BoredCount     int    `json:"bored_count"`
	AudioType      int    `json:"audio_type"`
	DrmVersion     int    `json:"drm_version"`
	SourceID       int    `json:"source_id"`
	SourceType     int    `json:"source_type"`
	SourceIcon     string `json:"source_icon"`
	SourceName     string `json:"source_name"`
	ListenProgress int    `json:"listen_progress"`
	ListenFinished bool   `json:"listen_finished"`
	DdArticleID    string `json:"dd_article_id"`
	AudioListIcon  string `json:"audio_list_icon"`
	ClassCourseID  string `json:"class_course_id"`
	ClassArticleID string `json:"class_article_id"`
	LogID          string `json:"log_id"`
	LogType        string `json:"log_type"`
	LogInterface   string `json:"log_interface"`
	Trackinfo      string `json:"trackinfo"`
	UsedDrm        int    `json:"used_drm"`
	IndexImg       string `json:"index_img"`
	Reader         string `json:"reader"`
	ReaderName     string `json:"reader_name"`
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
