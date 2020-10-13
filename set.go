package ouser

import (
	"errors"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
)

//PwdChange 更换用户密码
func (t *Ouser) PwdChange(p map[string]string) (interface{}, error) {
	pwd := p["pwd"]
	newPwd := p["newPwd"]
	if pwd == "" || newPwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	usr, _ := t.userCache(p)
	if sha(pwd) != usr.Pwd {
		return nil, errs(ErrUserLogin)
	}
	if cfg.OnlyGoogle {
		if err := t.checkPayPwd(p); err != nil {
			return nil, err
		}
	}
	return nil, t.userUpdateField(p["bsonid"], "pwd", sha(newPwd))
}

//SetNickName 用户设置显示名称
//Signed
//p "nickname" 用户昵称
func (t *Ouser) NickNameSet(p map[string]string) (interface{}, error) {
	id := p["bsonid"]
	nickname := p["nickname"]
	if id == "" || nickname == "" {
		return nil, errs(ErrParamsWrong)
	}
	err := t.userUpdateField(id, "nickname", nickname)
	return nil, err
}

//这个函数用户设置用户参数，id=用户的_id,然后需要更换的字段，需要更换的值
func (t *Ouser) userUpdateField(id, field string, value interface{}) error {
	if id == "" {
		return errors.New("wrongUserID")
	}
	_, err := t.mgo.C("user").UpdateOne(nil,
		bson.M{"_id": omongo.ID(id)},
		bson.M{"$set": bson.M{field: value}},
	)
	t.userCacheDelete(map[string]string{"bsonid": id})
	return err
}

//AvatarSet 设置用户头像
//Signed POST
//body里放入图片的二进制文件
func (t *Ouser) AvatarSet(p map[string]string) (interface{}, error) {
	bid := p["bsonid"]
	bt := []byte(p["body"])
	if len(bt) > 500*1024 || len(bt) < 1000 {
		return nil, errs(ErrAvatar)
	}
	cAvatar := t.mgo.C("avatar")
	avatarID := omongo.ID("")
	avatar := Avatar{avatarID, bid, bt}
	_, err := cAvatar.InsertOne(nil, avatar)
	if err != nil {
		return nil, err
	}
	err = t.userUpdateField(bid, "avatar", avatarID.Hex())
	return nil, err
}

//AvatarGet 获取用户头像
//返回用户头像的二进制
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

//PaypwdSet 设置支付密码
func (t *Ouser) PaypwdSet(p map[string]string) (interface{}, error) {
	if err := t.checkPrivateCode(p, smsPaypwdSet); err != nil {
		return nil, err
	}
	pwd := p["payPwd"]
	if pwd == "" {
		return nil, errs(ErrParamsWrong)
	}
	return nil, t.userUpdateField(p["bsonid"], "paypwd", sha(pwd))
}

//ContactChange 用户换绑手机
func (t *Ouser) ContactChange(p map[string]string) (interface{}, error) {
	if err := t.checkPrivateCode(p, smsPhoneChange); err != nil {
		return nil, err
	}
	field, err := t.checkPublicCode(p, smsPhoneBound)
	if err != nil {
		return nil, err
	}
	bid := omongo.ID(p["bsonid"])
	_, err = t.mgo.C("user").UpdateOne(nil,
		bson.M{"_id": bid}, bson.M{"$set": bson.M{field: p["contact"]}})
	if omongo.IsDuplicate(err) {
		return nil, errs(ErrUserExisted)
	}
	t.userCacheDelete(p)
	return nil, err
}

func (t *Ouser) checkPrivateCode(p map[string]string, checkType int) error {
	if cfg.OnlyGoogle {
		if t.g2faCheck(p) {
			return nil
		}
		return errs("wrongCode")
	}
	code := p["code"]
	bid := p["bsonid"]
	if checkType > 7100 || checkType < 7000 || code == "" || bid == "" {
		return errs(ErrParamsWrong)
	}
	key := joinSmsType(bid, checkType)
	if !ohttp.CodeCheck(key, code) {
		return errs(ErrCode)
	}
	return nil
}
