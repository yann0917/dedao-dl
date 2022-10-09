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
)

type HtmlEle struct {
	X       string `json:"x"`
	Y       string `json:"y"`
	ID      string `json:"id"`
	Width   string `json:"width"`
	Height  string `json:"height"`
	Offset  string `json:"offset"`
	Href    string `json:"href"`
	Name    string `json:"name"`
	Style   string `json:"style"`
	Content string `json:"content"`
	Class   string `json:"class"`
	Alt     string `json:"alt"`
	Len     string `json:"len"`
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

const footNoteImgW = 16

func Svg2Html(title string, svgContents []*SvgContent, toc []*EbookToc) (err error) {
	result := GenHeadHtml()
	for k, svgContent := range svgContents {

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

			// 锚点目录
			if k == 1 && len(toc) > 0 {
				result += GenTocHtml(toc)
			}
			// html 强制分页
			result += `
	<p style="page-break-after: always;">`
			for _, v := range keys {
				cont, id := "", ""
				if lineContent[v][0].ID != "" {
					id = lineContent[v][0].ID
				}

				for i, item := range lineContent[v] {
					// image class=epub-footnote 是注释图片
					style := item.Style
					style = strings.Replace(style, "fill", "color", -1)

					w, h := 0.0, 0.0
					w, _ = strconv.ParseFloat(item.Width, 64)
					h, _ = strconv.ParseFloat(item.Height, 64)

					if w > 900 {
						h = 900 * h / w
						w = 900
					}

					switch item.Name {
					case "image":
						img := `
	<img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
							// `" height="` + strconv.FormatFloat(h, 'f', 0, 64) +
							`" src="` + item.Href +
							`" alt="` + item.Alt +
							`" title="` + item.Alt
						if w < footNoteImgW {
							img += `" style="vertical-align:top; display:inline-block;" class="epub-footnote`
						}
						img += `"/>`
						if w < footNoteImgW {
							cont += img
						} else {
							result += img
						}

					case "text":
						if item.Content == "<" {
							item.Content = "&lt;"
						}
						if item.Content == ">" {
							item.Content = "&gt;"
						}
						cont += item.Content
					}

					if i == len(lineContent[v])-1 {
						matchH := false
						if strings.Contains(strings.Trim(svgContent.TocText, ""), strings.Trim(cont, "")) {
							matchH = true
						}
						if matchH {
							result += GenTocLevelHtml(svgContent.TocLevel, true)
						} else {
							result += `
	<p>`
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

	}
	result += `
</body>
</html>`

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

	// err = SaveFile("cover", "html", cover)
	// err = Html2PDF(title, "cover.html")
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
			result := GenHeadHtml()
			// html 强制分页
			for _, v := range keys {
				cont, id := "", ""
				if lineContent[v][0].ID != "" {
					id = lineContent[v][0].ID
				}

				for i, item := range lineContent[v] {
					// image class=epub-footnote 是注释图片
					style := item.Style
					style = strings.Replace(style, "fill", "color", -1)

					w, h := 0.0, 0.0
					w, _ = strconv.ParseFloat(item.Width, 64)
					h, _ = strconv.ParseFloat(item.Height, 64)

					if w > 900 {
						h = 900 * h / w
						w = 900
					}

					switch item.Name {
					case "image":
						img := `
	<img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
							`" src="` + item.Href +
							`" alt="` + item.Alt +
							`" title="` + item.Alt
						if w < footNoteImgW {
							img += `" style="vertical-align:top; display:inline-block;" class="epub-footnote`
						}
						img += `"/>`
						if w < footNoteImgW {
							cont += img
						}
						// create cover.html
						if k == 0 {
							cover = GenHeadHtml() + img + `</body></html>`
						}

						// filter cover content
						if k != 0 && w >= footNoteImgW {
							result += img
						}

					case "text":
						if item.Content == "<" {
							item.Content = "&lt;"
						}
						if item.Content == ">" {
							item.Content = "&gt;"
						}
						cont += item.Content
					}
					if i == len(lineContent[v])-1 {
						matchH := false
						if strings.Contains(strings.Trim(svgContent.TocText, ""), strings.Trim(cont, "")) {
							matchH = true
						}
						if matchH {
							result += GenTocLevelHtml(svgContent.TocLevel, true)
						} else {
							result += `
	<p>`
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
			result += `
</body>
</html>`
			buf.Write([]byte(result))
			buf.WriteString(`<P style="page-break-before: always">`)
		}
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
	pdfg.Cover.Input = "file://" + dir + "/" + coverPath

	pdfg.Dpi.Set(300)

	pdfg.TOC.Include = true
	pdfg.TOC.TocHeaderText.Set("目 录")
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

func GenHeadHtml() (result string) {
	result = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<style>
		@font-face { font-family: "FZFangSong-Z02";
		src:local("FZFangSong-Z02"),
		url("https://imgcdn.umiwi.com/ttf/fangzhengfangsong_gbk.ttf"); }
		@font-face { font-family: "FZKai-Z03";
		src:local("FZKai-Z03"),
		url("https://imgcdn.umiwi.com/ttf/fangzhengkaiti_gbk.ttf"); }
		@font-face { font-family: "PingFang SC";
		src:local("PingFang SC"); }
		@font-face { font-family: "Source Code Pro";
		src:local("Source Code Pro"),
		url("https://imgcdn.umiwi.com/ttf/0315911806889993935644188722660020367983.ttf"); }
		table, tr, td, th, tbody, thead, tfoot {page-break-inside: avoid !important;}
		img { page-break-inside: avoid; max-width: 100% !important;}
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
	switch level {
	case 0:
		if startTag {
			result = `<h1>`
		} else {
			result = `</h1>`
		}

	case 1:
		if startTag {
			result = `<h2>`
		} else {
			result = `</h2>`
		}

	case 2:
		if startTag {
			result = `<h3>`
		} else {
			result = `</h3>`
		}

	case 3:
		if startTag {
			result = `<h4>`
		} else {
			result = `</h4>`
		}
	default:
		if startTag {
			result = `<p>`
		} else {
			result = `</p>`
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

			ele.Style = attr["style"]
			ele.X = attr["x"]
			ele.Y = attr["y"]
			ele.Width = attr["width"]
			ele.Height = attr["height"]
			yInt, _ := strconv.ParseFloat(y, 64)

			// footnote image with text in one line
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
