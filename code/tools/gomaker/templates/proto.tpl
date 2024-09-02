
func init(){
{{range $st := .}} RegisterJson("{{$st.Type.Name}}", &{{$st.Type.Name}}{})
{{end}}
}
