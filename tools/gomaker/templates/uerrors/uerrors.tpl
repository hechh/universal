import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

{{range $v := .List}}
{{$funcName := $.TrimPrefix $v.Name (printf "%s_" $.Name)}}
func {{$funcName}}() *fbasic.UError {
	return fbasic.NewUError(1, {{$.PkgName}}.{{$v.Name}}, "{{$funcName}}")
}
{{end}}

