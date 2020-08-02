package ouser

import (
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
	return nil, err
}
func (t *Ouser) ShopList(p map[string]string) (interface{}, error) {
	c := t.mgo.C("shoplist")
	var result []ShopList
	err := c.FindAll(nil, bson.M{"uid": p["bsonid"]}).All(&result)
	return result, err
}
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
