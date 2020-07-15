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
	mgo = omongo.NewMongoDB("mongodb://127.0.0.1:27017", "user")
	//mgo = omongo.NewMongoDB("mongodb://root:root@172.31.39.207:13198,172.31.39.208:13198/?authSource=admin&replicaSet=nonomongo", "user")
	log.Println("Start to ensure users index")
	mgo.CreateIndexes("user", "user", []string{"u_user_1", "u_phone_1", "u_email_1"})
	log.Println("End ensure users index")

	h := ohttp.NewWithSession("127.0.0.1:6379", "")
	//登陆注册
	h.Add("/user/registerSimple", ouser.RegisterSimple)
	h.Add("/user/login", ouser.Login)
	h.AddAuth("/user/logout", ouser.Logout)
	//用户设置类
	h.AddAuth("/user/setNickname", ouser.SetNickName)
	//用户反馈类
	h.AddAuth("/user/feedback", ouser.Feedback)
	h.Run(os.Args[1])
}
