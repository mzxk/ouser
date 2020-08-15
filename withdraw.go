package ouser

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mzxk/ohttp"
	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Withdraw struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             //deposit或withdraw
	Bid         string
	Time        string  //时间
	Height      int64   //高度
	Txid        string  //交易id
	From        string  //从哪里来
	To          string  //到哪里去
	ToMemo      string  //转账备注（可能有特殊用途）
	Amount      float64 //数量
	Currency    string  //币种
	State       string  //状态名称
	VerifyRisk  bool    //风控是否通过
	VerifyAdmin bool    //人工是否通过
	FeeAmount   float64 //系统收取的手续费

	Confirm string //确认次数 通常为当前次数/需要次数
}

func (t *Ouser) WithdrawCheck(p map[string]string) (interface{}, error) {
	action := p["action"]
	id := p["id"]
	if id == "" {
		return nil, errors.New(ErrParamsWrong)
	}
	var wr Withdraw
	wrString := ohttp.RedisGet(id)
	err := json.Unmarshal([]byte(wrString), &wr)
	if err != nil {
		return nil, err
	}
	//确认用户提现，需要验证id，用户id和状态，设置用户确认为true以及状态为用户确认
	if action == "accept" {
		//扣款
		_, err = t.balance.New("withdraw", wr.Currency, wr.ID.Hex()).Lock(wr.Bid, wr.Amount).Run()
		if err != nil {
			return nil, err
		}
		//写入提现记录
		err = wr.Save(t.mgo)
		//写入出现问题，回退余额
		if err != nil {
			_, _ = t.balance.New("withdraw", wr.Currency, wr.ID.Hex()).UnLock(wr.Bid, wr.Amount).Run()
			return nil, err
		}
		return nil, err
	}
	return nil, nil
}

//WithdrawGet 获取提现记录
func (t *Ouser) WithdrawGet(p map[string]string) (interface{}, error) {
	var wr []Withdraw
	c := t.mgo.C("withdraw")
	err := c.FindAll(nil, bson.M{"bid": p["bsonid"]},
		options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100),
	).All(&wr)
	return wr, err
}

//Withdraw 这进行用户提现
func (t *Ouser) Withdraw(p map[string]string) (interface{}, error) {
	//判断验证码
	if err := t.checkPrivateCode(p, smsWithdraw); err != nil {
		return nil, err
	}
	//构建提款操作
	wr, err := t.withdrawCreate(p["bsonid"], p["currency"])
	if err != nil {
		return nil, err
	}
	wr, err = wr.Withdraw(p["to"], p["memo"], s2f(p["amount"]))
	if err != nil {
		return nil, err
	}
	js, _ := json.Marshal(wr)
	ohttp.RedisSet(wr.ID.Hex(), string(js), 750)

	//返回提现后余额
	return wr, nil
}
func (t *Ouser) withdrawCreate(bid, currency string) (*Withdraw, error) {
	if bid == "" || currency == "" {
		return nil, errors.New(ErrParamsWrong)
	}
	return &Withdraw{
		ID:       omongo.ID(""),
		Name:     "withdraw",
		Bid:      bid,
		Currency: currency,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
func (t *Withdraw) Withdraw(to, memo string, amount float64) (*Withdraw, error) {
	if to == "" || amount <= 0 {
		return nil, errors.New(ErrParamsWrong)
	}
	t.To = to
	t.ToMemo = memo
	t.Amount = amount
	return t, nil
}
func (t *Withdraw) Save(mgo *omongo.MongoDB) error {
	c := mgo.C("withdraw")
	_, err := c.InsertOne(nil, t)
	return err
}
