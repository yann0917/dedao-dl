package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

type PdfOption struct {
	FileName  string
	CoverPath string
	PageSize  string
	// Toc 预留：目录页已改由 Svg2Pdf 自绘 HTML（wkhtmltopdf 内置 TOC
	// 行距不可控，同样存在 1em 行距与页顶裁切问题）
	Toc bool
}

func (p *PdfOption) GenPdf(buf *bytes.Buffer) (err error) {
	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	page := wkhtmltopdf.NewPageReader(buf)
	page.FooterFontSize.Set(10)
	page.FooterRight.Set("[page]")
	page.DisableSmartShrinking.Set(true)

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.MarginTop.Set(15)
	pdfg.MarginBottom.Set(15)
	pdfg.MarginLeft.Set(15)
	pdfg.MarginRight.Set(15)

	// 封面单独渲染：wkhtmltopdf 的页边距是全局的，封面和正文同趟渲染
	// 必然带白边；先 0 边距渲染封面，再与正文合并实现满版。
	mainPDF := p.FileName
	var coverPDF string
	if p.CoverPath != "" {
		coverPDF, err = p.genCoverPdf(p.CoverPath)
		if err != nil {
			return fmt.Errorf("生成封面失败: %w", err)
		}
		defer os.Remove(coverPDF)

		tmp, err := os.CreateTemp("", "dedao-main-*.pdf")
		if err != nil {
			return err
		}
		tmp.Close()
		mainPDF = tmp.Name()
		defer os.Remove(mainPDF)
	}

	err = pdfg.Create()
	if err != nil {
		fmt.Printf("pdfg create err: %#v\n", err)
		return
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(mainPDF)
	if err != nil {
		fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
		return
	}

	if coverPDF != "" {
		if err = mergeCoverAndMain(coverPDF, mainPDF, p.FileName); err != nil {
			fmt.Printf("\033[31;1m%s\033[0m\n", "失败"+err.Error())
			return
		}
	}
	fmt.Printf("\033[32;1m%s\033[0m\n", "完成")
	if p.CoverPath != "" {
		err = os.Remove(p.CoverPath)
	}
	return
}

// mergeCoverAndMain 合并封面与正文。
// pdfcpu 的 MergeCreateFile 会把每个输入文件包成以文件名命名的顶层书签
// （dedao-cover / dedao-main，原大纲挂在 dedao-main 之下），因此合并后
// 需要用正文自身的大纲（页码整体偏移封面页数）替换掉书架大纲。
func mergeCoverAndMain(coverPDF, mainPDF, outFile string) error {
	coverPages, err := api.PageCountFile(context.Background(), coverPDF)
	if err != nil {
		return err
	}
	if err = api.MergeCreateFile(context.Background(), []string{coverPDF, mainPDF}, outFile, false, nil); err != nil {
		return err
	}

	// pdfcpu 合并不重映射命名目的地（/Dests）里的整数页号（0 基），
	// 正文页因封面整体后移，目录链接会全部落早一页，这里手动偏移。
	if err = shiftNamedDestPageInts(outFile, coverPages); err != nil {
		return err
	}

	mainFile, err := os.Open(mainPDF)
	if err != nil {
		return err
	}
	defer mainFile.Close()
	bms, err := api.Bookmarks(context.Background(), mainFile, nil)
	if err != nil || len(bms) == 0 {
		// 正文无大纲时无书架可替换，保持合并产物
		return nil
	}
	// pdfcpu 导出书签时对整数型目的地（wkhtmltopdf 的写法）按 PDF 规范返回 0 基页号，
	// 而导入 BookmarksJSON 按 1 基页号解析，因此偏移量 = 封面页数 + 1
	offsetBookmarks(bms, coverPages+1)

	buf, err := json.Marshal(map[string]any{"bookmarks": bms})
	if err != nil {
		return err
	}

	mergedFile, err := os.Open(outFile)
	if err != nil {
		return err
	}
	out, err := os.CreateTemp("", "dedao-outline-*.pdf")
	if err != nil {
		mergedFile.Close()
		return err
	}

	importErr := api.ImportBookmarks(context.Background(), mergedFile, bytes.NewReader(buf), out, true, nil)
	closeErr := errors.Join(mergedFile.Close(), out.Close())
	if importErr != nil {
		os.Remove(out.Name())
		return importErr
	}
	if closeErr != nil {
		os.Remove(out.Name())
		return closeErr
	}
	if err = os.Rename(out.Name(), outFile); err != nil {
		os.Remove(out.Name())
		return err
	}
	return nil
}

// offsetBookmarks 递归地把书签页码加上封面页数，使大纲指向合并后的正确页面。
func offsetBookmarks(bms []pdfcpu.Bookmark, offset int) {
	for i := range bms {
		bms[i].PageFrom += offset
		offsetBookmarks(bms[i].Kids, offset)
	}
}

// shiftNamedDestPageInts 把 inFile 的命名目的地（catalog /Dests，含名字树形式）
// 目的地数组首元素的整数页号（0 基）整体加上 offset，有修改时写回原文件。
func shiftNamedDestPageInts(inFile string, offset int) error {
	cc, err := api.ReadContextFile(context.Background(), inFile)
	if err != nil {
		return err
	}
	shifted := 0
	if destsIndRef, found := cc.RootDict["Dests"]; found {
		if destsObj, err := cc.Dereference(destsIndRef); err == nil {
			if destsDict, ok := destsObj.(types.Dict); ok {
				shifted, err = shiftDestsDict(cc, destsDict, offset)
				if err != nil {
					return err
				}
			}
		}
	}
	if shifted == 0 {
		return nil
	}
	tmp, err := os.CreateTemp("", "dedao-dests-*.pdf")
	if err != nil {
		return err
	}
	if err = api.WriteContextFile(context.Background(), cc, tmp.Name()); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Rename(tmp.Name(), inFile); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}

// shiftDestsDict 遍历 /Dests 字典（值 = 目的地数组或名字树字典），偏移其中的整数页号。
func shiftDestsDict(cc *model.Context, destsDict types.Dict, offset int) (int, error) {
	shifted := 0
	shiftValue := func(v types.Object) {
		obj, err := cc.Dereference(v)
		if err != nil {
			return
		}
		if arr, ok := obj.(types.Array); ok && len(arr) > 0 {
			if i, ok := arr[0].(types.Integer); ok {
				arr[0] = types.Integer(i.Value() + offset)
				shifted++
			}
		}
	}
	for _, v := range destsDict {
		obj, err := cc.Dereference(v)
		if err != nil {
			continue
		}
		switch o := obj.(type) {
		case types.Array:
			shiftValue(o)
		case types.Dict:
			// 名字树形式：/Names 为 [name, dest, name, dest...]，/Kids 为子树
			if namesObj, found := o["Names"]; found {
				if namesArr, err := cc.Dereference(namesObj); err == nil {
					if a, ok := namesArr.(types.Array); ok {
						for i := 1; i < len(a); i += 2 {
							shiftValue(a[i])
						}
					}
				}
			}
			if kidsObj, found := o["Kids"]; found {
				if kidsArr, err := cc.Dereference(kidsObj); err == nil {
					if a, ok := kidsArr.(types.Array); ok {
						for _, kid := range a {
							if kidObj, err := cc.Dereference(kid); err == nil {
								if kd, ok := kidObj.(types.Dict); ok {
									n, err := shiftDestsDict(cc, kd, offset)
									shifted += n
									if err != nil {
										return shifted, err
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return shifted, nil
}

// genCoverPdf 用 0 页边距单独渲染封面，输出临时 PDF 路径，由调用方负责删除。
func (p *PdfOption) genCoverPdf(coverHTMLPath string) (string, error) {
	f, err := os.Open(coverHTMLPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	page := wkhtmltopdf.NewPageReader(f)
	page.DisableSmartShrinking.Set(true)
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	if err = pdfg.Create(); err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "dedao-cover-*.pdf")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if err = pdfg.WriteFile(tmp.Name()); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}
