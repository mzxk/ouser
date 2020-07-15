package ouser

import (
	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
)

//SetNickName 用户设置显示名称
func (t *Ouser) NickNameSet(p map[string]string) (interface{}, error) {
	id := p["bsonid"]
	nickname := p["nickname"]
	if id == "" || nickname == "" {
		return nil, errs(ErrParamsWrong)
	}
	_, err := t.mgo.C("user").UpdateOne(nil,
		bson.M{"_id": omongo.ID(id)},
		bson.M{"$set": bson.M{"nickname": nickname}})
	return nil, err
}
func (t *Ouser) AvatarSet(p map[string]string) (interface{}, error) {
	bid := p["bsonid"]
	bt := []byte(p["body"])
	if len(bt) > 500*1024 || len(bt) < 1000 {
		return nil, errs(ErrAvatar)
	}
	cUser := t.mgo.C("user")
	cAvatar := t.mgo.C("avatar")
	avatarID := omongo.ID("")
	avatar := Avatar{avatarID, bid, bt}
	_, err := cAvatar.InsertOne(nil, avatar)
	if err != nil {
		return nil, err
	}
	_, err = cUser.UpdateOne(nil, bson.M{"_id": omongo.ID(bid)},
		bson.M{
			"$set": bson.M{
				"avatar": avatarID.Hex(),
			},
		},
	)
	return nil, err
}
func (t *Ouser) AvatarGet(p map[string]string) (interface{}, error) {
	id := p["id"]
	if len(id) != 24 {
		return nil, errs(ErrParamsWrong)
	}
	c := t.mgo.C("avatar")
	var result Avatar
	err := c.FindOne(nil, bson.M{"_id": omongo.ID(id)}).Decode(&result)
	return result, err
}
