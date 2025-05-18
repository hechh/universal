// Do not modify the generated code

{{/* 定义变量 */}}
{{$package := .Package}}
{{$dbname := .DbName}}
{{$pbname := .Name}}
{{$key := .Key}}
{{$ast := .Ast}}

package {{.Package}}

import (
	"fmt"

{{if .IsCache}} "forevernine.com/base/srvcore/framework" {{end}}
	"forevernine.com/planet/server/common/dao/domain"
	"forevernine.com/planet/server/common/dao/internal/manager"
	"forevernine.com/planet/server/common/dao/internal/service/redis"
	pb "forevernine.com/planet/server/common/pbclass"
	promith "forevernine.com/planet/server/common/prometheus/repository/old_promith"
)

const (
	TABLE_NAME = "{{$pbname}}"
	DBNAME = domain.{{$dbname}}
)

func new{{.Name}}() *pb.{{$pbname}} {
	return &pb.{{$pbname}}{
        {{range $field := $ast.Maps}} {{$field.Name}}: make(map[{{$field.KType.GetType}}]{{$field.VType.GetType}}),
        {{end}} {{range $field := $ast.Arrays}} {{$field.Name}}: make({{$field.Type.GetType}}, 0),
        {{end}}
	}
}

func GetRedisKey({{$key.Values.Arg}}) string {
	return fmt.Sprintf("{{$key.Format}}" {{if $key.Values}}, {{$key.Values.Val ""}} {{end}})
}

func Get({{$key.Values.Arg}}) (result *pb.{{$pbname}}, exist bool, err error) {
	var buf []byte
	if buf, err = redis.Get(DBNAME, GetRedisKey({{$key.Values.Val ""}})); err != nil {
		return
	}
	
	// 普罗米修斯上报
	promith.ReportRedis(len(buf))	

	exist = len(buf) > 0
	result = new{{$pbname}}()

	if err = result.Unmarshal(buf); err == nil {
		manager.CallGets(TABLE_NAME, result {{if $key.Values}}, {{$key.Values.Val ""}} {{end}})
	}
	return
}

func Set({{if $key.Values}} {{$key.Values.Arg}}, {{end}} data *pb.{{$pbname}}) (err error) {
	var buf []byte
	if buf, err = data.Marshal(); err != nil {
		return 
	}

	if err = redis.Set(DBNAME, GetRedisKey({{$key.Values.Val ""}}), buf); err == nil {
		manager.CallSets(TABLE_NAME, data {{if $key.Values}}, {{$key.Values.Val ""}} {{end}})
	}
	return
}

func Del({{$key.Values.Arg}}) (err error) {
	if err = redis.Del(DBNAME, GetRedisKey({{$key.Values.Val ""}})); err == nil {
		manager.CallDels(TABLE_NAME {{if $key.Values}}, {{$key.Values.Val ""}} {{end}})
	}
	return
}
