package actor

import (
	"reflect"
	"strings"
	"universal/library/baselib/uerror"
)

var (
	methods map[string]reflect.Method
)

func RegisterMethod(st interface{}) error {
	rType := reflect.TypeOf(st)
	if rType.Kind() != reflect.Ptr {
		return uerror.New(1, -1, "传入必须是指针类型")
	}

	// 类型名称
	name := rType.String()
	if index := strings.Index(name, "."); index != -1 {
		name = name[index+1:]
	}

	// 注册方法
	for i := 0; i < rType.NumMethod(); i++ {
		methond := rType.Method(i)
		fullName := name + "." + methond.Name

		// 判断是否已经存在
		if _, ok := methods[fullName]; ok {
			return uerror.New(1, -1, "方法已存在: %s", fullName)
		}
		methods[fullName] = methond
	}

	return nil
}

func GetMethod(name string) (reflect.Method, bool) {
	method, ok := methods[name]
	return method, ok
}
