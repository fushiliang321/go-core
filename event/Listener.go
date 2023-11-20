package event

import (
	"reflect"
	"sync"
)

type (
	EventName    = string
	EventProcess = func(Registered)
	Listen       struct {
		EventNames []EventName
		Process    EventProcess
	}
	eventListener struct {
		sync.RWMutex
		eventMap map[EventName][]EventProcess
	}
)

var globalEventListeners *eventListener

func init() {
	globalEventListeners = &eventListener{
		eventMap: map[EventName][]EventProcess{},
	}
}

// 添加事件监听
func (l *eventListener) Add(listen *Listen) {
	l.Lock()
	defer l.Unlock()
	for _, name := range listen.EventNames {
		if _, ok := l.eventMap[name]; !ok {
			l.eventMap[name] = []func(Registered){}
		}
		l.eventMap[name] = append(l.eventMap[name], listen.Process)
	}
}

// 移除事件监听
func (l *eventListener) Remove(listen *Listen) {
	l.Lock()
	defer l.Unlock()
	for _, name := range listen.EventNames {
		if _, ok := l.eventMap[name]; !ok {
			continue
		}
		for i, process := range l.eventMap[name] {
			if reflect.ValueOf(listen.Process).Pointer() != reflect.ValueOf(process).Pointer() {
				continue
			}

			switch len(l.eventMap[name]) {
			case 1:
				//切片只剩一个值，直接移除整个切片
				delete(l.eventMap, name)
			case i + 1:
				//移除切片最后一个值
				l.eventMap[name] = l.eventMap[name][:i]
			default:
				//移除切片指定位置的值
				l.eventMap[name] = append(l.eventMap[name][:i], l.eventMap[name][i+1:]...)
			}
			break
		}
	}
}

// 触发事件
func (l *eventListener) Trigger(reg *Registered) {
	l.RLock()
	defer l.RUnlock()
	_funs, ok := l.eventMap[reg.name]
	if !ok {
		return
	}
	var funs []EventProcess
	//防止解锁之后_funs的值马上被其他协程修改
	funs = append(funs, _funs...)
	//使用协程调用，避免fun内有添加或移除监听事件导致死锁
	go func() {
		for _, fun := range funs {
			fun(*reg)
		}
	}()
}

// 添加全局事件监听
func AddListener(listen *Listen) {
	globalEventListeners.Add(listen)
}

// 移除全局事件监听
func RemoveListener(listen *Listen) {
	globalEventListeners.Remove(listen)
}
