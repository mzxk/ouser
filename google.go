package ouser

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Google2fa struct{}

//G2faCreate 生成google验证码的key，这个接口没有限制，因为没有成本。
func (t *Ouser) G2faCreate(p map[string]string) (interface{}, error) {
	info, err := t.userCache(p)
	if err != nil {
		return nil, err
	}
	k, v := new(Google2fa).CreateKey(info.User, cfg.Name)
	ohttp.RedisSet(info.User+"google2fa", v, 15*60)
	return []string{k, v}, nil
}

//G2faAccept 设置用户的googlekey，这将验证支付密码，google密码和手机验证码
func (t *Ouser) G2faAccept(p map[string]string) (interface{}, error) {
	if err := t.checkPayPwd(p); err != nil {
		return nil, err
	}
	god := ohttp.RedisGet(p["bsonid"] + "google2fa")
	if new(Google2fa).Check2fa(p["googleCode"], god) {
		return nil, errs(ErrGoogle2fa)
	}
	if err := t.checkPrivateCode(p, smsGoogle2fa); err != nil {
		return nil, err
	}
	c := t.mgo.C("user")
	_, err := c.UpdateOne(nil,
		bson.M{"_id": omongo.ID(p["bsonid"])},
		bson.M{"$set": bson.M{"googlekey": god}})
	return nil, err
}

//Check2fa .
func (f *Google2fa) Check2fa(intuser, key string) bool {
	result, err := f.MakeGoogleAuthenticatorForNow(key)
	if err == nil && result == intuser {
		return true
	}
	now := time.Now().Unix()
	//当前时间前后10秒的验证码均予以通过
	result2, err := f.MakeGoogleAuthenticator(key, now-10)
	if err == nil && result2 == intuser {
		return true
	}
	result3, err := f.MakeGoogleAuthenticator(key, now+10)
	if err == nil && result3 == intuser {
		return true
	}
	return false
}

//CreateKey .
func (f *Google2fa) CreateKey(user, web string) (url string, key string) {
	key = f.createKey(user)
	url = "otpauth://totp/" + user + "?secret=" + key + "&issuer=" + web
	return url, key
}
func (f *Google2fa) createKey(s string) string {
	sha := sha1.New()
	sha.Write([]byte(s + time.Now().String() + "普通文本二维码的特点"))
	result := base32.StdEncoding.EncodeToString(sha.Sum(nil))
	return result
}

// MakeGoogleAuthenticator 获取key&t对应的验证码
// key 秘钥
// t 1970年的秒
func (f *Google2fa) MakeGoogleAuthenticator(key string, t int64) (string, error) {
	hs, e := f.hmacSha1(key, t/30)
	if e != nil {
		return "", e
	}
	snum := f.lastBit4byte(hs)
	d := snum % 1000000
	return fmt.Sprintf("%06d", d), nil
}

// MakeGoogleAuthenticatorForNow 获取key对应的验证码
func (f *Google2fa) MakeGoogleAuthenticatorForNow(key string) (string, error) {
	return f.MakeGoogleAuthenticator(key, time.Now().Unix())
}

func (f *Google2fa) lastBit4byte(hmacSha1 []byte) int32 {
	if len(hmacSha1) != sha1.Size {
		return 0
	}
	offsetBits := int8(hmacSha1[len(hmacSha1)-1]) & 0x0f
	p := (int32(hmacSha1[offsetBits]) << 24) | (int32(hmacSha1[offsetBits+1]) << 16) | (int32(hmacSha1[offsetBits+2]) << 8) | (int32(hmacSha1[offsetBits+3]) << 0)
	return p & 0x7fffffff
}

func (f *Google2fa) hmacSha1(key string, t int64) ([]byte, error) {
	decodeKey, err := base32.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	cData := make([]byte, 8)
	binary.BigEndian.PutUint64(cData, uint64(t))

	h1 := hmac.New(sha1.New, decodeKey)
	_, e := h1.Write(cData)
	if e != nil {
		return nil, e
	}
	return h1.Sum(nil), nil
}
