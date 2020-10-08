package ouser

import (
	"log"

	"github.com/mzxk/ohttp"
)

func (t *Ouser) AddressGet(p map[string]string) (interface{}, error) {

	uid := p["bsonid"]
	currency := p["currency"]
	if currency == "" {
		return nil, errs("WrongParams")
	}
	if cfg.AddressURL == "" {
		return [][]string{{currency, "服务未开放"}}, nil
	}
	chains := []string{currency}
	if currency == "usdt" {
		chains = []string{"eth", "trx"}
	}
	var result [][]string
	for _, v := range chains {
		result = append(result, []string{v, t.addressGet(cfg.Name, uid, v)})
	}
	return result, nil
}
func (t *Ouser) addressGet(web, uid, chain string) string {
	var result string
	rsp, err := ohttp.HTTP(cfg.AddressURL, map[string]interface{}{
		"web":   web,
		"user":  uid,
		"chain": chain,
	}).Get()
	if err != nil {
		return "服务未开放"
	}
	err = rsp.JSONSelf(&result)
	if err != nil {
		log.Println("ERROR Get Address", err)
	}
	return result
}
