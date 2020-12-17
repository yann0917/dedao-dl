package services

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-rod/rod"
	"github.com/yann0917/dedao-dl/utils"
)

var service *Service

func TestMain(m *testing.M) {
	co := &CookieOptions{
		GAT:           "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpZ2V0Z2V0LmNvbSIsImV4cCI6MTYwODMwMjY4OSwiaWF0IjoxNjA3ODcwNjg5LCJpc3MiOiJEREdXIEpXVCBNSURETEVXQVJFIiwibmJmIjoxNjA3ODcwNjg5LCJzdWIiOiIxNzk1MDMyNCIsImRldmljZV9pZCI6IjUzOGMzYzliYzk3MjI3ZmIzMTI5MGQ5ZjU0YWQ4N2NkIiwiZGV2aWNlX3R5cGUiOiJpZ2V0d2ViIn0.Vo_sORYNqr46IfnjwJyGpvQI8JeNIvt2cMjSos3awWkWwa9PiA8T6mARvH1GCfyX7EK6K5rNSnP9JBLWL-jFWQ",
		ISID:          "538c3c9bc97227fb31290d9f54ad87cd",
		SID:           "1ekdk2rsmmivs75orohjpkk6g2o49vpo",
		AcwTc:         "276082a916078690867185436edc12a1ca467545eedeab821efac43d6b6918",
		Iget:          "eyJzZWNyZXQiOiJYQndxVndSaTV1RENkejE5X2dtVEFYX1UiLCJfZXhwaXJlIjoxNjA4NDcxMjAyNjEwLCJfbWF4QWdlIjo2MDQ4MDAwMDB9",
		Token:         "alkcIlpm-1j-rusU7g3wL1TnF0QrL-oLPhIg",
		GuardDeviceID: "1epe8s675vg8Gg5c5dGvhQcgDfe9PD5FdKTkZc7",
	}
	service = NewService(co)
	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}
func TestGet(t *testing.T) {
	fmt.Println(dedaoCommURL.String())
}

func TestNewService(t *testing.T) {
	resp, _ := service.client.Request("GET", "/api/pc/user/info")
	defer resp.Body.Close()

	var user User

	handleJSONParse(resp.Body, &user)
	fmt.Println(user)
}

func TestUser(t *testing.T) {
	result, err := service.User()
	if err != nil {
		fmt.Printf("%#v \n", err)
	}
	fmt.Printf("%#v \n", result)
}

func TestCourseType(t *testing.T) {
	result, err := service.CourseType()
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestCourseList(t *testing.T) {
	result, err := service.CourseList("bauhinia", "study", 1, 10)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestCourseInfo(t *testing.T) {
	ID := "OY8PNZj5EavJq1aHO9Jn1eqGDdlgw7"
	result, err := service.CourseInfo(ID)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestArticleList(t *testing.T) {
	ID := "OY8PNZj5EavJq1aHO9Jn1eqGDdlgw7"
	result, err := service.ArticleList(ID)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestAudioByAlias(t *testing.T) {
	ID := "zDWMvqA3d2k94LZ9KQ0RVjnxyapBePZ7"
	result, err := service.AudioByAlias(ID)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestArticleDetail(t *testing.T) {
	token := "KWn/CP3W2txbAhtG26cVSr0YwlF3n7LCqzYAOHpyWw3+ft2hqSH+BqlOZTnBur2gXU0ByFmUQmz0tYVxepbdpTy81Gk="
	sign := "b23a426b357d1b83"
	appID := "1632426125495894021"
	result, err := service.ArticleDetail(token, sign, appID)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestConvertToStruct(t *testing.T) {
	rod.Try(func() {
		c := utils.Get("https://www.dedao.cn")
		if !strings.Contains(c, "ISID=") {
			t.Error("cookie should contain the ISID")
		}
		var cookie CookieOptions
		ConvertToStruct(c, &cookie)
		t.Log(c)
	})
}
