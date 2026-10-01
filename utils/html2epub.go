package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yann0917/dedao-dl/request"
)

type EpubOptions struct {
	Cover       string
	Title       string
	Author      string
	Description string
	Output      string
	ImagesDir   string
	FontsDir    string
	HTML        []HtmlContent
	Verbose     bool
	PTitle      map[int]string
	Toc         []EbookToc
}

type HtmlContent struct {
	Content   string
	ChapterID string
	Toc       []EbookToc
	TocLevel  int
	TocText   string
}

// HtmlToEpub EPUB 打包器。打包逻辑在 epub.go（手写 zip 组装），
// 这里只负责把各章节 HTML 清洗、图片本地化后交给 epubDoc。
type HtmlToEpub struct {
	EpubOptions
	DefaultCover []byte
	imgIdx       int
	downloads    map[string]string // src -> 本地文件
	refs         map[string]string // src -> epub 内部 href
	doc          *epubDoc
}

const kindleSafeCSS = `
body, div, p, li, span, a, blockquote, h1, h2, h3, h4, h5, h6 {
	font-family: STHeiti, STYuan, "Amazon Ember", Helvetica, Arial, sans-serif !important;
}
pre, code, kbd, samp {
	font-family: monospace !important;
}
p {
	margin: 1.5em 0;
}
`

var (
	cssFontFaceBlockRE = regexp.MustCompile(`(?is)@font-face\s*\{.*?\}`)
	cssFontFamilyRE    = regexp.MustCompile(`(?i)font-family\s*:[^;}{]+;?`)
	inlineFontFamilyRE = regexp.MustCompile(`(?i)font-family\s*:[^;]+;?`)
)

func (h *HtmlToEpub) Run() (err error) {
	if len(h.HTML) == 0 {
		return errors.New("no .html file given")
	}
	h.downloads = make(map[string]string)
	h.refs = make(map[string]string)
	h.doc = newEpubDoc(h.Title, h.Author, h.Description, kindleSafeCSS)

	if err = h.setCover(); err != nil {
		return
	}

	for _, html := range h.HTML {
		if err = h.add(html); err != nil {
			err = fmt.Errorf("parse %#v failed: %s", html, err)
			return
		}
	}

	h.doc.buildNavTree(h.Toc)
	if err = h.doc.write(h.Output); err != nil {
		return fmt.Errorf("cannot write output epub: %s", err)
	}
	return
}

func (h *HtmlToEpub) setCover() error {
	data := h.DefaultCover
	fallbackExt := filepath.Ext(h.Cover)
	if h.Cover != "" {
		b, err := os.ReadFile(h.Cover)
		if err != nil {
			return fmt.Errorf("can't read cover: %s", err)
		}
		data = b
	}
	return h.doc.setCover(data, fallbackExt)
}

func (h *HtmlToEpub) add(html HtmlContent) (err error) {
	// 封面页不进 spine，由 cover-image 承担（与旧 go-epub 流程行为一致）
	if html.ChapterID == "cover.xhtml" {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html.Content))
	if err != nil {
		return
	}
	h.sanitizeFontStyles(doc)

	h.saveImages(doc)
	doc.Find("img").
		Each(func(i int, img *goquery.Selection) {
			h.changeRef(img)
		})

	content, err := doc.Find("body").Html()
	if err != nil {
		return
	}
	h.doc.addChapter(html.ChapterID, epubXHTMLDoc(h.Title, content))
	return
}

func (h *HtmlToEpub) sanitizeFontStyles(doc *goquery.Document) {
	doc.Find("style").Each(func(i int, style *goquery.Selection) {
		text := style.Text()
		text = cssFontFaceBlockRE.ReplaceAllString(text, "")
		text = cssFontFamilyRE.ReplaceAllString(text, "")
		text = strings.TrimSpace(text)
		if text == "" {
			style.Remove()
			return
		}
		style.SetText(text)
	})

	doc.Find("[style]").Each(func(i int, s *goquery.Selection) {
		inline, ok := s.Attr("style")
		if !ok || inline == "" {
			return
		}
		inline = inlineFontFamilyRE.ReplaceAllString(inline, "")
		inline = strings.TrimSpace(inline)
		inline = strings.Trim(inline, ";")
		if inline == "" {
			s.RemoveAttr("style")
			return
		}
		s.SetAttr("style", inline)
	})
}

func (h *HtmlToEpub) saveImages(doc *goquery.Document) {
	tasks := request.NewDownloadTasks()
	doc.Find("img").Each(func(i int, img *goquery.Selection) {
		src, _ := img.Attr("src")
		if !strings.HasPrefix(src, "http") {
			return
		}

		if _, exist := h.downloads[src]; exist {
			return
		}

		localFile, err := localImageName(src, h.ImagesDir)
		if err != nil {
			log.Printf("parse %s fail: %s", src, err)
			return
		}

		tasks.Add(src, localFile)
		h.downloads[src] = localFile
	})
	request.Batch(tasks, 3, time.Minute*2).ForEach(func(t *request.DownloadTask) {
		if t.Err != nil {
			log.Printf("download %s fail: %s", t.Link, t.Err)
		}
	})
}

// localImageName 生成图片在本地的缓存文件名（内容寻址：src 的 MD5 + 原扩展名）。
func localImageName(src, dir string) (string, error) {
	uri, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("%s%s", MD5str(src), filepath.Ext(uri.Path))), nil
}

func (h *HtmlToEpub) changeRef(img *goquery.Selection) {
	img.RemoveAttr("loading")
	img.RemoveAttr("srcset")

	src, _ := img.Attr("src")
	if src == "" || strings.HasPrefix(src, "data:") {
		return
	}

	if ref, exist := h.refs[src]; exist {
		img.SetAttr("src", ref)
		return
	}

	var localFile string
	switch {
	case strings.HasPrefix(src, "http"):
		var exist bool
		localFile, exist = h.downloads[src]
		if !exist {
			log.Printf("local file of %s not exist", src)
			return
		}
	default:
		localFile = src
	}

	data, err := os.ReadFile(localFile)
	if err != nil {
		log.Printf("can't read image %s: %s", localFile, err)
		return
	}
	if mt := http.DetectContentType(data); !strings.HasPrefix(mt, "image/") {
		log.Printf("mime of %s is %s instead of images", src, mt)
		return
	}

	ext := strings.ToLower(filepath.Ext(localFile))
	if ext == "" {
		ext, _ = detectImageExtAndType(data, ".png")
	}

	internalRef := h.doc.addImage(data, ext)
	h.refs[src] = internalRef

	if h.Verbose {
		log.Printf("replace %s as %s", src, internalRef)
	}
	img.SetAttr("src", internalRef)
}

// getFontURLs TODO:
func (h *HtmlToEpub) getFontURLs(html HtmlContent) (downloads map[string]string, err error) {
	downloads = make(map[string]string)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html.Content))
	if err != nil {
		return
	}

	doc.Find("head>style").Each(func(i int, font *goquery.Selection) {
		fmt.Printf("%#v\n", font.Text())
		val, ok := font.Attr("font-family")
		fmt.Printf("%#v, %#v\n", val, ok)
		src, _ := font.Attr("url")
		if !strings.HasPrefix(src, "http") {
			return
		}

		if _, exist := downloads[src]; exist {
			return
		}

		localFile, err := localImageName(src, h.FontsDir)
		if err != nil {
			log.Printf("parse %s fail: %s", src, err)
			return
		}
		downloads[src] = localFile
	})

	return
}
