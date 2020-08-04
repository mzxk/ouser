package ouser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mzxk/obalance"
	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"github.com/mzxk/oval"
)

var userCache = oval.NewExpire() //用户信息的缓存

type Ouser struct {
	mgo        *omongo.MongoDB
	httpClient *ohttp.Server
	sms        ohttp.Sms
	balance    *obalance.Balance
}

func New(clt *ohttp.Server) *Ouser {
	Init()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	t := &Ouser{
		mgo:        omongo.NewMongoDB(cfg.MongoURL, "user"),
		httpClient: clt,
		sms:        ohttp.NewSms(cfg.Sms.Name, cfg.Sms.Key),
	}
	if cfg.BalanceURL != "" {
		t.balance = obalance.NewRemote(cfg.BalanceURL)
	} else {
		t.balance = obalance.NewLocal(cfg.RedisURL, cfg.RedisPwd)
	}
	return t
}

var idCreate = make(chan string, 5000)

func Init() {
	go func() {
		var begin int64 = 159651286600
		for {
			unix := time.Now().UnixNano()/10000 - begin
			idCreate <- fmt.Sprintf("%x", unix)
			time.Sleep(11 * time.Millisecond)
		}
	}()
	bt, err := ioutil.ReadFile("userConfig.json")
	if err != nil {
		cfg = &Config{
			MongoURL: "mongodb://127.0.0.1:27017",
			RedisURL: "127.0.0.1:6379",
		}
		cfg.Register.Limit = 5
		cfg.Login.Limit = 5

		js, _ := json.Marshal(cfg)
		_ = ioutil.WriteFile("userConfig.json", js, 0777)
		panic(err)
	}
	var cage Config
	err = json.Unmarshal(bt, &cage)
	if err != nil {
		panic(err)
	}
	cfg = &cage
	createIndex(cfg.MongoURL)
}
func createIndex(s string) {
	mgo := omongo.NewMongoDB(s, "user")
	defer mgo.MgoClient.Disconnect(nil)
	log.Println("Start to ensure users index")
	err := mgo.CreateIndexes("user", "user", []string{"u_user_1", "u_phone_1", "u_email_1"})
	err = mgo.CreateIndexes("user", "feedback", []string{"u_bid_1"})
	err = mgo.CreateIndexes("user", "shoplist", []string{"bid_1"})
	err = mgo.CreateIndexes("user", "shopitems", []string{"ban_1"})
	err = mgo.CreateIndexes("user", "withdraw", []string{"bid_1"})
	if err != nil {
		panic(err)
	}
	log.Println("End ensure users index")
}
