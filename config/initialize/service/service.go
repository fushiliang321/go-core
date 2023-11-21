package service

import "sync"

type (
	Service interface {
		Start(wg *sync.WaitGroup)
		PreServices() []string //前置服务
	}

	BaseStruct struct {
		preServices []string
	}
)

func (ser *BaseStruct) Start(wg *sync.WaitGroup) {}
func (ser *BaseStruct) PreServices() []string {
	return ser.preServices
}
