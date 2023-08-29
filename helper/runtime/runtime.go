package runtime

import (
	"github.com/fushiliang321/go-core/helper/logger"
	"runtime"
	"strings"
	"time"
)

// 当前运行时方法名
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	fn := runtime.FuncForPC(pc[0]).Name()
	countSplit := strings.Split(fn, ".")
	if len(countSplit) > 1 {
		return countSplit[1]
	}
	return ""
}

// 当前运行时包名
func RunPackageName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	fn := runtime.FuncForPC(pc[0]).Name()
	countSplit := strings.Split(fn, "/")
	splitLen := len(countSplit)
	if splitLen > 0 {
		countSplit = strings.Split(countSplit[splitLen-1], ".")
		if len(countSplit) > 0 {
			return countSplit[0]
		}
	}
	return ""
}

// 获取当前日期时间
func Time() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

// 获取当前文件名
//
//go:noinline
func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		logger.Warn("Can not get current file info")
	}
	return file
}
