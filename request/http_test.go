package request

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	resp, err := HTTPGet("https://www.baidu.com")
	if err != nil {

	}

	fmt.Println(string(resp))
}
