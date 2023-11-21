package service

import "sync"

type (
	Service interface {
		Start(wg *sync.WaitGroup)
		PreEvents() []string //前置事件
	}

	BaseStruct struct{}
)

func (ser *BaseStruct) Start(wg *sync.WaitGroup) {}
func (ser *BaseStruct) PreEvents() []string {
	return nil
}
