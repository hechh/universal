package uerrors

import (
	"universal/common/pb"
	"universal/framework/basic"
)

{{range $v := .List}}
{{$funcName := $.TrimPrefix $v.Name (printf "%s_" $.Name)}}
func {{$funcName}}() *basic.UError {
	return basic.NewUError(3, {{$.PkgName}}.{{$v.Name}}, "{{$funcName}}")
}
{{end}}

