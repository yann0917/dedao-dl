package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
	IsBold   bool   `json:"is_bold"`
	IsItalic bool   `json:"is_italic"`
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
	TocLevel   int
	TocHref    string
	TocText    string
}

const (
	footNoteImgW     = 20 // 脚注图片≈11x11px & 特殊字图片≈19x19
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
)

func Svg2Html(title string, svgContents []*SvgContent, toc []*EbookToc) (err error) {
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

func Svg2Pdf(title string, svgContents []*SvgContent) (err error) {

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
	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	cover := ""
	for k, svgContent := range svgContents {
		chapter, coverContent, err1 := OneByOneHtml(eBookTypePdf, k, svgContent, []*EbookToc{})
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

	page := wkhtmltopdf.NewPageReader(buf)
	page.FooterFontSize.Set(10)
	page.FooterRight.Set("[page]")
	page.DisableSmartShrinking.Set(true)

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	pdfg.Cover.EnableLocalFileAccess.Set(true)
	dir, _ := CurrentDir()

	pdfg.Cover.Input = "file://" + filepath.Join(dir, coverPath)

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
	err = os.Remove(dir + "/" + coverPath)
	return
}

func Svg2Epub(title string, svgContents []*SvgContent, opt EpubOptions) (err error) {
	var htmlAll []HtmlContent
	cover := ""
	for k, svgContent := range svgContents {
		chapter, coverUrl, err1 := OneByOneHtml(eBookTypeEpub, k, svgContent, []*EbookToc{})
		if err1 != nil {
			err = err1
			return
		}
		if k == 0 {
			cover = coverUrl
		}
		htmlAll = append(htmlAll, HtmlContent{
			Content:    chapter,
			ChapterID:  svgContent.ChapterID,
			PathInEpub: svgContent.PathInEpub,
			TocLevel:   svgContent.TocLevel,
			TocHref:    svgContent.TocHref,
			TocText:    svgContent.TocText,
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
func AllInOneHtml(svgContents []*SvgContent, toc []*EbookToc) (result string, err error) {
	result = GenHeadHtml()
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
func OneByOneHtml(eType string, index int, svgContent *SvgContent, toc []*EbookToc) (result, cover string, err error) {
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
		reader := strings.NewReader(content)

		element, err1 := svgparser.Parse(reader, false)
		if err1 != nil {
			err = err1
			return
		}

		lineContent := GenLineContentByElement(element)

		keys := make([]float64, 0, len(lineContent))
		for k := range lineContent {
			keys = append(keys, k)
		}
		sort.Float64s(keys)

		for _, v := range keys {
			cont, id, contWOTag := "", "", ""
			if lineContent[v][0].ID != "" {
				id = lineContent[v][0].ID
			}

			for i, item := range lineContent[v] {
				// image class=epub-footnote 是注释图片
				style := item.Style

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
					switch eType {
					case eBookTypeHtml, eBookTypePdf:
						img = `
	<img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
							`" src="` + item.Href +
							`" alt="` + item.Alt +
							`" title="` + item.Alt
						if w < footNoteImgW {
							img += `" style="vertical-align:top;" class="` + item.Class
						}
						img += `"/>`
					case eBookTypeEpub:
						img = `
	<img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
							`" src="` + item.Href +
							`" alt="` + item.Alt + `"/>`
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
					cont += item.Content
					if item.IsItalic {
						cont += `</i>`
					}
					if item.IsBold {
						cont += `</b>`
					}
					contWOTag += item.Content
					if i == len(lineContent[v])-1 {
						matchH := false
						tocText := strings.ReplaceAll(svgContent.TocText, " ", "")
						contWOTag = strings.ReplaceAll(contWOTag, "&nbsp;", "")
						if len(svgContent.TocText) > 0 &&
							len(contWOTag) > 1 &&
							strings.Contains(tocText, contWOTag) {
							matchH = true
						}
						if matchH {
							result += GenTocLevelHtml(svgContent.TocLevel, true)
						} else {
							result += `
	<p>`
						}
						if i > 1 && style == "" {
							style = lineContent[v][i-1].Style
						}
						result += `<span id="` + id + `" style="` + style + `">` + cont + `</span>`
						if matchH {
							result += GenTocLevelHtml(svgContent.TocLevel, false)
						} else {
							result += `</p>`
						}
					}
				}
			}
		}
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
		@font-face { font-family: "FZKai-Z03"; src:local("FZKai-Z03"), url("https://imgcdn.umiwi.com/ttf/fangzhengkaiti_gbk.ttf"); }
		@font-face { font-family: "PingFang SC"; src:local("PingFang SC"); }
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
func GenTocHtml(toc []*EbookToc) (result string) {
	if len(toc) == 0 {
		return
	}

	result = `
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

func GenLineContentByElement(element *svgparser.Element) (lineContent map[float64][]HtmlEle) {
	lineContent = make(map[float64][]HtmlEle)
	offset := ""
	for k, children := range element.Children {
		var ele HtmlEle
		attr := children.Attributes
		content := children.Content

		if y, ok := attr["y"]; ok {
			if children.Name == "text" {
				if content != "" {
					ele.Content = content
				} else {
					if children.Children != nil {
						for _, child := range children.Children {
							if child.Name == "a" {
								ele.Content += child.Content
							}
						}
					} else {
						ele.Content = "&nbsp;"
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
			ele.Y = attr["y"]
			ele.Width = attr["width"]
			ele.Height = attr["height"]

			// footnote image with text in one line
			yInt, _ := strconv.ParseFloat(y, 64)
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

			if _, ok := attr["href"]; ok && children.Name == "image" {
				ele.Href = attr["href"]
			} else {
				ele.Href = ""
			}
			if _, ok := attr["alt"]; ok && children.Name == "image" {
				ele.Alt = strings.ReplaceAll(attr["alt"], "\"", "&quot;")
			} else {
				ele.Alt = ""
			}
			ele.Name = children.Name

			if (children.Name == "text") ||
				children.Name == "image" {
				lineContent[yInt] = append(lineContent[yInt], ele)
			}
		}
	}
	return
}
