package utils

import (
	"fmt"
	"testing"
)

func TestM3u8URLs(t *testing.T) {
	URLs, err := M3u8URLs("https://m.igetget.com/ddmedia/public/v1/m3u8/3368680087879724/52/m.m3u8")

	if err != nil {
		t.Fatal("M3u8URLs test failed", err)
	}
	fmt.Println(URLs)
}
