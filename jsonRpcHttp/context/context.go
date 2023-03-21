package context

import "github.com/fushiliang321/go-core/context"

// 获取全部协程上下文数据
func GetAll() map[string]any {
	return context.GetAll()
}

// 获取协程上下文数据
func Get(key string) any {
	return context.Get(key)
}

// 设置协程上下文数据
func Set(key string, value any) {
	context.Set(key, value)
}

// 批量设置协程上下文数据
func SetBatch(values map[string]any) {
	context.SetBatch(values)
}

// 移除协程上下文数据
func Remove(key string) {
	context.Remove(key)
}
