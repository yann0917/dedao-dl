package utils

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixtureSvgPage 构造一页最小可解析的 SVG（含若干文本行），模拟 get_pages 解密后的页面。
func fixtureSvgPage(lines ...string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg" version="1.0" width="60000" height="200000">
` + strings.Join(lines, "\n") + `
</svg>`
}

func fixtureTextLine(id, offset, y, text, fontSize string) string {
	idAttr := ""
	if id != "" {
		idAttr = ` id="` + id + `"`
	}
	return `<text x="20000"` + idAttr + ` y="` + y + `" width="300.000000" top="` + y + `" height="30.000000" style="font-size:` +
		fontSize + `;fill:rgb(0, 0, 0);font-family:'PingFang SC';" offset="` + offset + `" len="4" newline="true">` + text + `</text>`
}

// TestSvg2EpubPackaging 离线验证手写 EPUB 打包：zip 结构、嵌套目录、锚点注入。
func TestSvg2EpubPackaging(t *testing.T) {
	tmp := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	page1 := fixtureSvgPage(
		fixtureTextLine("sigil_toc_id_0", "0", "100", "第一章 测试", "22px"),
		fixtureTextLine("sigil_toc_id_1", "10", "200", "第一节 测试小节", "20px"),
		fixtureTextLine("", "20", "300", "正文段落内容。", "16px"),
	)
	page2 := fixtureSvgPage(
		fixtureTextLine("sigil_toc_id_2", "400", "100", "第二章 测试", "22px"),
		fixtureTextLine("", "500", "200", "第二节 深入讨论", "20px"),
		fixtureTextLine("", "600", "300", "第二章正文。", "16px"),
	)
	// 本地图片文件：走 changeRef 的非 http 分支，离线验证图片内嵌与 src 重写
	// （index==0 的封面章会过滤图片，所以放在第二章）
	pngPath := filepath.Join(tmp, "fixture.png")
	pngBytes, _ := base64.StdEncoding.DecodeString(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==")
	if err := os.WriteFile(pngPath, pngBytes, 0644); err != nil {
		t.Fatal(err)
	}
	page2 = strings.Replace(page2, "</svg>",
		`<image x="20000" y="400" width="400.000000" height="100.000000" href="`+pngPath+`"/>
</svg>`, 1)
	svgContents := SvgContents{
		&SvgContent{Contents: []string{page1}, ChapterID: "Section001.xhtml", OrderIndex: 0},
		&SvgContent{Contents: []string{page2}, ChapterID: "Section002.xhtml", OrderIndex: 1},
	}
	toc := []EbookToc{
		{Href: "Section001.xhtml#sigil_toc_id_0", Level: 0, Offset: 0, Text: "第一章 测试", PlayOrder: 0},
		{Href: "Section001.xhtml#sigil_toc_id_1", Level: 1, Offset: 10, Text: "第一节 测试小节", PlayOrder: 1},
		{Href: "Section002.xhtml#sigil_toc_id_2", Level: 0, Offset: 400, Text: "第二章 测试", PlayOrder: 2},
		{Href: "Section002.xhtml#sigil_toc_id_3", Level: 1, Offset: 500, Text: "第二节 深入讨论", PlayOrder: 3},
	}

	opt := EpubOptions{Title: "测试书", Author: "测试作者", Description: "简介", Toc: toc}
	if err = Svg2Epub("测试书", svgContents, opt); err != nil {
		t.Fatalf("Svg2Epub 失败: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join("output", "Ebook", "*.epub"))
	if len(matches) != 1 {
		t.Fatalf("期望生成 1 个 epub，实际 %d", len(matches))
	}

	r, err := zip.OpenReader(matches[0])
	if err != nil {
		t.Fatalf("epub 无法打开: %v", err)
	}
	defer r.Close()

	// 1. mimetype 必须是第一个条目、不压缩、内容正确
	first := r.File[0]
	if first.Name != "mimetype" {
		t.Errorf("zip 第一个条目应为 mimetype，实际 %s", first.Name)
	}
	if first.Method != zip.Store {
		t.Errorf("mimetype 应为 Store 不压缩，实际 %d", first.Method)
	}
	fc, _ := first.Open()
	mt, _ := io.ReadAll(fc)
	fc.Close()
	if string(mt) != "application/epub+zip" {
		t.Errorf("mimetype 内容错误: %q", string(mt))
	}

	entries := map[string]string{}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("打开 %s 失败: %v", f.Name, err)
		}
		data, _ := io.ReadAll(rc)
		rc.Close()
		entries[f.Name] = string(data)
	}

	for _, name := range []string{
		"META-INF/container.xml",
		"OEBPS/content.opf",
		"OEBPS/nav.xhtml",
		"OEBPS/toc.ncx",
		"OEBPS/style.css",
		"OEBPS/Section001.xhtml",
		"OEBPS/Section002.xhtml",
	} {
		if _, ok := entries[name]; !ok {
			t.Errorf("缺少条目 %s", name)
		}
	}

	// 2. nav.xhtml：全部 4 个目录条目都在，且小节锚点带 fragment
	nav := entries["OEBPS/nav.xhtml"]
	for _, want := range []string{
		`href="Section001.xhtml#sigil_toc_id_0"`,
		`href="Section001.xhtml#sigil_toc_id_1"`,
		`href="Section002.xhtml#sigil_toc_id_2"`,
		`href="Section002.xhtml#sigil_toc_id_3"`,
	} {
		if !strings.Contains(nav, want) {
			t.Errorf("nav.xhtml 缺少目录链接 %s", want)
		}
	}
	// 嵌套结构：toc nav 内应有 3 个 ol（顶层 + 两个章节各一个子 ol）+ landmarks 1 个
	if got := strings.Count(nav, "<ol>"); got != 4 {
		t.Errorf("nav.xhtml 期望 4 个 <ol>（含 landmarks），实际 %d", got)
	}

	// 3. toc.ncx：嵌套深度与带锚点的 content src
	ncx := entries["OEBPS/toc.ncx"]
	if !strings.Contains(ncx, `<meta name="dtb:depth" content="2"/>`) {
		t.Errorf("toc.ncx 深度应为 2")
	}
	if !strings.Contains(ncx, `src="Section001.xhtml#sigil_toc_id_1"`) {
		t.Errorf("toc.ncx 缺少带锚点的 content src")
	}
	if got := strings.Count(ncx, "<navPoint "); got != 4 {
		t.Errorf("toc.ncx 期望 4 个 navPoint，实际 %d", got)
	}

	// 4. content.opf：manifest/spine 完整
	opf := entries["OEBPS/content.opf"]
	for _, want := range []string{
		`properties="nav"`,
		`href="Section001.xhtml" media-type="application/xhtml+xml"`,
		`href="Section002.xhtml" media-type="application/xhtml+xml"`,
		`<itemref idref="chapter001"/>`,
		`<itemref idref="chapter002"/>`,
		`toc="ncx"`,
	} {
		if !strings.Contains(opf, want) {
			t.Errorf("content.opf 缺少 %s", want)
		}
	}

	// 5. 锚点注入：章节 1 的两条目走 SVG 原生 id 信号（span id）；
	//    章节 2 的第二节无原生 id，走 offset 信号（p id）
	ch1 := entries["OEBPS/Section001.xhtml"]
	if !strings.Contains(ch1, `<span id="sigil_toc_id_0"`) || !strings.Contains(ch1, `<span id="sigil_toc_id_1"`) {
		t.Errorf("Section001.xhtml 缺少 id 主信号注入的锚点")
	}
	ch2 := entries["OEBPS/Section002.xhtml"]
	if !strings.Contains(ch2, `<span id="sigil_toc_id_2"`) {
		t.Errorf("Section002.xhtml 缺少章节标题锚点（id 主信号）")
	}
	if !strings.Contains(ch2, `<p id="sigil_toc_id_3"`) {
		t.Errorf("Section002.xhtml 缺少 offset 信号注入的锚点")
	}

	// 6. 图片内嵌：src 重写为 epub 内部引用，且 zip 中存在该图片
	ch2 = entries["OEBPS/Section002.xhtml"]
	if !strings.Contains(ch2, `src="images/image_000.png"`) {
		t.Errorf("Section002.xhtml 图片 src 未重写为内部引用")
	}
	if _, ok := entries["OEBPS/images/image_000.png"]; !ok {
		t.Errorf("zip 中缺少内嵌图片 image_000.png")
	}

	// 7. 所有 XML 部件可解析（保证 XHTML 良构），且章节文档都经过模板包装
	for _, name := range []string{"META-INF/container.xml", "OEBPS/content.opf", "OEBPS/nav.xhtml", "OEBPS/toc.ncx", "OEBPS/Section001.xhtml", "OEBPS/Section002.xhtml"} {
		if strings.HasPrefix(name, "OEBPS/") && strings.HasSuffix(name, ".xhtml") &&
			!strings.HasPrefix(entries[name], "<?xml") {
			t.Errorf("%s 缺少 XML 声明（模板未应用）", name)
		}
		dec := xml.NewDecoder(strings.NewReader(entries[name]))
		dec.Strict = false
		for {
			_, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("%s 不是良构 XML: %v", name, err)
				break
			}
		}
	}
}
