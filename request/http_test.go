package request

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGet(t *testing.T) {
	client := NewClient("https://www.baidu.com")
	resp, err := client.Request("GET", "")
	defer resp.Body.Close()
	if err != nil {

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	fmt.Println(string(body))
}
