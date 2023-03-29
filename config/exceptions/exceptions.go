package exceptions

import (
	"github.com/fushiliang321/go-core/exception/types"
)

type (
	Exceptions struct {
		Handlers []ExceptionHandler
	}
	ExceptionHandler interface {
		Handle(*types.Runtime)
	}
)

var exceptions = &Exceptions{
	Handlers: []ExceptionHandler{},
}

func Set(config *Exceptions) {
	exceptions = config
}

func Get() *Exceptions {
	return exceptions
}
