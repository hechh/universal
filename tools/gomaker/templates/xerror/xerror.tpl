// Do not modify the generated code

package xerrors


import (
	"forevernine.com/base/srvcore/libs/terror"
)

{{range $v := .}}
func {{$v.FName}}() *terror.Terror {
	return terror.New({{$v.Code}}, "{{$v.ErrMsg}}", nil)
}
{{end}}
