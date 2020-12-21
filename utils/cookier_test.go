package utils

import (
	"strings"
	"testing"

	"github.com/go-rod/rod"
)

func TestGet(t *testing.T) {
	rod.Try(func() {
		c := Get("https://www.dedao.cn")
		if !strings.Contains(c, "ISID=") {
			t.Error("cookie should contain the ISID")
		}
		t.Log(c)
	})

}
