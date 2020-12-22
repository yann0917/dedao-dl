package services

import "time"

// Catelog ebook catalog
type Catelog struct {
	Level     int    `json:"level"`
	Text      string `json:"text"`
	Href      string `json:"href"`
	PlayOrder int    `json:"playOrder"`
}

// Press ebook press info
type Press struct {
	Name  string `json:"name"`
	Brief string `json:"brief"`
}

// EbookDetail ebook detail
type EbookDetail struct {
	ID                  int           `json:"id"`
	Title               string        `json:"title"`
	Style               int           `json:"style"`
	Cover               string        `json:"cover"`
	Count               int           `json:"count"`
	Price               string        `json:"price"`
	Status              int           `json:"status"`
	OperatingTitle      string        `json:"operating_title"`
	OtherShareTitle     string        `json:"other_share_title"`
	OtherShareSummary   string        `json:"other_share_summary"`
	AuthorInfo          string        `json:"author_info"`
	BookAuthor          string        `json:"book_author"`
	PublishTime         string        `json:"publish_time"`
	CatalogList         []Catelog     `json:"catalog_list"`
	BookIntro           string        `json:"book_intro"`
	BSpecialPrice       string        `json:"b_special_price"`
	CurrentPrice        string        `json:"current_price"`
	IsBuy               bool          `json:"is_buy"`
	IsTrial             bool          `json:"is_trial"`
	IsTtsSwitch         bool          `json:"is_tts_switch"`
	LogID               string        `json:"log_id"`
	LogType             string        `json:"log_type"`
	OriginalPrice       string        `json:"original_price"`
	AuthorList          []string      `json:"author_list"`
	CanTrialRead        bool          `json:"can_trial_read"`
	TrialReadProportion string        `json:"trial_read_proportion"`
	WithVideo           bool          `json:"with_video"`
	Enid                string        `json:"enid"`
	BOverseasPurchase   int           `json:"b_overseas_purchase"`
	RankName            string        `json:"rank_name"`
	RankNum             int           `json:"rank_num"`
	IsVipBook           int           `json:"is_vip_book"`
	IsOnBookshelf       bool          `json:"is_on_bookshelf"`
	ProductScore        string        `json:"product_score"`
	ReadTime            int           `json:"read_time"`
	ReadNumber          []interface{} `json:"read_number"`
	Press               Press         `json:"press"`
	DoubanScore         string        `json:"douban_score"`
	ClassifyName        string        `json:"classify_name"`
	ClassifyID          int           `json:"classify_id"`
	AddStudylistDdURL   string        `json:"add_studylist_dd_url"`
}

// EbookBlock ebook block
type EbookBlock struct {
	ChapterID   string `json:"chapterId"`
	SectionID   string `json:"sectionID"`
	EndOffset   int    `json:"endOffset"`
	StartOffset int    `json:"startOffset"`
}

// EbookOrders ebook orders
type EbookOrders struct {
	ChapterID  string `json:"chapterId"`
	PathInEpub string `json:"pathInEpub"`
}

// EbookToc ebook toc
type EbookToc struct {
	Href      string `json:"href"`
	Level     int    `json:"level"`
	PlayOrder int    `json:"playOrder"`
	Offset    int    `json:"offset"`
	Text      string `json:"text"`
}

// EbookInfo ebook info
type EbookInfo struct {
	BookInfo struct {
		EbookBlock [][]EbookBlock `json:"block"`
		Orders     []EbookOrders  `json:"orders"`
		Toc        []EbookToc     `json:"toc"`
	} `json:"bookInfo"`
}

// EbookDetail get ebook detail
func (s *Service) EbookDetail(enid string) (detail *EbookDetail, err error) {
	cacheFile := "ebookDetail:" + enid
	x, ok := detail.getCache(cacheFile)
	if ok {
		detail = x.(*EbookDetail)
		return
	}
	body, err := s.reqEbookDetail(enid)
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

// EbookReadToken get ebook read token
func (s *Service) EbookReadToken(enid string) (t *Token, err error) {
	body, err := s.reqEbookReadToken(enid)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &t); err != nil {
		return
	}
	return
}

// EbookInfo get ebook info
//
// include book block, book TOC, epubPath etc
func (s *Service) EbookInfo(token string) (info *EbookInfo, err error) {
	body, err := s.reqEbookInfo(token)
	defer body.Close()
	if err != nil {
		return
	}
	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

func (c *EbookDetail) getCacheKey() string {
	return "ebookDetail"
}

func (c *EbookDetail) getCache(fileName string) (interface{}, bool) {
	err := LoadCacheFile(fileName)
	if err != nil {
		return nil, false
	}
	x, ok := Cache.Get(cacheKey(c))
	return x, ok
}

func (c *EbookDetail) setCache(fileName string) error {
	Cache.Set(cacheKey(c), c, 1*time.Hour)
	err := SaveCacheFile(fileName)
	return err
}
