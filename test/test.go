package main

import (
	"fmt"

	"github.com/mzxk/ohttp"
)

var host = "http://127.0.0.1:8090"

func main() {
	err := setNickName("a7099690ae", "10711c2e48b2b06af", "狗老猫")
	fmt.Println(err)
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
