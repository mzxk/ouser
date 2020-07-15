package ouser

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	User     string             //用户名
	Pwd      string             `json:"-"` //密码s
	RegIP    string             `json:"-"` //注册ip
	Referrer string             `json:"-"` //推荐人

	NickName        string //用户自定义昵称
	Avatar          string //头像ID
	DeliveryAddress string //用户实际地址

	Email string //用户email
	Phone string //用户手机号码

	RealName string //用户真名
	Verify   int64  //用户认证等级

	IsLocked bool     //用户是否锁定
	Group    []string //用户组

	RateLimit int64 `json:"-"` //用户接口全局限制

	GoogleKey string `json:"-"` //用户的googlekey
}

//Feedback 用户反馈的结构
type Feedback struct {
	Bid  string
	Text []FeedbackText
}

//FeedbackText 用户反馈的数组
type FeedbackText struct {
	Time  string
	Admin bool
	Text  string
}

//Avatar 用户头像
type Avatar struct {
	ID     primitive.ObjectID `bson:"_id"`
	Bid    string             `json:"-"`
	Avatar []byte
}
