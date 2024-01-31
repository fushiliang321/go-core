package context

import "github.com/timandy/routine"

type Context struct {
	Local routine.ThreadLocal[map[string]any]
}

func NewInstance() *Context {
	return &Context{
		Local: routine.NewInheritableThreadLocal[map[string]any](),
	}
}

// 获取全部协程上下文数据
func (ctx *Context) GetAll() map[string]any {
	if data := ctx.Local.Get(); data != nil {
		return data
	}
	return nil
}

// 获取协程上下文数据
func (ctx *Context) Get(key string) any {
	if data := ctx.Local.Get(); data != nil {
		if value, ok := data[key]; ok {
			return value
		}
	}
	return nil
}

// 设置协程上下文数据
func (ctx *Context) Set(key string, value any) {
	data := ctx.Local.Get()
	if data == nil {
		data = map[string]any{}
	}
	data[key] = value
	ctx.Local.Set(data)
}

// 批量设置协程上下文数据
func (ctx *Context) SetBatch(values map[string]any) {
	if values == nil {
		return
	}
	data := ctx.Local.Get()
	if data == nil {
		data = map[string]any{}
	}
	for key, value := range values {
		data[key] = value
	}
	ctx.Local.Set(data)
}

// 移除协程上下文数据
func (ctx *Context) Remove(key string) {
	data := ctx.Local.Get()
	if data == nil {
		return
	}
	delete(data, key)
	if len(data) > 0 {
		ctx.Local.Set(data)
		return
	}
	ctx.Local.Remove()
}
