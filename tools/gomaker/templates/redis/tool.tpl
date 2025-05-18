// Do not modify the generated code

package tools

import (
	pb "forevernine.com/planet/server/common/pbclass"
	"forevernine.com/planet/server/common/dao/internal/manager"
{{range $val := .}} "forevernine.com/planet/server/common/dao/repository/redis/{{$val.PackageName}}"
{{end}}
)

func init() {
{{range $val := .}} {{if $val.IsHash}} manager.RegisterTool({{$val.PackageName}}.DBNAME, "{{$val.Desc}}", "{{$val.Args}}", &pb.{{$val.PBName}}{}, {{$val.PackageName}}.GetRedisKey, {{$val.PackageName}}.GetRedisField) {{else}} manager.RegisterTool({{$val.PackageName}}.DBNAME, "{{$val.Desc}}", "{{$val.Args}}", &pb.{{$val.PBName}}{}, {{$val.PackageName}}.GetRedisKey, nil) {{end}}
{{end}}
}
