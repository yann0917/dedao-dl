package utils

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func Md2Pdf(path, title string, md []byte) (err error) {
	title = FileName(title, "pdf")
	filePreName := filepath.Join(path, title)
	fileName, err := FilePath(filePreName, "", false)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)

	h := mdToHTML(md)
	article := genHeadHtml() + string(h) + `
</body>
</html>`
	buf.Write([]byte(article))
	pdf := PdfOption{
		FileName: fileName,
		PageSize: "A4",
		Toc:      false,
	}
	fmt.Printf("正在生成文件：【\033[37;1m%s\033[0m】 ", title)
	err = pdf.GenPdf(buf)
	return
}

func genHeadHtml() (result string) {
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
		body {font-family: PingFang SC,Arial,sans-serif,Source Code Pro;color: #333;text-align: left;line-height: 1.8;}
		em {font-style: normal;}
		h2>code { background-color: rgb(255, 96, 2);padding: 0.5%;border-radius: 10%;color: white;}
		p>em {color: rgb(255, 96, 2);}
	</style>
</head>
<body>
`
	return
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
