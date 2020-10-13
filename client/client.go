package ouserClient

import (
	"fmt"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/ouser"
)

type Client struct {
	Key   string
	Value string
	Url   string
}

func New(url string, group map[string]string) *Client {
	t := &Client{Url: url}
	//登陆
	user := Readline("Input User...")
	pwd := ReadPwd("InputPwd...")
	rsp, err := ohttp.HTTP(t.Url+"/user/login", map[string]interface{}{"user": user, "pwd": pwd}).Get()
	if err != nil {
		panic(err)
	}
	var result map[string]string
	err = rsp.JSONSelf(&result)
	if err != nil {
		panic(err)
	}
	t.Key = result["Key"]
	t.Value = result["Value"]
	//确认用户组
	var rlt2 ouser.User
	err2 := t.sign("/user/info", nil, &rlt2)
	if err2 != nil {
		panic(err2)
	}

	for k, v := range group {
		if rlt2.Group[k] != v {
			fmt.Println(group, rlt2.Group)
			panic("notRightGroup")
		}
	}
	return t
}
func (t *Client) UserInfo(userID string) (result ouser.User, err error) {
	err = t.sign("/user/admin/info", map[string]interface{}{"userID": userID}, &result)
	return
}
func (t *Client) CheckGP(userID, payPwd, g2fa string) (bool, bool) {
	var result []bool
	_ = t.sign("/user/admin/checkGP", map[string]interface{}{
		"payPwd": payPwd, "googleCode": g2fa, "userID": userID,
	}, &result)
	if len(result) == 2 {
		return result[0], result[1]
	}
	return false, false
}
func (t *Client) sign(url string, p map[string]interface{}, result interface{}) error {
	rsp2, err2 := ohttp.HTTPSign(t.Url+url, p, t.Key, t.Value).Get()
	if err2 != nil {
		return err2
	}
	return rsp2.JSONSelf(result)
}
