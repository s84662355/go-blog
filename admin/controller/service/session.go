package service

import (
	libsession "GoBlog/lib/session"
	"GoBlog/logic/model"
	"encoding/json"
	"errors"
	_ "fmt"
)

type Session struct {
	s       libsession.HttpSessionInterface
	userKey string
	token   string
}

func CreateSession(token string) Session {
	s := Session{}
	c := libsession.CreateHttpSession(token)
	s.s = &c
	s.userKey = "userKey"
	s.token = token
	return s
}

func (l *Session) Get(key string) interface{} {
	return l.s.Get(key)
}

func (l *Session) Set(key string, data interface{}) bool {
	if !l.s.Set(key, data) {
		return false
	}
	return true
}

func (l *Session) GetSystemUser() (model.SystemUser, error) {
	userData := l.s.Get(l.userKey)
	//fmt.Println(userData)
	if userData != nil {
		res := model.SystemUser{}
		err := json.Unmarshal([]byte(userData.(string)), &res)
		//fmt.Println(res)
		if err != nil {
			return res, err
		}
		return res, nil

	}
	return model.SystemUser{}, errors.New("system user")
}

func (l *Session) SetSystemUser(data model.SystemUser) bool {
	str, err := json.Marshal(data)
	if err != nil {
		return false
	}
	if !l.s.Set(l.userKey, str) {
		return false
	}
	return true
}

func (l *Session) LoginOut() {
	l.s.Del()
}
