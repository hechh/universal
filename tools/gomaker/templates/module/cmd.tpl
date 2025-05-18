syntax = "proto3";
package forevernine.com.planet.proto;
option  go_package = "common/pbclass";

enum CMD {
{{range $v := .}}	{{$v.Name}} = {{$v.Value}};
{{end}} }
