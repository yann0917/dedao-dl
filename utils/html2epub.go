package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"github.com/gabriel-vasile/mimetype"
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
}

type HtmlContent struct {
	Content    string
	ChapterID  string
	PathInEpub string
	TocLevel   int
	TocHref    string
	TocText    string
}

type HtmlToEpub struct {
	EpubOptions
	DefaultCover []byte
	book         *epub.Epub
	imgIdx       int
}

func (h *HtmlToEpub) Run() (err error) {
	if len(h.HTML) == 0 {
		return errors.New("no .html file given")
	}
	h.PTitle = make(map[int]string)
	return h.run()
}
func (h *HtmlToEpub) run() (err error) {
	err = h.genBook()
	if err != nil {
		return
	}

	for _, html := range h.HTML {
		err = h.add(html)
		if err != nil {
			err = fmt.Errorf("parse %#v failed: %s", html, err)
			return
		}
	}

	err = h.book.Write(h.Output)
	if err != nil {
		return fmt.Errorf("cannot write output epub: %s", err)
	}

	return
}

func (h *HtmlToEpub) genBook() error {
	h.book = epub.NewEpub(h.Title)
	h.book.SetAuthor(h.Author)
	h.book.SetDescription(h.Description)
	return h.setCover()
}

func (h *HtmlToEpub) setCover() (err error) {
	if h.Cover == "" {
		temp, err := os.CreateTemp("", "html-to-epub")
		if err != nil {
			return fmt.Errorf("can't create tempfile: %s", err)
		}
		_, err = temp.Write(h.DefaultCover)
		if err != nil {
			return fmt.Errorf("can't write tempfile: %s", err)
		}
		_ = temp.Close()

		h.Cover = temp.Name()
	}

	m, err := mimetype.DetectFile(h.Cover)
	if err != nil {
		return fmt.Errorf("can't detect cover mime type %s", err)
	}
	cover, err := h.book.AddImage(h.Cover, "cover"+m.Extension())
	if err != nil {
		return fmt.Errorf("can't add cover %s", err)
	}
	h.book.SetCover(cover, "")

	return
}

func (h *HtmlToEpub) add(html HtmlContent) (err error) {
	refs := make(map[string]string)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html.Content))
	if err != nil {
		return
	}

	images := h.saveImages(doc)
	doc.Find("img").
		Each(func(i int, img *goquery.Selection) {
			h.changeRef(html.Content, img, refs, images)
		})
	content, err := doc.Find("body").Html()
	if err != nil {
		return
	}

	// FIXME: bug
	switch html.TocLevel {
	case 0, 1:
		h.PTitle[html.TocLevel], err = h.book.AddSection(content, html.TocText, html.ChapterID, "")
		if err != nil {
			return
		}
	case 2, 3, 4, 5, 6:
		h.PTitle[html.TocLevel], err = h.book.AddSubSection(h.PTitle[html.TocLevel-1], content, html.TocText, html.ChapterID, "")
		if err != nil {
			return
		}
	}
	return
}

func (h *HtmlToEpub) saveImages(doc *goquery.Document) map[string]string {
	downloads := make(map[string]string)

	tasks := request.NewDownloadTasks()
	doc.Find("img").Each(func(i int, img *goquery.Selection) {
		src, _ := img.Attr("src")
		if !strings.HasPrefix(src, "http") {
			return
		}

		localFile, exist := downloads[src]
		if exist {
			return
		}

		uri, err := url.Parse(src)
		if err != nil {
			log.Printf("parse %s fail: %s", src, err)
			return
		}
		_ = os.MkdirAll(h.ImagesDir, 0766)
		localFile = filepath.Join(h.ImagesDir, fmt.Sprintf("%s%s", MD5str(src), filepath.Ext(uri.Path)))

		tasks.Add(src, localFile)
		downloads[src] = localFile
	})
	request.Batch(tasks, 3, time.Minute*2).ForEach(func(t *request.DownloadTask) {
		if t.Err != nil {
			log.Printf("download %s fail: %s", t.Link, t.Err)
		}
	})

	return downloads
}

// TODO:
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

		localFile, exist := downloads[src]
		if exist {
			return
		}

		uri, err := url.Parse(src)
		if err != nil {
			log.Printf("parse %s fail: %s", src, err)
			return
		}
		_ = os.MkdirAll(h.FontsDir, 0766)
		localFile = filepath.Join(h.FontsDir, fmt.Sprintf("%s%s", MD5str(src), filepath.Ext(uri.Path)))

		downloads[src] = localFile
	})

	return
}

func (h *HtmlToEpub) changeRef(htmlFile string, img *goquery.Selection, refs, downloads map[string]string) {
	img.RemoveAttr("loading")
	img.RemoveAttr("srcset")

	src, _ := img.Attr("src")

	internalRef, exist := refs[src]
	if exist {
		img.SetAttr("src", internalRef)
		return
	}

	var localFile string
	switch {
	case strings.HasPrefix(src, "data:"):
		return
	case strings.HasPrefix(src, "http"):
		localFile, exist = downloads[src]
		if !exist {
			log.Printf("local file of %s not exist", src)
			return
		}
	default:
		fd, err := h.openLocalFile(htmlFile, src)
		if err != nil {
			log.Printf("local ref %s not found: %s", src, err)
			return
		}
		_ = fd.Close()
		localFile = fd.Name()
	}

	// check mime
	fmime, err := mimetype.DetectFile(localFile)
	{
		if err != nil {
			log.Printf("can't detect image mime of %s: %s", src, err)
			return
		}
		if !strings.HasPrefix(fmime.String(), "image") {
			log.Printf("mime of %s is %s instead of images", src, fmime.String())
			return
		}
	}

	// add image
	internalName := fmt.Sprintf("image_%03d", h.imgIdx)
	{
		h.imgIdx += 1
		if !strings.HasSuffix(internalName, fmime.Extension()) {
			internalName += fmime.Extension()
		}
		internalRef, err = h.book.AddImage(localFile, internalName)
		if err != nil {
			log.Printf("can't add image %s: %s", localFile, err)
			return
		}
		refs[src] = internalRef
	}

	if h.Verbose {
		log.Printf("replace %s as %s", src, localFile)
	}

	img.SetAttr("src", internalRef)
}

func (h *HtmlToEpub) openLocalFile(htmlFile string, ref string) (fd *os.File, err error) {
	fd, err = os.Open(ref)
	if err == nil {
		return
	}

	// compatible with evernote's exported htmls
	dirname := strings.TrimSuffix(htmlFile, filepath.Ext(htmlFile))
	name := filepath.Base(ref)
	fd, err = os.Open(filepath.Join(dirname+"_files", name))
	if err == nil {
		return
	}
	fd, err = os.Open(filepath.Join(dirname+".resources", name))
	if err == nil {
		return
	}
	if strings.HasSuffix(ref, ".") {
		return h.openLocalFile(htmlFile, strings.TrimSuffix(ref, "."))
	}

	return
}
