package context

import "github.com/timandy/routine"

var local = routine.NewInheritableThreadLocal()

// 获取协程上下文数据
func Get(key string) any {
	data := local.Get()
	if data != nil {
		if data, ok := data.(map[string]any); ok {
			if value, ok := data[key]; ok {
				return value
			}
		}
	}
	return nil
}

// 设置协程上下文数据
func Set(key string, value any) {
	_map := map[string]any{}
	data := local.Get()
	if data != nil {
		if data, ok := data.(map[string]any); ok {
			_map = data
		}
	}
	_map[key] = value
	local.Set(_map)
}

// 移除协程上下文数据
func Remove(key string) {
	data := local.Get()
	if data != nil {
		if data, ok := data.(map[string]any); ok {
			delete(data, key)
			if len(data) > 0 {
				local.Set(data)
				return
			}
		}
		local.Remove()
	}
}
