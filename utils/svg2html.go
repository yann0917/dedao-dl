package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/yann0917/dedao-dl/request"
)

type HtmlEle struct {
	X        string `json:"x"`
	Y        string `json:"y"`
	ID       string `json:"id"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Offset   string `json:"offset"`
	Href     string `json:"href"`
	Name     string `json:"name"`
	Style    string `json:"style"`
	Content  string `json:"content"`
	Class    string `json:"class"`
	Alt      string `json:"alt"`
	Len      string `json:"len"`
	Newline  bool   `json:"newline"`
	IsBold   bool   `json:"is_bold"`
	IsItalic bool   `json:"is_italic"`
	IsFn     bool   `json:"is_fn"`  // footnote: sup tag
	IsSub    bool   `json:"is_sub"` // sub tag
	Fn       struct {
		Href  string `json:"href"`
		Style string `json:"style"`
	} `json:"fn"`
	TextAlign string `json:"text_align"` // left; center; right
}

type EbookToc struct {
	Href      string `json:"href"`
	Level     int    `json:"level"`
	PlayOrder int    `json:"playOrder"`
	Offset    int    `json:"offset"`
	Text      string `json:"text"`
}

type SvgContent struct {
	Contents   []string
	ChapterID  string
	PathInEpub string
	OrderIndex int
}

type SvgContents []*SvgContent

func (a SvgContents) Len() int           { return len(a) }
func (a SvgContents) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SvgContents) Less(i, j int) bool { return a[i].OrderIndex < a[j].OrderIndex } // 从小到大排序

const (
	footNoteImgW     = 20 // 脚注图片≈11x11px & 特殊字图片≈19x19
	footNoteImgH     = 20 // 行内图片高度=20
	svgShapePath     = "path"
	svgShapePolygon  = "polygon"
	svgShapePolyline = "polyline"
	svgShapeLine     = "line"
	svgShapeRect     = "rect"
	svgShapeEllipse  = "ellipse"
	svgShapeCircle   = "circle"

	eBookTypeHtml = "html"
	eBookTypePdf  = "pdf"
	eBookTypeEpub = "epub"

	reqEbookPageWidth = 60000
)

// 假设整本书脚注跳转符号相同
var fnA, fnB = "", ""

var tocLevel map[string]int

func Svg2Html(title string, svgContents []*SvgContent, toc []EbookToc) (err error) {
	tocLevel = make(map[string]int, len(toc))
	for _, ebookToc := range toc {
		tocLevel[ebookToc.Text] = ebookToc.Level
	}
	// fmt.Println(tocLevel)

	result, err := AllInOneHtml(svgContents, toc)
	if err != nil {
		return err
	}
	path, err := Mkdir(OutputDir, "Ebook")
	if err != nil {
		return err
	}

	fileName, err := FilePath(filepath.Join(path, FileName(title, "")), "html", false)
	if err != nil {
		return err
	}
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", fileName)
	if err = WriteFileWithTrunc(fileName, result); err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		return
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
	return
}

func Svg2Pdf(title string, svgContents []*SvgContent, toc []EbookToc) (err error) {

	path, err := Mkdir(OutputDir, "Ebook")
	if err != nil {
		return err
	}
	filePreName := filepath.Join(path, FileName(title, ""))
	fileName, err := FilePath(filePreName, "pdf", false)
	if err != nil {
		return err
	}
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", fileName)
	buf := new(bytes.Buffer)
	cover := ""
	tocLevel = make(map[string]int, len(toc))
	for _, ebookToc := range toc {
		tocLevel[ebookToc.Text] = ebookToc.Level
	}

	for k, svgContent := range svgContents {
		chapter, coverContent, err1 := OneByOneHtml(eBookTypePdf, k, svgContent, toc)
		if err1 != nil {
			err = err1
			return
		}
		if k == 0 {
			cover = coverContent
		}
		buf.Write([]byte(chapter))
		buf.WriteString(`<P style="page-break-before: always">`)
	}

	// write cover into cover.html file
	coverPath, _ := FilePath(filepath.Join(path, FileName("cover", "")), "html", false)
	if err = WriteFileWithTrunc(coverPath, cover); err != nil {
		return
	}

	err = genPdf(buf, fileName, coverPath)
	return
}

func Svg2Epub(title string, svgContents []*SvgContent, opt EpubOptions) (err error) {
	var htmlAll []HtmlContent
	cover := ""
	tocLevel = make(map[string]int, len(opt.Toc))
	chapterToc := make(map[string][]EbookToc, len(opt.Toc))
	for _, ebookToc := range opt.Toc {
		tocLevel[ebookToc.Text] = ebookToc.Level
		tagArr := strings.Split(ebookToc.Href, "#")
		// footnote jump back and forth
		if len(tagArr) > 0 {
			chapterToc[tagArr[0]] = append(chapterToc[tagArr[0]], ebookToc)
		}
	}
	// fmt.Println(chapterToc)

	for k, svgContent := range svgContents {
		chapter, coverUrl, err1 := OneByOneHtml(eBookTypeEpub, k, svgContent, opt.Toc)
		if err1 != nil {
			err = err1
			return
		}
		if k == 0 {
			cover = coverUrl
		}
		htmlAll = append(htmlAll, HtmlContent{
			Content:   chapter,
			ChapterID: svgContent.ChapterID,
			Toc:       chapterToc[svgContent.ChapterID],
		})
	}

	path, err := Mkdir(OutputDir, "Ebook")
	if err != nil {
		return err
	}

	fileName, err := FilePath(filepath.Join(path, FileName(title, "")), "epub", false)

	imageDir, err := Mkdir(OutputDir, "Ebook", "images")
	if err != nil {
		return err
	}
	opt.ImagesDir = imageDir

	h2e := HtmlToEpub{
		EpubOptions: opt,
	}

	if coverByte, err := request.HTTPGet(cover); err == nil {
		h2e.DefaultCover = coverByte
	}

	h2e.HTML = htmlAll
	h2e.Output = fileName
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", fileName)
	if err = h2e.Run(); err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")

	return err
}

func genPdf(buf *bytes.Buffer, fileName, coverPath string) (err error) {
	pdfg, _ := wkhtmltopdf.NewPDFGenerator()

	page := wkhtmltopdf.NewPageReader(buf)
	page.FooterFontSize.Set(10)
	page.FooterRight.Set("[page]")
	page.DisableSmartShrinking.Set(true)

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	pdfg.Cover.EnableLocalFileAccess.Set(true)

	dir, _ := CurrentDir()
	coverPath = filepath.Join(dir, coverPath)

	if runtime.GOOS == "windows" {
		pdfg.Cover.Input = coverPath
	} else {
		pdfg.Cover.Input = "file://" + coverPath
	}

	pdfg.Dpi.Set(300)

	pdfg.TOC.Include = true
	pdfg.TOC.TocHeaderText.Set("目 录")
	pdfg.TOC.HeaderFontSize.Set(18)

	pdfg.TOC.TocLevelIndentation.Set(15)
	pdfg.TOC.TocTextSizeShrink.Set(0.9)
	pdfg.TOC.DisableDottedLines.Set(false)
	pdfg.TOC.EnableTocBackLinks.Set(true)

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.MarginTop.Set(15)
	pdfg.MarginBottom.Set(15)
	pdfg.MarginLeft.Set(15)
	pdfg.MarginRight.Set(15)

	err = pdfg.Create()
	if err != nil {
		fmt.Printf("pdfg create err: %#v\n", err)
		return
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(fileName)
	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		return
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
	err = os.Remove(coverPath)
	return
}

func SaveFile(title, ext, content string) (err error) {
	path, err := Mkdir(OutputDir, "Ebook")
	if err != nil {
		return err
	}

	fileName, err := FilePath(filepath.Join(path, FileName(title, "")), ext, false)
	if err != nil {
		return err
	}
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", fileName)
	if err = WriteFileWithTrunc(fileName, content); err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		return
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
	return
}

// AllInOneHtml generate ebook content all in one html file
func AllInOneHtml(svgContents []*SvgContent, toc []EbookToc) (result string, err error) {
	result = GenHeadHtml()
	fnA, fnB = ParseBookFnDelimiter(svgContents)
	for k, svgContent := range svgContents {
		chapter, _, err1 := OneByOneHtml(eBookTypeHtml, k, svgContent, toc)
		if err1 != nil {
			err = err1
			return
		}
		result += chapter
	}
	result += `
</body>
</html>`
	return
}

// OneByOneHtml one by one generate chapter html
// eType: html/pdf/epub, index: []*SvgContent index, svgContent: one chapter content
func OneByOneHtml(eType string, index int, svgContent *SvgContent, toc []EbookToc) (result, cover string, err error) {
	switch eType {
	case eBookTypeHtml:
		// 锚点目录
		if index == 1 && len(toc) > 0 {
			result += GenTocHtml(toc)
		}
		// html 强制分页
		result += `
	<p style="page-break-after: always;">`

	case eBookTypePdf, eBookTypeEpub:
		result += GenHeadHtml()
	}

	for _, content := range svgContent.Contents {
		result += `
<div id="` + svgContent.ChapterID + `">`
		reader := strings.NewReader(content)

		element, err1 := svgparser.Parse(reader, false)
		if err1 != nil {
			err = err1
			return
		}

		lineContent := GenLineContentByElement(svgContent.ChapterID, element)

		keys := make([]float64, 0, len(lineContent))
		for k := range lineContent {
			keys = append(keys, k)
		}
		sort.Float64s(keys)

		for _, v := range keys {
			cont, id, contWOTag, firstX := "", "", "", 0.0
			if lineContent[v][0].ID != "" {
				id = lineContent[v][0].ID
			}

			lineStyle, currentSpanStyle := "", ""
			hasUncloseSpan := false

			for i, item := range lineContent[v] {
				// image class=epub-footnote 是注释图片
				style := item.Style

				if i == 0 {
					firstX, _ = strconv.ParseFloat(item.X, 64)
					lastIndex := len(lineContent[v]) - 1
					if lineContent[v][lastIndex].Name != "image" {
						lineStyle = lineContent[v][lastIndex].Style
					} else if lastIndex-1 >= 0 {
						lineStyle = lineContent[v][lastIndex-1].Style
					} else {
						lineStyle = item.Style
					}
				}
				centerL := (reqEbookPageWidth / 2) * 0.9
				centerH := (reqEbookPageWidth / 2) * 1.1
				rightL := (reqEbookPageWidth) * 0.9

				w, h := 0.0, 0.0
				w, _ = strconv.ParseFloat(item.Width, 64)
				h, _ = strconv.ParseFloat(item.Height, 64)

				if w > 900 {
					h = 900 * h / w
					w = 900
				}

				switch item.Name {
				case "image":
					img := ""
					if firstX >= centerL && firstX <= centerH {
						style = style + "display: block;text-align:center;"
					} else if firstX >= rightL {
						style = style + "display: block;text-align:right;"
					}
					switch eType {
					case eBookTypeHtml, eBookTypePdf:
						img = `
	<img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
							`" src="` + item.Href +
							`" alt="` + item.Alt +
							`" title="` + item.Alt + `"/>`
						if len(style) > 0 {
							img = `<div style="` + style + `">` + img + `</div>`
						}
						if (w < footNoteImgW || h < footNoteImgH) && len(item.Class) > 0 {
							img = `
	<sup><img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
								`" src="` + item.Href +
								`" alt="` + item.Alt +
								`" title="` + item.Alt +
								`" class="` + item.Class +
								`"/></sup>`
						}
					case eBookTypeEpub:
						img = `
	<img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
							`" src="` + item.Href +
							`" alt="` + item.Alt + `"/>`
						if len(style) > 0 {
							img = `<div style="` + style + `">` + img + `</div>`
						}
						if w < footNoteImgW {
							// epub popup comment
							if len(item.Class) > 0 {
								footnoteId := "footnote-" + strconv.Itoa(index) + "-" + strconv.Itoa(i)
								img = `
	<sup><a epub:type="noteref" href="#` + footnoteId + `"> <img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
									`" src="` + item.Href +
									`" alt="` + item.Alt +
									`" class="` + item.Class + `"/></a></sup>`
								result += `<aside epub:type="footnote" id="` + footnoteId + `">` + item.Alt + `</aside>`
							}
						}
					}

					switch eType {
					case eBookTypePdf:
						// create cover.html
						if index == 0 {
							cover = GenHeadHtml() + img + `</body></html>`
						}
					case eBookTypeEpub:
						// get cover url
						cover = item.Href
					}

					if w < footNoteImgW {
						cont += img
					}

					switch eType {
					case eBookTypeHtml:
						if w >= footNoteImgW {
							result += img
						}
					case eBookTypePdf, eBookTypeEpub:
						// filter cover content
						if index != 0 && w >= footNoteImgW {
							result += img
						}
					}

				case "text":
					if hasUncloseSpan && item.Style != currentSpanStyle {
						cont += `</span>`
						hasUncloseSpan = false
					}

					keepStyle := item.Style != lineStyle
					if keepStyle && !hasUncloseSpan {
						cont += `<span style="` + item.Style + `">`
						currentSpanStyle = item.Style
						hasUncloseSpan = true
					}

					if firstX >= centerL && firstX <= centerH {
						style = style + "display: block;text-align:center;"
					} else if firstX >= rightL {
						style = style + "display: block;text-align:right;"
					}
					if item.Content == "<" {
						item.Content = "&lt;"
					}
					if item.Content == ">" {
						item.Content = "&gt;"
					}
					if item.IsBold {
						cont += `<b>`
					}
					if item.IsItalic {
						cont += `<i>`
					}
					if item.IsFn {
						cont += `<sup>`
					}
					if item.IsSub {
						cont += `<sub>`
					}
					if item.Fn.Href != "" {
						cont += `<a id=` + item.ID + ` href=` + item.Fn.Href
						if item.Fn.Style != "" {
							cont += ` style="` + item.Fn.Style + `"`
						}
						cont += `>`
					}
					cont += item.Content
					if item.Fn.Href != "" {
						cont += `</a>`
					}
					if item.IsFn {
						cont += `</sup>`
					}
					if item.IsSub {
						cont += `</sub>`
					}
					if item.IsItalic {
						cont += `</i>`
					}
					if item.IsBold {
						cont += `</b>`
					}
					contWOTag += item.Content
				}
				if i == len(lineContent[v])-1 {
					matchH := false
					contWOTag = strings.ReplaceAll(contWOTag, "&nbsp;", "")

					level := 0
					for k, v := range tocLevel {
						if strings.Contains(strings.ReplaceAll(k, " ", ""), contWOTag) {
							matchH, level = true, v
							break
						}
					}
					if contWOTag != "" {
						if matchH {
							result += `
</div>`
							result += `<div class='header` + strconv.Itoa(level) + `'>` + GenTocLevelHtml(level, true)
						} else {
							result += `
	<p>`
						}
					}
					if i > 1 && item.Name == "image" {
						style = lineContent[v][i-1].Style
					}
					if cont != "" {
						if id != "" && style != "" {
							result += `<span id="` + id + `" style="` + style + `">`
						} else {
							if id != "" {
								result += `<span id="` + id + `">`
							}
							if style != "" {
								result += `<span style="` + style + `">`
							}
						}
						result += cont + `</span>`
					}
					if contWOTag != "" {
						if matchH {
							result += GenTocLevelHtml(level, false) + `</div>
<div class="part">`
						} else {
							result += `</p>`
						}
					}
				}
			}
		}
		result += `</div>`
		switch eType {
		case eBookTypePdf, eBookTypeEpub:
			result += `
</body>
</html>`
		}
	}
	return
}

