package ouser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"github.com/mzxk/oval"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCache sync.Map //用户信息的缓存

type Ouser struct {
	mgo        *omongo.MongoDB
	httpClient *ohttp.Server
	sms        ohttp.Sms
}

func New(clt *ohttp.Server) *Ouser {
	return &Ouser{
		omongo.NewMongoDB(cfg.MongoURL, "user"),
		clt,
		ohttp.NewSms(cfg.Sms.Name, cfg.Sms.Key),
	}
}

//Logout 用户登出
func (t *Ouser) Logout(p map[string]string) (interface{}, error) {
	ohttp.DeleteSession(p["key"])
	return nil, nil
}
func (t *Ouser) SmsLogin(p map[string]string) (interface{}, error) {
	//判断是否是正常的联系方式
	contact := p["contact"]
	code := p["code"]
	if contact == "" || code == "" {
		return nil, errs(ErrParamsWrong)
	}
	isEmail := strings.Contains(contact, "@")
	if isEmail {
		return nil, errs(ErrParamsWrong)
	}
	if ohttp.CodeCheck(getType(contact, "login"), code) {
		usr, err := t.login("phone", contact)
		//代表不存在
		if err == mongo.ErrNoDocuments {
			//新建一个用户
			usrNew := User{
				ID:       omongo.ID(""),
				User:     "",
				RegIP:    p["ip"],
				Referrer: p["referrer"],
				Phone:    contact,
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
		return getLoginToken(usr.ID)
	}
	return nil, errs(ErrUserLogin)
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
	return t.userGet(p)
}

//这个函数将输入用户id，来获取用户信息
//如果缓存中存在，那么直接返回
//如果缓存中不存在，那么重新从库中读取
func (t *Ouser) userGet(p map[string]string) (*User, error) {
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
		userCache.Store(usr.ID.Hex(), &usr)
		oval.UnLimited("getUserByID" + id)
	}
	return &usr, err
}

//这个函数调用可能存在的字段和名字来获取用户信息
//比如 	"user" , 	"username"
//		"phone",	"16600001111"
//		"email",	"abc@bcd.com"
//同时如果获取成功，会存入userCache
func (t *Ouser) login(field, value string) (*User, error) {
	var usr User
	err := t.mgo.C("user").FindOne(nil, bson.M{field: value}).Decode(&usr)
	if err == nil && usr.User != "" {
		userCache.Store(usr.ID.Hex(), &usr)
	}
	return &usr, err
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
		User:     user,
		Pwd:      sha(pwd),
		RegIP:    p["ip"],
		Referrer: p["referrer"],
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
