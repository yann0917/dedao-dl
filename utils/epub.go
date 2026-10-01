package utils

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// epubDoc 手写 EPUB 打包器：不依赖 go-epub（已停止维护），直接按 EPUB 3 规范组装 zip。
// 相比 go-epub 的关键差异：目录（nav.xhtml / toc.ncx）由 bookInfo.Toc 全量构建，
// 支持任意层级嵌套，且小节条目可以链接到章节文件内部的锚点（Section001.xhtml#sigil_toc_id_1）。
type epubDoc struct {
	title       string
	author      string
	description string
	language    string
	uid         string
	css         string

	chapters   []*epubChapter
	assets     []*epubAsset
	cover      *epubAsset
	chapterSet map[string]int // 章节 href -> chapters 下标，用于目录链接解析

	nav      []*epubTocNode
	navDepth int
}

type epubChapter struct {
	id      string
	href    string // 相对 OEBPS/ 的文件名（保留 ChapterID 原名，目录与脚注链接按原名解析）
	content string // 完整 XHTML 文档
}

type epubAsset struct {
	id         string
	href       string // 相对 OEBPS/
	mediaType  string
	properties string // "cover-image" 等
	data       []byte
}

// epubTocNode 目录树节点；href 为空表示纯文本条目（目标文件不在书中）。
type epubTocNode struct {
	text     string
	href     string
	children []*epubTocNode
}

func newEpubDoc(title, author, description, css string) *epubDoc {
	return &epubDoc{
		title:       title,
		author:      author,
		description: description,
		language:    "zh-CN",
		uid:         epubUUID(title),
		css:         css,
		chapterSet:  make(map[string]int),
	}
}

// epubUUID 由书名派生确定性 UUID（v5 风格），重复导出同一本书时产物稳定。
func epubUUID(seed string) string {
	h := sha1.Sum([]byte(seed))
	b := h[:]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// addChapter 注册一个章节，返回实际使用的文件名。ChapterID 原样保留为文件名，
// 因为 TOC 的 href 前缀与脚注的跨章链接都以这个名字引用章节。
func (d *epubDoc) addChapter(rawName, content string) string {
	name := epubChapterFileName(rawName, len(d.chapters))
	if _, dup := d.chapterSet[name]; dup {
		name = fmt.Sprintf("%s_%d", name, len(d.chapters)+1)
	}
	d.chapters = append(d.chapters, &epubChapter{
		id:      fmt.Sprintf("chapter%03d", len(d.chapters)+1),
		href:    name,
		content: content,
	})
	d.chapterSet[name] = len(d.chapters) - 1
	return name
}

func epubChapterFileName(rawName string, idx int) string {
	name := path.Base(strings.ReplaceAll(strings.TrimSpace(rawName), "\\", "/"))
	if name == "" || name == "." || name == "/" {
		name = fmt.Sprintf("section_%03d.xhtml", idx+1)
	}
	return name
}

// addImage 注册一张图片，返回章节内引用用的相对 href（images/xxx）。
func (d *epubDoc) addImage(data []byte, ext string) string {
	if ext == "" || !strings.HasPrefix(ext, ".") {
		ext = "." + strings.TrimPrefix(ext, ".")
	}
	name := fmt.Sprintf("image_%03d%s", len(d.assets), ext)
	asset := &epubAsset{
		id:        fmt.Sprintf("image%03d", len(d.assets)),
		href:      path.Join("images", name),
		mediaType: epubImageMediaType(ext, data),
		data:      data,
	}
	d.assets = append(d.assets, asset)
	return asset.href
}

// setCover 注册封面图；data 为空时跳过（保持无封面的兼容行为）。
func (d *epubDoc) setCover(data []byte, fallbackExt string) error {
	if len(data) == 0 {
		return nil
	}
	ext, _ := detectImageExtAndType(data, fallbackExt)
	d.cover = &epubAsset{
		id:         "cover-image",
		href:       path.Join("images", "cover"+ext),
		mediaType:  epubImageMediaType(ext, data),
		properties: "cover-image",
		data:       data,
	}
	return nil
}

// detectImageExtAndType 通过内容嗅探图片扩展名，嗅探失败时回退到原始扩展名。
func detectImageExtAndType(data []byte, fallbackExt string) (ext, mediaType string) {
	if fallbackExt != "" && !strings.HasPrefix(fallbackExt, ".") {
		fallbackExt = "." + fallbackExt
	}
	m := http.DetectContentType(data)
	if strings.HasPrefix(m, "image/") {
		switch {
		case strings.Contains(m, "png"):
			return ".png", m
		case strings.Contains(m, "jpeg"):
			return ".jpg", m
		case strings.Contains(m, "gif"):
			return ".gif", m
		case strings.Contains(m, "webp"):
			return ".webp", m
		case strings.Contains(m, "svg"):
			return ".svg", m
		case strings.Contains(m, "bmp"):
			return ".bmp", m
		}
		return fallbackExt, m
	}
	return fallbackExt, ""
}

func epubImageMediaType(ext string, data []byte) string {
	if _, m := detectImageExtAndType(data, ext); m != "" {
		return m
	}
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".bmp":
		return "image/bmp"
	}
	return "application/octet-stream"
}

