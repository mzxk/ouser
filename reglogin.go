package ouser

import (
	"fmt"
	"log"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"github.com/mzxk/oval"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Logout 用户登出
func (t *Ouser) Logout(p map[string]string) (interface{}, error) {
	ohttp.DeleteSession(p["key"])
	return nil, nil
}

//SmsResetPwd 这将重制用户登录密码，输入参数为 contact=联系方式 ， pwd=密码 ， contactType=验证方式
func (t *Ouser) SmsResetPwd(p map[string]string) (interface{}, error) {
	contactType, err := t.checkPublicCode(p, smsResetPwd)
	if err != nil {
		return nil, err
	}
	contact := p["contact"]
	pwd := p["pwd"]
	if pwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	_, err = t.mgo.C("user").UpdateOne(nil,
		bson.M{contactType: contact},
		bson.M{"$set": bson.M{
			"pwd": sha(pwd),
		}})
	return nil, err
}

//这个函数将统一的验证公开sms输入的合理性，比如code是否存在，contact是否正常，同时验证code是否正确
//默认的，这将联系方式设置为手机短信，如果前端有额外的参数contactType=email，才会设置成email
func (t *Ouser) checkPublicCode(p map[string]string, checkType int) (string, error) {
	if cfg.OnlyGoogle {
		return "nil", errs("SmsLoginClosed")
	}
	//判断是否是正常的联系方式
	contact := p["contact"]
	code := p["codePublic"]
	field := "phone"
	if p["contactType"] == "email" {
		field = "email"
	}
	if contact == "" || code == "" {
		return field, errs(ErrParamsWrong)
	}

	if ohttp.CodeCheck(joinSmsType(contact, checkType), code) == false {
		return field, errs(ErrCode)
	}
	return field, nil
}

//SmsLogin 使用短信验证码登录，如果账户不存在，就新建一个
//可选字段有 		referrer 	推荐人
//可选字段			paypwd：	支付密码
//					pwd:		登录密码
func (t *Ouser) SmsLogin(p map[string]string) (interface{}, error) {
	contactType, err := t.checkPublicCode(p, smsLogin)
	if err != nil {
		return nil, err
	}
	contact := p["contact"]
	usr, err := t.userGetByField(contactType, contact)
	//代表不存在
	if err == mongo.ErrNoDocuments {
		//新建一个用户
		usrNew := User{
			User:     contact,
			ID:       omongo.ID(""),
			Uid:      <-idCreate,
			RegIP:    p["ip"],
			Referrer: t.getReferrerID(p["referrer"]),
			Paypwd:   sha(p["paypwd"]),
			Pwd:      sha(p["pwd"]),
		}
		if contactType == "phone" {
			usrNew.Phone = contact
		} else {
			usrNew.Email = contact
		}
		//插入用户
		_, err = t.mgo.C("user").InsertOne(nil, usrNew)
		//有极低的概率这里会有重复的用户名
		if err != nil && omongo.IsDuplicate(err) {
			fmt.Println(err)
			return nil, errs(ErrUserExisted)
		}
		//这里代表插入成功，返回key和value
		if err == nil {
			return getLoginToken(usrNew.ID)
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if p["resetPwd"] != "" {
		_ = t.userUpdateField(usr.ID.Hex(), "pwd", sha(p["resetPwd"]))
	}
	return getLoginToken(usr.ID)

}

//Login 用户登陆
func (t *Ouser) Login(p map[string]string) (interface{}, error) {
	user := p["user"]
	pwd := p["pwd"]

	if user == "" || pwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	var usr User
	err := t.mgo.C("user").FindOne(nil, bson.M{"user": user}).Decode(&usr)
	if err != nil {
		return nil, errs(ErrUserLogin)
	}
	if sha(pwd) != usr.Pwd {
		if oval.Limited(user+"login", 60, cfg.Login.Limit) {
			return nil, errs(ErrLimit)
		}
		return nil, errs(ErrUserLogin)
	}
	oval.UnLimited(user + "login")
	return getLoginToken(usr.ID)
}

func (t *Ouser) UserInfo(p map[string]string) (interface{}, error) {
	return t.userCache(p)
}

//这个函数将输入用户id，来获取用户信息
//如果缓存中存在，那么直接返回
//如果缓存中不存在，那么重新从库中读取
func (t *Ouser) userCache(p map[string]string) (*User, error) {
	id := p["bsonid"]
	if id == "" {
		return nil, errs(ErrParamsWrong)
	}
	if u, ok := userCache.Load(id); ok {
		return u.(*User), nil
	}
	//这里的限制是为了避免恶意用户大量调用
	if oval.Limited("getUserByID"+id, 60, 5) {
		return nil, errs(ErrLimit)
	}
	var usr User
	err := t.mgo.C("user").FindOne(nil, bson.M{"_id": omongo.ID(id)}).Decode(&usr)
	if err == nil && usr.User != "" {
		userCache.Store(usr.ID.Hex(), &usr, 7200)
		oval.UnLimited("getUserByID" + id)
	}
	return &usr, err
}

//在有更新后清楚一下用户缓存
func (t *Ouser) userCacheDelete(p map[string]string) {
	userCache.Delete(p["bsonid"])
}

//这个函数调用可能存在的字段和名字来获取用户信息
//比如 	"user" , 	"username"
//		"phone",	"16600001111"
//		"email",	"abc@bcd.com"
//同时如果获取成功，会存入userCache
func (t *Ouser) userGetByField(field, value string) (*User, error) {
	var usr User
	err := t.mgo.C("user").FindOne(nil, bson.M{field: value}).Decode(&usr)
	if err == nil && usr.User != "" {
		userCache.Store(usr.ID.Hex(), &usr, 7200)
	}
	return &usr, err
}

//这个函数用户获得推荐人实际的ID
func (t *Ouser) getReferrerID(referrer string) string {
	//当推荐人为空，返回空
	if referrer == "" {
		return ""
	}
	usr, e := t.userGetByField("uid", referrer)
	if e != nil || usr == nil {
		log.Println("GetReferrerError:", referrer, e)
		return ""
	}
	return usr.ID.Hex()
}

//RegisterSimple 用户注册
func (t *Ouser) RegisterSimple(p map[string]string) (interface{}, error) {
	if cfg.Register.SimpleClosed {
		return nil, errs(ErrClosedMethod)
	}
	if oval.Limited(p["ip"]+"reg", 60, cfg.Register.Limit) {
		return nil, errs(ErrLimit)
	}
	user := p["user"]
	pwd := p["pwd"]

	if user == "" || pwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	c := t.mgo.C("user")
	usr := User{
		ID:       omongo.ID(""),
		Uid:      <-idCreate,
		User:     user,
		Pwd:      sha(pwd),
		RegIP:    p["ip"],
		Paypwd:   sha(p["payPwd"]),
		Referrer: t.getReferrerID(p["referrer"]),
	}
	rst, err := c.Upsert(nil, bson.M{"user": usr.User}, bson.M{"$setOnInsert": usr})
	if err != nil {
		return nil, err
	}
	if rst.UpsertedID == nil {
		return nil, errs(ErrUserExisted)
	}

	return getLoginToken(usr.ID)
}

//通过通过用户id获取token
func getLoginToken(id primitive.ObjectID) (struLogin, error) {
	k, v, err := ohttp.AddSession(id.Hex())
	return struLogin{k, v}, err
}

type struLogin struct {
	Key   string
	Value string
}
