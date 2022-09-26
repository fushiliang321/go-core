package consul

import (
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	api.Config
}

var consul = &Consul{}

func Set(config *Consul) {
	consul = config
}

func Get() *Consul {
	return consul
}
