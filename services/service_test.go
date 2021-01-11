package services

import (
	"fmt"
	"os"
	"testing"
)

var service *Service

func TestMain(m *testing.M) {
	// cookie := utils.Get(baseURL)
	cookie := "ISID=538c3c9bc97227fb31290d9f54ad87cd; _sid=1ekdk2rsmmivs75orohjpkk6g2o49vpo; _guard_device_id=1epe8s675vg8Gg5c5dGvhQcgDfe9PD5FdKTkZc7; token=iXBjd6WM-ar7lw2242KOjqm4sQPm-lJViRnw; iget=eyJzZWNyZXQiOiJzVWdwSVJwc3c3QzVoY21DYVhRX1RlSnoiLCJfZXhwaXJlIjoxNjEwMTk2MzkwMzI4LCJfbWF4QWdlIjo2MDQ4MDAwMDB9; GAT=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpZ2V0Z2V0LmNvbSIsImV4cCI6MTYxMDQ1MjkwNiwiaWF0IjoxNjEwMDIwOTA2LCJpc3MiOiJEREdXIEpXVCBNSURETEVXQVJFIiwibmJmIjoxNjEwMDIwOTA2LCJzdWIiOiIxNzk1MDMyNCIsImRldmljZV9pZCI6IjUzOGMzYzliYzk3MjI3ZmIzMTI5MGQ5ZjU0YWQ4N2NkIiwiZGV2aWNlX3R5cGUiOiJpZ2V0d2ViIn0.GH0o6ZYOCt0tSI4GzUnJnRkvJqEBOKpnI4FrFADfyN0GD0oE4-0ZwhxQ7qbXfWYCtSRN_iJ8PiZ3XqJyijhtmA; acw_tc=276082a916100771826244684e8b9b53f794448d7df554acf58703e1319d18"
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
	result, err := service.ArticleList(ID, 30)
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

func TestArticleInfo(t *testing.T) {
	enid := "R2Mo65zY4QZ3VnmvraKqEdNAa98jGB"
	result, err := service.ArticleInfo(enid)
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

func TestEbookVIPInfo(t *testing.T) {
	result, err := service.EbookVIPInfo()
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestTopicAll(t *testing.T) {
	result, err := service.TopicAll(0, 10, false)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestTopicDetail(t *testing.T) {
	id := "4qpo7LxOynVXemeY6pALW1JXrlwG6E"
	result, err := service.TopicDetail(id)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}

func TestTopicNotesList(t *testing.T) {
	id := "4qpo7LxOynVXemeY6pALW1JXrlwG6E"
	result, err := service.TopicNotesList(id)
	if err != nil {
		fmt.Printf("err:=%#v \n", err)
	}
	fmt.Printf("result:=%v \n", result)
}
