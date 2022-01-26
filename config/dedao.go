package config

import (
	"github.com/yann0917/dedao-dl/services"
)

var (
	// BaseURL dedao url
	BaseURL = "https://www.dedao.cn"
)

// User dedao user info
type User struct {
	UIDHazy string `json:"uid_hazy"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
}

// Dedao geek time info
type Dedao struct {
	User
	services.CookieOptions
}

// New dedao service
func (d *Dedao) New() *services.Service {
	return services.NewService(&d.CookieOptions)
}
