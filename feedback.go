package ouser

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//FeedbackPull 用户反馈
func (t *Ouser) FeedbackPull(p map[string]string) (interface{}, error) {
	txt := p["text"]
	if len(txt) < 5 {
		return nil, errs(ErrParamsWrong)
	}
	c := t.mgo.C("feedback")
	_, err := c.Upsert(nil, bson.M{"bid": p["bsonid"]}, bson.M{"$push": bson.M{"text": FeedbackText{
		Time:  time.Now().Format("2006-01-02 15:04:05"),
		Text:  txt,
		Type:  p["type"],
		Title: p["title"],
	}}})
	return nil, err
}

//FeedbackList ß用户反馈
func (t *Ouser) FeedbackList(p map[string]string) (interface{}, error) {
	var result Feedback
	err := t.mgo.C("feedback").FindOne(nil,
		bson.M{"bid": p["bsonid"]},
		options.FindOne().SetProjection(bson.M{
			"text": bson.M{
				"$slice": bson.A{-10, 10},
			},
		}),
	).Decode(&result)
	return result, err
}
