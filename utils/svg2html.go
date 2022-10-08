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

func Svg2Html(title string, svgContents []*SvgContent) (err error) {
	result := GenHeadHtml()
	cover := ""
	for k, svgContent := range svgContents {

		for _, content := range svgContent.Contents {
			reader := strings.NewReader(content)

			element, err1 := svgparser.Parse(reader, false)
			if err1 != nil {
				err = err1
				return
			}

			lineContent := make(map[float64][]HtmlEle)
			offset := ""

			for _, children := range element.Children {
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
						yInt, _ := strconv.ParseFloat(y, 64)
						lineContent[yInt] = append(lineContent[yInt], ele)
					}
				}
			}

			keys := make([]float64, 0, len(lineContent))
			for k := range lineContent {
				keys = append(keys, k)
			}
			sort.Float64s(keys)

			// 锚点目录
			// if i == 1 && len(toc) > 0 {
			// 	result += GenTocHtml(toc)
			// }
			// html 强制分页
			result += `
	<p style="page-break-after: always;">`
			for _, v := range keys {
				cont, id := "", ""
				if lineContent[v][0].ID != "" {
					id = lineContent[v][0].ID
				}

				for i, item := range lineContent[v] {
					// TODO： image class=epub-footnote 是注释图片
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
						if w < 18 {
							img += `" style="vertical-align:top;" class="epub-footnote`
						}
						img += `"/>`

						// create cover.html
						if k == 0 {
							cover = GenHeadHtml() + img + `</body></html>`
						}

						// filter cover content
						if k != 0 {
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
						if i == len(lineContent[v])-1 {
							matchH := false
							if strings.Contains(strings.Trim(svgContent.TocText, ""), strings.Trim(cont, "")) {
								matchH = true
							}
							if matchH {
								switch svgContent.TocLevel {
								case 0:
									result += `
	<h1>`
								case 1:
									result += `
	<h2>`
								case 2:
									result += `
	<h3>`
								case 3:
									result += `
	<h4>`
								}
							} else {
								result += `
	<p>`
							}
							result += `<span id="` + id + `" style="` + style + `">` + cont + `</span>`
							if matchH {
								switch svgContent.TocLevel {
								case 0:
									result += `</h1>`
								case 1:
									result += `</h2>`
								case 2:
									result += `</h3>`
								case 3:
									result += `</h4>`

								}
							} else {
								result += `</p>`
							}
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

	err = SaveFile("cover", "html", cover)
	err = Html2PDF(title, "cover.html")
	return
}

func Html2PDF(filename, cover string) (err error) {
	path, err := Mkdir(OutputDir, "Ebook")
	if err != nil {
		return err
	}

	filePreName := filepath.Join(path, FileName(filename, ""))

	fileName, err := FilePath(filePreName, "html", false)
	if err != nil {
		return err
	}

	pdfg, _ := wkhtmltopdf.NewPDFGenerator()

	htmlfile, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("htmlfile err: %#v\n", err)
		return
	}

	fileName, err = FilePath(filePreName, "pdf", false)
	if err != nil {
		return err
	}
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", fileName)

	// see https://wkhtmltopdf.org/usage/wkhtmltopdf.txt
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlfile))
	page.FooterFontSize.Set(10)
	page.FooterRight.Set("[page]")
	page.DisableSmartShrinking.Set(true)

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	coverPath, _ := FilePath(filepath.Join(path, FileName(cover, "")), "", false)
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

	// debug wkhtmltopdf
	// errBuf := new(bytes.Buffer)
	// pdfg.SetStderr(errBuf)
	// done := false
	// defer func() { done = true }()
	// go func() {
	// 	for !done {
	// 		time.Sleep(500 * time.Millisecond)
	// 		fmt.Println(errBuf.String())
	// 	}
	// }()

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
func GenTocHtml(toc []EbookToc) (result string) {
	if len(toc) > 0 {
		result = `
			<p style="page-break-after: always;">
			<p><span style="font-size:28px;font-weight: bold;color:rgb(0, 0, 0);font-family:'PingFang SC';">目 录</span></p>`
		for _, ebookToc := range toc {
			style := "font-size:20px;color:rgb(0, 0, 0);font-family:'PingFang SC';text-decoration: none;"
			if ebookToc.Level == 0 {
				style = "font-size:24px;font-weight: bold;color:rgb(0, 0, 0);font-family:'PingFang SC';text-decoration: none;"
			}
			href := strings.Split(ebookToc.Href, "#")
			text := strings.Repeat("&nbsp;", ebookToc.Level*4) + ebookToc.Text
			result += `
			<p><a href="#` + href[1] + `" style="` + style + `">` + text + `</a></p>`
		}

	}
	return
}
