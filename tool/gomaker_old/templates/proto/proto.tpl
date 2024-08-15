import (
	"corps/pb"
)

func init(){
{{range $st := .}} RegisterProto(&pb.{{$st.Name}}{}, &{{$st.Name}}{})
{{end}}
}