func GenHeadHtml() (result string) {
	result = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<style>
		@font-face { font-family: "FZFangSong-Z02"; src:local("FZFangSong-Z02"), url("https://imgcdn.umiwi.com/ttf/fangzhengfangsong_gbk.ttf"); }
		@font-face { font-family: "FZKai-Z03"; src:local("FZFangSong-Z02S"), url("https://imgcdn.umiwi.com/ttf/0315911813008928624065681028886857980055.ttf"); }
		@font-face { font-family: "FZKai-Z03"; src:local("FZKai-Z03"), url("https://imgcdn.umiwi.com/ttf/fangzhengkaiti_gbk.ttf"); }
		@font-face { font-family: "PingFang SC"; src:local("PingFang SC"); }
		@font-face { font-family: "DeDaoJinKai"; src:local("DeDaoJinKai"), url("https://imgcdn.umiwi.com/ttf/dedaojinkaiw03.ttf");}
		@font-face { font-family: "Source Code Pro"; src:local("Source Code Pro"), url("https://imgcdn.umiwi.com/ttf/0315911806889993935644188722660020367983.ttf"); }
		table, tr, td, th, tbody, thead, tfoot {page-break-inside: avoid !important;}
		img { page-break-inside: avoid; max-width: 100% !important;}
		img.epub-footnote { padding-right:5px;}
	</style>
</head>
<body>`
	return
}

// GenTocHtml generate toc html anchor
func GenTocHtml(toc []EbookToc) (result string) {
	if len(toc) == 0 {
		return
	}

	result = `
<div id="toc">
		<p style="page-break-after: always;">
		<p><span style="font-size:24px;font-weight: bold;color:rgb(0, 0, 0);font-family:'PingFang SC';">目 录</span></p>`
	for _, ebookToc := range toc {
		style := "font-size:18px;color:rgb(0, 0, 0);font-family:'PingFang SC';text-decoration: none;"
		if ebookToc.Level == 0 {
			style = "font-size:20px;font-weight: bold;color:rgb(0, 0, 0);font-family:'PingFang SC';text-decoration: none;"
		}
		href := strings.Split(ebookToc.Href, "#")
		text := strings.Repeat("&nbsp;", ebookToc.Level*4) + ebookToc.Text
		if len(href) > 1 {
			result += `
		<p><a href="#` + href[1] + `" style="` + style + `">` + text + `</a></p>`
		} else {
			result += `
		<p><a style="` + style + `">` + text + `</a></p>`
		}
	}
	result += `
</div>`

	return
}

func GenTocLevelHtml(level int, startTag bool) (result string) {
	sTag := map[int]string{0: `<h1>`, 1: `<h2>`, 2: `<h3>`, 3: `<h4>`, 4: `<h5>`, 5: `<h6>`}
	eTag := map[int]string{0: `</h1>`, 1: `</h2>`, 2: `</h3>`, 3: `</h4>`, 4: `</h5>`, 5: `</h6>`}
	if startTag {
		if tag, ok := sTag[level]; ok {
			result = tag
		}
	} else {
		if tag, ok := eTag[level]; ok {
			result = tag
		}
	}
	return
}

func GenLineContentByElement(chapterID string, element *svgparser.Element) (lineContent map[float64][]HtmlEle) {
	lineContent = make(map[float64][]HtmlEle)
	offset := ""
	lastY, lastTop, lastH, lastNewLine := "", "", "", false

	for k, children := range element.Children {
		var ele HtmlEle
		attr := children.Attributes
		content := children.Content

		if _, ok := attr["y"]; ok {
			if children.Name == "text" {
				if content != "" {
					ele.Content = content
				} else {
					if children.Children != nil {
						for _, child := range children.Children {
							if child.Name == "a" {
								ele.Content += child.Content
								attrC := child.Attributes
								if href, ok := attrC["href"]; ok {
									// href="/OEBPS/Text/chapter_00001.xhtml#abc123
									hrefArr := strings.Split(href, "/")
									href = hrefArr[len(hrefArr)-1:][0]
									tagArr := strings.Split(href, "#")
									// footnote jump back and forth
									if len(tagArr) > 1 {
										if strings.Contains(tagArr[1], fnA) {
											ele.Fn.Href = "#" + tagArr[0] + "_" + strings.Replace(tagArr[1], fnA, fnB, -1)
										} else {
											ele.Fn.Href = "#" + tagArr[0] + "_" + strings.Replace(tagArr[1], fnB, fnA, -1)
										}
										attr["id"] = chapterID + "_" + tagArr[1]
									} else {
										ele.Fn.Href = "#" + tagArr[0]
										attr["id"] = chapterID
									}
									ele.Fn.Style = attrC["style"]
								}
							}
						}
					} else {
						ele.Content = "&nbsp;"
					}
				}
				ele.Newline = parseAttrNewline(attr)
				if _, ok := attr["top"]; ok {
					topInt, _ := strconv.ParseFloat(attr["top"], 64)
					heightInt, _ := strconv.ParseFloat(attr["height"], 64)
					lenInt, _ := strconv.ParseFloat(attr["len"], 64)
					lastTopInt, _ := strconv.ParseFloat(lastTop, 64)
					lastHInt, _ := strconv.ParseFloat(lastH, 64)

					// 中文字符 len=3, FIXME: 英文字符无法根据 len 区分是否是下标
					if !lastNewLine && heightInt < lastHInt && heightInt <= 20 && lenInt < 3 {
						if topInt < lastTopInt {
							ele.IsFn = true
						} else {
							ele.IsSub = true
						}

						attr["style"] = ""
					} else {
						lastTop = attr["top"]
						lastH = attr["height"]
					}
				}
			} else {
				ele.Content = ""
			}
			ele.Len = attr["len"]
			if class, ok := attr["class"]; ok {
				ele.Class = class
			}

			if style, ok := attr["style"]; ok {
				style = strings.Replace(style, "fill", "color", -1)
				ele.Style = style
				if strings.Contains(style, "font-weight: bold;") {
					ele.IsBold = true
				}
				if strings.Contains(style, "font-style: oblique") ||
					strings.Contains(style, "font-style: italic") {
					ele.IsItalic = true
				}
			}
			ele.X = attr["x"]

			if ele.IsFn || ele.IsSub {
				ele.Y = lastY
			} else {
				ele.Y = attr["y"]
				if children.Name == "text" {
					lastY = attr["y"]
				}
			}

			ele.Width = attr["width"]
			ele.Height = attr["height"]

			// footnote image with text in one line
			yInt, _ := strconv.ParseFloat(ele.Y, 64)
			w, _ := strconv.ParseFloat(ele.Width, 64)
			if children.Name == "image" && w < footNoteImgW {
				attrPre := element.Children[k-1].Attributes
				yInt, _ = strconv.ParseFloat(attrPre["y"], 64)
				ele.Y = attrPre["y"]
			}
			// id &offset 设置标题 margin-left
			if _, ok := attr["id"]; ok {
				ele.ID = attr["id"]
				if _, ok := attr["offset"]; ok {
					offset = attr["offset"]
				}
			}
			ele.Offset = offset
			ele.Href = parseAttrHref(attr)
			ele.Alt = parseAttrAlt(attr)
			ele.Name = children.Name

			if (children.Name == "text") ||
				children.Name == "image" {
				lineContent[yInt] = append(lineContent[yInt], ele)
			}
			lastNewLine = ele.Newline
		}
	}
	return
}

func parseAttrHref(attr map[string]string) string {
	if href, ok := attr["href"]; ok {
		return href
	}
	return ""
}

func parseAttrAlt(attr map[string]string) string {
	if alt, ok := attr["alt"]; ok {
		return strings.ReplaceAll(alt, "\"", "&quot;")
	}
	return ""
}

func parseAttrNewline(attr map[string]string) bool {
	if newline, ok := attr["newline"]; ok && newline == "true" {
		return true
	}
	return false
}

func ParseBookFnDelimiter(svgContents []*SvgContent) (fnA, fnB string) {
	fn := make(map[string]struct{})
outer:
	for _, svgContent := range svgContents {
		for _, content := range svgContent.Contents {
			reader := strings.NewReader(content)
			element, err := svgparser.Parse(reader, false)
			if err != nil {
				continue
			}
			a, b := parseFootNoteDelimiter(element)
			if a != "" {
				fn[a] = struct{}{}
			}
			if b != "" {
				fn[b] = struct{}{}
			}
			if a != "" && b != "" {
				break outer
			}
		}
	}
	keys := make([]string, 0, len(fn))
	for k := range fn {
		keys = append(keys, k)
	}
	if len(keys) < 1 {
		return
	} else if len(keys) == 1 {
		fnA = keys[0]
	} else if len(keys) == 2 {
		fnA = keys[0]
		fnB = keys[1]
	}
	return
}

func parseFootNoteDelimiter(element *svgparser.Element) (a, b string) {
outer:
	for _, children := range element.Children {
		if children.Name == "text" &&
			children.Content == "" &&
			children.Children != nil {
			for _, child := range children.Children {
				if child.Name == "a" {
					attr := child.Attributes
					if href, ok := attr["href"]; ok {
						// href="/OEBPS/Text/chapter_00001.xhtml#abc123
						hrefArr := strings.Split(href, "/")
						href = hrefArr[len(hrefArr)-1:][0]
						tagArr := strings.Split(href, "#")
						reg := regexp.MustCompile(`([a-zA-Z_-]+)`)
						var params []string
						if len(tagArr) > 1 {
							params = reg.FindStringSubmatch(tagArr[1])
						} else {
							params = reg.FindStringSubmatch(tagArr[0])
						}
						if len(params) > 1 {
							if a == "" {
								a = params[0]
							} else {
								if a != params[0] {
									b = params[0]
									break outer
								}
							}
						}
					}
				}
			}
		}
	}
	return
}
