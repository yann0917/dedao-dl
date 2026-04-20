package app

import "github.com/yann0917/dedao-dl/services"

// User 获取当前登录用户信息
func User() (*services.User, error) {
	return getService().User()
}

// EbookVIPInfo 获取电子书 VIP 信息
func EbookVIPInfo() (*services.EbookVIPInfo, error) {
	return getService().EbookVIPInfo()
}

// OdobVIPInfo 获取每天听本书 VIP 信息
func OdobVIPInfo() (*services.OdobVipUser, error) {
	return getService().OdobVIPInfo()
}

// Recent 获取用户最近学习情况
func Recent(maxID int64, pageSize int, productType, uidHazy string, filterProductType bool) (*services.RecentResponse, error) {
	service := getService()
	if uidHazy == "" {
		user, err := service.User()
		if err != nil {
			return nil, err
		}
		uidHazy = user.UIDHazy
	}

	req := services.RecentRequest{
		FilterProductType: filterProductType,
		MaxID:             maxID,
		PageSize:          pageSize,
		ProductType:       productType,
		UID:               nil,
		UIDHazy:           uidHazy,
	}
	return service.Recent(req)
}
