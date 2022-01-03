package services

import "github.com/pkg/errors"

// User user info
type User struct {
	Nickname        string  `json:"nickname"`
	Avatar          string  `json:"avatar"`
	TodayStudyTime  int     `json:"today_study_time"`
	StudySerialDays int     `json:"study_serial_days"`
	IsV             int     `JSON:"is_v"`
	VIPUser         VIPUser `json:"vip_user"`
	IsTeacher       int     `json:"is_teacher"`
	UIDHazy         string  `json:"uid_hazy"`
}

// VIPUser vip info
type VIPUser struct {
	Info string `json:"info"`
	Stat int    `json:"stat"`
}

// User get user info
func (s *Service) User() (user *User, err error) {
	body, err := s.reqUser()
	if err != nil {
		err = errors.Wrap(err, "request user err")
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &user); err != nil {
		return
	}
	return
}

// Token token
type Token struct {
	Token string `json:"token"`
}

// Token get token
func (s *Service) Token() (t *Token, err error) {
	body, err := s.reqToken()
	if err != nil {
		err = errors.Wrap(err, "request token err")
		return
	}
	defer body.Close()
	if err = handleJSONParse(body, &t); err != nil {
		return
	}
	return
}
