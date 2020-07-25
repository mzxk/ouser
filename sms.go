package ouser

import (
	"strings"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/oval"
)

//SmsPublic 发送短信验证码，公开接口，只能用于注册，登录，找回密码，修改手机号
//* params
//* type          发送代码        6001｜6002｜6005 注册｜重置密码｜绑定手机
//* contact       联系方式        手机号码｜信箱     后端自动识别
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
	tp := int(s2i(p["type"]))
	if tp != 6001 && tp != 6002 && tp != 6005 {
		return nil, errs(ErrParamsWrong)
	}

	if oval.Limited(getType("limitPublicSms", p["ip"]), 60, 5) || oval.Limited(getType("limitPublicSms", contact), 60, 1) {
		return nil, errs(ErrLimit)
	}
	//获取一个 code
	code, err := ohttp.CodeGet(joinSmsType(contact, tp), 600, 3)
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

//SmsPrivate 这是个只能发送用户自身的短信,
func (t *Ouser) SmsPrivate(p map[string]string) (interface{}, error) {
	return t.contactPrivate(p, true)
}

//SmsPrivate 这是个只能发送用户自身的短信,
func (t *Ouser) contactPrivate(p map[string]string, isPhone bool) (interface{}, error) {
	//确认是不是正常的id
	tp := int(s2i(p["type"]))
	if tp > 6100 || tp < 6000 {
		return nil, errs(ErrParamsWrong)
	}
	//获取用户缓存
	usr, err := t.userCache(p)
	if err != nil {
		return nil, err
	}
	//确认用户联系方式是否正确
	contact := usr.Phone
	if !isPhone {
		contact = usr.Email
	}
	if contact == "" {
		return nil, errs(ErrContact)
	}
	//用户短信保存的key
	codeKey := joinSmsType(usr.ID.Hex(), tp)
	//用户短信是否超限
	if oval.Limited(getType("limitPrivateSms", usr.ID.Hex()), 60, 5) || oval.Limited(getType("limitPrivateSms", contact),
		60,
		1) {
		return nil, errs(ErrLimit)
	}
	//生成用户验证码
	code, err := ohttp.CodeGet(codeKey, 600, 3)
	if err != nil {
		return nil, err
	}
	if isPhone {
		//判断发送模版是否存在，如果不存在，使用默认模版
		modelID, ok := cfg.Sms.IDs[tp]
		if !ok {
			modelID = cfg.Sms.IDefault
		}
		//发送短信
		err = t.sms.Send(usr.Phone, modelID, code)
		return nil, err
	}
	//TODO 发送email
	return nil, errs(ErrParamsWrong)
}

const (
	smsLogin = iota + 6001
	smsResetPwd
	smsPhoneBound
)
const (
	smsPaypwdSet = iota + 7001
	smsPhoneChange
	smsWithdraw
)
