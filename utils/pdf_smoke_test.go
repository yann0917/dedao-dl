package utils

import (
	"context"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// TestSvg2PdfRealDataSmoke 用本地真实解密页面生成 PDF（与线上同一代码路径），
// 用于人工检查排版（页顶截断、行距等）。本机无 wkhtmltopdf 时跳过。
func TestSvg2PdfRealDataSmoke(t *testing.T) {
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		t.Skipf("本机无 wkhtmltopdf，跳过 PDF 冒烟: %v", err)
	}
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
	for i, e := range entriesOnDisk {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".xhtml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("读取 %s 失败: %v", e.Name(), err)
		}
		content := imageRE.ReplaceAllString(string(data), "")
		if i == 0 {
			// 第一章注入一张本地图片，强制触发封面渲染与合并路径
			pngPath := filepath.Join(t.TempDir(), "cover.png")
			pngBytes, _ := base64.StdEncoding.DecodeString(
				"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==")
			if err := os.WriteFile(pngPath, pngBytes, 0644); err != nil {
				t.Fatal(err)
			}
			img := `<image x="0" y="57114" width="60000" height="85752" href="` + pngPath + `"/>`
			content = strings.Replace(content, "</svg>", img+"\n</svg>", 1)
		}
		svgContents = append(svgContents, &SvgContent{
			Contents:   []string{content},
			ChapterID:  e.Name(),
			OrderIndex: i,
		})
		seen := make(map[string]bool)
		for _, m := range sigilIDRE.FindAllStringSubmatch(string(data), -1) {
			id := m[1]
			if seen[id] {
				continue
			}
			seen[id] = true
			toc = append(toc, EbookToc{
				Href:      e.Name() + "#" + id,
				Level:     1,
				Text:      strings.TrimSuffix(e.Name(), ".xhtml") + " " + id,
				PlayOrder: len(toc),
			})
		}
	}

	if err = Svg2Pdf("冒烟测试书", svgContents, toc); err != nil {
		t.Fatalf("Svg2Pdf 失败: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join("output", "Ebook", "*.pdf"))
	if len(matches) != 1 {
		t.Fatalf("期望生成 1 个 pdf，实际 %d", len(matches))
	}
	st, err := os.Stat(matches[0])
	if err != nil || st.Size() < 1024 {
		t.Fatalf("pdf 产物异常: %v (%v)", matches[0], err)
	}

	// 大纲必须是正文自身的单一目录：顶层不允许出现 pdfcpu 合并产生的
	// 书架条目（dedao-cover / dedao-main）
	pdfFile, err := os.Open(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	defer pdfFile.Close()
	bms, err := api.Bookmarks(context.Background(), pdfFile, nil)
	if err != nil {
		t.Fatalf("读取大纲失败: %v", err)
	}
	if len(bms) == 0 {
		t.Fatalf("大纲为空")
	}
	for _, bm := range bms {
		if strings.HasPrefix(bm.Title, "dedao-") {
			t.Errorf("大纲顶层出现合并器书架条目: %q", bm.Title)
		}
	}
	first := bms[0]
	t.Logf("PDF 已生成: %s (%d bytes)，大纲顶层 %d 条，首条 %q@p%d",
		matches[0], st.Size(), len(bms), first.Title, first.PageFrom)
}
