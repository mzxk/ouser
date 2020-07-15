package ouser

import (
	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
)

//用户反馈
func Feedback(p map[string]string) (interface{}, error) {
	ohttp.DeleteSession(p["key"])
	return nil, nil
}

//用户设置显示名称
func SetNickName(p map[string]string) (interface{}, error) {
	id := p["bsonid"]
	nickname := p["nickname"]
	if id == "" || nickname == "" {
		return nil, errs(ErrParamsWrong)
	}
	_, err := mgo.C("user").UpdateOne(nil,
		bson.M{"_id": omongo.ID(id)},
		bson.M{"$set": bson.M{"nickname": nickname}})
	return nil, err
}
