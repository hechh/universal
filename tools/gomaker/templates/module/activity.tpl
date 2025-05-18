syntax = "proto3";
package forevernine.com.planet.proto;
option  go_package = "common/pbclass";

// 活动模块
enum ActivityType {
{{range $v := .}}	{{$v.Name}} = {{$v.Value}};
{{end}} }
