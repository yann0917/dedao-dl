package services

import (
	"net/http"
	"net/url"

	"github.com/yann0917/dedao-dl/request"
)

var (
	dedaoCommURL = &url.URL{
		Scheme: "https",
		Host:   "dedao.cn",
	}
	baseURL = "https://www.dedao.cn"
)

//Service dedao service
type Service struct {
	client *request.HTTPClient
}

//NewService new service
func NewService(gat, isid, sid, acwTc, iget, token, guardDeviceID string) *Service {
	client := request.NewClient(baseURL)
	client.ResetCookieJar()
	cookies := []*http.Cookie{}
	cookies = append(cookies, &http.Cookie{
		Name:   "GAT",
		Value:  gat,
		Domain: "." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "ISID",
		Value:  isid,
		Domain: "." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_guard_device_id",
		Value:  guardDeviceID,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_sid",
		Value:  sid,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "acw_tc",
		Value:  acwTc,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "iget",
		Value:  iget,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "token",
		Value:  token,
		Domain: "www." + dedaoCommURL.Host,
	})
	client.Client.Jar.SetCookies(dedaoCommURL, cookies)

	return &Service{client: client}
}

//Cookies get cookies string
func (s *Service) Cookies() map[string]string {
	cookies := s.client.Client.Jar.Cookies(dedaoCommURL)

	cstr := map[string]string{}

	for _, cookie := range cookies {
		cstr[cookie.Name] = cookie.Value
	}

	return cstr
}
