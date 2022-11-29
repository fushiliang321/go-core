package helper

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// 获取文件内的所有变量名
func GetFileVariateNameAll(FilePath string) (names []string) {
	cache, ok := cacheData["GetFileVariateNameAll"]
	mapData := map[string][]string{}
	if ok {
		if mapData1, ok := cache.(map[string][]string); ok {
			mapData = mapData1
			if names, ok := mapData[FilePath]; ok {
				return names
			}
		}
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, FilePath, nil, 0)
	if err != nil {
		fmt.Printf("err = %s", err)
	}
	for s, object := range f.Scope.Objects {
		if object.Kind == ast.Var {
			names = append(names, s)
		}
	}
	mapData[FilePath] = names
	cacheData["GetFileVariateNameAll"] = mapData
	return
}