// buildNavTree 由 bookInfo.Toc 构建嵌套目录树。
// 层级取 toc.Level（0 为最顶层），层级跳档时收敛到上一层，避免悬空节点；
// 目录未覆盖到的章节文件在顶层补空名条目，保证每个 spine 文件都能从目录到达。
func (d *epubDoc) buildNavTree(toc []EbookToc) {
	var stack []*epubTocNode
	place := func(depth int, node *epubTocNode) {
		if depth < 0 {
			depth = 0
		}
		if depth > len(stack) {
			depth = len(stack)
		}
		if depth == 0 {
			d.nav = append(d.nav, node)
		} else {
			parent := stack[depth-1]
			parent.children = append(parent.children, node)
		}
		stack = append(stack[:depth], node)
		if d.navDepth < depth+1 {
			d.navDepth = depth + 1
		}
	}

	referenced := make(map[string]bool)
	for _, t := range toc {
		node := &epubTocNode{text: t.Text, href: d.resolveTocHref(t.Href)}
		place(t.Level, node)
		if node.href != "" {
			referenced[strings.SplitN(node.href, "#", 2)[0]] = true
		}
	}
	for _, ch := range d.chapters {
		if !referenced[ch.href] {
			place(0, &epubTocNode{href: ch.href})
		}
	}
}

// resolveTocHref 把 TOC 的 href（ChapterID#fragment）解析成书中实际存在的文件链接。
// 解析不到目标文件时返回空串，调用方渲染为无链接的纯文本。
func (d *epubDoc) resolveTocHref(href string) string {
	if href == "" {
		return ""
	}
	file, frag := href, ""
	if i := strings.Index(href, "#"); i >= 0 {
		file, frag = href[:i], href[i+1:]
	}
	target := ""
	for _, candidate := range []string{file, path.Base(file), file + ".xhtml"} {
		if candidate == "" {
			continue
		}
		if _, ok := d.chapterSet[candidate]; ok {
			target = candidate
			break
		}
	}
	if target == "" {
		return ""
	}
	if frag != "" {
		return target + "#" + frag
	}
	return target
}

