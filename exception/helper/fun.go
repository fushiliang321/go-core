package helper

import (
	"fmt"
	"runtime"
)

func Trace(skip int) (data []string) {
	pcs := make([]uintptr, 10)
	if skip < 0 {
		skip = 2
	}
	n := runtime.Callers(skip, pcs)
	for i := 0; i < n; i++ {
		pc := pcs[i]
		fn := runtime.FuncForPC(pc)
		fname := fn.Name()
		file, line := fn.FileLine(pc)
		data = append(data, fmt.Sprintf("%s:%d %s", file, line, fname))
	}
	return
}
