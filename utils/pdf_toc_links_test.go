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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// TestPdfTocLinkDestOffsets 回归：封面与正文合并后，目录链接的命名目的地
// （/Dests 里的 0 基整数页号）必须整体偏移封面页数，否则点击会落早一页。
func TestPdfTocLinkDestOffsets(t *testing.T) {
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		t.Skipf("本机无 wkhtmltopdf，跳过: %v", err)
	}
	const debugEnid = "donM9vjLM8m6d5YQ7lvJGVz2Xgaqb3pZYYo0nBRyejENkP1KZo9rD4OApxpxkY7K"
	dir := filepath.Join("..", "output", "debug", debugEnid)
	entriesOnDisk, err := os.ReadDir(dir)
	if err != nil {
		t.Skipf("本地无真实调试数据: %v", err)
	}

	sigilIDRE := regexp.MustCompile(`id="(sigil[^"]+)"`)
	imageRE := regexp.MustCompile(`(?i)<image[^>]*>`)
	build := func(withCover bool, outDir string) string {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			t.Fatal(err)
		}
		var svgContents SvgContents
		var toc []EbookToc
		for i, e := range entriesOnDisk {
			if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".xhtml") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			content := imageRE.ReplaceAllString(string(data), "")
			if i == 0 && withCover {
				pngPath := filepath.Join(outDir, "cover.png")
				pngBytes, _ := base64.StdEncoding.DecodeString(
					"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==")
				if err := os.WriteFile(pngPath, pngBytes, 0644); err != nil {
					t.Fatal(err)
				}
				img := `<image x="0" y="57114" width="60000" height="85752" href="` + pngPath + `"/>`
				content = strings.Replace(content, "</svg>", img+"\n</svg>", 1)
			}
			svgContents = append(svgContents, &SvgContent{Contents: []string{content}, ChapterID: e.Name(), OrderIndex: i})
			seen := map[string]bool{}
			for _, m := range sigilIDRE.FindAllStringSubmatch(string(data), -1) {
				if seen[m[1]] {
					continue
				}
				seen[m[1]] = true
				toc = append(toc, EbookToc{Href: e.Name() + "#" + m[1], Level: 1, Text: m[1], PlayOrder: len(toc)})
			}
		}
		if err := os.MkdirAll(outDir, 0755); err != nil {
			t.Fatal(err)
		}
		oldWd, _ := os.Getwd()
		if err := os.Chdir(outDir); err != nil {
			t.Fatal(err)
		}
		err = Svg2Pdf("回归", svgContents, toc)
		_ = os.Chdir(oldWd)
		if err != nil {
			t.Fatal(err)
		}
		matches, _ := filepath.Glob(filepath.Join(outDir, "output", "Ebook", "*.pdf"))
		if len(matches) != 1 {
			t.Fatalf("pdf 数量 %d", len(matches))
		}
		return matches[0]
	}

	tmp := t.TempDir()
	// 无封面单趟产物：目的地为基准值（0 基页号）
	mainOnly := build(false, filepath.Join(tmp, "main"))
	// 带封面合并产物：目的地应 = 基准 + 封面页数
	merged := build(true, filepath.Join(tmp, "merged"))

	const coverPageCount = 1 // 单页封面
	mainDests := namedDestPageInts(t, mainOnly)
	mergedDests := namedDestPageInts(t, merged)
	if len(mainDests) == 0 || len(mergedDests) == 0 {
		t.Fatalf("目的地为空: main=%d merged=%d", len(mainDests), len(mergedDests))
	}

	bad := 0
	for anchor, base := range mainDests {
		got, ok := mergedDests[anchor]
		if !ok {
			t.Errorf("合并产物缺少目的地 %s", anchor)
			continue
		}
		if got != base+coverPageCount {
			bad++
			if bad <= 3 {
				t.Errorf("锚点 %s 目的地未正确偏移: 期望 %d，实际 %d", anchor, base+coverPageCount, got)
			}
		}
	}
	if bad > 0 {
		t.Fatalf("共 %d 个目的地偏移错误", bad)
	}
	t.Logf("回归通过：%d 个命名目的地全部偏移 +%d 页", len(mainDests), coverPageCount)
}

// namedDestPageInts 解析 PDF 的 catalog /Dests，返回「锚点名 -> 目的地整数页号」。
func namedDestPageInts(t *testing.T, pdfPath string) map[string]int {
	cc, err := api.ReadContextFile(context.Background(), pdfPath)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]int{}
	destsIndRef, found := cc.RootDict["Dests"]
	if !found {
		return out
	}
	destsObj, err := cc.Dereference(destsIndRef)
	if err != nil {
		t.Fatal(err)
	}
	destsDict, ok := destsObj.(types.Dict)
	if !ok {
		return out
	}
	for name, v := range destsDict {
		obj, err := cc.Dereference(v)
		if err != nil {
			continue
		}
		arr, ok := obj.(types.Array)
		if !ok || len(arr) == 0 {
			continue
		}
		i, ok := arr[0].(types.Integer)
		if !ok {
			continue
		}
		// 名称形如 file:///...html#sigil_toc_id_1（#xx 转义已被 pdfcpu 解码），
		// 只取 .html# 之后的片段为锚点键（其余是含随机字节的内部名字）
		if idx := strings.LastIndex(name, ".html#"); idx >= 0 && idx+6 < len(name) {
			out[name[idx+6:]] = i.Value()
		}
	}
	return out
}