// write 组装并写出 EPUB 文件。mimetype 必须是 zip 第一个条目且不压缩（规范要求）。
func (d *epubDoc) write(outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	writeEntry := func(name string, method uint16, data []byte) error {
		h := &zip.FileHeader{Name: name, Method: method, Modified: time.Now()}
		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}

	if err = writeEntry("mimetype", zip.Store, []byte("application/epub+zip")); err != nil {
		return err
	}
	if err = writeEntry("META-INF/container.xml", zip.Deflate, []byte(d.containerXML())); err != nil {
		return err
	}
	if err = writeEntry("OEBPS/content.opf", zip.Deflate, []byte(d.contentOPF())); err != nil {
		return err
	}
	if err = writeEntry("OEBPS/nav.xhtml", zip.Deflate, []byte(d.navXHTML())); err != nil {
		return err
	}
	if err = writeEntry("OEBPS/toc.ncx", zip.Deflate, []byte(d.tocNCX())); err != nil {
		return err
	}
	if d.css != "" {
		if err = writeEntry("OEBPS/style.css", zip.Deflate, []byte(d.css)); err != nil {
			return err
		}
	}
	for _, a := range d.assets {
		if err = writeEntry(path.Join("OEBPS", a.href), zip.Deflate, a.data); err != nil {
			return err
		}
	}
	if d.cover != nil {
		if err = writeEntry(path.Join("OEBPS", d.cover.href), zip.Deflate, d.cover.data); err != nil {
			return err
		}
	}
	for _, ch := range d.chapters {
		if err = writeEntry(path.Join("OEBPS", ch.href), zip.Deflate, []byte(ch.content)); err != nil {
			return err
		}
	}
	return zw.Close()
}

func (d *epubDoc) containerXML() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>
`
}

func (d *epubDoc) contentOPF() string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="pub-id" xml:lang="` + d.language + `">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="pub-id">urn:uuid:` + d.uid + `</dc:identifier>
    <dc:title>` + xmlEscapeText(d.title) + `</dc:title>`)
	if d.author != "" {
		sb.WriteString("\n    <dc:creator>" + xmlEscapeText(d.author) + "</dc:creator>")
	}
	if d.description != "" {
		sb.WriteString("\n    <dc:description>" + xmlEscapeText(d.description) + "</dc:description>")
	}
	sb.WriteString("\n    <dc:language>" + d.language + "</dc:language>")
	sb.WriteString("\n    <meta property=\"dcterms:modified\">" + time.Now().UTC().Format("2006-01-02T15:04:05Z") + "</meta>")
	if d.cover != nil {
		sb.WriteString("\n    <meta name=\"cover\" content=\"" + d.cover.id + "\"/>")
	}

	sb.WriteString("\n  </metadata>\n  <manifest>")
	sb.WriteString("\n    <item id=\"nav\" href=\"nav.xhtml\" media-type=\"application/xhtml+xml\" properties=\"nav\"/>")
	sb.WriteString("\n    <item id=\"ncx\" href=\"toc.ncx\" media-type=\"application/x-dtbncx\"/>")
	if d.css != "" {
		sb.WriteString("\n    <item id=\"css\" href=\"style.css\" media-type=\"text/css\"/>")
	}
	if d.cover != nil {
		sb.WriteString("\n    <item id=\"" + d.cover.id + "\" href=\"" + d.cover.href + "\" media-type=\"" + d.cover.mediaType + "\" properties=\"cover-image\"/>")
	}
	for _, a := range d.assets {
		sb.WriteString("\n    <item id=\"" + a.id + "\" href=\"" + a.href + "\" media-type=\"" + a.mediaType + "\"/>")
	}
	for _, ch := range d.chapters {
		sb.WriteString("\n    <item id=\"" + ch.id + "\" href=\"" + ch.href + "\" media-type=\"application/xhtml+xml\"/>")
	}

	sb.WriteString("\n  </manifest>\n  <spine toc=\"ncx\">")
	for _, ch := range d.chapters {
		sb.WriteString("\n    <itemref idref=\"" + ch.id + "\"/>")
	}
	sb.WriteString("\n  </spine>\n</package>\n")
	return sb.String()
}

func (d *epubDoc) navXHTML() string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="` + d.language + `" xml:lang="` + d.language + `">
<head>
  <meta charset="utf-8"/>
  <title>目录</title>
  <link rel="stylesheet" type="text/css" href="style.css"/>
</head>
<body>
  <nav epub:type="toc" id="toc">
    <h1>目录</h1>
