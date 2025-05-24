package templ

const stringTpl = `
{{/* 定义全局变量  */}}
{{$pkg := .Pkg}}
{{$pb := .Name}} 
{{$dbname := .DbName}}
{{$format := .Format}}

/*
* 本代码由pbtool工具生成，请勿手动修改
*/

package {{$pkg}}

import (
	"fmt"
	"universal/common/dao/internal/manager"
	"universal/common/pb"
	"universal/framework/library/uerror"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	DBNAME = "{{$dbname}}"
)

func GetKey({{.Args}}) string {
	return fmt.Sprintf("{{$format}}", {{range $i, $item := .Keys}} {{if $i}}, {{end}} {{$item.Name}} {{end}})
}

func Get({{.Args}}) (*pb.{{$pb}}, error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return nil, uerror.New(1, -1,"redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Values}})

	// 加载数据
	str, err := client.Get(key)
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

func Set({{.Args}}, data *pb.{{$pb}}) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1,"redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Values}})
	
	// 存储数据
	return client.Set(key, buf)
}


func SetEX({{.Args}}, data *pb.{{$pb}}, ttl time.Duration) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1,"redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Values}})
	
	// 存储数据
	return client.SetEX(key, buf, ttl)
}

func Del({{.Args}}) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1,"redis数据库不存在: %s", DBNAME)
	}
	key := GetKey({{.Values}})

	// 删除数据
	return client.Del(key)
}

`
