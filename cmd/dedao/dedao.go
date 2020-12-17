package dedao

import (
	"github.com/yann0917/dedao-dl/services"
)

//User dedao user info
type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

//Dedao geek time info
type Dedao struct {
	User
	services.CookieOptions
}

//New new dedao service
func (d *Dedao) New() *services.Service {
	return services.NewService(&d.CookieOptions)
}
