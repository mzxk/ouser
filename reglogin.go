package ouser

import (
	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"github.com/mzxk/oval"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ouser struct {
	mgo        *omongo.MongoDB
	httpClient *ohttp.Server
}

func New(mgo *omongo.MongoDB, clt *ohttp.Server) *Ouser {
	return &Ouser{mgo, clt}
}

//Logout 用户登出
func (t *Ouser) Logout(p map[string]string) (interface{}, error) {
	ohttp.DeleteSession(p["key"])
	return nil, nil
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
		if oval.Limited(user+"login", 60, 5) {
			return nil, errs(ErrLimit)
		}
		return nil, errs(ErrUserLogin)
	}
	oval.UnLimited(user + "login")
	return getLoginToken(usr.ID)
}

//RegisterSimple 用户注册
func (t *Ouser) RegisterSimple(p map[string]string) (interface{}, error) {
	if oval.Limited(p["ip"]+"reg", 60, 5) {
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
