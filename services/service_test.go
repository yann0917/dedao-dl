package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/yann0917/dedao-dl/utils"
)

var service *Service

func TestMain(m *testing.M) {
	cookie := utils.Get(baseURL)
	co := &CookieOptions{}
	ParseCookies(cookie, &co)
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

func TestToken(t *testing.T) {
	result, err := service.Token()
	if err != nil {
		fmt.Printf("%#v \n", err)
	}
	fmt.Printf("%#v \n", result)
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

func TestEbookDetail(t *testing.T) {
	enid := "DLnMGAEG7gKLyYmkAbPaEXxD8BM4J0LMedWROrpdZn19VNzv2o5e6lqjQQ1poxqy"
	result, err := service.EbookDetail(enid)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestEbookReadToken(t *testing.T) {
	enid := "DLnMGAEG7gKLyYmkAbPaEXxD8BM4J0LMedWROrpdZn19VNzv2o5e6lqjQQ1poxqy"
	result, err := service.EbookReadToken(enid)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}
