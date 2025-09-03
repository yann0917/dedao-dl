package services

// CourseList 课程列表V2数据
type CourseList struct {
	List          []Course `json:"list"`            // 课程列表
	Total         int      `json:"total"`           // 总数量
	IsMore        int      `json:"is_more"`         // 是否有更多
	HasSingleBook bool     `json:"has_single_book"` // 是否有单本书籍
}

// CourseV2 课程V2数据结构
type Course struct {
	Enid             string           `json:"enid"`
	ID               int              `json:"id"`
	Type             int              `json:"type"` // 13-单本,1013-名家讲书合集
	ClassType        int              `json:"class_type"`
	ClassID          int              `json:"class_id"`
	HasExtra         bool             `json:"has_extra"`
	ClassFinished    bool             `json:"class_finished"`
	Title            string           `json:"title"`
	Intro            string           `json:"intro"`
	Author           string           `json:"author"`
	Icon             string           `json:"icon"`
	CreateTime       int              `json:"create_time"`
	LastRead         string           `json:"last_read"`
	Progress         int              `json:"progress"`
	Duration         int              `json:"duration"`
	CourseNum        int              `json:"course_num"`
	PublishNum       int              `json:"publish_num"`
	LogID            string           `json:"log_id"`
	LogType          string           `json:"log_type"`
	IsTop            int              `json:"is_top"`
	LastActionTime   int              `json:"last_action_time"`
	IsNew            int              `json:"is_new"`
	IsFinished       int              `json:"is_finished"`
	Size             string           `json:"size"`
	DdURL            string           `json:"dd_url"`
	AssetsType       int              `json:"assets_type"`
	DrmToken         string           `json:"drm_token"`
	AudioDetail      Audio            `json:"audio_detail"`
	ProductPrice     int              `json:"product_price"`
	Price            string           `json:"price"`
	ProductIntro     string           `json:"product_intro"`
	HasPlayAuth      bool             `json:"has_play_auth"`
	ExtInfo          interface{}      `json:"ext_info"`
	Status           int              `json:"status"`
	DdExtURL         string           `json:"dd_ext_url"`
	IsCollected      bool             `json:"is_collected"`
	WendaExtInfo     WendaExtInfo     `json:"wenda_ext_info"`
	ClassExtReview   ClassExtReview   `json:"class_ext_review"`
	PlanStatus       int              `json:"plan_status"`
	IsSphere         int              `json:"is_sphere"`
	InstituteExtInfo InstituteExtInfo `json:"institute_ext_info"`
	TrainingCampInfo TrainingCampInfo `json:"training_camp_ext_info"`
	LastReadInfo     string           `json:"last_read_info"`
	OdobGroupExtInfo OdobGroupExtInfo `json:"odob_group_ext_info"`
	EbookExtInfo     EbookExtInfo     `json:"ebook_ext_info"`
	IsVideoOdob      bool             `json:"is_video_odob"`
	IsShowPlayLater  bool             `json:"is_show_play_later_button"`
	IsAddPlayLater   bool             `json:"is_add_play_later"`
	EnPackageID      string           `json:"en_package_id,omitempty"`
	IsGroup          bool             `json:"is_group"`
	GroupID          int              `json:"group_id"`
	IsSelfBuildGroup bool             `json:"is_self_build_group"`
	GroupType        int              `json:"group_type"`
	LabelID          int              `json:"label_id"`
}

// WendaExtInfo 问答扩展信息
type WendaExtInfo struct {
	AnswerID int `json:"answer_id"`
}

// ClassExtReview 课程评价扩展信息
type ClassExtReview struct {
	ShowText    string `json:"show_text"`
	ReviewDone  bool   `json:"review_done"`
	IsSatisfied bool   `json:"is_satisfied"`
	Detail      string `json:"detail"`
}

// InstituteExtInfo 机构扩展信息
type InstituteExtInfo struct {
	CourseID   int `json:"course_id"`
	CourseType int `json:"course_type"`
}

// TrainingCampInfo 训练营扩展信息
type TrainingCampInfo struct {
	CourseID          int `json:"course_id"`
	CourseType        int `json:"course_type"`
	CourseSalesNature int `json:"course_sales_nature"`
}

// OdobGroupExtInfo 每天听本书分组扩展信息
type OdobGroupExtInfo struct {
	GroupType            int      `json:"group_type"` // 1-名家讲书
	GroupTypeName        string   `json:"group_type_name"`
	OdobTotalDuration    int      `json:"odob_total_duration"`
	OdobAliasList        []string `json:"odob_alias_list"` // 音频 alias_id
	ProgressPercent      int      `json:"progress_percent"`
	ProgressLearnedCount int      `json:"progress_learned_count"`
	PublishTime          string   `json:"publish_time"`
	PublishStatus        int      `json:"publish_status"`
	AudioDetail          Audio    `json:"audio_detail"`
}

