package main

import (
	"errors"
	"log"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

func FTPTest() {
	driver, err := file.NewDriver("/home/leaf/")
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(&server.Options{
		Driver:    driver,
		Auth:      &UserAuth{},
		Perm:      server.NewSimplePerm("root", "root"),
		RateLimit: 0,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type UserAuth struct{}

func (ua *UserAuth) CheckPasswd(ctx *server.Context, username, password string) (bool, error) {
	if username == "leaf" && password == "test" {
		return true, nil
	}
	return false, errors.New("用户密码错误")
}
