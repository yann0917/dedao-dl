package services

type RankSquareBaseInfo struct {
	Header  RankSquareHeader `json:"header"`
	NavList []RankNavGroup   `json:"nav_list"`
}

type RankSquareHeader struct {
	IndexImage string `json:"index_image"`
}

type RankNavGroup struct {
	PType        int            `json:"ptype"`
	Name         string         `json:"name"`
	RankTypeList []RankTypeInfo `json:"rank_type_list"`
}

type RankTypeInfo struct {
	RankID   int    `json:"rank_id"`
	Title    string `json:"title"`
	SubTitle string `json:"sub_title"`
}

type RankListData struct {
	List []RankBoard `json:"list"`
}

type RankBoard struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	SubTitle        string         `json:"sub_title"`
	Count           int            `json:"count"`
	RNImg           string         `json:"rn_img"`
	RNLeftTitle     string         `json:"rn_left_title"`
	RNRightTitle    string         `json:"rn_right_title"`
	RankType        int            `json:"rank_type"`
	BackgroundColor string         `json:"background_color"`
	ResourceType    int            `json:"resource_type"`
	Tips            string         `json:"tips"`
	ShareImg        string         `json:"share_img"`
	DDURL           string         `json:"dd_url"`
	LogID           int            `json:"log_id"`
	LogType         string         `json:"log_type"`
	List            []RankItemCard `json:"list"`
}

type RankItemCard struct {
	ProductID         int                 `json:"product_id"`
	ProductType       int                 `json:"product_type"`
	OriginType        int                 `json:"origin_type"`
	OriginID          int                 `json:"origin_id"`
	TypeName          string              `json:"type_name"`
	IndexImg          string              `json:"index_img"`
	SquareImg         string              `json:"square_img"`
	Logo              string              `json:"logo"`
	Title             string              `json:"title"`
	RecommendTitle    string              `json:"recommend_title"`
	RecommendIntro    string              `json:"recommend_intro"`
	Intro             string              `json:"intro"`
	Author            string              `json:"author"`
	Status            int                 `json:"status"`
	LogID             int                 `json:"log_id"`
	LogType           string              `json:"log_type"`
	TrackInfo         string              `json:"track_info"`
	DDURL             string              `json:"dd_url"`
	Metrics           string              `json:"metrics"`
	Resource          RankItemResource    `json:"resource"`
	HotIntro          RankHotIntro        `json:"hot_intro"`
	CollectionInfo    RankCollectionInfo  `json:"collection_info"`
	ProgressIntro     RankProgressIntro   `json:"progress_intro"`
	CostIntro         RankCostIntro       `json:"cost_intro"`
	AuthorityIntro    RankAuthorityIntro  `json:"authority_intro"`
	AuthorInfo        RankAuthorInfo      `json:"author_info"`
	CornerLabelInfo   RankCornerLabelInfo `json:"corner_label_info"`
	CornerImgVertical string              `json:"corner_img_vertical"`
	OnlineTime        string              `json:"online_time"`
	IsToday           bool                `json:"is_today"`
	IsShowNewbook     bool                `json:"is_show_newbook"`
	RecommendEbook    bool                `json:"recommend_ebook"`
	DoubanInfo        RankDoubanInfo      `json:"douban_info"`
	LearnCountInfo    RankLearnCountInfo  `json:"learn_count_info"`
}

type RankItemResource struct {
	AudioID      string `json:"audio_id"`
	ResourceID   int    `json:"resource_id"`
	ResourceType int    `json:"resource_type"`
	DDURL        string `json:"dd_url"`
}

type RankHotIntro struct {
	Number         int    `json:"number"`
	Intro          string `json:"intro"`
	Score          string `json:"score"`
	UserScoreCount int    `json:"user_score_count"`
}

type RankCollectionInfo struct {
	Number int    `json:"number"`
	Intro  string `json:"intro"`
}

type RankProgressIntro struct {
	Intro       string `json:"intro"`
	Progress    int    `json:"progress"`
	MaxProgress int    `json:"max_progress"`
	IsFinish    int    `json:"is_finish"`
	Uint        string `json:"uint"`
}

type RankCostIntro struct {
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	CouponPrice   string `json:"coupon_price"`
	Tag           string `json:"tag"`
	TagColor      int    `json:"tag_color"`
}

type RankAuthorityIntro struct {
	IsBuy                bool `json:"is_buy"`
	IsVip                bool `json:"is_vip"`
	IsVipExpired         bool `json:"is_vip_expired"`
	HasAuthority         bool `json:"has_authority"`
	ConsumedNum          int  `json:"consumed_num"`
	FreeMaximum          int  `json:"free_maximum"`
	InBookRack           bool `json:"in_book_rack"`
	InLimitFreeList      bool `json:"in_limit_free_list"`
	InVipPool            bool `json:"in_vip_pool"`
	FreeBeginTime        int  `json:"free_begin_time"`
	FreeEndTime          int  `json:"free_end_time"`
	IsCanBuy             bool `json:"is_can_buy"`
	BSellingChannelGroup int  `json:"b_selling_channel_group"`
	IsSubscribe          bool `json:"is_subscribe"`
	IsRedPacketTry       bool `json:"is_red_packet_try"`
}

type RankAuthorInfo struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Intro  string `json:"intro"`
	IsAuth int    `json:"is_auth"`
	UID    int    `json:"uid"`
	DDURL  string `json:"dd_url"`
}

type RankCornerLabelInfo struct {
	LeftTop     string `json:"left_top"`
	LeftBottom  string `json:"left_bottom"`
	RightTop    string `json:"right_top"`
	RightBottom string `json:"right_bottom"`
}

type RankDoubanInfo struct {
	Score      string `json:"score"`
	UseDouban  bool   `json:"use_douban"`
	DoubanIcon string `json:"douban_icon"`
}

type RankLearnCountInfo struct {
	LearnCount     int    `json:"learn_count"`
	LearnCountDesc string `json:"learn_count_desc"`
}

func (s *Service) RankSquareBaseInfo() (resp *RankSquareBaseInfo, err error) {
	response, err := s.client.R().
		SetBody(map[string]any{}).
		Post(mobileBaseURL + "/native/api/rankSquare/baseInfo")
	reader, err := handleHTTPResponse(response, err)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if err = handleJSONParse(reader, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) RankListData(rankType int) (resp *RankListData, err error) {
	response, err := s.client.R().
		SetBody(map[string]any{
			"rank_type": rankType,
		}).
		Post(mobileBaseURL + "/native/api/rankListData")
	reader, err := handleHTTPResponse(response, err)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if err = handleJSONParse(reader, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
