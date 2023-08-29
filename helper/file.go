package helper

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
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

type byModTime []*os.FileInfo

func (fis byModTime) Len() int {
	return len(fis)
}

func (fis byModTime) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis byModTime) Less(i, j int) bool {
	return (*fis[i]).ModTime().Before((*fis[j]).ModTime())
}

// 指定目录下的文件按时间大小排序，从远到近
func SortFile(path string) ([]*os.FileInfo, error) {
	dirEntrys, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files byModTime
	files = make(byModTime, len(dirEntrys))
	j := 0
	for _, dirEntry := range dirEntrys {
		if dirEntry.IsDir() {
			continue
		}
		info, err := dirEntry.Info()
		if err != nil {
			continue
		}
		files[j] = &info
		j++
	}
	files = files[:j]

	sort.Sort(files)
	return files, nil
}
