package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// 在线阅读按批加载：每次 /reader/chunk 只发一次 get_pages 请求（readerPageSize 页），
// 前端把最多 readerBatchRequests 次请求作为一批，读到接近已加载末尾时再来一批，
// 避免一次性下载整书触发风控。实际请求节奏由 services 包的全局限流器兜底。
//
// 与下载路径不同，阅读器把浏览器可视区宽高与阅读器设置（字号/行距/字体）传给
// get_pages，服务端按该版式分页，返回的页面即整屏内容，前端直接渲染 SVG
// （下载路径固定超大画布是为了整章拉取后重排）。
const (
	readerBatchRequests = 5
	readerPageSize      = 20
	readerMinDim        = 200
	readerMaxDim        = 8000
	readerMinFontSize   = 12
	readerMaxFontSize   = 28
)

// 阅读器可选字体：白名单防注入，来源是页面 SVG @font-face 里出现过的字体族
var readerFontFamilies = map[string]bool{
	"pingfang":       true,
	"FZKai-Z03":      true,
	"FZFangSong-Z02": true,
	"DeDaoJinKai":    true,
}

type readerChapter struct {
	ChapterID  string `json:"chapterId"`
	Title      string `json:"title"`
	OrderIndex int    `json:"orderIndex"`
}

type readerInfoData struct {
	Chapters []readerChapter `json:"chapters"`
}

