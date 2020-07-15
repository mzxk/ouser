package ouser

import (
	"errors"

	"github.com/mzxk/omongo"
)

const (
	ErrParamsWrong = "params error" //交了错误的参数
	ErrUserExisted = "user existed"
	ErrUserLogin   = "wrong user or password"
	ErrLimit       = "out of limit"
)

func errs(s string) error {
	return errors.New(s)
}

var mgo = omongo.NewMongoDB("mongodb://root:root@172.31.39.207:13198,172.31.39.208:13198/?authSource = admin&replicaSet = nonomongo", "user")
