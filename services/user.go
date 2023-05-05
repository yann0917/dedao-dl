package services

// User user info
type User struct {
	Nickname        string  `json:"nickname"`
	Avatar          string  `json:"avatar"`
	TodayStudyTime  int     `json:"today_study_time"`
	StudySerialDays int     `json:"study_serial_days"`
	IsV             int     `JSON:"is_v"`
	VIPUser         VIPUser `json:"vip_user"`
	IsTeacher       int     `json:"is_teacher"`
	UIDHazy         string  `json:"uid_hazy"`
}

// VIPUser vip info
type VIPUser struct {
	Info string `json:"info"`
	Stat int    `json:"stat"`
}
type OdobVipUser struct {
	Card []struct {
		AvailableBuy     bool   `json:"available_buy"`
		CanGift          int    `json:"can_gift"`
		CardType         int    `json:"card_type"`
		ContractTemplate string `json:"contract_template"`
		Cover            string `json:"cover"`
		CreatedAt        string `json:"created_at"`
		Days             int    `json:"days"`
		Description      string `json:"description"`
		Device           string `json:"device"`
		DiscountTip      string `json:"discount_tip"`
		Displayed        int    `json:"displayed"`
		FirstMonthPrice  string `json:"first_month_price"`
		ID               int    `json:"id"`
		IosSwitch        bool   `json:"ios_switch"`
		IsFirstSubscribe bool   `json:"is_first_subscribe"`
		IsSubscribed     int    `json:"is_subscribed"`
		MaxNum           int    `json:"max_num"`
		Name             string `json:"name"`
		OriginPrice      string `json:"origin_price"`
		Price            string `json:"price"`
		PriceDesc        string `json:"price_desc"`
		Ptype            int64  `json:"ptype"`
		Rights           []struct {
			Name  string `json:"name"`
			Right bool   `json:"right"`
		} `json:"rights"`
		Selected        int    `json:"selected"`
		Sku             string `json:"sku"`
		SortNo          int    `json:"sort_no"`
		Status          int    `json:"status"`
		StyleType       int    `json:"style_type"`
		SubscribeDesc   string `json:"subscribe_desc"`
		SubscribedFlag  string `json:"subscribed_flag"`
		Type            int    `json:"type"`
		UpdatedAt       string `json:"updated_at"`
		UsedCardCount   int    `json:"used_card_count"`
		UserMaxBuyCount int    `json:"user_max_buy_count"`
		WelfareInfo     string `json:"welfare_info"`
	} `json:"card"`
	User struct {
		Avatar            string `json:"avatar"`
		AvatarS           string `json:"avatar_s"`
		BeginTime         int    `json:"begin_time"`
		CanGivenCount     int    `json:"can_given_count"`
		CardID            int    `json:"card_id"`
		DdURL             string `json:"dd_url"`
		EndTime           int64  `json:"end_time"`
		EnterpriseEndTime int    `json:"enterprise_end_time"`
		ErrTips           string `json:"err_tips"`
		ExpireTime        int    `json:"expire_time"`
		GivenCardCount    int    `json:"given_card_count"`
		IsExpire          bool   `json:"is_expire"`
		IsVip             bool   `json:"is_vip"`
		Nickname          string `json:"nickname"`
		PriceDesc         string `json:"price_desc"`
		SavePrice         string `json:"save_price"`
		Slogan            string `json:"slogan"`
		SurplusTime       int64  `json:"surplus_time"`
		TotalCount        int    `json:"total_count"`
		UID               int    `json:"uid"`
		VStateValue       int    `json:"v_state_value"`
		WeekCount         int    `json:"week_count"`
	}
}

// User get user info
func (s *Service) User() (user *User, err error) {
	body, err := s.reqUser()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &user); err != nil {
		return
	}
	return
}

// Token token
type Token struct {
	Token string `json:"token"`
}

// Token get token
func (s *Service) Token() (t *Token, err error) {
	body, err := s.reqToken()
	if err != nil {
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &t); err != nil {
		return
	}
	return
}