type readerPage struct {
	SVG    string  `json:"svg"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type readerChunkData struct {
	ChapterID string                   `json:"chapterId"`
	Index     int                      `json:"index"`
	Layout    services.EbookPageLayout `json:"layout"`
	IsEnd     bool                     `json:"isEnd"`
	NextIndex int                      `json:"nextIndex"`
	Pages     []readerPage             `json:"pages"`
	Cached    bool                     `json:"cached"`
}

// read token 与登录态绑定且整本下载期间复用，这里按 enid 做进程内缓存，
// 避免每批都多打一次 read/token 接口
var (
	readerTokenMutex sync.Mutex
	readerTokens     = map[string]readerTokenEntry{}
)

type readerTokenEntry struct {
	token    string
	expireAt time.Time
}

func registerEbookReaderRoutes(group *gin.RouterGroup) {
	ebook := group.Group("/ebook")
	ebook.GET("/reader/info", getEbookReaderInfo)
	ebook.GET("/reader/chunk", getEbookReaderChunk)
}

func getReaderToken(service *services.Service, enid string) (string, error) {
	readerTokenMutex.Lock()
	if entry, found := readerTokens[enid]; found && time.Now().Before(entry.expireAt) {
		readerTokenMutex.Unlock()
		return entry.token, nil
	}
	readerTokenMutex.Unlock()

	t, err := service.EbookReadToken(enid)
	if err != nil {
		return "", err
	}

	readerTokenMutex.Lock()
	readerTokens[enid] = readerTokenEntry{token: t.Token, expireAt: time.Now().Add(time.Hour)}
	readerTokenMutex.Unlock()
	return t.Token, nil
}

func invalidateReaderToken(enid string) {
	readerTokenMutex.Lock()
	delete(readerTokens, enid)
	readerTokenMutex.Unlock()
}

// decryptEbookPage 解密 get_pages 返回的单页 SVG 内容（AES-128-CBC + PKCS7），
// 密钥与 cmd/app.DecryptAES 相同
func decryptEbookPage(contents string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(contents)
	if err != nil || len(ciphertext) == 0 || len(ciphertext)%16 != 0 {
		return ""
	}
	key := []byte("3e4r06tjkpjcevlbslr3d96gdb5ahbmo")
	iv := []byte("6fd89a1b3a7f48fb")
	plaintext, err := utils.AES128Decrypt(ciphertext, key, iv)
	if err != nil {
		return ""
	}
	return string(plaintext)
}

var (
	// 页面 SVG 根标签的 width/height 只是请求布局的回显，内容实际占位可能小于画布
	// （最后一页尤甚），也可能个别元素越出画布。只统计 text/image 的几何属性
	// （装饰性 rect 会横跨整个画布，计入会放大包围盒），前端据此设置 viewBox，
	// 保证按请求宽高渲染时既不裁内容也不留大片空白
	svgElementRe  = regexp.MustCompile(`<(text|image)\b[^>]*>`)
	svgGeomAttrRe = regexp.MustCompile(`\b(x|y|width|height)="(-?[0-9.]+)"`)
)

// svgPageBounds 计算单页 SVG 内容包围盒（宽高），带少量右边距/下边距
func svgPageBounds(svg string) (float64, float64) {
	var maxX, maxY float64
	for _, el := range svgElementRe.FindAllString(svg, -1) {
		var x, y, w, h float64
		for _, attr := range svgGeomAttrRe.FindAllStringSubmatch(el, -1) {
			v, err := strconv.ParseFloat(attr[2], 64)
			if err != nil {
				continue
			}
			switch attr[1] {
			case "x":
				x = v
			case "y":
				y = v
			case "width":
				w = v
			case "height":
				h = v
			}
		}
		if x+w > maxX {
			maxX = x + w
		}
		if y+h > maxY {
			maxY = y + h
		}
	}
	if maxX <= 0 || maxY <= 0 {
		// 空页/异常页给一个最小的可读比例
		return 600, 120
	}
	return maxX + 20, maxY + 30
}

func buildReaderPages(contents []string) []readerPage {
	pages := make([]readerPage, 0, len(contents))
	for _, content := range contents {
		if content == "" {
			continue
		}
		width, height := svgPageBounds(content)
		pages = append(pages, readerPage{SVG: content, Width: width, Height: height})
	}
	return pages
}

func getEbookReaderInfo(c *gin.Context) {
	enid := c.Query("enid")
	if enid == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 参数")
		return
	}

	service := config.Instance.ActiveUserService()
	token, err := getReaderToken(service, enid)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	info, err := service.EbookInfo(token)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	// Toc 里同一章可能有多个小节（href 带 #frag），与下载进度一致取首个标题
	titleByChapter := make(map[string]string)
	for _, ebookToc := range info.BookInfo.Toc {
		href := ebookToc.Href
		if i := strings.Index(href, "#"); i >= 0 {
			href = href[:i]
		}
		if href == "" {
			continue
		}
		if _, found := titleByChapter[href]; !found {
			titleByChapter[href] = ebookToc.Text
		}
	}

	chapters := make([]readerChapter, 0, len(info.BookInfo.Orders))
	for i, order := range info.BookInfo.Orders {
		chapters = append(chapters, readerChapter{
			ChapterID:  order.ChapterID,
			Title:      titleByChapter[order.ChapterID],
			OrderIndex: i,
		})
	}

	ok(c, readerInfoData{Chapters: chapters})
}

func readQueryDim(c *gin.Context, key string, fallback int) int {
	v := readQueryInt(c, key, fallback)
	if v < readerMinDim {
		return fallback
	}
	if v > readerMaxDim {
		return readerMaxDim
	}
	return v
}

func readQueryFloat(c *gin.Context, key string, fallback float64) float64 {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

// readerLayoutFromQuery 从请求参数解析阅读器版式：
// w/h 浏览器可视区，fs 字号，fl 行距（em），ff 字体
func readerLayoutFromQuery(c *gin.Context) (services.EbookPageLayout, bool) {
	width := readQueryDim(c, "w", 0)
	height := readQueryDim(c, "h", 0)
	if width <= 0 || height <= 0 {
		return services.EbookPageLayout{}, false
	}

	fontSize := readQueryInt(c, "fs", 16)
	if fontSize < readerMinFontSize {
		fontSize = readerMinFontSize
	}
	if fontSize > readerMaxFontSize {
		fontSize = readerMaxFontSize
	}

	lineHeight := readQueryFloat(c, "fl", 2)
	if lineHeight < 1.2 {
		lineHeight = 1.2
	}
	if lineHeight > 3 {
		lineHeight = 3
	}

	fontFamily := c.Query("ff")
	if !readerFontFamilies[fontFamily] {
		fontFamily = "pingfang"
	}

	return services.EbookPageLayout{
		Width:      width,
		Height:     height,
		FontSize:   fontSize,
		LineHeight: lineHeight,
		FontFamily: fontFamily,
	}, true
}

func getEbookReaderChunk(c *gin.Context) {
	enid := c.Query("enid")
	chapterID := c.Query("chapterId")
	if enid == "" || chapterID == "" {
		fail(c, http.StatusBadRequest, "缺少 enid 或 chapterId 参数")
		return
	}
	index := readQueryInt(c, "index", 0)
	if index < 0 {
		index = 0
	}
	// 浏览器可视区与排版设置决定分页：任一参数变化分页都不同
	layout, okLayout := readerLayoutFromQuery(c)
	if !okLayout {
		fail(c, http.StatusBadRequest, "缺少 w 或 h 参数")
		return
	}

	// 分片缓存按版式+页序分键：回看/重进同一页零请求
	if chunk, found := services.LoadReaderChunk(enid, chapterID, index, readerPageSize, layout); found {
		okGzip(c, readerChunkData{
			ChapterID: chapterID,
			Index:     index,
			Layout:    layout,
			IsEnd:     chunk.IsEnd,
			NextIndex: index + readerPageSize,
			Pages:     buildReaderPages(chunk.Pages),
			Cached:    true,
		})
		return
	}

	service := config.Instance.ActiveUserService()
	token, err := getReaderToken(service, enid)
	if err != nil {
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	pageList, err := service.EbookPagesSized(chapterID, token, index, readerPageSize, 0, layout)
	if err != nil {
		// 账号/权限类业务错误多半是 token 失效，清掉缓存下次重取
		var be *services.BusinessError
		if errors.As(err, &be) {
			invalidateReaderToken(enid)
		}
		fail(c, http.StatusBadGateway, err.Error())
		return
	}

	var decrypted []string
	isEnd := false
	if pageList != nil {
		isEnd = pageList.IsEnd
		for _, item := range pageList.Pages {
			decrypted = append(decrypted, decryptEbookPage(item.Svg))
		}
	}

	if err := services.SaveReaderChunk(enid, chapterID, index, readerPageSize, layout, &services.ReaderChunk{
		Pages: decrypted,
		IsEnd: isEnd,
	}); err != nil {
		println("警告: 无法缓存章节分片 " + chapterID + ": " + err.Error())
	}

	okGzip(c, readerChunkData{
		ChapterID: chapterID,
		Index:     index,
		Layout:    layout,
		IsEnd:     isEnd,
		NextIndex: index + readerPageSize,
		Pages:     buildReaderPages(decrypted),
	})
}

// okGzip 带 gzip 压缩的成功响应。原始页面 SVG 文本量不小，压缩后可缩一个数量级
func okGzip(c *gin.Context, data interface{}) {
	payload, err := json.Marshal(envelope{Code: 0, Msg: "", Data: data})
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(payload); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := zw.Close(); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Encoding", "gzip")
	c.Data(http.StatusOK, "application/json; charset=utf-8", buf.Bytes())
}
