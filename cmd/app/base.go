package app

import (
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/services"
)

func getService() *services.Service {
	return config.Instance.ActiveUserService()
}
