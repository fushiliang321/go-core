package exception

import (
	"fmt"
	"github.com/fushiliang321/go-core/config/exceptions"
	"github.com/fushiliang321/go-core/exception/helper"
	"github.com/fushiliang321/go-core/exception/types"
	"github.com/fushiliang321/go-core/logger"
)

func Listener(mark string, err any) {
	if err == nil {
		return
	}
	go func() {
		if e := recover(); e != nil {
			logger.Error("exception listener error", e)
		}
	}()
	handle(&types.Runtime{
		Msg:   fmt.Sprint(err),
		Mark:  mark,
		Trace: helper.Trace(3),
	})
}

func handle(runtime *types.Runtime) {
	config := exceptions.Get()
	for _, handler := range config.Handlers {
		handler.Handle(runtime)
	}
}
