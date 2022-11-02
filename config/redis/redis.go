package redis

type Redis struct {
	Host     string
	Port     int
	Password string
	Db       int
	Options  map[string]any
}

var data = &Redis{}

func Set(config *Redis) {
	data = config
}

func Get() *Redis {
	return data
}
