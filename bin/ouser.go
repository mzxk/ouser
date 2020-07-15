package main

import (
	"log"
	"os"

	"github.com/mzxk/ouser"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
)

var mgo *omongo.MongoDB

func main() {
	//mgo = omongo.NewMongoDB("mongodb://192.168.1.3:27017", "user")
	mgo = omongo.NewMongoDB("mongodb://root:root@172.31.39.207:13198,172.31.39.208:13198/?authSource=admin&replicaSet=nonomongo", "user")
	log.Println("Start to ensure users index")
	err := mgo.CreateIndexes("user", "user", []string{"u_user_1", "u_phone_1", "u_email_1"})
	err = mgo.CreateIndexes("user", "feedback", []string{"u_bid_1"})
	if err != nil {
		panic(err)
	}
	log.Println("End ensure users index")

	log.Println("Reg Router....")
	h := ohttp.NewWithSession("127.0.0.1:6379", "")
	o := ouser.New(mgo, h)
	//登陆注册
	h.Add("/user/registerSimple", o.RegisterSimple) //简单注册
	h.Add("/user/login", o.Login)                   //简单登录
	h.AddAuth("/user/logout", o.Logout)             //登出
	//用户设置类
	h.AddAuth("/user/nicknameSet", o.NickNameSet) //设置显示名
	//用户反馈类
	h.AddAuth("/user/feedbackPull", o.FeedbackPull) //用户反馈
	h.AddAuth("/user/feedbackGet", o.FeedbackList)  //用户反馈列表，读取
	//头像
	h.AddAuth("/user/avatarSet", o.AvatarSet)
	h.Add("/user/avatarGet", o.AvatarGet)
	h.Run(os.Args[1])
}