`)
	d.writeNavList(&sb, d.nav, "    ")
	sb.WriteString(`  </nav>
  <nav epub:type="landmarks" id="landmarks" hidden="hidden">
    <h2>Guide</h2>
    <ol>
      <li><a epub:type="toc" href="#toc">目录</a></li>`)
	if len(d.chapters) > 0 {
		sb.WriteString("\n      <li><a epub:type=\"bodymatter\" href=\"" + xmlEscapeAttr(d.chapters[0].href) + "\">正文</a></li>")
	}
	sb.WriteString(`
    </ol>
  </nav>
</body>
</html>
`)
	return sb.String()
}

func (d *epubDoc) writeNavList(sb *strings.Builder, nodes []*epubTocNode, indent string) {
	if len(nodes) == 0 {
		return
	}
	sb.WriteString(indent + "<ol>\n")
	for _, n := range nodes {
		sb.WriteString(indent + "\t<li>")
		if n.href != "" {
			sb.WriteString(`<a href="` + xmlEscapeAttr(n.href) + `">` + xmlEscapeText(n.text) + `</a>`)
		} else {
			sb.WriteString(xmlEscapeText(n.text))
		}
		if len(n.children) > 0 {
			sb.WriteString("\n")
			d.writeNavList(sb, n.children, indent+"\t\t")
			sb.WriteString(indent + "\t")
		}
		sb.WriteString("</li>\n")
	}
	sb.WriteString(indent + "</ol>\n")
}

func (d *epubDoc) tocNCX() string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="` + d.language + `">
  <head>
    <meta name="dtb:uid" content="urn:uuid:` + d.uid + `"/>
    <meta name="dtb:depth" content="` + fmt.Sprintf("%d", d.navDepth) + `"/>
    <meta name="dtb:totalPageCount" content="0"/>
    <meta name="dtb:maxPageNumber" content="0"/>
  </head>
  <docTitle><text>` + xmlEscapeText(d.title) + `</text></docTitle>`)
	if d.author != "" {
		sb.WriteString("\n  <docAuthor><text>" + xmlEscapeText(d.author) + "</text></docAuthor>")
	}
	sb.WriteString("\n  <navMap>")
	playOrder := 0
	var writePoints func(nodes []*epubTocNode, indent string)
	writePoints = func(nodes []*epubTocNode, indent string) {
		for _, n := range nodes {
			playOrder++
			sb.WriteString("\n" + indent + `<navPoint id="navpoint-` + fmt.Sprintf("%d", playOrder) + `" playOrder="` + fmt.Sprintf("%d", playOrder) + `">`)
			sb.WriteString("\n" + indent + "\t<navLabel><text>" + xmlEscapeText(n.text) + "</text></navLabel>")
			sb.WriteString("\n" + indent + "\t<content src=\"" + xmlEscapeAttr(n.href) + "\"/>")
			writePoints(n.children, indent+"\t\t")
			sb.WriteString("\n" + indent + "</navPoint>")
		}
	}
	writePoints(d.nav, "    ")
	sb.WriteString("\n  </navMap>\n</ncx>\n")
	return sb.String()
}

// epubXHTMLDoc 用章节正文（body 内部 HTML）组装完整 XHTML 文档。
// 内层包一层容器 div：goquery 序列化的 body 内部 HTML 可能有多个顶层兄弟节点，
// 而 XHTML 要求单一根元素，否则部分阅读器/校验器会判定文档损坏。
func epubXHTMLDoc(title, bodyInner string) string {
	return `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" lang="zh-CN" xml:lang="zh-CN">
<head>
  <meta charset="utf-8"/>
  <title>` + xmlEscapeText(title) + `</title>
  <link rel="stylesheet" type="text/css" href="style.css"/>
</head>
<body>
<div class="chapter-body">
` + bodyInner + `
</div>
</body>
</html>
`
}

func xmlEscapeText(s string) string {
	return xmlEscape(s, true)
}

func xmlEscapeAttr(s string) string {
	return xmlEscape(s, true)
}

func xmlEscape(s string, _ bool) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
