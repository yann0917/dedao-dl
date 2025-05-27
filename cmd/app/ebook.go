package app

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

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
	// fmt.Printf("%#v\n", info.BookInfo.Pages)
	// fmt.Printf("%#v\n", info.BookInfo.EbookBlock)
	// fmt.Printf("%#v\n", info.BookInfo.Toc)
	// fmt.Printf("%#v\n", info.BookInfo.Orders)
	wgp := utils.NewWaitGroupPool(5)
	for i, order := range info.BookInfo.Orders {
		wgp.Add()
		go func(i int, order services.EbookOrders) {
			defer func() {
				wgp.Done()
			}()
			index, count, offset := 0, 20, 0
			svgList, err1 := generateEbookPages(enID, order.ChapterID, token.Token, index, count, offset)
			if err1 != nil {
				err = err1
				return
			}

			svgContent = append(svgContent, &utils.SvgContent{
				Contents:   svgList,
				ChapterID:  order.ChapterID,
				OrderIndex: i,
			})
		}(i, order)
	}
	wgp.Wait()
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
		return
	}

	// // 创建用于保存原始SVG内容的目录
	// debugDir, err1 := utils.Mkdir(utils.OutputDir, "Debug", "SVG", enid)
	// if err1 != nil {
	// 	fmt.Printf("创建调试目录失败: %v\n", err1)
	// 	// 继续执行，不要因为调试目录创建失败而中断主流程
	// }

	for _, item := range pageList.Pages {
		desContents := DecryptAES(item.Svg)
		svgList = append(svgList, desContents)

		// // 保存原始SVG内容到文件，用于调试
		// if debugDir != "" {
		// 	debugFileName := fmt.Sprintf("%s_%d_%d.svg", chapterID, index, i)
		// 	debugFilePath, err3 := utils.FilePath(filepath.Join(debugDir, debugFileName), "", false)
		// 	if err3 != nil {
		// 		fmt.Printf("创建调试文件路径失败: %v\n", err3)
		// 		continue
		// 	}
		// 	// 写入文件
		// 	if err2 := utils.WriteFileWithTrunc(debugFilePath, desContents); err2 != nil {
		// 		fmt.Printf("保存调试SVG文件失败: %v\n", err2)
		// 	} else {
		// 		fmt.Printf("已保存调试SVG文件: %s\n", debugFilePath)
		// 	}
		// }
	}

	if !pageList.IsEnd {
		index += count
		count = 20
		fmt.Printf("下载章节 %s 的更多页面 (索引: %d)\n", chapterID, index)
		list, err1 := generateEbookPages(enid, chapterID, token, index, count, offset)
		if err1 != nil {
			err = err1
			return
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
