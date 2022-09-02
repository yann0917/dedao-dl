package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func Svg2Html(title string, contents []string) (err error) {

	result := `<!DOCTYPE html>
	<html>
	<head>
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
		</style>
	</head>
	<body>`

	for _, content := range contents {
		reader := strings.NewReader(content)

		element, err1 := svgparser.Parse(reader, false)
		if err1 != nil {
			fmt.Println(err)
			return err1
		}

		lineContent := make(map[float64][]HtmlEle)
		offset := ""

		for _, children := range element.Children {
			// fmt.Printf("%#v\n", children)
			var ele HtmlEle
			attr := children.Attributes
			content := children.Content

			if y, ok := attr["y"]; ok {
				if children.Name == "text" {
					if content == "" && attr["len"] == "1" {
						ele.Content = " "
					} else {
						ele.Content = content
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

		// html 强制分页
		result += `
		<div style="page-break-after: always;">`
		for _, v := range keys {
			cont := ""
			result += `
		<p>`
			id := ""
			if lineContent[v][0].ID != "" {
				id = lineContent[v][0].ID
			}

			for i, item := range lineContent[v] {
				// TODO： image class=epub-footnote 是注释图片

				style := item.Style
				style = strings.Replace(style, "fill", "color", -1)
				// if id != "" {
				// 	style += " margin-left:" + item.Offset + "px;"
				// }
				w, h := 0.0, 0.0
				w, _ = strconv.ParseFloat(item.Width, 64)
				h, _ = strconv.ParseFloat(item.Height, 64)
				// 1240x1754
				if w > 1240 {
					w = 1240
				}
				if h > 1754 {
					h = 1754
				}
				switch item.Name {
				case "image":
					result += `<img width="` + strconv.FormatFloat(w, 'f', 0, 64) + `" height="` + strconv.FormatFloat(h, 'f', 0, 64) + `" src="` + item.Href + `"/>`
				case "text":
					cont += item.Content
					if i == len(lineContent[v])-1 {
						result += `<span id="` + id + `" style="` + style + `">` + cont + `</span>`
					}
				}
			}
			result += `</p>`
		}
		result += `</div>`
	}
	result += `</body>
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
	err = Html2PDF(title)
	return
}

func Html2PDF(filename string) (err error) {
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

	htmlfile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("htmlfile err: %#v\n", err)
		return
	}

	fileName, err = FilePath(filePreName, "pdf", false)
	if err != nil {
		return err
	}
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", fileName)

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlfile)))
	pdfg.Dpi.Set(300)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.MarginTop.Set(20)
	pdfg.MarginBottom.Set(20)
	pdfg.MarginLeft.Set(20)

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

	fileName, err := FilePath(filepath.Join(path, FileName(title, ext)), ext, false)
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
