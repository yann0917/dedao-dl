package app

import (
	"errors"

	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

const (
	defaultListOrder = "study"
)

// CourseType 课程分类
func CourseType() (list *services.CourseCategoryList, err error) {
	list, err = getService().CourseType()
	return
}

// CourseList 已购课程列表
func CourseList(category string) (list *services.CourseList, err error) {
	return CourseListWithOptions(category, defaultListOrder, 0, 0)
}

// CourseListWithOptions 已购课程列表（支持排序和分页）
// page/limit 同时大于 0 时按单页拉取，否则拉取全量。
func CourseListWithOptions(category, order string, page, limit int) (list *services.CourseList, err error) {
	if page > 0 && limit > 0 {
		list, err = getService().CourseList(category, order, page, limit)
	} else {
		list, err = getService().CourseListAll(category, order)
	}
	if err != nil {
		return
	}
	if list == nil {
		err = errors.New("已购书架为空")
		return
	}
	cacheCourseList(category, list.List)
	return
}

func cacheCourseList(category string, courses []services.Course) {
	switch category {
	case CateCourse:
		for _, course := range courses {
			if course.IsGroup || course.ClassID <= 0 {
				continue
			}
			config.Instance.SetCourseCache(category, course.ClassID, course)
		}
	case CateAudioBook, CateEbook, CateAce:
		for _, course := range courses {
			if course.IsGroup || course.ID <= 0 {
				continue
			}
			config.Instance.SetCourseCache(category, course.ID, course)
		}
	}
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

// GetGroupItems fetches all items within a specific group.
// Returns an error if the group is empty or doesn't exist.
// 获取分组内的所有项目
func GetGroupItems(category string, groupID int) (list *services.CourseList, err error) {
	return GetGroupItemsWithOptions(category, defaultListOrder, groupID, 0, 0)
}

// GetGroupItemsWithOptions 获取分组项目（支持排序和分页）
// page/limit 同时大于 0 时按单页拉取，否则拉取全量。
func GetGroupItemsWithOptions(category, order string, groupID, page, limit int) (list *services.CourseList, err error) {
	if page > 0 && limit > 0 {
		list, err = getService().CourseGroupList(category, order, groupID, page, limit)
	} else {
		list, err = getService().CourseGroupListAll(category, order, groupID)
	}
	if err != nil {
		return
	}
	if list == nil {
		err = errors.New("分组为空或不存在")
		return
	}
	cacheCourseList(category, list.List)
	return
}
