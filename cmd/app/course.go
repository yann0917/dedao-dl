package app

import (
	"github.com/yann0917/dedao-dl/services"
)

// CourseType 课程分类
func CourseType() (list *services.CourseCategoryList, err error) {
	list, err = getService().CourseType()
	return
}

// CourseList 已购课程列表
func CourseList(category string) (list *services.CourseList, err error) {
	list, err = getService().CourseListAll(category, "study")
	return
}

// CourseInfo 已购课程列表
func CourseInfo(id int) (info *services.CourseInfo, err error) {
	courseDetail, err := getService().CourseDetail(CateCourse, id)
	if err != nil {
		return
	}
	enID := courseDetail.Enid
	info, err = getService().CourseInfo(enID)
	if err != nil {
		return
	}
	return
}
