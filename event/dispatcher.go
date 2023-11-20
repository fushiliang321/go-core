package event

import (
	"sync"
)

type coreServiceEventLog struct {
	sync.RWMutex
	events map[string]bool
}

var (
	CoreServiceEventLog = coreServiceEventLog{
		events: map[string]bool{},
	} //核心服务事件记录
	_dispatch func(reg *Registered) //事件调用方法
)

func init() {
	_dispatch = beforeServerStartDispatch
}

// 判断指定事件是否已经触发过
func (l *coreServiceEventLog) AlreadyTriggered(name string) bool {
	_, ok := CoreServiceEventLog.events[name]
	return ok
}

// 核心服务全部启动前的调用方法
func beforeServerStartDispatch(reg *Registered) {
	if reg.name == AfterServerStart {
		_dispatch = afterServerStartDispatch
	}
	CoreServiceEventLog.Lock()
	defer CoreServiceEventLog.Unlock()

	CoreServiceEventLog.events[reg.name] = true
	globalEventListeners.Trigger(reg)
}

// 核心服务全部启动后的调用方法
func afterServerStartDispatch(reg *Registered) {
	globalEventListeners.Trigger(reg)
}

// 事件调用
func Dispatch(reg *Registered) {
	_dispatch(reg)
}
