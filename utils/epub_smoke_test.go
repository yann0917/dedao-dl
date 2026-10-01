package utils

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestSvg2EpubRealDataSmoke 用本地 output/debug 下的真实解密页面做端到端冒烟
// （覆盖 cover.xhtml 跳过、全书脚注锚点、大量章节的 nav/ncx 构建）。
// 数据目录不存在时跳过（CI / 其他机器无此数据）。
func TestSvg2EpubRealDataSmoke(t *testing.T) {
	const debugEnid = "donM9vjLM8m6d5YQ7lvJGVz2Xgaqb3pZYYo0nBRyejENkP1KZo9rD4OApxpxkY7K"
	dir := filepath.Join("..", "output", "debug", debugEnid)
	entriesOnDisk, err := os.ReadDir(dir)
	if err != nil {
		t.Skipf("本地无真实调试数据 %s，跳过冒烟: %v", dir, err)
	}

	sigilIDRE := regexp.MustCompile(`id="(sigil[^"]+)"`)
	imageRE := regexp.MustCompile(`(?i)<image[^>]*>`)

	var svgContents SvgContents
	var toc []EbookToc
	names := make([]string, 0)
	for i, e := range entriesOnDisk {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".xhtml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("读取 %s 失败: %v", e.Name(), err)
		}
		// 冒烟离线进行：去掉 <image>，避免下载 CDN 图片
		content := imageRE.ReplaceAllString(string(data), "")
		svgContents = append(svgContents, &SvgContent{
			Contents:   []string{content},
			ChapterID:  e.Name(),
			OrderIndex: i,
		})
		names = append(names, e.Name())

		seen := make(map[string]bool)
		first := true
		for _, m := range sigilIDRE.FindAllStringSubmatch(string(data), -1) {
			id := m[1]
			if seen[id] {
				continue
			}
			seen[id] = true
			level := 1
			if first {
				level = 0
				first = false
			}
			toc = append(toc, EbookToc{
				Href:      e.Name() + "#" + id,
				Level:     level,
				Text:      fmt.Sprintf("%s · %s", strings.TrimSuffix(e.Name(), ".xhtml"), id),
				PlayOrder: len(toc),
			})
		}
	}
	if len(svgContents) < 5 {
		t.Fatalf("调试数据异常，仅 %d 章", len(svgContents))
	}

	tmp := t.TempDir()
	oldWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	opt := EpubOptions{Title: "冒烟测试书", Author: "作者", Toc: toc}
	if err = Svg2Epub("冒烟测试书", svgContents, opt); err != nil {
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

	if r.File[0].Name != "mimetype" || r.File[0].Method != zip.Store {
		t.Errorf("mimetype 必须是 zip 首个未压缩条目")
	}

	files := make(map[string]string)
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("打开 %s 失败: %v", f.Name, err)
		}
		data, _ := io.ReadAll(rc)
		rc.Close()
		files[f.Name] = string(data)
	}

	// 封面不进 spine，其余章节都应打包
	wantChapters := len(names) - 1
	gotChapters := 0
	for name := range files {
		if strings.HasPrefix(name, "OEBPS/") &&
			(strings.HasSuffix(name, ".xhtml") && name != "OEBPS/nav.xhtml") {
			gotChapters++
		}
	}
	if gotChapters != wantChapters {
		t.Errorf("期望 %d 个章节文件（cover.xhtml 跳过），实际 %d", wantChapters, gotChapters)
	}

	// nav 应包含每一个真实 sigil 锚点链接
	nav := files["OEBPS/nav.xhtml"]
	missing := 0
	for _, tE := range toc {
		frag := strings.SplitN(tE.Href, "#", 2)
		if len(frag) == 2 && !strings.Contains(nav, `href="`+tE.Href+`"`) {
			missing++
		}
	}
	if missing > 0 {
		t.Errorf("nav.xhtml 缺少 %d 个目录锚点链接", missing)
	}

	// 顶层还应有无目录覆盖的章节的兜底条目（如 Copyright.xhtml）
	if !strings.Contains(nav, `href="Copyright.xhtml"`) {
		t.Errorf("nav.xhtml 缺少未覆盖章节的兜底条目（Copyright.xhtml）")
	}

	// nav 与 ncx 必须是严格良构的 XML；所有章节文档必须经过模板包装
	for name, content := range files {
		if strings.HasPrefix(name, "OEBPS/") && strings.HasSuffix(name, ".xhtml") &&
			!strings.HasPrefix(content, "<?xml") {
			t.Errorf("%s 缺少 XML 声明（模板未应用）", name)
		}
	}
	for _, name := range []string{"OEBPS/nav.xhtml", "OEBPS/toc.ncx", "OEBPS/content.opf"} {
		dec := xml.NewDecoder(strings.NewReader(files[name]))
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

	// ncx navPoint 数应与 nav 目录（toc nav，不含 landmarks）条目数一致
	navToc := strings.SplitN(nav, `<nav epub:type="landmarks"`, 2)[0]
	navEntryCount := strings.Count(navToc, "<li>")
	ncxPointCount := strings.Count(files["OEBPS/toc.ncx"], "<navPoint ")
	if navEntryCount != ncxPointCount {
		t.Errorf("nav 条目数 %d 与 ncx navPoint 数 %d 不一致", navEntryCount, ncxPointCount)
	}
	t.Logf("冒烟完成：%d 章，%d 个目录条目", len(names), navEntryCount)
}
