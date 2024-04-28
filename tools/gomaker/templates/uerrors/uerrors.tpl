import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

{{range $v := .List}}
{{$funcName := $.TrimPrefix $v.Name (printf "%s_" $.Name)}}
func {{$funcName}}(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, {{$.PkgName}}.{{$v.Name}}, args...)
}
{{end}}

