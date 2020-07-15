package ouser

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	User     string             //用户名
	Pwd      string             `json:"-"` //密码s
	RegIP    string             `json:"-"` //注册ip
	Referrer string             `json:"-"` //推荐人

	NickName        string //用户自定义昵称
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
