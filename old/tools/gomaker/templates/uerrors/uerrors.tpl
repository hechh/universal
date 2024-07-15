import (
	"universal/common/pb"
	"universal/framework/common/uerror"
)

{{$type := printf "%s_" .Type.Name}}
{{$pkg := .Type.PkgName}}

{{range $v := .List}}
{{$funcName := $.TrimPrefix $v.Name $type}}
func {{$funcName}}(args ...interface{}) *uerror.UError {
	return uerror.NewUError(1, int32({{$pkg}}.{{$v.Name}}), args...)
}
{{end}}

