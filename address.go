package ouser

import "github.com/mzxk/ohttp"

func (t *Ouser) AddressGet(p map[string]string) (interface{}, error) {
	uid := p["bsonid"]
	//TODO 未实装多币种
	//coin := p["coin"]
	//if coin == "usdt" {
	//
	//}
	var result string
	rsp, err := ohttp.HTTP(cfg.AddressURL, map[string]interface{}{
		"web":   cfg.Name,
		"user":  uid,
		"chain": "eth",
	}).Get()
	if err != nil {
		return nil, err
	}
	err = rsp.JSONSelf(&result)
	if err != nil {
		return nil, err
	}

	return [][]string{
		{"eth", result},
	}, nil
}
