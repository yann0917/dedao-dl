package services

import (
	"github.com/yann0917/dedao-dl/utils"
)

type QrCodeResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		QrCode       string `json:"qrcode"`
		QrCodeString string `json:"qrCodeString"`
	} `json:"data"`
}

type CheckLoginResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		Status int `json:"status"` // 1-扫码成功,2-过期
	} `json:"data"`
}

// LoginAccessToken get login access token
func (s *Service) LoginAccessToken() (token string, err error) {
	token, err = s.reqGetLoginAccessToken()
	if err != nil {
		return
	}

	return
}

func (s *Service) GetQrcode(token string) (resp *QrCodeResp, err error) {
	resp, err = s.reqGetQrcode(token)
	if err != nil {
		return
	}
	content := "https://m.igetget.com/oauth/qrcode/v2/authorize?token=" + resp.Data.QrCodeString
	obj := utils.New()
	obj.Get(content).Print()

	return
}

func (s *Service) CheckLogin(token, qrcode string) (check *CheckLoginResp, cookie string, err error) {
	check, cookie, err = s.reqCheckLogin(token, qrcode)
	if err != nil {
		return
	}

	return
}
