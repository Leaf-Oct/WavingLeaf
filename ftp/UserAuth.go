package ftp

import (
	"errors"

	"goftp.io/server/v2"
)

type UserAuth struct{}

func (ua *UserAuth) CheckPasswd(ctx *server.Context, username, password string) (bool, error) {
	// TODO 待完善成接数据库的逻辑
	if username == "leaf" && password == "test" {
		return true, nil
	}
	return false, errors.New("用户密码错误")
}
