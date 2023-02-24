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

var operators = []string{"=", "<", "<=", ">", ">=", "<>", "<=>", "!=", "like", "not like", "in", "not in", "between", "not between"}
var operatorMap = map[string]byte{}

func init() {
	for _, operator := range operators {
		operatorMap[operator] = 0
	}
}

func filter(str string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9_$>./[/]\"-]+")
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

func (m *Model[t]) Where(where any, args ...interface{}) *Model[t] {
	if where == nil {
		return m
	}
	if args == nil {
		if w, ok := where.(map[string]any); ok {
			return m._where(&w)
		} else if w, ok := where.(*map[string]any); ok {
			return m._where(w)
		}
	}
	m.Db = m.Db.Where(where, args)
	return m
}

func (m *Model[t]) _where(where *map[string]any) *Model[t] {
	var operator string
	for k, v := range *where {
		var value any
		operator = "="
		switch v.(type) {
		case []any:
			vArr := v.([]any)
			operator = getOperator(vArr[0].(string))
			value = vArr[1]
		case []string:
			vArr := v.([]string)
			operator = getOperator(vArr[0])
			value = vArr[1]
		default:
			value = v
		}
		k = filter(k)
		switch operator {
		case "not between", "between":
			betweenValues := value.([]any)
			m.Db.Where(fmt.Sprintf("%s %s ? AND ?", k, operator), betweenValues[0], betweenValues[1])
		default:
			m.Db.Where(fmt.Sprintf("%s %s ?", k, operator), value)
		}
	}
	return m
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

func (m *Model[t]) Distinct(args ...interface{}) *Model[t] {
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
