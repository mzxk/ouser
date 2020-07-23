# ouser
## 签名流程
```
Sign的方式：
		1. 登陆，获取后端返回的key和value。
				例如：key = foo , value = bar
				
		2. 在正常的请求上url，增加两个参数：
			nonce	时间戳 精确到秒
			key		login后后端返回的key，value之一的key
				例如：
					原始请求是http://localhost/user/changePwd?oldpwd=old&newpwd=new
					增加后变成http://localhost/user/changePwd?oldpwd=old&newpwd=new&key=foo&nonce=158754624
		
		3：计算sign
			使用整个requestURI 不包含域名和端口的部分 加上 value 进行sha512计算 ， 然后截取[58:98]
				例如：
					needsign:="/user/changePwd?oldpwd=old&newpwd=new&key=foo&nonce=158754624"+"bar"
					sign=sha512(needsign)[58:98]
					"bf11206c32dc27301f681da41fc3570937c1e2cd"
			然后在请求request里增加一个header，名字为key,值为上面的sign

		请注意：
			如果需要加密的字符串里含有中文字符,那么一定要在url被转码成%ab的格式后在进行签名，直接使用中文拼接进行的签名是无法通过检查的
```
## 返回值
```
    {
        "m":"ok",               //接口正常，返回ok，不然返回具体原因
        "r":null,               //接口正常的返回内容，参考每个接口
        "t":1598542145          //服务器时间戳(秒)
    }
```

## 发送验证码 1
### 不登录发送的验证码 1-1
* /user/smsPublic
* params
    * type          发送代码        6001｜6002｜6005 注册登录｜重置密码｜绑定手机
    * contact       联系方式        手机号码｜信箱     后端自动识别
    
### 登录后发送验证码 1-2
* /user/smsPrivate
* 只能发送用户自己的手机号码
* params
    * type          发送代码        具体如下
* 发送代码
    7001    设置支付密码
    7002    修改手机号码
## 注册/登录/找回密码

### 注册 2 

#### 简单注册 2-1
* /user/registerSimple
* params
    * user          用户名      
    * pwd           密码
    * referrer      推荐人      可选
    * paypwd        支付密码
* return
    * Key           用户的key
    * value         用户的value
    
#### 手机号码注册/登录 2-2
* /user/smsLogin
* 登录注册是一体的，如果用户不存在，直接注册，存在则登录
* 首先发送一个公开接口的验证码 6001
* params
    * contact       手机号码 必须和发送验证码的手机号码一致
    * codePublic    手机收到的验证码
    * contactType   是手机还是信箱     phone|email 如果留空默认手机
    * pwd           密码          可选
    * referrer      推荐人        可选
    * paypwd        支付密码       可选
* return
    * Key           用户的key
    * value         用户的value
    
#### 用户简单登录 2-3
* /user/login
* params
    * user      用户名
    * pwd       密码
* return
    * Key           用户的key
    * value         用户的value
    
#### 用户登出
* /user/logout
* 需要签名

### 设置类 3

#### 用户修改手机号 3-1
* 先发送私有接口的验证码
* 发送共有接口的验证码
* /user/contactChange
* 需要签名
* params
    * contact       新的手机号码
    * contactType   email｜phone
    * code          用户老手机号码收到的验证码
    * codePublic    用户新手机号码的验证码

#### 用户设置支付密码 3-2
* 发送私有接口的验证码
* 需要签名
* /user/paypwdSet
* params
    * paypwd        新的支付密码
    * code          用户收到的验证码
    
   