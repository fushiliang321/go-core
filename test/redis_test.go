package test

import (
	"fmt"
	"github.com/fushiliang321/go-core"
	redis2 "github.com/fushiliang321/go-core/config/redis"
	"github.com/fushiliang321/go-core/helper/file"
	"github.com/fushiliang321/go-core/helper/logger"
	logger2 "github.com/fushiliang321/go-core/logger"
	"github.com/fushiliang321/go-core/redis"
	"reflect"
	"strings"
	"testing"
)

const b = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	redis2.Set(&redis2.Redis{
		Host:     "192.168.31.10",
		Port:     6379,
		Password: "zct123",
	})
}

func TestSet(t *testing.T) {
	redis.Set("test", "asdfasdf")
}
func TestGet(t *testing.T) {
	res, _ := redis.Get[string]("test")
	fmt.Println(*res, reflect.TypeOf(res).String())
}

func TestRandString(t *testing.T) {
	//fmt.Println(helper.RandString(70, helper.RangeLetter0))
	//str := helper.RandString(70, helper.RangeLetter0)
	//fmt.Println(str)
	//for i := 0; i < 10000000; i++ {
	//	//helper.RandString(70)
	//	//helper.RandString(70, helper.RangeLetter0)
	//	json.Marshal(str)
	//}
	for i := 0; i < 10000000; i++ {
		s := "data->content"
		jsonFieldNameTransition(&s)
	}
}

// json字段名称转换 ->格式转为 ->'$.'
func jsonFieldNameTransition(filedName *string) {
	i := strings.Index(*filedName, "->")
	if i < 1 {
		return
	}
	builder := strings.Builder{}
	builder.WriteString((*filedName)[:i+2])
	builder.WriteString("'$")
	builder.WriteString(strings.Replace((*filedName)[i:], "->", ".", -1))

	//filedNameSplit := strings.Split(*filedName, "->")
	//builder.WriteString(filedNameSplit[0])
	//builder.WriteString("->'$")
	//
	//for i := 1; i < len(filedNameSplit); i++ {
	//	builder.WriteString(".")
	//	builder.WriteString(filedNameSplit[i])
	//}
	builder.WriteString("'")
	*filedName = builder.String()
}

func Test1(t *testing.T) {
	s := logger2.Service{}
	s.Start(nil)
	file, err := file.SortFile("D:\\代码")
	if err != nil {
		return
	}

	logger.Warn("ok")
	logger.Warn(len(file))
	logger.Debug("debug")
	logger.Info("info")
	for _, info := range file {
		logger.Warn(1)
		logger.Info(info)
	}
}

func TestStart(t *testing.T) {
	core.Start()
}
