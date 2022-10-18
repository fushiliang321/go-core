package context

import "github.com/fushiliang321/go-core/context"

var local = context.NewInstance()

// 获取全部协程上下文数据
func GetAll() map[string]any {
	return local.GetAll()
}

// 获取协程上下文数据
func Get(key string) any {
	return local.Get(key)
}

// 设置协程上下文数据
func Set(key string, value any) {
	local.Set(key, value)
}

// 批量设置协程上下文数据
func SetBatch(values map[string]any) {
	local.SetBatch(values)
}

// 移除协程上下文数据
func Remove(key string) {
	local.Remove(key)
}
