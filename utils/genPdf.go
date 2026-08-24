package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PdfOption struct {
	FileName  string
	CoverPath string
	PageSize  string
	Toc       bool
}

func (p *PdfOption) GenPdf(buf *bytes.Buffer) (err error) {
	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	page := wkhtmltopdf.NewPageReader(buf)
	page.FooterFontSize.Set(10)
	page.FooterRight.Set("[page]")
	page.DisableSmartShrinking.Set(true)

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	if p.CoverPath != "" {
		pdfg.Cover.EnableLocalFileAccess.Set(true)
		dir, err := os.Getwd()
		if err != nil {
			pdfg.Cover.EnableLocalFileAccess.Set(false)
		}
		dir = filepath.Join(dir, p.CoverPath)
		if runtime.GOOS == "windows" {
			// Windows: 反斜杠路径会被 wkhtmltopdf 当作 URL 协议(c:)解析失败，
			// 统一转为正斜杠 file:// URL
			dir = strings.ReplaceAll(dir, "\\", "/")
			pdfg.Cover.Input = "file:///" + dir
		} else {
			pdfg.Cover.Input = "file://" + dir
		}
	}

	pdfg.Dpi.Set(300)
	if p.Toc {
		pdfg.TOC.Include = true
		pdfg.TOC.TocHeaderText.Set("目 录")
		pdfg.TOC.HeaderFontSize.Set(18)

		pdfg.TOC.TocLevelIndentation.Set(15)
		pdfg.TOC.TocTextSizeShrink.Set(0.9)
		pdfg.TOC.DisableDottedLines.Set(false)
		pdfg.TOC.EnableTocBackLinks.Set(true)
	}

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
	err = pdfg.WriteFile(p.FileName)
	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		return
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
	if p.CoverPath != "" {
		err = os.Remove(p.CoverPath)
	}
	return
}
