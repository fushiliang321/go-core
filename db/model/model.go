package model

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strings"
)

type Model[t any] struct {
	Db *gorm.DB
}

type PaginateData[t any] struct {
	Page     int `json:"page"`
	LastPage int `json:"last_page"`
	Limit    int `json:"limit"`
	Total    int `json:"total"`
	Lists    []t `json:"lists"`
}

var operatorMap = map[string]byte{"=": 0, "<": 0, "<=": 0, ">": 0, ">=": 0, "<>": 0, "<=>": 0, "!=": 0, "like": 0, "not like": 0, "in": 0, "not in": 0, "between": 0, "not between": 0}

func filter(str string) string {
	reg, err := regexp.Compile(`[^a-zA-Z0-9_*.\-\[\]>]+`)
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(str, "")
}

func (m *Model[t]) DB() *gorm.DB {
	return m.Db
}

func (m *Model[t]) SetDB(db *gorm.DB) *Model[t] {
	m.Db = db
	return m
}

func (m *Model[t]) Where(where any, args ...any) *Model[t] {
	if where == nil {
		return m
	}
	if args == nil {
		switch w := where.(type) {
		case map[string]any:
			return m._where(w)
		case *map[string]any:
			return m._where(*w)
		}
	}
	m.Db = m.Db.Where(where, args...)
	return m
}

func (m *Model[t]) _where(where map[string]any) *Model[t] {
	var (
		operator string
		k        string
		v        any
	)

	for k, v = range where {
		var value any
		operator = "="
		switch vArr := v.(type) {
		case []any:
			operator = getOperator(vArr[0].(string))
			value = vArr[1]
		case []string:
			operator = getOperator(vArr[0])
			value = vArr[1]
		default:
			value = v
		}
		k = filter(k)
		jsonFieldNameTransition(&k)
		switch operator {
		case "not between", "between":
			betweenValues := value.([]any)
			m.Db.Where(fmt.Sprintf("`%s` %s ? AND ?", k, operator), betweenValues[0], betweenValues[1])
		default:
			m.Db.Where(fmt.Sprintf("`%s` %s ?", k, operator), value)
		}
	}
	return m
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
	builder.WriteString("'")
	*filedName = builder.String()
}

func getOperator(str string) (operator string) {
	operator = strings.TrimSpace(str)
	if _, ok := operatorMap[operator]; !ok {
		operator = "="
	}
	return
}

func (m *Model[t]) Order(order string) *Model[t] {
	if order != "" {
		m.Db = m.Db.Order(order)
	}
	return m
}

func (m *Model[t]) Unscoped() *Model[t] {
	m.Db.Statement.Unscoped = true
	return m
}

func (m *Model[t]) Distinct(args ...any) *Model[t] {
	m.Db = m.Db.Distinct(args)
	return m
}

func (m *Model[t]) Select(fields []string) *Model[t] {
	if fields != nil && len(fields) > 0 {
		m.Db = m.Db.Select(fields)
	}
	return m
}

func (m *Model[t]) Offset(offset int) *Model[t] {
	m.Db = m.Db.Offset(offset)
	return m
}

func (m *Model[t]) Limit(offset int) *Model[t] {
	m.Db = m.Db.Limit(offset)
	return m
}

func (m *Model[t]) Omit(columns ...string) *Model[t] {
	m.Db = m.Db.Omit(columns...)
	return m
}
