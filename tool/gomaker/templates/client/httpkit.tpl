
import (
	"net/http"
)


func init() {
{{range $st := .}} http.HandleFunc("/api/{{$st.Type.Name}}", handle)
{{end}}
}

