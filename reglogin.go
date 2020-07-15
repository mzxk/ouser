package ouser

import (
	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Logout 用户登出
func Logout(p map[string]string) (interface{}, error) {
	ohttp.DeleteSession(p["key"])
	return nil, nil
}

//Login 用户登陆
func Login(p map[string]string) (interface{}, error) {
	//TODO limit IP
	user := p["user"]
	pwd := p["pwd"]

	if user == "" || pwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	var usr User
	err := mgo.C("user").FindOne(nil, bson.M{"user": user}).Decode(&usr)
	if err != nil {
		return nil, errs(ErrUserLogin)
	}
	if sha(pwd) != usr.Pwd {
		return nil, errs(ErrUserLogin)
	}
	return getLogin(usr.ID)
}

//RegisterSimple 用户注册
func RegisterSimple(p map[string]string) (interface{}, error) {
	//TODO limit IP
	user := p["user"]
	pwd := p["pwd"]

	if user == "" || pwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	c := mgo.C("user")
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

	return getLogin(usr.ID)
}

//通过通过用户id获取token
func getLogin(id primitive.ObjectID) (struLogin, error) {
	k, v, err := ohttp.AddSession(id.Hex())
	return struLogin{k, v}, err
}

type struLogin struct {
	Key   string
	Value string
}