// EbookExtInfo 电子书扩展信息
type EbookExtInfo struct {
	IsTtsSwitch bool `json:"is_tts_switch"`
}

// OutsideDetail 名家讲书课程详情
type OutsideDetail struct {
	SPU          OutsideSPU    `json:"spu"`           // 课程SPU信息
	Items        []OutsideItem `json:"items"`         // 课程包含的音频列表
	Count        int           `json:"count"`         // 总数量
	CurrentCount int           `json:"current_count"` // 当前数量
	NowTime      string        `json:"now_time"`      // 当前时间
	HasMore      bool          `json:"has_more"`      // 是否有更多
	UpMore       bool          `json:"up_more"`       // 是否有上一页
	DownMore     bool          `json:"down_more"`     // 是否有下一页
}

// OutsideSPU 名家讲书课程SPU信息
type OutsideSPU struct {
	BizType              int                 `json:"biz_type"`                // 业务类型
	PType                int                 `json:"ptype"`                   // 产品类型
	PID                  int                 `json:"pid"`                     // 产品ID
	PTypePID             string              `json:"ptype_pid"`               // 产品类型-产品ID
	Title                string              `json:"title"`                   // 标题
	Summary              string              `json:"summary"`                 // 摘要
	Intro                string              `json:"intro"`                   // 介绍
	Icon                 string              `json:"icon"`                    // 图标
	LessonsNum           int                 `json:"lessons_num"`             // 课程数量
	Department           []string            `json:"department"`              // 部门
	FirstItemPublishTime string              `json:"first_item_publish_time"` // 首个课程发布时间
	CreateTime           string              `json:"create_time"`             // 创建时间
	UpdateTime           string              `json:"update_time"`             // 更新时间
	SPUCode              string              `json:"spu_code"`                // SPU代码
	Creator              string              `json:"creator"`                 // 创建者
	DDID                 int                 `json:"ddid"`                    // 得到ID
	ShareTitle           string              `json:"share_title"`             // 分享标题
	ShareSummary         string              `json:"share_summary"`           // 分享摘要
	ShareURL             string              `json:"share_url"`               // 分享URL
	DDURL                string              `json:"ddurl"`                   // 得到URL
	AuthorUID            int                 `json:"author_uid"`              // 作者UID
	SKUInfo              OutsideSKUInfo      `json:"sku_info"`                // SKU信息
	RelationInfo         OutsideRelationInfo `json:"relation_info"`           // 关系信息
	GroupTypeID          int                 `json:"group_type_id"`           // 分组类型ID
	GroupTypeName        string              `json:"group_type_name"`         // 分组类型名称
	PublishTime          string              `json:"publish_time"`            // 发布时间
	Extra                OutsideSPUExtra     `json:"extra"`                   // 额外信息
}

// OutsideSKUInfo 名家讲书课程SKU信息
type OutsideSKUInfo struct {
	SKUCode         string `json:"sku_code"`         // SKU代码
	PTypePID        string `json:"ptype_pid"`        // 产品类型-产品ID
	Status          int    `json:"status"`           // 状态
	ChannelType     string `json:"channel_type"`     // 渠道类型
	PublishTime     string `json:"publish_time"`     // 发布时间
	SPUCode         string `json:"spu_code"`         // SPU代码
	PType           int    `json:"ptype"`            // 产品类型
	IsSale          int    `json:"is_sale"`          // 是否在售
	ProductFeatures int    `json:"product_features"` // 产品特性
	Remark          string `json:"remark"`           // 备注
	Price           string `json:"price"`            // 价格
}

// OutsideRelationInfo 名家讲书课程关系信息
type OutsideRelationInfo struct {
	ParentPK string `json:"parent_pk"` // 父级PK
	Status   int    `json:"status"`    // 状态
	SortID   int    `json:"sort_id"`   // 排序ID
	Title    string `json:"title"`     // 标题
	SKUCode  string `json:"sku_code"`  // SKU代码
	Price    string `json:"price"`     // 价格
}

