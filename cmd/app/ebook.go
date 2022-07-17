package app

import (
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
		index, count, offset := 0, 20, 20000
		pageList, err1 := getService().EbookPages(order.ChapterID, token.Token, index, count, offset)
		if err1 != nil {
			err = err1
			return
		}

		for _, item := range pageList.Pages {
			svgContent = append(svgContent, item.Svg)
		}

		if !pageList.IsEnd {
			index = count
			count += 20
			pageList, err1 = getService().EbookPages(order.ChapterID, token.Token, index, count, offset)
			if err1 != nil {
				err = err1
				return
			}
			for _, item := range pageList.Pages {
				svgContent = append(svgContent, item.Svg)
			}
		}
	}
	if len(svgContent) > 0 {
		// utils.WriteFile(order.ChapterID+".svg", strings.Join(svgContent, "\n"))
		utils.Svg2Html(title, svgContent)
	}
	return
}
