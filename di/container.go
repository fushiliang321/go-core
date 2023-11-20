package di

import (
	"errors"
	"fmt"
	"reflect"
)

var container = map[any]any{}

func getKey[T any]() (string, reflect.Kind) {
	var t *T
	ref := reflect.TypeOf(t).Elem()
	return ref.PkgPath() + "/" + ref.Name(), ref.Kind()
}

func Set[T any](entrys ...T) error {
	var entry T
	key, kind := getKey[T]()
	fmt.Println(key, kind)
	switch kind {
	case reflect.Interface:
		if len(entrys) == 0 {
			return errors.New("entry missing")
		}
		entry = entrys[0]
	case reflect.Struct:
		if len(entrys) == 0 {
			entry = *new(T)
		} else {
			entry = entrys[0]
		}
	}
	fmt.Println(entry)
	container[key] = entry
	return nil
}

func Get[T any]() T {
	var t T
	fmt.Println(container[t])

	key, kind := getKey[T]()
	fmt.Println(key, kind)

	return container[key].(T)
}
