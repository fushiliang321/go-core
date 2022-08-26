package exception

import (
	"fmt"
	"gitee.com/zvc/go-core/config/exceptions"
	"gitee.com/zvc/go-core/exception/helper"
	"gitee.com/zvc/go-core/exception/types"
	"log"
)

func Listener(mark string, err any) {
	if err == nil {
		return
	}
	go func() {
		if e := recover(); e != nil {
			log.Println("exception listener error", e)
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
