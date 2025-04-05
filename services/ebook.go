package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yann0917/dedao-dl/utils"
)

const (
	maxRetries     = 3
	initialBackoff = 3 * time.Second
	retryDelay     = 1 * time.Second
)

// Add new cache-related types
type PageCache struct {
	ChapterID string    `json:"chapter_id"`
	Pages     []string  `json:"pages"`
	Timestamp time.Time `json:"timestamp"`
}

// Make cache functions public
func GetCachePath(enid, chapterID string) string {
	cacheDir := filepath.Join("output", ".cache", "pages", enid)
	os.MkdirAll(cacheDir, 0755)
	return filepath.Join(cacheDir, fmt.Sprintf("%s.json", chapterID))
}

func LoadFromCache(enid, chapterID string) ([]string, bool) {
	cachePath := GetCachePath(enid, chapterID)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, false
	}

	var cache PageCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, false
	}

	// Optional: Check if cache is too old
	if time.Since(cache.Timestamp) > 24*time.Hour {
		os.Remove(cachePath)
		return nil, false
	}

	return cache.Pages, true
}

func SaveToCache(enid, chapterID string, pages []string) error {
	cache := PageCache{
		ChapterID: chapterID,
		Pages:     pages,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(GetCachePath(enid, chapterID), data, 0644)
}

func ClearBookCache(enid string) error {
	cacheDir := filepath.Join("output", ".cache", "pages", enid)
	return os.RemoveAll(cacheDir)
}

func ClearAllCache() error {
	cacheDir := filepath.Join("output", ".cache", "pages")
	return os.RemoveAll(cacheDir)
}

func withRetry[T any](operation func() (T, error), chapterID string) (result T, err error) {
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		result, err = operation()
		if err == nil {
			return result, nil
		}

		// Print detailed error information
		fmt.Printf("\nAttempt %d/%d failed for chapter %s:\n", i+1, maxRetries, chapterID)
		fmt.Printf("Error: %v\n", err)

		if i < maxRetries-1 {
			fmt.Printf("Retrying in %v...\n", backoff)
			time.Sleep(backoff)
			// Double the backoff for next attempt
			backoff *= 2
		} else {
			fmt.Printf("Max retries reached for chapter %s, giving up.\n", chapterID)
		}
	}
	return result, err
}

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

type EbookInfoPage struct {
	Cid         string `json:"cid"`
	EndOffset   int    `json:"end_offset"`
	PageNum     int    `json:"page_num"`
	StartOffset int    `json:"start_offset"`
}

type EbookPage struct {
	IsEnd bool `json:"is_end"`
	Pages []struct {
		BeginOffset           int64  `json:"begin_offset"`
		EndOffset             int64  `json:"end_offset"`
		IsFirst               bool   `json:"is_first"`
		IsLast                bool   `json:"is_last"`
		Svg                   string `json:"svg"`
		ViewHeighToChapterTop int64  `json:"view_heigh_to_chapter_top"`
	} `json:"pages"`
}

// EbookInfo ebook info
type EbookInfo struct {
	BookInfo struct {
		EbookBlock [][]EbookBlock   `json:"block"`
		Orders     []EbookOrders    `json:"orders"`
		Toc        []utils.EbookToc `json:"toc"`
		Pages      []EbookInfoPage  `json:"pages"`
	} `json:"bookInfo"`
}

// EbookVIPInfo ebook vip info
type EbookVIPInfo struct {
	UID                int    `json:"uid"`
	Nickname           string `json:"nickname"`
	Slogan             string `json:"slogan"`
	Avatar             string `json:"avatar"`
	AvatarS            string `json:"avatar_s"`
	MonthCount         int    `json:"month_count"`
	TotalCount         int    `json:"total_count"`
	FinishedCount      int    `json:"finished_count"`
	SavePrice          string `json:"save_price"`
	IsVip              bool   `json:"is_vip"`
	BeginTime          int    `json:"begin_time"`
	EndTime            int    `json:"end_time"`
	EnterpriseEndTime  int    `json:"enterprise_end_time"`
	ExpireTime         int    `json:"expire_time"`
	SurplusTime        int    `json:"surplus_time"`
	IsExpire           bool   `json:"is_expire"`
	CardID             int    `json:"card_id"`
	CardType           int    `json:"card_type"`
	PriceDesc          string `json:"price_desc"`
	IsBuyMonthDiscount bool   `json:"is_buy_month_discount"`
	MonthDiscountPrice int    `json:"month_discount_price"`
	DdURL              string `json:"dd_url"`
	ErrTips            string `json:"err_tips"`
}

// EbookDetail get ebook detail
func (s *Service) EbookDetail(enid string) (detail *EbookDetail, err error) {
	operation := func() (*EbookDetail, error) {
		body, err := s.reqEbookDetail(enid)
		if err != nil {
			return nil, err
		}
		defer body.Close()
		var d *EbookDetail
		if err = handleJSONParse(body, &d); err != nil {
			return nil, err
		}
		return d, nil
	}
	return withRetry(operation, "book-detail")
}

// EbookReadToken get ebook read token
func (s *Service) EbookReadToken(enid string) (t *Token, err error) {
	operation := func() (*Token, error) {
		body, err := s.reqEbookReadToken(enid)
		if err != nil {
			return nil, err
		}
		defer body.Close()
		var token *Token
		if err = handleJSONParse(body, &token); err != nil {
			return nil, err
		}
		return token, nil
	}
	return withRetry(operation, "read-token")
}

// EbookInfo get ebook info
// include book block, book TOC, epubPath etc
func (s *Service) EbookInfo(token string) (info *EbookInfo, err error) {
	operation := func() (*EbookInfo, error) {
		body, err := s.reqEbookInfo(token)
		if err != nil {
			return nil, err
		}
		defer body.Close()
		var i *EbookInfo
		if err = handleJSONParse(body, &i); err != nil {
			return nil, err
		}
		return i, nil
	}
	return withRetry(operation, "book-info")
}

// EbookVIPInfo get ebook vip info
func (s *Service) EbookVIPInfo() (info *EbookVIPInfo, err error) {
	operation := func() (*EbookVIPInfo, error) {
		body, err := s.reqEbookVIPInfo()
		if err != nil {
			return nil, err
		}
		defer body.Close()
		var i *EbookVIPInfo
		if err = handleJSONParse(body, &i); err != nil {
			return nil, err
		}
		return i, nil
	}
	return withRetry(operation, "vip-info")
}

func (s *Service) EbookPages(chapterID, token string, index, count, offset int) (pages *EbookPage, err error) {
	operation := func() (*EbookPage, error) {

		body, err := s.reqEbookPages(chapterID, token, index, count, offset)
		if err != nil {
			return nil, err
		}
		defer body.Close()

		var p *EbookPage
		if err = handleJSONParse(body, &p); err != nil {
			return nil, err
		}
		return p, nil
	}
	return withRetry(operation, chapterID)
}
