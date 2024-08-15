
import (
	"net/http"
)


func init() {
{{range $st := .}} http.HandleFunc("/api/{{$st.Name}}", handle)
//http.HandleFunc("/html/api/{{$st.Name}}", htmlApiHandle)
{{end}}
}

