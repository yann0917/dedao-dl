package app

import "github.com/yann0917/dedao-dl/services"

// TopicAll 推荐话题列表
func TopicAll() (list *services.TopicAll, err error) {
	list, err = getService().TopicAll(0, 10, true)
	if err != nil {
		return
	}
	return
}

// TopicDetail Topic Detail
func TopicDetail(id string) (detail *services.TopicDetail, err error) {
	detail, err = getService().TopicDetail(id)
	if err != nil {
		return
	}
	return
}

// TopicNotesList Topic NotesList
func TopicNotesList(id string) (list *services.NotesList, err error) {
	list, err = getService().TopicNotesList(id)
	if err != nil {
		return
	}
	return
}
