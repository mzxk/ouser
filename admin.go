package ouser

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *Ouser) AdminGet(p map[string]string) (interface{}, error) {
	usr, err := t.userCache(p)
	if err != nil {
		return nil, err
	}
	if usr.Group["AdminRead"] != "all" {
		return nil, errs("WrongGroup")
	}
	query := p["query"]
	db := p["db"]
	coll := p["coll"]
	if db == "" || coll == "" || query == "" {
		return nil, errs("WrongDbCollQuery")
	}
	var find bson.M
	err = json.Unmarshal([]byte(query), &find)
	if err != nil {
		return nil, err
	}
	c := t.mgo.CDb(db, coll)
	var result []bson.M
	err = c.FindAll(nil, find).All(&result)
	return result, err
}
