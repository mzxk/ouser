package ouser

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/mzxk/omongo"
)

type Config struct {
	Name       string
	MongoURL   string
	RedisURL   string
	RedisPwd   string
	BalanceURL string
	Register   struct {
		SimpleClosed bool
		Limit        int64
	}
	Login struct {
		Limit int64
	}
	Sms struct {
		Name     string
		Key      string
		Value    string
		IDefault string
		IDs      map[int]string
	}
}

var cfg *Config

func Init() {
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
	if err != nil {
		panic(err)
	}
	log.Println("End ensure users index")
}
