package templ

const hashTpl = `
{{/* 定义全局变量  */}}
{{$pkg := .Pkg}}
{{$pb := .Name}} 
{{$dbname := .DbName}}
{{$kformat := .KFormat}}
{{$fformat := .FFormat}}

/*
* 本代码由pbtool工具生成，请勿手动修改
*/

package {{$pkg}}

import (
	"fmt"
	"universal/common/dao/internal/manager"
	"universal/common/pb"
	"universal/framework/library/uerror"

	"github.com/golang/protobuf/proto"
)

const (
	DBNAME = "{{$dbname}}"
)

func GetKey({{.Kargs}}) string {
	return fmt.Sprintf("{{$kformat}}", {{range $i, $item := .Keys}} {{if $i}}, {{end}} {{$item.Name}} {{end}})
}

func GetField({{.Fargs}}) string {
	return fmt.Sprintf("{{$fformat}}", {{range $i, $item := .Fields}} {{if $i}}, {{end}} {{$item.Name}} {{end}})
}

func HGetAll({{.Kargs}}) (ret map[string]*pb.{{$pb}}, err error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		err = uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
		return
	}
	key := GetKey({{.Kvalues}})

	// 加载数据
	kvs, err := client.HGetAll(key)
	if err != nil {
		return
	}

	// 解析数据
	ret = make(map[string]*pb.{{$pb}})
	for k, item := range kvs {
		data := &pb.{{$pb}}{}
		if err := proto.Unmarshal([]byte(item), data); err != nil {
			return nil, err
		}
		ret[k] = data
	}
	return
}

func HMSet({{.Kargs}} {{if .Keys}},{{end}} data map[string]*pb.{{$pb}}) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Kvalues}})

	// 设置数据
	vals := []interface{}{}
	for k, v := range data {
		buf, err := proto.Marshal(v)
		if err != nil {
			return err
		}
		vals = append(vals, k, buf)
	}
	return client.HMSet(key, vals...)
}

func HGet({{.Args}}) (*pb.{{$pb}}, error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return nil, uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Kvalues}})
	field := GetField({{.Fvalues}})

	// 加载数据
	str, err := client.HGet(key, field)
	if err != nil {
		return nil, err
	}

	// 解析数据
	data := &pb.{{$pb}}{}
	if err := proto.Unmarshal([]byte(str), data); err != nil {
		return nil, err
	}
	return data, nil
}

func HSet({{.Args}}, data *pb.{{$pb}}) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Kvalues}})
	field := GetField({{.Fvalues}})

	// 序列化数据
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// 设置数据
	return client.HSet(key, field, buf)
}

func HDel({{.Kargs}} {{if .Keys}},{{end}} fields ...string) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Kvalues}})
	return client.HDel(key, fields...)
}

`
