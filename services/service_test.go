package services

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

var service *Service

func TestMain(m *testing.M) {
	// cookie := utils.Get(baseURL)
	cookie := "aliyungf_tc=92ca2d34f067361044be841c6c559599f5b3411f030fac6be52bc59e06cfb77c; ISID=b013cac8ffcabe56f77760d8bf6a47e6; csrfToken=2_EoW1CDysblBqri4ZEGtLH_; token=2_EoW1CDysblBqri4ZEGtLH_; _guard_device_id=1g7vrlkr0C2VzYWMB05NqrSDNUvbqqEotsv6YKt; Hm_lvt_be36b12b82a5f4eaa42c23989d277bb0=1657539221,1657716818,1657722320,1657848773; _sid=1g7vrllhlkomvv71g5e193aw5we817wf; acw_tc=2f6fc12616578553056474156e8e6d9076552d85101eff7b108bcc4bdb9595; GAT=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpZ2V0Z2V0LmNvbSIsImV4cCI6MTY2MDQ0NzMyNCwiaWF0IjoxNjU3ODU1MzI0LCJpc3MiOiJEREdXIEpXVCBNSURETEVXQVJFIiwibmJmIjoxNjU3ODU1MzI0LCJzdWIiOiIxNzk1MDMyNCIsImRldmljZV9pZCI6ImIwMTNjYWM4ZmZjYWJlNTZmNzc3NjBkOGJmNmE0N2U2IiwiZGV2aWNlX3R5cGUiOiJpZ2V0d2ViIn0.GpToOpMNjqRG8dbP8bHlE0BIK4AmW-jmGOwmWu7cORvLuSzS9GA1OpovuIMCIprLFp33xDiWlyMaPqdM2oNmfA; iget=eyJ1c2VyIjp7InVzZXJJZCI6MTc5NTAzMjQsIm5pY2tuYW1lIjoiWWFibyIsImF2YXRhciI6Imh0dHBzOi8vcGljY2RuLnVtaXdpLmNvbS91cGxvYWRlci9pbWFnZS9hdmF0YXIvMjAyMTEyMjAxMS8xNzYwOTA0MzUwMTE0MjUwOTIwIiwiYXZhdGFyU2hvcnQiOiJodHRwczovL3BpY2Nkbi51bWl3aS5jb20vdXBsb2FkZXIvaW1hZ2UvYXZhdGFyLzIwMjExMjIwMTEvMTc2MDkwNDM1MDExNDI1MDkyMD94LW9zcy1wcm9jZXNzPWltYWdlL3Jlc2l6ZSxwXzI1Iiwic3RhdHVzIjoyLCJwaG9uZSI6IjE3NjgyMzMwOTE3IiwibG9naW5UaW1lIjoiMjAyMi0wNy0xNSAxMToyMjowNCJ9LCJfZXhwaXJlIjoxNjU4NDYwMTI1MDQwLCJfbWF4QWdlIjo2MDQ4MDAwMDB9; Hm_lpvt_be36b12b82a5f4eaa42c23989d277bb0=1657855331"
	// cookie := ""
	co := &CookieOptions{}
	ParseCookies(cookie, &co)
	service = NewService(co)
	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

func TestGetLoginAccessToken(t *testing.T) {
	result, err := service.reqGetLoginAccessToken()
	if err != nil {
		fmt.Printf("%#v \n", err)
	}
	fmt.Printf("%#v \n", result)
}

func TestGet(t *testing.T) {
	fmt.Println(dedaoCommURL.String())
}

func TestNewService(t *testing.T) {
	resp, err := service.client.R().Get("/api/pc/user/info")
	if err != nil {
		fmt.Printf("%#v \n", err)
	}
	var user User
	data := resp.Body()

	reader := bytes.NewReader(data)
	result := io.NopCloser(reader)
	handleJSONParse(result, &user)
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
	result, err := service.ArticleList(ID, "", 30)
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
	result, err := service.ArticleInfo(enid, 1)
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
