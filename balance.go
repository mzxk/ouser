package ouser

//BalanceGet 获取用户余额
func (t *Ouser) BalanceGet(p map[string]string) (interface{}, error) {
	return t.balance.GetBalance(p["bsonid"])
}
