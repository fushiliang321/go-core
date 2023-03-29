package context

import "github.com/timandy/routine"

type Context struct {
	Local routine.ThreadLocal
}

func NewInstance() *Context {
	return &Context{
		Local: routine.NewInheritableThreadLocal(),
	}
}

// 获取全部协程上下文数据
func (ctx *Context) GetAll() map[string]any {
	if data := ctx.Local.Get(); data != nil {
		if data, ok := data.(map[string]any); ok {
			return data
		}
	}
	return nil
}

// 获取协程上下文数据
func (ctx *Context) Get(key string) any {
	if data := ctx.Local.Get(); data != nil {
		if data, ok := data.(map[string]any); ok {
			if value, ok := data[key]; ok {
				return value
			}
		}
	}
	return nil
}

// 设置协程上下文数据
func (ctx *Context) Set(key string, value any) {
	_map := map[string]any{}
	if data := ctx.Local.Get(); data != nil {
		if data, ok := data.(map[string]any); ok {
			_map = data
		}
	}
	_map[key] = value
	ctx.Local.Set(_map)
}

// 批量设置协程上下文数据
func (ctx *Context) SetBatch(values map[string]any) {
	_map := map[string]any{}
	if data := ctx.Local.Get(); data != nil {
		if data, ok := data.(map[string]any); ok {
			_map = data
		}
	}
	for key, value := range values {
		_map[key] = value
	}
	ctx.Local.Set(_map)
}

// 移除协程上下文数据
func (ctx *Context) Remove(key string) {
	if data := ctx.Local.Get(); data != nil {
		if data, ok := data.(map[string]any); ok {
			delete(data, key)
			if len(data) > 0 {
				ctx.Local.Set(data)
				return
			}
		}
		ctx.Local.Remove()
	}
}
