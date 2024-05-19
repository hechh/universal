import (
	"universal/common/pb"
	"universal/framework/common/uerror"
)

{{range $v := .List}}
{{$funcName := $.TrimPrefix $v.Name (printf "%s_" $.Name)}}
func {{$funcName}}(args ...interface{}) *uerror.UError {
	return uerror.NewUError(1, int32({{$.PkgName}}.{{$v.Name}}), args...)
}
{{end}}

