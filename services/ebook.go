package services

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/yann0917/dedao-dl/utils"
)

const (
	maxRetries     = 3
	initialBackoff = 3 * time.Second
	retryDelay     = 1 * time.Second
	// 缓存前缀
	ebookPageCachePrefix = "ebook:page:"
	// 缓存过期时间：24小时
	ebookPageCacheTTL = 24 * time.Hour
	// 请求间隔时间范围（秒）
	minRequestInterval = 1
	maxRequestInterval = 3
	// 触发反爬虫后的冷却时间（秒）
	cooldownTime = 60
	// 最大连续失败次数，超过此数认为触发了反爬虫
	maxConsecutiveFailures = 3
	// 全局请求令牌桶大小
	tokenBucketSize = 5
	// 令牌产生速率（秒/个）
	tokenRefillRate = 0.5
)

// requestLimiter 请求限流器
type requestLimiter struct {
	tokens         int        // 当前可用令牌数
	maxTokens      int        // 最大令牌数
	refillRate     float64    // 令牌填充速率（个/秒）
	lastRefillTime time.Time  // 上次填充时间
	mutex          sync.Mutex // 互斥锁
}

// newRequestLimiter 创建新的请求限流器
func newRequestLimiter(maxTokens int, refillRate float64) *requestLimiter {
	return &requestLimiter{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

// getToken 获取一个请求令牌，如果没有可用令牌则等待
func (r *requestLimiter) getToken() time.Duration {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 计算从上次填充到现在应该填充的令牌数
	now := time.Now()
	elapsedTime := now.Sub(r.lastRefillTime).Seconds()
	newTokens := int(elapsedTime * r.refillRate)

	if newTokens > 0 {
		// 填充令牌，但不超过最大值
		r.tokens = min(r.tokens+newTokens, r.maxTokens)
		r.lastRefillTime = now
	}

	// 如果没有令牌，计算等待时间
	if r.tokens <= 0 {
		// 计算需要等待多久才能获得一个令牌
		waitTime := time.Duration((1.0 / r.refillRate) * float64(time.Second))
		return waitTime
	}

	// 消耗一个令牌
	r.tokens--
	// 添加小的随机抖动，使请求不那么规律
	jitter := time.Duration(rand.Float64() * 200 * float64(time.Millisecond))
	return jitter
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 全局请求控制相关变量
var (
	lastRequestTime     time.Time
	consecutiveFailures int
	antispiderCooldown  bool
	antispiderMutex     sync.Mutex

	// 全局请求限流器
	globalLimiter = newRequestLimiter(tokenBucketSize, tokenRefillRate)
)

// 全局请求控制：检查是否需要等待并设置等待时间
func waitForNextRequest() {
	antispiderMutex.Lock()

	// 如果在冷却期，检查是否已经过了冷却时间
	if antispiderCooldown {
		if time.Since(lastRequestTime) > cooldownTime*time.Second {
			antispiderCooldown = false
			consecutiveFailures = 0
			antispiderMutex.Unlock()
		} else {
			// 仍在冷却期，需要额外等待
			waitTime := cooldownTime*time.Second - time.Since(lastRequestTime)
			antispiderMutex.Unlock()
			fmt.Printf("处于反爬虫冷却期，等待 %.1f 秒...\n", waitTime.Seconds())
			time.Sleep(waitTime)
			return
		}
	} else {
		antispiderMutex.Unlock()
	}

	// 从令牌桶获取令牌，可能会有等待时间
	waitTime := globalLimiter.getToken()
	if waitTime > 0 {
		time.Sleep(waitTime)
	}

	// 更新最后请求时间
	antispiderMutex.Lock()
	lastRequestTime = time.Now()
	antispiderMutex.Unlock()
}

// 记录请求失败
func recordRequestFailure(err error) {
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()

	// 检查错误是否可能是反爬虫引起的
	if strings.Contains(err.Error(), "403") ||
		strings.Contains(err.Error(), "forbidden") ||
		strings.Contains(err.Error(), "too many requests") ||
		strings.Contains(err.Error(), "429") {
		// 直接进入冷却期
		fmt.Println("检测到可能的反爬虫限制，进入冷却期")
		consecutiveFailures = maxConsecutiveFailures
		antispiderCooldown = true
		lastRequestTime = time.Now()
		return
	}

	// 增加连续失败计数
	consecutiveFailures++

	// 如果连续失败次数超过阈值，启动冷却期
	if consecutiveFailures >= maxConsecutiveFailures {
		fmt.Printf("连续请求失败 %d 次，可能触发了反爬虫机制，进入冷却期\n", consecutiveFailures)
		antispiderCooldown = true
		lastRequestTime = time.Now()
	}
}

// 记录请求成功
func recordRequestSuccess() {
	antispiderMutex.Lock()
	defer antispiderMutex.Unlock()

	// 重置连续失败计数
	consecutiveFailures = 0
}

// 使用 BadgerDB 实现缓存相关函数
func getEbookPageCacheKey(enid, chapterID string) string {
	return fmt.Sprintf("%s%s:%s", ebookPageCachePrefix, enid, chapterID)
}

// LoadFromCache 从缓存加载页面数据
func LoadFromCache(enid, chapterID string) ([]string, bool) {
	// 获取 BadgerDB 实例
	db, err := utils.GetBadgerDB(utils.GetDefaultBadgerDBPath())
	if err != nil {
		return nil, false
	}

	var pages []string
	cacheKey := getEbookPageCacheKey(enid, chapterID)

	err = db.Get(cacheKey, &pages)
	if err != nil {
		return nil, false
	}

	return pages, true
}

// SaveToCache 保存页面数据到缓存
func SaveToCache(enid, chapterID string, pages []string) error {
	// 获取 BadgerDB 实例
	db, err := utils.GetBadgerDB(utils.GetDefaultBadgerDBPath())
	if err != nil {
		return err
	}

	cacheKey := getEbookPageCacheKey(enid, chapterID)
	return db.SetWithTTL(cacheKey, pages, ebookPageCacheTTL)
}

// ClearBookCache 清除指定电子书的所有缓存
func ClearBookCache(enid string) error {
	// 获取 BadgerDB 实例
	db, err := utils.GetBadgerDB(utils.GetDefaultBadgerDBPath())
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("%s%s:", ebookPageCachePrefix, enid)
	return db.DeleteWithPrefix(prefix)
}

// ClearAllCache 清除所有电子书缓存
func ClearAllCache() error {
	// 获取 BadgerDB 实例
	db, err := utils.GetBadgerDB(utils.GetDefaultBadgerDBPath())
	if err != nil {
		return err
	}

	return db.DeleteWithPrefix(ebookPageCachePrefix)
}

func withRetry[T any](operation func() (T, error), chapterID string) (result T, err error) {
	backoff := initialBackoff
	var zero T

	for i := 0; i < maxRetries; i++ {
		result, err = operation()
		if err == nil {
			return result, nil
		}

		// 检查错误是否是反爬虫相关
		isAntiSpider := strings.Contains(err.Error(), "反爬虫") ||
			strings.Contains(err.Error(), "403") ||
			strings.Contains(err.Error(), "429") ||
			strings.Contains(err.Error(), "too many requests") ||
			strings.Contains(err.Error(), "forbidden")

		// 打印详细错误信息
		fmt.Printf("\n尝试 %d/%d 失败，章节 %s:\n", i+1, maxRetries, chapterID)
		fmt.Printf("错误: %v\n", err)

		if i < maxRetries-1 {
			// 如果是反爬虫错误，使用更长的退避时间
			if isAntiSpider {
				backoff = backoff * 3 // 更激进的退避
				fmt.Printf("检测到可能的反爬虫限制，使用更长的等待时间\n")
			}

			fmt.Printf("将在 %v 后重试...\n", backoff)
			time.Sleep(backoff)

			// 指数退避策略
			backoff = backoff * 2
		} else {
			fmt.Printf("达到最大重试次数 %d，章节 %s 放弃获取。\n", maxRetries, chapterID)

			// 对于反爬虫错误，建议用户稍后再试
			if isAntiSpider {
				fmt.Printf("建议等待一段时间后再尝试下载此章节，以避免触发反爬虫机制。\n")
			}
		}
	}

	return zero, fmt.Errorf("获取章节 %s 失败，错误: %v", chapterID, err)
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

// EbookPages 使用默认选项的 EbookPages
func (s *Service) EbookPages(chapterID, token string, index, count, offset int) (pages *EbookPage, err error) {
	operation := func() (*EbookPage, error) {
		// 在请求之前检查并等待合适的时间间隔
		// 使用全局令牌桶限流器来平衡并发请求速率
		waitForNextRequest()

		// 请求API获取数据
		body, err := s.reqEbookPages(chapterID, token, index, count, offset)
		if err != nil {
			// 记录请求失败
			recordRequestFailure(err)
			return nil, err
		}
		defer body.Close()

		var p *EbookPage
		if err = handleJSONParse(body, &p); err != nil {
			// 记录请求失败
			recordRequestFailure(err)
			return nil, err
		}

		// 请求成功，记录成功状态
		recordRequestSuccess()

		// 检查返回的页面数据是否为空
		if p == nil || len(p.Pages) == 0 {
			// 如果返回的页面为空，可能是触发了反爬虫
			recordRequestFailure(fmt.Errorf("返回的页面数据为空，可能触发了反爬虫"))
		}

		return p, nil
	}

	return withRetry(operation, chapterID)
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
