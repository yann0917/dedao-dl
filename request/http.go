package request

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/imroc/req/v3"
)

var (
	// UserAgent UserAgent
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
)

// HTTPClient http client
type HTTPClient struct {
	req.Client
}

// NewClient new HTTPClient
func NewClient(baseURL string) *req.Client {
	c := req.C()
	c = c.SetBaseURL(baseURL)
	return c
}

// HTTPGet http get request
func HTTPGet(url string) (body []byte, err error) {
	r, err := req.Get(url)
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	return
}

// Get http get request
func Get(url string) (io.ReadCloser, error) {
	client := NewClient(url)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// Headers return the HTTP Headers of the url
func Headers(url string) (http.Header, error) {
	client := NewClient(url)
	resp, err := client.R().Get("")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Header, nil
}

// Size get size of the url
func Size(url string) (int, error) {
	h, err := Headers(url)
	if err != nil {
		return 0, err
	}
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not present")
	}
	size, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return size, nil
}
