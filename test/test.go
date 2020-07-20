package main

import (
	"github.com/mzxk/ohttp"
)

var host = "http://127.0.0.1:6666"

func main() {

}
func smsLogin(phone, code string) (result struReg, err error) {
	rst, err := ohttp.HTTP(host+"/user/smsLogin", map[string]interface{}{
		"contact": phone,
		"code":    code,
	}).Get()
	if err != nil {
		return
	}
	err = rst.JSONSelf(&result)
	return
}
func smsPublic(phone, tp string) (result interface{}, err error) {
	rst, err := ohttp.HTTP(host+"/user/smsPublic", map[string]interface{}{
		"type":    "login",
		"contact": phone,
	}).Get()
	if err != nil {
		return result, err
	}
	err = rst.JSON(&result)
	return
}
func feedbackGet(key, value string) (result interface{}, err error) {
	rst, err := ohttp.HTTPSign(host+"/user/feedbackGet", nil, key, value).Get()
	if err != nil {
		panic(err)
	}
	err = rst.JSONSelf(&result)
	return
}
func feedbackPull(key, value, text string) error {
	rst, err := ohttp.HTTPSign(host+"/user/feedbackPull", map[string]interface{}{
		"text": text,
	}, key, value).Get()
	if err != nil {
		panic(err)
	}
	var result interface{}
	err = rst.JSONSelf(&result)
	return err
}
func setNickName(key, value, name string) error {
	rst, err := ohttp.HTTPSign(host+"/user/setNickname", map[string]interface{}{
		"nickname": name,
	}, key, value).Get()
	if err != nil {
		panic(err)
	}
	var result interface{}
	err = rst.JSONSelf(&result)
	return err
}
func login(user, pwd string) (struReg, error) {
	rst, err := ohttp.HTTP(host+"/user/login", map[string]interface{}{
		"user": user,
		"pwd":  pwd,
	}).Get()
	if err != nil {
		panic(err)
	}
	var result struReg
	err = rst.JSONSelf(&result)
	return result, err
}
func regSimple(user, pwd string) (struReg, error) {
	rst, err := ohttp.HTTP(host+"/user/registerSimple", map[string]interface{}{
		"user": user,
		"pwd":  pwd,
	}).Get()
	if err != nil {
		panic(err)
	}
	var result struReg
	err = rst.JSONSelf(&result)
	return result, err
}

type struReg struct {
	Key   string
	Value string
}
