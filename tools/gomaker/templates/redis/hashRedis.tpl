// Do not modify the generated code

{{/* 定义变量 */}}
{{$package := .Package}}
{{$dbname := .DbName}}
{{$pbname := .Name}}
{{$key := .Key}}
{{$field := .Field}}
{{$total :=  .Key.Values.Join .Field.Values}}

package {{$package}}

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

func new{{$pbname}}() *pb.{{$pbname}} {
	return &pb.{{$pbname}}{
        {{range $field := .Ast.Maps}} {{$field.Name}}: make(map[{{$field.KType.GetType}}]{{$field.VType.GetType}}),
        {{end}} {{range $field := .Ast.Arrays}} {{$field.Name}}: make({{$field.Type.GetType}}, 0),
        {{end}}
	}
}

func GetRedisKey({{$key.Values.Arg}}) string {
	return fmt.Sprintf("{{$key.Format}}" {{if $key.Values}}, {{$key.Values.Val ""}} {{end}})
}

func GetRedisField({{$field.Values.Arg}}) string {
	return fmt.Sprintf("{{$field.Format}}" {{if $field.Values}}, {{$field.Values.Val ""}} {{end}})
}

func HGet({{$total.Arg}}) (result *pb.{{$pbname}}, exist bool, err error) {
	buf, err := redis.HGet(DBNAME, GetRedisKey({{$key.Values.Val ""}}), GetRedisField({{$field.Values.Val ""}}))
	if err != nil {
		return
	}

	// 普罗米修斯上报
	promith.ReportRedis(len(buf))	
	
	exist = len(buf) > 0
	result = new{{$pbname}}()
	if err = result.Unmarshal(buf); err == nil {
		manager.CallGets(TABLE_NAME, result {{if $total}}, {{$total.Val ""}} {{end}})
	}
	return
}

func HSet({{if $total}} {{$total.Arg}}, {{end}} data *pb.{{$pbname}}) (err error) {
	var buf []byte
	if buf, err = data.Marshal(); err != nil {
		return
	}
	
	if err = redis.HSet(DBNAME, GetRedisKey({{$key.Values.Val ""}}), GetRedisField({{$field.Values.Val ""}}), buf); err == nil {
		manager.CallSets(TABLE_NAME, data {{if $total}}, {{$total.Val ""}} {{end}})
	}
	return
}

func HDel({{$total.Arg}}) (err error) {
	if err = redis.HDel(DBNAME, GetRedisKey({{$key.Values.Val ""}}), GetRedisField({{$field.Values.Val ""}})); err == nil {
		manager.CallDels(TABLE_NAME, {{$total.Val ""}})
	}
	return
}

{{if .UUID}} {{$utotal := .Key.Values.Join .Field.Values .UUID.Values}}
{{if eq .UUID.Field "struct"}}
func GetByUUID({{$utotal.Arg}}) (result *pb.{{$pbname}}, isnew bool, err error) {
	result, _, err = HGet({{$total.Val ""}})
	if err != nil {
		return
	}
	{{$tmp := index .UUID.Values 0}}
	if result.{{$tmp.Name}} != {{$tmp.Name}} {
		isnew = true
		result = new{{.Name}}()
		result.{{$tmp.Name}} = {{$tmp.Name}}	
	}
	return
}
{{end}}

{{if eq .UUID.Field "list"}}
func newSub() *pb.{{.SubAst.Type.Name}} {
	return &pb.{{.SubAst.Type.Name}}{
        {{range $field := .SubAst.Maps}} {{$field.Name}}: make(map[{{$field.KType.GetType}}]{{$field.VType.GetType}}),
        {{end}} {{range $field := .SubAst.Arrays}} {{$field.Name}}: make({{$field.Type.GetType}}, 0),
        {{end}}
	}
}
{{$zero := index .UUID.Values 0}}
{{$one := index .UUID.Values 1}}
func GetByUUID({{if $total}} {{$total.Arg}}, {{end}} {{$one.Name}} {{$one.Type}}) (datahis *pb.{{$pbname}}, data *pb.{{$zero.Type}}, isnew bool, err error) {
	datahis, _, err = HGet({{$total.Val ""}})
	if err != nil {
		return
	}
	for _, item := range datahis.{{$zero.Name}} {
		if item.{{$one.Name}} == {{$one.Name}} {
			data = item
			return
		}
	}
	
	isnew = true
	data = newSub()
	data.{{$one.Name}} = {{$one.Name}}
	datahis.{{$zero.Name}} = append(datahis.{{$zero.Name}}, data)
	for i := len(datahis.{{$zero.Name}})-1; i >=1; i-- {
		tmp := datahis.{{$zero.Name}}[i]
		datahis.{{$zero.Name}}[i] = datahis.{{$zero.Name}}[i-1]
		datahis.{{$zero.Name}}[i-1] = tmp 
	}
	if len(datahis.{{$zero.Name}}) > 2 {
		datahis.{{$zero.Name}} = datahis.{{$zero.Name}}[0:2]
	}
	return
}
{{end}}
{{end}}
