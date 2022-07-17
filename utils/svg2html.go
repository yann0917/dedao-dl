package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type HtmlEle struct {
	X       string `json:"x"`
	Y       string `json:"y"`
	Width   string `json:"width"`
	Height  string `json:"height"`
	Href    string `json:"href"`
	Name    string `json:"name"`
	Style   string `json:"style"`
	Content string `json:"content"`
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

		lineContent := make(map[string][]HtmlEle)

		for _, children := range element.Children {

			var ele HtmlEle
			attr := children.Attributes
			content := children.Content

			if y, ok := attr["y"]; ok {
				// cont += content
				// if newline, ok := attr["newline"]; ok && newline == "true" && children.Name == "text" {
				// 	cont += "<br/>"
				// 	ele.Content = cont
				// 	cont = ""
				// } else {
				// 	ele.Content = ""
				// }
				if children.Name == "text" {
					ele.Content = content
				} else {
					ele.Content = ""
				}

				ele.Style = attr["style"]
				ele.X = attr["x"]
				ele.Y = attr["y"]
				ele.Width = attr["width"]
				ele.Height = attr["height"]

				if _, ok := attr["href"]; ok && children.Name == "image" {
					ele.Href = attr["href"]
				} else {
					ele.Href = ""
				}

				ele.Name = children.Name

				if (children.Name == "text" && ele.Content != "") ||
					children.Name == "image" {
					lineContent[y] = append(lineContent[y], ele)
				}
			}
		}

		keys := make([]string, 0, len(lineContent))
		for k := range lineContent {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		cont := ""
		for _, item := range lineContent {
			result += `
		<p>`
			for i, v := range item {
				switch v.Name {
				case "image":
					result += `<img width="` + v.Width + `" height="` + v.Height + `" src="` + v.Href + `"/>`
				case "text":
					cont += v.Content
					if i == len(item)-1 {
						result += `<span style="` + v.Style + `">` + cont + `</span>`
					}
				}
			}
			cont = ""
			result += `</p>`
		}
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

	if err = WriteFileWithTrunc(fileName, result); err != nil {
		return
	}
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

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlfile)))
	pdfg.Dpi.Set(300)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.MarginBottom.Set(40)
	pdfg.MarginLeft.Set(30)

	err = pdfg.Create()
	if err != nil {
		fmt.Printf("pdfg create err: %#v\n", err)
		return
		// log.Fatal(err)
	}

	fileName, err = FilePath(filePreName, "pdf", false)
	if err != nil {
		return err
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(fileName)
	if err != nil {
		fmt.Printf("pdfg WriteFile err: %#v\n", err)
		return
		// log.Fatal(err)
	}
	return
}
