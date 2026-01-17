package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/yann0917/dedao-dl/utils"
)

var (
	dedaoCommURL = &url.URL{
		Scheme: "https",
		Host:   "dedao.cn",
	}
	baseURL   = "https://www.dedao.cn"
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
)

// Response dedao success response
type Response struct {
	H respH `json:"h"`
	C respC `json:"c"`
}

type respH struct {
	C   int    `json:"c"`
	E   string `json:"e"`
	S   int    `json:"s"`
	T   int    `json:"t"`
	Apm string `json:"apm"`
}

// respC response content
type respC []byte

func (r *respC) UnmarshalJSON(data []byte) error {
	*r = data

	return nil
}

func (r respC) String() string {
	return string(r)
}

// Service dedao service
type Service struct {
	client *resty.Client
}

// CookieOptions dedao cookie options
type CookieOptions struct {
	GAT           string `json:"gat"`
	ISID          string `json:"isid"`
	Iget          string `json:"iget"`
	Token         string `json:"token"`
	GuardDeviceID string `json:"_guard_device_id"`
	SID           string `json:"_sid"`
	AcwTc         string `json:"acw_tc"`
	AliyungfTc    string `json:"aliyungf_tc"`
	CookieStr     string `json:"cookieStr"`
}

// NewService new service
func NewService(co *CookieOptions) *Service {
	cookies := []*http.Cookie{}
	cookies = append(cookies, &http.Cookie{
		Name:   "GAT",
		Value:  co.GAT,
		Domain: "." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "ISID",
		Value:  co.ISID,
		Domain: "." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_guard_device_id",
		Value:  co.GuardDeviceID,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_sid",
		Value:  co.SID,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "acw_tc",
		Value:  co.AcwTc,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "iget",
		Value:  co.Iget,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "token",
		Value:  co.Token,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "aliyungf_tc",
		Value:  co.AliyungfTc,
		Domain: "www." + dedaoCommURL.Host,
	})
	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(baseURL).
		SetCookies(cookies).
		SetHeader("User-Agent", UserAgent)

	return &Service{client: client}
}

func (r *Response) isSuccess() bool {
	return r.H.C == 0
}

func handleHTTPResponse(resp *resty.Response, err error) (io.ReadCloser, error) {
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", resp.Request.URL, err)
	}

	// Check content type for HTML error pages
	contentType := resp.Header().Get("Content-Type")
	if strings.Contains(contentType, "text/html") && resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("temporary: HTTP %d from %s - Server returned HTML error page (will retry)",
			resp.StatusCode(), resp.Request.URL)
	}

	// Permanent errors that shouldn't be retried
	switch resp.StatusCode() {
	case http.StatusNotFound:
		return nil, fmt.Errorf("404 NotFound from %s", resp.Request.URL)
	case http.StatusBadRequest:
		return nil, fmt.Errorf("400 BadRequest from %s", resp.Request.URL)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("401 Unauthorized from %s", resp.Request.URL)
	case 496:
		return nil, fmt.Errorf("496 NoCertificate from %s", resp.Request.URL)
	}

	// Temporary errors that should be retried
	if resp.StatusCode() == http.StatusBadGateway {
		return nil, fmt.Errorf("temporary: 502 Bad Gateway from %s - Backend service error", resp.Request.URL)
	}

	data := resp.Body()
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	reader := bytes.NewReader(data)
	result := &responseWrapper{
		ReadCloser: io.NopCloser(reader),
		rawData:    dataCopy,
	}
	return result, nil
}

type responseWrapper struct {
	io.ReadCloser
	rawData []byte
}

func (r *responseWrapper) GetRawResponse() []byte {
	return r.rawData
}

func handleJSONParse(reader io.Reader, v interface{}) error {
	result := new(Response)

	// Read the entire body
	rawData, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// // Log raw response if it's not valid JSON
	// if !json.Valid(rawData) {
	// 	// Write to error log file
	// 	logFile, _ := os.OpenFile(
	// 		fmt.Sprintf("error_log_%s.txt", time.Now().Format("2006-01-02_15-04-05")),
	// 		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
	// 		0666,
	// 	)
	// 	defer logFile.Close()

	// 	fmt.Fprintf(logFile, "\n=== Invalid JSON Response ===\n%s\n=== End Response ===\n", string(rawData))
	// }

	err = utils.UnmarshalReader(bytes.NewReader(rawData), &result)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w\nraw response: %s", err, string(rawData))
	}

	if !result.isSuccess() {
		return fmt.Errorf("service error: %s\nraw response: %s", result.H.E, string(rawData))
	}

	err = utils.UnmarshalJSON(result.C, v)
	if err != nil {
		return fmt.Errorf("unmarshal content error: %w\nraw response: %s", err, string(rawData))
	}

	return nil
}

// ParseCookies parse cookie string to cookie options
func ParseCookies(cookie string, v interface{}) (err error) {
	if cookie == "" {
		return errors.New("cookie is empty")
	}
	list := strings.Split(cookie, ";")
	cookieM := make(map[string]string, len(list))
	for _, item := range list {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			if parts[1] != "" {
				cookieM[strings.TrimSpace(parts[0])] = parts[1]
			}
		}
	}

	// 创建大小写不敏感的 map（为了兼容 mapstructure 的行为）
	cookieMInsensitive := make(map[string]string)
	for k, v := range cookieM {
		cookieMInsensitive[strings.ToLower(k)] = v
	}

	// 使用反射将 map 的值赋给结构体
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to struct")
	}

	elem := value.Elem()
	structType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := structType.Field(i)

		// 获取 json tag
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// 处理逗号后面的选项
		jsonName := strings.Split(tag, ",")[0]
		if jsonName == "" {
			jsonName = fieldType.Name
		}

		// 查找 map 中的值（大小写不敏感）
		if mapValue, ok := cookieMInsensitive[strings.ToLower(jsonName)]; ok {
			if field.Kind() == reflect.String {
				field.SetString(mapValue)
			}
		}
	}

	return nil
}
