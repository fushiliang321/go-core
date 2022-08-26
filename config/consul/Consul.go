package consul

import (
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	api.Config
	Consumers []string
	Services  []any
}

var consul = &Consul{
	Consumers: []string{},
	Services:  []any{},
}

func Set(config *Consul) {
	consul = config
}

func Get() *Consul {
	return consul
}
