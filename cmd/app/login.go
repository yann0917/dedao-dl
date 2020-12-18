package app

import (
	"fmt"

	"github.com/yann0917/dedao-dl/config"
)

//Who get current login user
func Who() {
	activeUser := config.Instance.ActiveUser()
	fmt.Printf("当前帐号 uid: %s, 用户名: %s, 头像地址: %s \n", activeUser.UIDHazy, activeUser.Name, activeUser.Avatar)
}
