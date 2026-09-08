package utils

import (
	"bytes"
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/JoshVarga/svgparser"
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
	// AnchorName / AnchorChapter 记录注释引用 <a href> 的原始锚点名与所属章节，
	// 用于在渲染前把「正文角标」与「注释正文」配对，统一成 fn-ref / fn-note 锚点。
	AnchorName    string `json:"anchor_name"`
	AnchorChapter string `json:"anchor_chapter"`
	TextAlign     string `json:"text_align"` // left; center; right
}

type SvgRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Rx     string
	Ry     string
	Style  string
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
	footNoteImgW = 20 // 脚注图片≈11x11px & 特殊字图片≈19x19
	footNoteImgH = 20 // 行内图片高度=20

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

	// 整书级注释配对：正文角标 ↔ 注释正文 双向跳转
	chapters, bErr := buildBookChapters(svgContents)
	if bErr != nil {
		return bErr
	}
	assignFootnoteAnchors(chapters)

	for k, svgContent := range svgContents {
		chapter, coverContent, err1 := OneByOneHtml(eBookTypePdf, k, svgContent, toc, chapters[k].pages)
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
	pdf := PdfOption{
		FileName:  fileName,
		CoverPath: coverPath,
		PageSize:  "A4",
		Toc:       true,
	}
	err = pdf.GenPdf(buf)
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

	// 整书级注释配对：正文角标 ↔ 注释正文 双向跳转
	chapters, bErr := buildBookChapters(svgContents)
	if bErr != nil {
		err = bErr
		return
	}
	assignFootnoteAnchors(chapters)

	for k, svgContent := range svgContents {
		chapter, coverUrl, err1 := OneByOneHtml(eBookTypeEpub, k, svgContent, opt.Toc, chapters[k].pages)
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

// footnoteExpandScript 电子书脚注点击展开脚本：
// 点击脚注小图标，在图标旁显示注释浮层，浮层内容可选中复制。
const footnoteExpandScript = `<script>
(function () {
  if (window.__fnPopupBound) return;
  window.__fnPopupBound = true;
  var popup = document.createElement('div');
  popup.className = 'fn-popup';
  popup.style.display = 'none';
  document.body.appendChild(popup);
  function close() { popup.style.display = 'none'; }
  document.addEventListener('click', function (e) {
    var trigger = e.target.closest ? e.target.closest('.fn-trigger') : null;
    if (trigger) {
      e.preventDefault();
      if (popup.style.display === 'block') { close(); return; }
      popup.textContent = trigger.getAttribute('data-footnote') || '';
      popup.style.display = 'block';
      var rect = trigger.getBoundingClientRect();
      var pw = popup.offsetWidth, ph = popup.offsetHeight;
      var left = rect.left + rect.width / 2 - pw / 2;
      var top = rect.bottom + 6;
      if (left < 8) left = 8;
      if (left + pw > window.innerWidth - 8) left = window.innerWidth - pw - 8;
      if (top + ph > window.innerHeight - 8) top = window.innerHeight - ph - 8;
      if (top < 8) top = rect.top - ph - 6;
      popup.style.left = left + 'px';
      popup.style.top = top + 'px';
      return;
    }
    if (popup.contains(e.target)) return;
    close();
  });
})();
</script>`

// AllInOneHtml generate ebook content all in one html file
func AllInOneHtml(svgContents []*SvgContent, toc []EbookToc) (result string, err error) {
	result = GenHeadHtml()
	fnA, fnB = ParseBookFnDelimiter(svgContents)
	// 整书级注释配对：正文角标 ↔ 注释正文 双向跳转
	chapters, bErr := buildBookChapters(svgContents)
	if bErr != nil {
		return "", bErr
	}
	assignFootnoteAnchors(chapters)
	for k, svgContent := range svgContents {
		chapter, _, err1 := OneByOneHtml(eBookTypeHtml, k, svgContent, toc, chapters[k].pages)
		if err1 != nil {
			err = err1
			return
		}
		result += chapter
	}
	result += footnoteExpandScript + `
</body>
</html>`
	// 使用同样的方法处理反转义，保留已转义的HTML标签
	result = preserveEscapedHtmlTags(result)
	return
}

// locateTocAnchorID 在当前渲染行上定位 TOC 锚点
// offset 主信号：用 TOC 的 Offset（字节偏移）匹配本行 SVG text 的最小 offset，
// 定位到章节内标题所在行（含二级/三级目录）；匹配成功推进 offsetIdx。
// 文本兜底：对未提供有效 Offset 的目录项，按标题文本匹配定位。
// 返回锚点 id（TOC Href 中 # 后部分），无匹配返回空串。
func locateTocAnchorID(items []HtmlEle, offsetEntries []EbookToc, offsetIdx *int, textEntries []EbookToc, textMatched []bool, lineText string) string {
	// 计算本行最小的字节偏移（SVG text 的 offset 属性）
	lineMinOffset := -1
	for _, item := range items {
		if item.Offset == "" {
			continue
		}
		if o, err := strconv.Atoi(item.Offset); err == nil {
			if lineMinOffset == -1 || o < lineMinOffset {
				lineMinOffset = o
			}
		}
	}

	// offset 主信号：按字节偏移定位下一个未消费的目录项
	if *offsetIdx < len(offsetEntries) && lineMinOffset >= 0 {
		entry := offsetEntries[*offsetIdx]
		if lineMinOffset >= entry.Offset {
			*offsetIdx++
			if tagArr := strings.Split(entry.Href, "#"); len(tagArr) > 1 {
				return tagArr[1]
			}
		}
	}

	// 文本兜底：仅对未提供有效 offset 的目录项，按标题文本匹配
	if len([]rune(lineText)) >= 2 {
		for i, e := range textEntries {
			if textMatched[i] {
				continue
			}
			norm := strings.ReplaceAll(e.Text, " ", "")
			if norm != "" && strings.Contains(norm, lineText) {
				textMatched[i] = true
				if tagArr := strings.Split(e.Href, "#"); len(tagArr) > 1 {
					return tagArr[1]
				}
			}
		}
	}

	return ""
}

// OneByOneHtml one by one generate chapter html
// eType: html/pdf/epub, index: []*SvgContent index, svgContent: one chapter content
func OneByOneHtml(eType string, index int, svgContent *SvgContent, toc []EbookToc, prepared []chapterPage) (result, cover string, err error) {
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

	// 按章节分组 TOC 项，用于在章节内定位锚点（含二级/三级目录）
	// 有有效 Offset 的目录项走字节偏移匹配（主信号），无 Offset 的走标题文本匹配（兜底）
	offsetEntries := make([]EbookToc, 0)
	textEntries := make([]EbookToc, 0)
	for _, t := range toc {
		tagArr := strings.Split(t.Href, "#")
		if len(tagArr) > 0 && tagArr[0] == svgContent.ChapterID {
			if t.Offset > 0 {
				offsetEntries = append(offsetEntries, t)
			} else {
				textEntries = append(textEntries, t)
			}
		}
	}
	sort.Slice(offsetEntries, func(i, j int) bool { return offsetEntries[i].Offset < offsetEntries[j].Offset })
	offsetIdx := 0
	textMatched := make([]bool, len(textEntries))

	// 视图渲染：优先使用调用方传入（已整书解析并配对）的 pages，
	// 否则回退到单章解析（此时不做跨章节配对，仅供 PDF 等兜底）。
	var pages []chapterPage
	if prepared != nil {
		pages = prepared
	} else {
		var err1 error
		pages, err1 = parseChapterPages(svgContent)
		if err1 != nil {
			err = err1
			return
		}
	}

	for _, pg := range pages {
		result += `
<div id="` + svgContent.ChapterID + `">`
		lineContent := pg.lineContent
		rects := pg.rects
		keys := pg.keys

		activeRectIdx := -1
		for _, v := range keys {
			rectIdx := matchRectForLine(rects, v)
			if rectIdx != activeRectIdx {
				if activeRectIdx >= 0 {
					result += `
</div>`
				}
				if rectIdx >= 0 {
					result += `
<div style="` + buildRectWrapperStyle(rects[rectIdx]) + `">`
				}
				activeRectIdx = rectIdx
			}
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
					// lineStyle 代表本行“正文”的样式：从行尾向前查找非注释引用、非图片
					// 的文本元素，避免注释角标（Fn.Href 非空）的蓝色小字号样式污染整行正文。
					lineStyle = findLineStyle(lineContent[v])
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
							if eType == eBookTypeHtml {
								// HTML：点击展开浮层，注释文本存入 data-footnote，由脚本读取展示，可选中复制
								img = `
	<sup class="fn-trigger" data-footnote="` + item.Alt + `"><img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
									`" src="` + item.Href +
									`" alt="` + item.Alt +
									`" class="` + item.Class +
									`"/></sup>`
							} else {
								// PDF：保留原生 title 悬浮提示
								img = `
	<sup><img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
									`" src="` + item.Href +
									`" alt="` + item.Alt +
									`" title="` + item.Alt +
									`" class="` + item.Class +
									`"/></sup>`
							}
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
	<sup><a class="duokan-footnote" epub:type="noteref" href="#` + footnoteId + `"> <img width="` + strconv.FormatFloat(w, 'f', 0, 64) +
									`" src="` + item.Href +
									`" alt="` + item.Alt +
									`" zy-footnote="` + item.Alt +
									`" class="` + item.Class + ` zhangyue-footnote qqreader-footnote"/></a></sup>`
								result += `<aside epub:type="footnote" id="` + footnoteId +
									`"><ol class="duokan-footnote-content" style="list-style:none;padding:0px;margin:0px;"><li class="duokan-footnote-item" id="` +
									footnoteId + `"></a>` + item.Alt + `</li></ol></aside>`
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
						cont += "</span>"
						hasUncloseSpan = false
					}

					if item.Style != lineStyle && !hasUncloseSpan {
						cont += fmt.Sprintf(`<span style="%s">`, item.Style)
						currentSpanStyle = item.Style
						hasUncloseSpan = true
					}

					if firstX >= centerL && firstX <= centerH {
						style += "display: block;text-align:center;"
					} else if firstX >= rightL {
						style += "display: block;text-align:right;"
					}

					item.Content = html.EscapeString(item.Content)

					tags := []struct {
						condition bool
						open      string
						close     string
					}{
						{item.IsBold, "<b>", "</b>"},
						{item.IsItalic, "<i>", "</i>"},
						{item.IsFn, "<sup>", "</sup>"},
						{item.IsSub, "<sub>", "</sub>"},
					}

					for _, tag := range tags {
						if tag.condition {
							cont += tag.open
						}
					}

					if item.Fn.Href != "" {
						if item.ID != "" {
							cont += fmt.Sprintf(`<a id="%s" href="%s"`, item.ID, item.Fn.Href)
						} else {
							// 同一注释项拆分出的非首段只挂 href，不重复 id
							cont += fmt.Sprintf(`<a href="%s"`, item.Fn.Href)
						}
						if item.Fn.Style != "" {
							cont += fmt.Sprintf(` style="%s"`, item.Fn.Style)
						}
						cont += ">"
					}

					cont += item.Content

					if item.Fn.Href != "" {
						cont += "</a>"
					}

					for i := len(tags) - 1; i >= 0; i-- {
						if tags[i].condition {
							cont += tags[i].close
						}
					}

					contWOTag += item.Content
				}
				if i == len(lineContent[v])-1 {
					matchH := false
					contWOTag = html.UnescapeString(contWOTag)

					level := 0
					for k, v := range tocLevel {
						contWOTagMatch := strings.ReplaceAll(contWOTag, "&nbsp;", "")
						if strings.Contains(strings.ReplaceAll(k, " ", ""), contWOTagMatch) {
							matchH, level = true, v
							break
						}
					}
					if contWOTag != "" {
						// 定位本行对应的 TOC 锚点（offset 主信号 + 文本兜底）
						anchorID := locateTocAnchorID(lineContent[v], offsetEntries, &offsetIdx, textEntries, textMatched, contWOTag)
						// 行内已有同名 id（SVG 自带锚点）时避免生成重复 id
						if anchorID != "" && id == anchorID {
							anchorID = ""
						}
						if matchH {
							result += `
</div>`
							result += `<div class='header` + strconv.Itoa(level) + `'`
							if anchorID != "" {
								result += ` id="` + anchorID + `"`
							}
							result += `>` + GenTocLevelHtml(level, true)
						} else {
							result += `
	<p`
							if anchorID != "" {
								result += ` id="` + anchorID + `"`
							}
							result += `>`
						}
					}
					if i > 1 && item.Name == "image" {
						style = lineContent[v][i-1].Style
					}
					if cont != "" {
						// 保留每个元素的原始样式
						// 外层包裹整行：若行末是注释角标，用本行“正文”样式 lineStyle，
						// 避免角标的小蓝样式污染整行正文（正文变蓝/变小）；
						// 行末为正常正文/图片时保留 style（含 display:block/居中）。
						wrapStyle := style
						if item.Fn.Href != "" && lineStyle != "" {
							wrapStyle = lineStyle
						}
						if id != "" && wrapStyle != "" {
							result += `<span id="` + id + `" style="` + wrapStyle + `">`
						} else {
							if id != "" {
								result += `<span id="` + id + `">`
							}
							if wrapStyle != "" {
								// 确保样式正确应用
								result += `<span style="` + wrapStyle + `">`
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
		if activeRectIdx >= 0 {
			result += `
</div>`
		}
		result += `</div>`
		switch eType {
		case eBookTypePdf, eBookTypeEpub:
			result += `
</body>
</html>`
		}
	}

	// 防止将已转义的HTML标签还原
	// 比如 &lt;script&gt; 变成 <script>
	// 仅对非HTML标签的实体进行反转义，如 &nbsp; &quot; 等
	result = preserveEscapedHtmlTags(result)

	return
}

// preserveEscapedHtmlTags 保留所有转义的HTML标签，只处理特定的实体符号
func preserveEscapedHtmlTags(content string) string {
	// 首先找出并保存所有转义的HTML标签
	// 通用模式：&lt;任何标签&gt; 和 &lt;/任何标签&gt;
	preservedMap := make(map[string]string)
	result := content

	// 1. 匹配所有转义的开始标签 &lt;tag...&gt;
	openTagPattern := regexp.MustCompile(`&lt;[a-zA-Z][^&]*&gt;`)
	openTags := openTagPattern.FindAllString(result, -1)
	for i, match := range openTags {
		placeholder := fmt.Sprintf("__PRESERVED_OPEN_TAG_%d__", i)
		preservedMap[placeholder] = match
		result = strings.ReplaceAll(result, match, placeholder)
	}

	// 2. 匹配所有转义的结束标签 &lt;/tag&gt;
	closeTagPattern := regexp.MustCompile(`&lt;/[a-zA-Z][^&]*&gt;`)
	closeTags := closeTagPattern.FindAllString(result, -1)
	for i, match := range closeTags {
		placeholder := fmt.Sprintf("__PRESERVED_CLOSE_TAG_%d__", i)
		preservedMap[placeholder] = match
		result = strings.ReplaceAll(result, match, placeholder)
	}

	// 3. 匹配所有转义的自闭合标签 &lt;tag...&gt;
	selfCloseTagPattern := regexp.MustCompile(`&lt;[a-zA-Z][^&]*/&gt;`)
	selfCloseTags := selfCloseTagPattern.FindAllString(result, -1)
	for i, match := range selfCloseTags {
		placeholder := fmt.Sprintf("__PRESERVED_SELFCLOSE_TAG_%d__", i)
		preservedMap[placeholder] = match
		result = strings.ReplaceAll(result, match, placeholder)
	}

	// 需要处理的特定实体符号列表
	entities := map[string]string{
		"&nbsp;":   " ",  // 不间断空格
		"&ensp;":   " ",  // 半角空格
		"&emsp;":   " ",  // 全角空格
		"&quot;":   "\"", // 双引号
		"&apos;":   "'",  // 单引号
		"&amp;":    "&",  // 和号
		"&mdash;":  "—",  // 破折号
		"&ndash;":  "–",  // 连字符
		"&hellip;": "…",  // 省略号
		"&copy;":   "©",  // 版权符号
		"&reg;":    "®",  // 注册商标
		"&trade;":  "™",  // 商标
		"&deg;":    "°",  // 度数符号
		"&plusmn;": "±",  // 正负号
	}

	// 替换特定的实体符号
	for entity, replacement := range entities {
		result = strings.ReplaceAll(result, entity, replacement)
	}

	// 恢复所有转义的HTML标签
	for placeholder, original := range preservedMap {
		result = strings.ReplaceAll(result, placeholder, original)
	}

	return result
}

func GenHeadHtml() (result string) {
	result = `<!DOCTYPE html>
<html lang="zh-CN" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
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
		p { margin: 1.5em 0; }
		img { page-break-inside: avoid; max-width: 100% !important;}
		img.epub-footnote { margin-right:5px;display: inline;font-size: 12px;}
		/* 脚注点击展开浮层：内容可选中复制 */
		.fn-trigger { cursor: pointer; }
		.fn-popup { position: fixed; z-index: 9999; background: #fff; border: 1px solid #eee; box-shadow: 0 2px 12px rgba(0,0,0,.18); padding: 10px 14px; border-radius: 6px; max-width: 320px; font-size: 14px; line-height: 1.7; color: #333; user-select: text; white-space: normal; word-break: break-word; }
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

func GenTocLevelHtml(level int, startTag bool) string {
	tags := map[bool]map[int]string{
		true:  {0: "<h1>", 1: "<h2>", 2: "<h3>", 3: "<h4>", 4: "<h5>", 5: "<h6>"},
		false: {0: "</h1>", 1: "</h2>", 2: "</h3>", 3: "</h4>", 4: "</h5>", 5: "</h6>"},
	}

	if tag, ok := tags[startTag][level]; ok {
		return tag
	}
	return ""
}

func GenLineContentByElement(chapterID string, element *svgparser.Element) (lineContent map[float64][]HtmlEle) {
	lineContent = make(map[float64][]HtmlEle)
	offset := ""
	lastY, lastTop, lastH, lastName := "", "", "", ""

	for k, children := range element.Children {
		var ele HtmlEle
		attr := children.Attributes
		content := children.Content
		ele.Newline = parseAttrNewline(attr)
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
									if len(tagArr) > 1 {
										// 记录注释引用的原始锚点信息，渲染前由 assignFootnoteAnchors
										// 把「正文角标」与「注释正文」配对，统一生成 fn-ref / fn-note 锚点。
										ele.AnchorChapter = tagArr[0]
										ele.AnchorName = tagArr[1]
										// 兜底：未配对前先构造成同章节锚点，保证链接可展示、可点击。
										ele.Fn.Href = "#" + tagArr[0] + "_" + tagArr[1]
									} else {
										ele.Fn.Href = "#" + tagArr[0]
									}
									// 不再用 chapterID 前缀覆盖原生 id，保留 SVG 原始 id（注释正文首段自带的）。
									ele.Fn.Style = attrC["style"]
								}
							}
						}
					} else {
						ele.Content = "&nbsp;"
					}
				}

				if _, ok := attr["top"]; ok {
					topInt, _ := strconv.ParseFloat(attr["top"], 64)
					heightInt, _ := strconv.ParseFloat(attr["height"], 64)
					lenInt, _ := strconv.ParseFloat(attr["len"], 64)
					lastTopInt, _ := strconv.ParseFloat(lastTop, 64)
					lastHInt, _ := strconv.ParseFloat(lastH, 64)

					// 判断是否可能是上标或下标，如果是，忽略 newline 设置
					isPossibleSuperOrSub := (heightInt < lastHInt*0.8) && // 提高高度比例要求，确保只有明显更小的文字才被识别为上标
						(children.Name == lastName || (content != "" && len(content) <= 3 && isNumericOrMathSymbol(content)))

					// 检查是否有字体大小指示为上标
					fontSizeIsSmaller := false
					if style, ok := attr["style"]; ok {
						// 仅当字体明确小于16px时才视为可能的上标
						if strings.Contains(style, "font-size:11px") ||
							strings.Contains(style, "font-size:12px") ||
							strings.Contains(style, "font-size:13px") {
							fontSizeIsSmaller = true
						}
					}

					// 使用更严格的条件组合判断上标
					if fontSizeIsSmaller {
						isPossibleSuperOrSub = true
					}

					// 检查Y坐标是否不同，通常上标的Y坐标会比基线的Y坐标小
					yInt, _ := strconv.ParseFloat(attr["y"], 64)
					lastYInt, _ := strconv.ParseFloat(lastY, 64)
					// 上下标必须紧贴前一个文本元素；跨行的数字开头正文不能按上下标处理。
					if isPossibleSuperOrSub && yInt != 0 && lastYInt != 0 && lastHInt > 0 {
						yDiff := yInt - lastYInt
						if yDiff < 0 {
							yDiff = -yDiff
						}
						if yDiff > lastHInt*0.5 {
							isPossibleSuperOrSub = false
						}
					}
					if yInt != 0 && lastYInt != 0 && yInt < lastYInt &&
						(lastYInt-yInt > 2) { // 至少需要有明显的Y轴差异
						// Y坐标比前一个元素小，很可能是上标
						isPossibleSuperOrSub = true
					}

					if isPossibleSuperOrSub {
						// 如果满足上标条件，则优先识别上标

						// 使用更严格的条件判断是否真的是上标或下标
						isLikelyPower := (content != "" && len(content) <= 3 && isNumericOrMathSymbol(content))
						// 对于特殊情况，如N³或2ⁿ，需要额外判断
						if lenInt <= 5 && fontSizeIsSmaller {
							isLikelyPower = true
						}

						if isLikelyPower {
							// 根据相对位置判断是上标还是下标
							if topInt < lastTopInt {
								ele.IsFn = true
								// 上标和下标元素不应该被视为新行
								ele.Newline = false
							} else {
								ele.IsSub = true
								ele.Newline = false
							}
							// 不要清空原始样式，会导致字体样式丢失
							// attr["style"] = ""
						} else {
							lastTop = attr["top"]
							lastH = attr["height"]
						}
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

			// 小图标通常跟正文同行，但如果前一个元素已经显式换行，
			// 则应当跟随下一行文本，而不是挂到上一行末尾。
			yInt, _ := strconv.ParseFloat(ele.Y, 64)
			w, _ := strconv.ParseFloat(ele.Width, 64)
			if children.Name == "image" && w < footNoteImgW {
				aligned := false
				if k > 0 {
					attrPre := element.Children[k-1].Attributes
					prevHasExplicitNewline := attrPre["newline"] == "true"
					if !prevHasExplicitNewline {
						if prevY, ok := attrPre["y"]; ok {
							yInt, _ = strconv.ParseFloat(prevY, 64)
							ele.Y = prevY
							aligned = true
						}
					}
				}
				if !aligned && k+1 < len(element.Children) {
					attrNext := element.Children[k+1].Attributes
					if nextY, ok := attrNext["y"]; ok {
						yInt, _ = strconv.ParseFloat(nextY, 64)
						ele.Y = nextY
					}
				}
			}
			// 捕获 id 和 offset，用于后续 TOC 锚点定位
			ele.ID = attr["id"]
			if _, ok := attr["offset"]; ok {
				offset = attr["offset"]
			}
			ele.Offset = offset
			ele.Href = parseAttrHref(attr)
			ele.Alt = parseAttrAlt(attr)
			ele.Name = children.Name

			if (children.Name == "text") ||
				children.Name == "image" {
				lineContent[yInt] = append(lineContent[yInt], ele)
			}
			lastName = children.Name
		}
	}
	for y := range lineContent {
		sort.SliceStable(lineContent[y], func(i, j int) bool {
			xi, errI := strconv.ParseFloat(lineContent[y][i].X, 64)
			xj, errJ := strconv.ParseFloat(lineContent[y][j].X, 64)
			if errI != nil || errJ != nil {
				return lineContent[y][i].X < lineContent[y][j].X
			}
			return xi < xj
		})
	}
	return
}

// findLineStyle 从行尾向前查找本行“正文”样式：
// 跳过注释角标（AnchorName 非空）与图片，避免角标蓝色小字号污染整行正文；
// 若整行都是角标/图片，则回退到行末元素样式。
func findLineStyle(line []HtmlEle) string {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i].AnchorName == "" && line[i].Name != "image" {
			return line[i].Style
		}
	}
	if len(line) > 0 {
		return line[len(line)-1].Style
	}
	return ""
}

// chapterPage 是一章内单页的解析产物；bookChapter 是一章的全部页。
type chapterPage struct {
	chapterID   string
	lineContent map[float64][]HtmlEle
	rects       []SvgRect
	keys        []float64
}

type bookChapter struct {
	chID  string
	pages []chapterPage
}

// parseChapterPages 解析一章的全部页面，返回每页的 lineContent/rects/keys，
// 供整书注释配对与渲染使用（正文角标与章尾注释可能落在不同页/不同章节）。
func parseChapterPages(svgContent *SvgContent) ([]chapterPage, error) {
	pages := make([]chapterPage, 0, len(svgContent.Contents))
	for _, content := range svgContent.Contents {
		processedContent := preprocessSvgContent(content)
		valid := NewValidUTF8Reader(strings.NewReader(processedContent))
		validReader := []byte(processedContent)
		_, _ = valid.Read(validReader)
		element, err1 := svgparser.Parse(bytes.NewReader(validReader), false)
		if err1 != nil {
			return nil, err1
		}
		lineContent := GenLineContentByElement(svgContent.ChapterID, element)
		rects := GenRectContentByElement(element)
		keys := make([]float64, 0, len(lineContent))
		for k := range lineContent {
			keys = append(keys, k)
		}
		sort.Float64s(keys)
		pages = append(pages, chapterPage{
			chapterID:   svgContent.ChapterID,
			lineContent: lineContent,
			rects:       rects,
			keys:        keys,
		})
	}
	return pages, nil
}

// buildBookChapters 解析整本书的所有章节，返回每章的全部页（含各页 lineContent），
// 供 assignFootnoteAnchors 进行整书级注释配对，再逐章交由 OneByOneHtml 渲染。
func buildBookChapters(svgContents []*SvgContent) ([]bookChapter, error) {
	chapters := make([]bookChapter, 0, len(svgContents))
	for _, sc := range svgContents {
		pages, err := parseChapterPages(sc)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, bookChapter{chID: sc.ChapterID, pages: pages})
	}
	return chapters, nil
}

// assignFootnoteAnchors 在渲染前把【整本书】的注释引用配对，统一生成互跳锚点。
// 判定规则（与文字形式无关，只看结构）：
//   - 带 <a href="#anchor"> 的 text 说明它引用了一段注释（正文角标 / 注释正文段）。
//   - 原始 id 非空的段视为「注释正文」首段（注释正文首段自带 id，如 id="ch1"）；
//     无原始 id 的段视为「正文角标」或注释正文的后续段（[、1、] 拆分的非首段）。
//   - 当「正文角标组的 AnchorName == 注释正文组首段的原始 id」时，二者构成一对，
//     分配 fn-ref-<角标所在章节>_<序号>（正文角标）与 fn-note-<注释所在章节>_<序号>（注释正文），
//     两条锚点互相指向，从而实现正文角标 ↔ 注释正文 的双向跳转。
//
// 非配对段保持原始 href；同一注释项拆分出的多段只有首段挂 id，避免重复 id。
func assignFootnoteAnchors(chapters []bookChapter) {
	type fnAnchorPos struct {
		chapter int
		page    int
		lineKey float64
		index   int
	}
	type fnGroup struct {
		pos   []fnAnchorPos
		rawID string // 该组首个带原始 id 的段（注释正文首段）
	}
	groups := make(map[string]*fnGroup)
	var order []string
	for ci, ch := range chapters {
		for pi, pg := range ch.pages {
			for k, line := range pg.lineContent {
				for i, ele := range line {
					if ele.AnchorName == "" {
						continue
					}
					g, ok := groups[ele.AnchorName]
					if !ok {
						g = &fnGroup{}
						groups[ele.AnchorName] = g
						order = append(order, ele.AnchorName)
					}
					g.pos = append(g.pos, fnAnchorPos{chapter: ci, page: pi, lineKey: k, index: i})
					if g.rawID == "" && ele.ID != "" {
						g.rawID = ele.ID
					}
				}
			}
		}
	}
	if len(groups) == 0 {
		return
	}

	// 以「注释正文首段的原始 id」作桥，建立 rawID -> 注释正文组 的索引
	idToGroup := make(map[string]string)
	for key, g := range groups {
		if g.rawID != "" {
			idToGroup[g.rawID] = key
		}
	}

	// g 为要写入的组；ownID 为该组自己的锚点 id；targetHref 为该组 href 指向的锚点 id
	write := func(g *fnGroup, ownID, targetHref string) {
		for i, p := range g.pos {
			if i == 0 {
				chapters[p.chapter].pages[p.page].lineContent[p.lineKey][p.index].ID = ownID
			}
			chapters[p.chapter].pages[p.page].lineContent[p.lineKey][p.index].Fn.Href = "#" + targetHref
		}
	}

	seq := 0
	for _, key := range order {
		g := groups[key]
		target := groups[idToGroup[key]]
		if target == nil || target == g {
			// 未配对（或自引用），交给兜底
			continue
		}
		seq++
		// 角标组锚点前缀用其所在章节；注释正文组用其所在章节（两者可能不同章节）。
		refID := fmt.Sprintf("fn-ref-%s_%d", chapters[g.pos[0].chapter].chID, seq)
		noteID := fmt.Sprintf("fn-note-%s_%d", chapters[target.pos[0].chapter].chID, seq)
		write(g, refID, noteID)      // 正文角标组：id=refID，href=#noteID
		write(target, noteID, refID) // 注释正文组：id=noteID，href=#refID
	}

	// 兜底：未配对且无原始 id 的段，补一个唯一 id，避免渲染出空的 id 属性
	for _, key := range order {
		g := groups[key]
		if len(g.pos) == 0 {
			continue
		}
		first := g.pos[0]
		if chapters[first.chapter].pages[first.page].lineContent[first.lineKey][first.index].ID == "" {
			chapters[first.chapter].pages[first.page].lineContent[first.lineKey][first.index].ID = "fn-orphan-" + chapters[first.chapter].chID + "_" + key
		}
	}
}

func GenRectContentByElement(element *svgparser.Element) (rects []SvgRect) {
	for _, child := range element.Children {
		if child.Name != "rect" {
			continue
		}
		attr := child.Attributes
		x, errX := strconv.ParseFloat(attr["x"], 64)
		y, errY := strconv.ParseFloat(attr["y"], 64)
		w, errW := strconv.ParseFloat(attr["width"], 64)
		h, errH := strconv.ParseFloat(attr["height"], 64)
		if errX != nil || errY != nil || errW != nil || errH != nil {
			continue
		}
		rects = append(rects, SvgRect{
			X:      x,
			Y:      y,
			Width:  w,
			Height: h,
			Rx:     attr["rx"],
			Ry:     attr["ry"],
			Style:  strings.ReplaceAll(attr["style"], "fill", "background"),
		})
	}
	sort.SliceStable(rects, func(i, j int) bool {
		return rects[i].Y < rects[j].Y
	})
	return
}

func matchRectForLine(rects []SvgRect, lineY float64) int {
	for i, rect := range rects {
		if lineY >= rect.Y && lineY <= rect.Y+rect.Height {
			return i
		}
	}
	return -1
}

func buildRectWrapperStyle(rect SvgRect) string {
	styles := []string{
		rect.Style,
		"padding: 0.2em 1.2em;",
		"margin: 1em 0;",
	}
	if rect.Rx != "" {
		styles = append(styles, "border-radius:"+rect.Rx+"px;")
	} else if rect.Ry != "" {
		styles = append(styles, "border-radius:"+rect.Ry+"px;")
	}
	return strings.Join(styles, "")
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
	// 如果元素有明显的上标特征，即使有 newline="true" 也应该忽略
	if newline, ok := attr["newline"]; ok && newline == "true" {
		// 检查是否可能是上标
		if _, topOk := attr["top"]; topOk {
			if height, heightOk := attr["height"]; heightOk {
				heightInt, _ := strconv.ParseFloat(height, 64)
				// 使用更精确的高度阈值
				if heightInt <= 16 { // 典型上标的高度通常小于16
					return false
				}
			}

			// 检查字体大小，是否明显小于正常文本
			if style, sizeOk := attr["style"]; sizeOk {
				// 仅当字体明确小于16px时才视为可能的上标
				if strings.Contains(style, "font-size:11px") ||
					strings.Contains(style, "font-size:12px") ||
					strings.Contains(style, "font-size:13px") {
					return false
				}
			}

			// 检查内容长度，上标通常很短
			if len, lenOk := attr["len"]; lenOk {
				lenInt, _ := strconv.ParseFloat(len, 64)
				if lenInt <= 5 { // 上标通常只有几个字符
					return false
				}
			}
		}
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

func preprocessSvgContent(content string) string {
	// 主要目的是处理实体符号，如 &nbsp; 等
	// 我们不需要删除 HTML 标签，因为它们已经是转义形式（如 &lt;script&gt;）
	// 在后续的 preserveEscapedHtmlTags 中会处理它们

	// 保留内容不变，转义操作在 preserveEscapedHtmlTags 中完成
	return content
}

// isNumericOrMathSymbol 检查字符串是否只包含数字或数学符号
func isNumericOrMathSymbol(s string) bool {
	for _, r := range s {
		// 检查是否为数字或常见数学符号（+, -, *, /, ^, etc.）
		if !unicode.IsDigit(r) && !strings.ContainsRune("+-*/^()[]{}.,", r) {
			return false
		}
	}
	return true
}
