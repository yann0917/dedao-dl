package app

import (
	"fmt"

	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// EbookDetail 电子书详情
func EbookDetail(id int) (detail *services.EbookDetail, err error) {
	courseDetail, err := getService().CourseDetail(CateEbook, id)
	if err != nil {
		return
	}

	enID := courseDetail.Enid
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

func EbookPage(title, enID string) (pages *services.EbookPage, err error) {
	token, err1 := getService().EbookReadToken(enID)
	if err1 != nil {
		err = err1
		return
	}

	info, err1 := getService().EbookInfo(token.Token)
	if err1 != nil {
		err = err1
		return
	}
	var svgContent []string
	for _, order := range info.BookInfo.Orders {

		index, count, offset := 0, 20, 0
		svgList, err1 := generateEbookPages(order.ChapterID, token.Token, index, count, offset)
		if err1 != nil {
			err = err1
			return
		}
		svgContent = append(svgContent, svgList...)

	}
	if len(svgContent) > 0 {
		utils.Svg2Html(title, svgContent)
	}
	return
}

func generateEbookPages(chapterID, token string, index, count, offset int) (svgList []string, err error) {
	fmt.Printf("chapterID:%#v\n", chapterID)
	pageList, err := getService().EbookPages(chapterID, token, index, count, offset)
	if err != nil {
		return
	}

	for _, item := range pageList.Pages {
		svgList = append(svgList, item.Svg)
	}
	fmt.Printf("IsEnd:%#v\n", pageList.IsEnd)
	if !pageList.IsEnd {
		index = count
		count += 20
		list, err1 := generateEbookPages(chapterID, token, index, count, offset)
		if err1 != nil {
			err = err1
			return
		}

		svgList = append(svgList, list...)
	}
	// utils.WriteFileWithTrunc(chapterID+".svg", strings.Join(svgList, "\n"))
	return
}
