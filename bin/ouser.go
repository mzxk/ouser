package main

import (
	"log"
	"os"

	"github.com/mzxk/ouser"

	"github.com/mzxk/ohttp"
)

func main() {
	log.Println("Reg Router....")
	hh := ohttp.NewWithSession("127.0.0.1:6379", "")
	o := ouser.New(hh)
	h := hh.Group("/user")
	//登陆注册
	h.Add("/registerSimple", o.RegisterSimple) //简单注册
	h.Add("/login", o.Login)                   //简单登录
	h.AddAuth("/logout", o.Logout)             //登出
	//用户设置类
	h.AddAuth("/nicknameSet", o.NickNameSet) //设置显示名
	//用户反馈类
	h.AddAuth("/feedbackPull", o.FeedbackPull) //用户反馈
	h.AddAuth("/feedbackGet", o.FeedbackList)  //用户反馈列表，读取
	//头像
	h.AddAuth("/avatarSet", o.AvatarSet) //设置头像
	h.Add("/avatarGet", o.AvatarGet)     //获取头像

	h.AddAuth("/info", o.UserInfo) //获取用户信息

	//短信类
	h.Add("/smsPublic", o.SmsPublic) //发送短信，这将只能调用注册，登录，找回密码
	h.Add("/smsLogin", o.SmsLogin)   //使用验证码登录 ， 如果不存在，这将新注册账号
	log.Println("Reg Router Done!")

	hh.Run(os.Args[1])
}
