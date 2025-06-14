package app

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// EbookDetail 电子书详情
func EbookDetail(id int) (detail *services.EbookDetail, err error) {
	courseDetail, err := CourseDetail(CateEbook, id)
	if err != nil {
		return
	}
	enID := courseDetail.Enid
	detail, err = getService().EbookDetail(enID)

	return
}

func EbookDetailByEnID(enID string) (detail *services.EbookDetail, err error) {
	detail, err = getService().EbookDetail(enID)

	return
}

func EbookInfo(enID string) (info *services.EbookInfo, err error) {
	token, err1 := getService().EbookReadToken(enID)
	if err1 != nil {
		err = err1
		return
	}

	info, err = getService().EbookInfo(token.Token)
	return
}

func EbookPage(enID string) (info *services.EbookInfo, svgContent utils.SvgContents, err error) {
	token, err1 := getService().EbookReadToken(enID)
	if err1 != nil {
		err = err1
		return
	}

	info, err = getService().EbookInfo(token.Token)
	if err != nil {
		return
	}

	// 使用 channel 来收集结果和错误
	type result struct {
		content *utils.SvgContent
		err     error
	}

	resultChan := make(chan result, len(info.BookInfo.Orders))
	var mu sync.Mutex
	var svgContents utils.SvgContents

	wgp := utils.NewWaitGroupPool(5)
	for i, order := range info.BookInfo.Orders {
		wgp.Add()
		go func(i int, order services.EbookOrders) {
			defer wgp.Done()

			index, count, offset := 0, 20, 0
			svgList, err1 := generateEbookPages(enID, order.ChapterID, token.Token, index, count, offset)

			if err1 != nil {
				resultChan <- result{nil, fmt.Errorf("章节 %s 处理失败: %v", order.ChapterID, err1)}
				return
			}

			content := &utils.SvgContent{
				Contents:   svgList,
				ChapterID:  order.ChapterID,
				OrderIndex: i,
			}

			resultChan <- result{content, nil}
		}(i, order)
	}

	// 等待所有 goroutine 完成
	go func() {
		wgp.Wait()
		close(resultChan)
	}()

	// 收集结果
	var firstError error
	for res := range resultChan {
		if res.err != nil {
			if firstError == nil {
				firstError = res.err
			}
			fmt.Printf("错误: %v\n", res.err)
			continue
		}

		if res.content != nil {
			mu.Lock()
			svgContents = append(svgContents, res.content)
			mu.Unlock()
		}
	}

	// 如果有错误且没有成功的内容，返回错误
	if firstError != nil && len(svgContents) == 0 {
		err = firstError
		return
	}

	// 如果有部分成功，继续处理，但打印警告
	if firstError != nil {
		fmt.Printf("警告: 部分章节处理失败，但继续处理成功的章节\n")
	}

	svgContent = svgContents
	return
}

func generateEbookPages(enid, chapterID, token string, index, count, offset int) (svgList []string, err error) {
	// Try to load from cache first
	if cachedPages, found := services.LoadFromCache(enid, chapterID); found {
		fmt.Printf("使用缓存内容：%s\n", chapterID)
		return cachedPages, nil
	}

	fmt.Printf("下载章节 %s\n", chapterID)
	pageList, err := getService().EbookPages(chapterID, token, index, count, offset)
	if err != nil {
		return nil, fmt.Errorf("获取章节页面失败: %v", err)
	}

	// 添加 nil 检查
	if pageList == nil {
		return nil, fmt.Errorf("章节 %s 返回的页面数据为空", chapterID)
	}

	if pageList.Pages == nil {
		return nil, fmt.Errorf("章节 %s 的页面列表为空", chapterID)
	}

	if len(pageList.Pages) == 0 {
		fmt.Printf("警告: 章节 %s 没有页面内容\n", chapterID)
		return []string{}, nil
	}

	for _, item := range pageList.Pages {
		desContents := DecryptAES(item.Svg)
		svgList = append(svgList, desContents)
	}

	if !pageList.IsEnd {
		index += count
		count = 20
		fmt.Printf("下载章节 %s 的更多页面 (索引: %d)\n", chapterID, index)
		list, err1 := generateEbookPages(enid, chapterID, token, index, count, offset)
		if err1 != nil {
			return nil, fmt.Errorf("获取更多页面失败: %v", err1)
		}

		svgList = append(svgList, list...)
	} else {
		fmt.Printf("章节 %s 下载完成 (共 %d 页)\n", chapterID, len(svgList))
	}

	// Save to cache
	if err := services.SaveToCache(enid, chapterID, svgList); err != nil {
		fmt.Printf("警告: 无法缓存章节 %s: %v\n", chapterID, err)
	} else {
		fmt.Printf("已缓存章节 %s\n", chapterID)
	}

	return
}

// PKCS7Unpad 实现PKCS7去填充
func PKCS7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// DecryptAES 实现AES - CBC解密
func DecryptAES(contents string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(contents)
	if err != nil {
		fmt.Println("Base64解码错误:", err)
		return ""
	}

	key := []byte("3e4r06tjkpjcevlbslr3d96gdb5ahbmo")
	iv := []byte("6fd89a1b3a7f48fb")

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	blockSize := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = PKCS7Unpad(plaintext)
	return string(plaintext)
}