// OutsideSPUExtra 名家讲书课程SPU额外信息
type OutsideSPUExtra struct {
	Enid                 string `json:"enid"`                   // 加密ID
	BookShelfStatus      bool   `json:"book_shelf_status"`      // 书架状态
	GroupType            int    `json:"group_type"`             // 分组类型
	GroupTypeName        string `json:"group_type_name"`        // 分组类型名称
	IntroText            string `json:"intro_text"`             // 介绍文本
	LatestLearningID     int    `json:"latest_learning_id"`     // 最新学习ID
	NewestIntro          string `json:"newest_intro"`           // 最新介绍
	OdobConsumerNum      int    `json:"odob_consumer_num"`      // 每天听本书消费者数量
	OdobFreeMaximum      int    `json:"odob_free_maximum"`      // 每天听本书免费最大数量
	ProgressPercent      int    `json:"progress_percent"`       // 进度百分比
	QcgID                int    `json:"qcg_id"`                 // QCG ID
	RnLearnCount         int    `json:"rn_learn_count"`         // 学习计数
	RnLearnCountDesc     string `json:"rn_learn_count_desc"`    // 学习计数描述
	SubscribeStatus      int    `json:"subscribe_status"`       // 订阅状态
	LiveTag              string `json:"live_tag"`               // 直播标签
	TeacherAvatar        string `json:"teacher_avatar"`         // 老师头像
	TeacherDDURL         string `json:"teacher_ddurl"`          // 老师得到URL
	TeacherIntro         string `json:"teacher_intro"`          // 老师介绍
	TeacherName          string `json:"teacher_name"`           // 老师名称
	ProgressLearnedCount int    `json:"progress_learned_count"` // 进度学习计数
}

// OutsideItem 名家讲书课程子项目信息
type OutsideItem struct {
	BizType              int                 `json:"biz_type"`                // 业务类型
	PType                int                 `json:"ptype"`                   // 产品类型
	PID                  int                 `json:"pid"`                     // 产品ID
	PTypePID             string              `json:"ptype_pid"`               // 产品类型-产品ID
	Title                string              `json:"title"`                   // 标题
	Summary              string              `json:"summary"`                 // 摘要
	Intro                string              `json:"intro"`                   // 介绍
	Icon                 string              `json:"icon"`                    // 图标
	LessonsNum           int                 `json:"lessons_num"`             // 课程数量
	Department           []string            `json:"department"`              // 部门
	FirstItemPublishTime string              `json:"first_item_publish_time"` // 首个课程发布时间
	CreateTime           string              `json:"create_time"`             // 创建时间
	UpdateTime           string              `json:"update_time"`             // 更新时间
	SPUCode              string              `json:"spu_code"`                // SPU代码
	Creator              string              `json:"creator"`                 // 创建者
	DDID                 int                 `json:"ddid"`                    // 得到ID
	ShareTitle           string              `json:"share_title"`             // 分享标题
	ShareSummary         string              `json:"share_summary"`           // 分享摘要
	ShareURL             string              `json:"share_url"`               // 分享URL
	DDURL                string              `json:"ddurl"`                   // 得到URL
	AuthorUID            int                 `json:"author_uid"`              // 作者UID
	SKUInfo              OutsideSKUInfo      `json:"sku_info"`                // SKU信息
	RelationInfo         OutsideRelationInfo `json:"relation_info"`           // 关系信息
	GroupTypeID          int                 `json:"group_type_id"`           // 分组类型ID
	GroupTypeName        string              `json:"group_type_name"`         // 分组类型名称
	PublishTime          string              `json:"publish_time"`            // 发布时间
	Extra                OutsideItemExtra    `json:"extra"`                   // 额外信息
}

// OutsideItemExtra 名家讲书课程子项目额外信息
type OutsideItemExtra struct {
	Enid                string         `json:"enid"`                  // 加密ID
	Duration            int            `json:"duration"`              // 时长
	GroupType           int            `json:"group_type"`            // 分组类型
	GroupTypeName       string         `json:"group_type_name"`       // 分组类型名称
	IsEbookVIP          bool           `json:"is_ebook_vip"`          // 是否电子书VIP
	IsOdobVIP           bool           `json:"is_odob_vip"`           // 是否每天听本书VIP
	LiveInfoDetail      LiveInfoDetail `json:"live_info_detail"`      // 直播信息详情
	AudioAliasID        string         `json:"audio_alias_id"`        // 音频别名ID
	Equity              bool           `json:"equity"`                // 权益
	OdobAudioDetail     Audio          `json:"odob_audio_detail"`     // 每天听本书音频详情
	ProgressMaxProgress int            `json:"progress_max_progress"` // 进度最大值
	ProgressLearned     bool           `json:"progress_learned"`      // 是否已学习
	ProgressIsFinish    bool           `json:"progress_is_finish"`    // 是否已完成
	PublishTime         int            `json:"publish_time"`          // 发布时间
}

