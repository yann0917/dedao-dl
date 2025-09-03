package app

import (
	"errors"

	"github.com/yann0917/dedao-dl/config"
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
	if err != nil {
		return
	}
	if list == nil {
		err = errors.New("已购书架为空")
		return
	}
	switch category {
	case CateCourse:
		for _, course := range list.List {
			config.Instance.SetCourseCache(category, course.ClassID, course)
		}
	case CateAudioBook, CateEbook:
		for _, course := range list.List {
			config.Instance.SetCourseCache(category, course.ID, course)
		}
	}

	return
}

// CourseInfo 已购课程详情
func CourseInfo(id int) (info *services.CourseInfo, err error) {
	course := config.Instance.GetCourseCache(CateCourse, id)
	enID := ""
	if course != nil {
		enID = course.Enid
	}
	if enID == "" {
		courseDetail, err1 := CourseDetail(CateCourse, id)
		if err1 != nil {
			err = err1
			return
		}
		enID = courseDetail.Enid
	}
	info, err = getService().CourseInfo(enID)
	if err != nil {
		return
	}
	return
}

// CourseDetail 已购课程详情
func CourseDetail(category string, id int) (course *services.Course, err error) {
	course = config.Instance.GetCourseCache(category, id)
	if course != nil && course.Enid != "" {
		return
	}

	// 如果获取不到或 Enid 为空，则从服务器获取
	detail, err1 := getService().CourseDetail(category, id)
	if err1 != nil {
		err = err1
		return
	}
	course = detail
	return
}
