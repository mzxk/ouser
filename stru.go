package ouser

import "go.mongodb.org/mongo-driver/bson/primitive"

//User 最主要的用户结构
type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"-"`
	Uid      string             `json:"ID"`
	User     string             `bson:"user,omitempty"` //用户名
	Pwd      string             `json:"-"`              //密码s
	Paypwd   string             `json:"-"`              //支付密码
	RegIP    string             `json:"-"`              //注册ip
	Referrer string             `json:"-"`              //推荐人

	NickName        string //用户自定义昵称
	Avatar          string //头像ID
	DeliveryAddress string //用户实际地址

	Email string `bson:"email,omitempty"` //用户email
	Phone string `bson:"phone,omitempty"` //用户手机号码

	RealName string //用户真名
	Verify   int64  //用户认证等级

	Locked         bool  //用户是否锁定
	LockedBalance  bool  //用户余额锁定
	LockedWithdraw int64 //用户提现锁定，代表锁定到的unix时间，在此时间之前不允许提现

	Group map[string]string //用户组

	RateLimit int64 `json:"-"` //用户接口全局限制

	GoogleKey string `json:"-"` //用户的google key
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
	Title string
	Text  string
	Type  string
}

//Avatar 用户头像
type Avatar struct {
	ID     primitive.ObjectID `bson:"_id"`
	Bid    string             `json:"-"`
	Avatar []byte
}
