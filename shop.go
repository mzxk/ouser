package ouser

import (
	"log"
	"time"

	"github.com/mzxk/omongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopItem struct {
	ID            primitive.ObjectID `bson:"_id"`
	Title         string             //标题
	SubTitle      string             //副标题
	Price         float64            //价格
	OriginalPrice float64            //原始价格
	ImgURL        string             //图片地址
	Tag           []string           //商品标签
	Currency      string             //计价单位
	Rebate        float64            //返利
	Discount      map[string]float64 `json:"-"` //折扣码
	Detail        interface{}        //商品详情
	discountPrice float64
}
type ShopList struct {
	ID           primitive.ObjectID `bson:"_id"`
	Uid          string
	ItemID       string
	Title        string
	Price        float64
	Amount       int64
	Currency     string
	TotalPrice   float64
	Discount     float64
	DiscountCode string
}

func (t *Ouser) ShopItemsGet(_ map[string]string) (interface{}, error) {
	c := t.mgo.C("shopitems")
	var result []ShopItem
	err := c.FindAll(nil, bson.M{}).All(&result)
	return result, err
}
func (t *Ouser) ShopBuy(p map[string]string) (interface{}, error) {
	item, err := t.ShopDiscount(p)
	if err != nil {
		return nil, err
	}
	var amount int64 = 1
	amt := s2i(p["amount"])
	if amt > 1 {
		amount = amt
	}
	shopItem := item.(ShopItem)
	sl := ShopList{
		ID:           omongo.ID(""),
		Uid:          p["bsonid"],
		ItemID:       shopItem.ID.Hex(),
		Title:        shopItem.Title,
		Price:        shopItem.Price,
		Discount:     shopItem.discountPrice,
		DiscountCode: p["discountCode"],
		Amount:       amount,
		Currency:     shopItem.Currency,
		TotalPrice:   float64(amount) * shopItem.Price,
	}
	_, err = t.balance.New("buy", shopItem.Currency, sl.ID.Hex()).DecrAvail(p["bsonid"], sl.Price).Run()
	if err != nil {
		return nil, err
	}
	c := t.mgo.C("shoplist")
	_, err = c.InsertOne(nil, sl)
	t.shopRebate(p, sl, shopItem.Rebate)
	return nil, err
}

//ShopList 获取当前用户的购买列表
func (t *Ouser) ShopList(p map[string]string) (interface{}, error) {
	c := t.mgo.C("shoplist")
	var result []ShopList
	err := c.FindAll(nil, bson.M{"uid": p["bsonid"]}).All(&result)
	return result, err
}

//ShopDiscount 获取折扣价格，需要输入的有商品id和折扣码
//id-discountCode
func (t *Ouser) ShopDiscount(p map[string]string) (interface{}, error) {
	id := p["id"]
	code := p["discountCode"]
	if id == "" {
		return nil, errs(ErrParamsWrong)
	}
	c := t.mgo.C("shopitems")
	var result ShopItem
	err := c.FindOne(nil, bson.M{"_id": omongo.ID(id)}).Decode(&result)
	if err != nil {
		return nil, err
	}
	//
	if code == "" {
		return result, nil
	}
	if result.Discount != nil && result.Discount[code] != 0 {
		dis := result.Discount[code]
		var reply float64
		if dis < 0 { //如果折扣比0低，那么实际价格就是当前价格-折扣
			reply = result.Price + dis
		}
		if dis > 0 && dis < 1 { //如果折扣比0高比1低，那么就是当前价格*折扣
			reply = dis * result.Price
		}
		//上面低减法可能造成价格低于0的情况,只要reply比0大，那么当前价格就是reply
		if reply > 0 {
			result.discountPrice = result.Price - reply
			result.Price = reply
		}
	}
	return result, nil
}

type Rebate struct {
	Uid      string
	Amount   map[string]float64
	Invitees int64
	List     []RebateList
}

func (t *Ouser) shopRebate(p map[string]string, sl ShopList, re float64) {
	//返现比例为0，直接返回
	if re < 0.0001 {
		return
	}
	re = sl.TotalPrice * re
	usr, err := t.userCache(p)
	if err != nil || usr == nil {
		log.Println("???", p["bsonid"], err)
	}
	if usr.Referrer == "" {
		return
	}
	lst := RebateList{
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		ID:       sl.ID.Hex(),
		Name:     "购物返利",
		Uid:      usr.ID.Hex(),
		Currency: sl.Currency,
		Amount:   re,
	}
RE:
	m := t.mgo.C("rebate")
	_, err = m.Upsert(nil,
		bson.M{"uid": usr.Referrer},
		bson.M{
			"$inc":  bson.M{"amount." + sl.Currency: re},
			"$push": bson.M{"list": lst},
		})
	if err != nil {
		log.Println(err)
		time.Sleep(1 * time.Second)
		goto RE
	}
	ba := t.balance.New("rebate", lst.Currency, sl.ID.Hex())
	_, err = ba.IncrAvail(usr.Referrer, re).Run()
	if err != nil {
		log.Println(err)
	}
}

type RebateList struct {
	Time     string  //时间
	ID       string  //增加这笔记录的事件ID
	Name     string  //显示名称
	Uid      string  //这笔来源用户ID
	Currency string  //币种
	Amount   float64 //数量
}

func (t *Ouser) RebateGet(p map[string]string) (interface{}, error) {
	m := t.mgo.C("rebate")
	var result Rebate
	err := m.FindOne(nil, bson.M{"uid": p["bsonid"]}).Decode(&result)
	return result, err
	//r := &Rebate{
	//	Uid:      "test",
	//	Amount:   map[string]float64{"T": 100.0},
	//	Invitees: 10,
	//	List: []RebateList{{
	//		Time:     "2020-01-01 12:33:45",
	//		ID:       "idididididididid",
	//		Name:     "用户购买矿机",
	//		Uid:      "uiduiduid",
	//		Currency: "T",
	//		Amount:   100.0,
	//	}},
	//}
	//return r, nil
}
