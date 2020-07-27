package ouser

import (
	"errors"
)

const (
	ErrParamsWrong  = "params error" //交了错误的参数
	ErrUserExisted  = "existed"
	ErrUserLogin    = "wrong user or password"
	ErrLimit        = "out of limit"
	ErrAvatar       = "wrong Avatar Size"
	ErrClosedMethod = "method closed"
	ErrCode         = "wrong contact or code"
	ErrPayPwdNeed   = "must set pay password"
	ErrContact      = "contact error"
	ErrWrongPwd     = "password wrong"
	ErrGoogle2fa    = "google code wrong"
)

func errs(s string) error {
	return errors.New(s)
}
