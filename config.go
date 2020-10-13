package ouser

type Config struct {
	Name       string
	MongoURL   string
	RedisURL   string
	RedisPwd   string
	BalanceURL string
	AddressURL string
	OnlyGoogle bool
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
