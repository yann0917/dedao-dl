package app

import "github.com/yann0917/dedao-dl/services"

// SearchSuggest 搜索建议
func SearchSuggest(query string, searchType int) (resp *services.SearchSuggestResp, err error) {
	resp, err = getService().SearchSuggest(query, searchType)
	if err != nil {
		return
	}
	return
}