// LiveInfoDetail 直播信息详情
type LiveInfoDetail struct {
	AlertTips            string           `json:"alert_tips"`             // 提醒提示
	AliasID              string           `json:"alias_id"`               // 别名ID
	AppointmentStatus    int              `json:"appointment_status"`     // 预约状态
	Author               string           `json:"author"`                 // 作者
	CanWatch             bool             `json:"can_watch"`              // 是否可以观看
	DDURL                string           `json:"dd_url"`                 // 得到URL
	Duration             int              `json:"duration"`               // 时长
	Endtime              int              `json:"endtime"`                // 结束时间
	HasBuy               bool             `json:"has_buy"`                // 是否已购买
	HideOnlineNumber     int              `json:"hide_online_number"`     // 隐藏在线人数
	HomeImg              string           `json:"home_img"`               // 首页图片
	ID                   int              `json:"id"`                     // ID
	Intro                string           `json:"intro"`                  // 介绍
	InviteCount          int              `json:"invite_count"`           // 邀请计数
	IsPrivilegeLive      bool             `json:"is_privilege_live"`      // 是否特权直播
	LdFlv                string           `json:"ld_flv"`                 // 低清FLV
	LiveCoverM3U8        string           `json:"live_cover_m3u8"`        // 直播封面M3U8
	LiveDurationText     string           `json:"live_duration_text"`     // 直播时长文本
	LivePrivilegeTips    string           `json:"live_privilege_tips"`    // 直播特权提示
	LivePrivilegeType    int              `json:"live_privilege_type"`    // 直播特权类型
	LivePV               int              `json:"live_pv"`                // 直播PV
	LiveType             string           `json:"live_type"`              // 直播类型
	LiveViewers          int              `json:"live_viewers"`           // 直播观看者
	LogID                int              `json:"log_id"`                 // 日志ID
	LogType              string           `json:"log_type"`               // 日志类型
	Logo                 string           `json:"logo"`                   // 标志
	OnlineNum            int              `json:"online_num"`             // 在线人数
	PlaybackStatus       int              `json:"playback_status"`        // 回放状态
	PrivilegeLiveTag     string           `json:"privilege_live_tag"`     // 特权直播标签
	PrivilegeProduct     PrivilegeProduct `json:"privilege_product"`      // 特权产品
	PrivilegeStatus      int              `json:"privilege_status"`       // 特权状态
	PublishStatus        int              `json:"publish_status"`         // 发布状态
	PVNum                int              `json:"pv_num"`                 // PV数量
	ReservationNum       int              `json:"reservation_num"`        // 预约数量
	RoomID               int              `json:"room_id"`                // 房间ID
	ShowPV               int              `json:"show_pv"`                // 显示PV
	ShowSubscribeSummary int              `json:"show_subscribe_summary"` // 显示订阅摘要
	Starttime            int              `json:"starttime"`              // 开始时间
	StarttimeDesc        string           `json:"starttime_desc"`         // 开始时间描述
	Status               int              `json:"status"`                 // 状态
	SubscribeSummary     string           `json:"subscribe_summary"`      // 订阅摘要
	TimeReport           int              `json:"time_report"`            // 时间报告
	Title                string           `json:"title"`                  // 标题
	Type                 int              `json:"type"`                   // 类型
	VideoCoverM3U8       string           `json:"video_cover_m3u8"`       // 视频封面M3U8
	VideoCoverMediaID    int              `json:"video_cover_media_id"`   // 视频封面媒体ID
	VideoDuration        string           `json:"video_duration"`         // 视频时长
	WebPCMediaToken      string           `json:"web_pc_media_token"`     // Web PC媒体令牌
}

// PrivilegeProduct 特权产品
type PrivilegeProduct struct {
	AliasID         string `json:"alias_id"`          // 别名ID
	AliasName       string `json:"alias_name"`        // 别名名称
	DDURL           string `json:"ddurl"`             // 得到URL
	HasPrivilege    bool   `json:"has_privilege"`     // 是否有特权
	ID              int    `json:"id"`                // ID
	IsPrivilegeLive bool   `json:"is_privilege_live"` // 是否特权直播
	LiveID          int    `json:"live_id"`           // 直播ID
	LiveStartTime   int    `json:"live_start_time"`   // 直播开始时间
	ProductGroupID  int    `json:"product_group_id"`  // 产品分组ID
	ProductID       int    `json:"product_id"`        // 产品ID
	ProductSubID    int    `json:"product_sub_id"`    // 产品子ID
	ProductType     int    `json:"product_type"`      // 产品类型
	Title           string `json:"title"`             // 标题
}

// TopicPkgOdobDetails 名家讲书每天听本书音频集合详情
type TopicPkgOdobDetails struct {
	OdobAudioDetailList []Audio `json:"odob_audio_detail_list"` // 每天听本书音频详情列表
}
