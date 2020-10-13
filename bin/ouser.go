package main

import (
	"log"
	"os"

	"github.com/mzxk/ouser"

	"github.com/mzxk/ohttp"
)

func main() {
	log.Println("Reg Router....")
	hh := ohttp.NewWithSession(os.Getenv("redisURL"), os.Getenv("redisPWD"))
	o := ouser.New(hh)
	h := hh.Group("/user")

	//登陆注册
	h.Add("/registerSimple", o.RegisterSimple) //简单注册
	h.Add("/login", o.Login)                   //简单登录
	h.AddAuth("/logout", o.Logout)             //登出

	//用户设置类
	h.AddAuth("/nicknameSet", o.NickNameSet)     //设置显示名
	h.AddAuth("/paypwdSet", o.PaypwdSet)         //设置支付密码
	h.AddAuth("/contactChange", o.ContactChange) //用户换绑手机
	h.AddAuth("/pwdChange", o.PwdChange)

	//用户反馈类
	h.AddAuth("/feedbackPull", o.FeedbackPull) //用户反馈
	h.AddAuth("/feedbackGet", o.FeedbackList)  //用户反馈列表，读取

	//头像
	h.AddAuth("/avatarSet", o.AvatarSet) //设置头像
	h.Add("/avatarGet", o.AvatarGet)     //获取头像

	//信息获取类
	h.AddAuth("/info", o.UserInfo) //获取用户信息

	//短信类
	h.Add("/smsPublic", o.SmsPublic)       //发送短信，这将只能调用注册，登录，找回密码
	h.Add("/smsLogin", o.SmsLogin)         //使用验证码 登录 ， 如果不存在，这将新注册账号 , 同时重置密码也是这个接口。
	h.AddAuth("/smsPrivate", o.SmsPrivate) //发送用户本身的短信接口

	//账务类
	h.AddAuth("/withdraw", o.Withdraw)            //用户提现
	h.AddAuth("/withdrawAccept", o.WithdrawCheck) //用户确认提款
	h.AddAuth("/withdrawGet", o.WithdrawGet)
	h.AddAuth("/addressGet", o.AddressGet)
	h.AddAuth("/balanceGet", o.BalanceGet)

	//google验证码
	h.AddAuth("/g2faCreate", o.G2faCreate)
	h.AddAuth("/g2faAccept", o.G2faAccept)

	//购物类
	h.Add("/shopItemsGet", o.ShopItemsGet) //获取列表
	h.Add("/shopDiscount", o.ShopDiscount) //查看购买折扣
	h.AddAuth("/shopBuy", o.ShopBuy)       //购买
	h.AddAuth("/shopList", o.ShopList)     //查看自己的购买列表

	h.AddAuth("/rebateGet", o.RebateGet)

	//用户类
	h.AddAuth("/admin/get", o.AdminGet)
	h.AddAuth("/admin/info", o.AdminUserInfo)
	h.AddAuth("/admin/checkGP", o.AdminCheckGP)
	log.Println("Reg Router Done!")

	hh.Run(os.Args[1])
}
