package ouser

import (
	"strings"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/oval"
)

//SmsPublic 发送短信验证码，公开接口，只能用于注册，登录，找回密码，修改手机号
//参数
//	type : login|reset
//	contact : email或者手机号码
func (t *Ouser) SmsPublic(p map[string]string) (interface{}, error) {
	//判断是否是正常的联系方式，手机或者信箱
	contact := p["contact"]
	if contact == "" { //TODO 这里应该正则判断手机号码和email
		return nil, errs(ErrParamsWrong)
	}
	isEmail := strings.Contains(contact, "@")
	if isEmail {
		return nil, errs(ErrParamsWrong)
	}
	//首先判断是否可以正常发送短信 这里有email.的问题
	if cfg.Sms.Name == "" {
		return nil, errs(ErrClosedMethod)
	}
	//判断是否是允许的type
	tp := p["type"]
	if tp != "register" && tp != "login" && tp != "reset" {
		return nil, errs(ErrParamsWrong)
	}

	if oval.Limited(getType("limitPublicSms", p["ip"]), 60, 5) || oval.Limited(getType("limitPublicSms", contact), 60, 1) {
		return nil, errs(ErrLimit)
	}
	//获取一个 code
	code, err := ohttp.CodeGet(getType(contact, tp), 600, 3)
	if err != nil {
		return nil, err
	}
	//手机的发送方式
	if !isEmail {
		//判断发送模版是否存在，如果不存在，使用默认模版
		modelID, ok := cfg.Sms.IDs[tp]
		if !ok {
			modelID = cfg.Sms.IDefault
		}
		//发送短信
		err = t.sms.Send(contact, modelID, code)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
