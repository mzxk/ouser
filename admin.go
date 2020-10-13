package ouser

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *Ouser) AdminGet(p map[string]string) (interface{}, error) {
	if err := t.adminCheckGroup(p, "adminRead", "all"); err != nil {
		return nil, err
	}
	query := p["query"]
	db := p["db"]
	coll := p["coll"]
	if db == "" || coll == "" || query == "" {
		return nil, errs("WrongDbCollQuery")
	}
	var find bson.M
	err := json.Unmarshal([]byte(query), &find)
	if err != nil {
		return nil, err
	}
	c := t.mgo.CDb(db, coll)
	var result []bson.M
	err = c.FindAll(nil, find).All(&result)
	return result, err
}
func (t *Ouser) AdminUserInfo(p map[string]string) (interface{}, error) {
	if err := t.adminCheckGroup(p, "userInfo", "get"); err != nil {
		return nil, err
	}
	find := map[string]string{"bsonid": p["userID"]}
	return t.UserInfo(find)
}

//userID代表用户名，参数googleKey和参数payPwd，返回为bool数组，第一个确认支付密码是否正常，第二个确认google是否正常
func (t *Ouser) AdminCheckGP(p map[string]string) (interface{}, error) {
	if err := t.adminCheckGroup(p, "userCheck", "get"); err != nil {
		return nil, err
	}
	p["bsonid"] = p["userID"]
	result := make([]bool, 2)
	result[0] = t.checkPayPwd(p) == nil
	result[1] = t.g2faCheck(p)
	return result, nil
}
func (t *Ouser) adminCheckGroup(p map[string]string, groupType, groupValue string) error {
	usr, err := t.userCache(p)
	if err != nil {
		return err
	}
	if usr.Group[groupType] == groupValue {
		return nil
	}
	return errs("WrongGroup")
}
