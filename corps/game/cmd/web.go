package cmd

import (
	"corps/common"
	"net/http"
)

// http://localhost:8080/gm?cmd=cpus()
func cmdHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cmd := r.FormValue("cmd")
	if cmd != "" {
		common.ParseConsole(g_Cmd, (cmd))
	}
}

func InitWeb() {
	/*go func() {
		http.HandleFunc("/gm", cmdHandle)
		err := http.ListenAndServe(world.Web_Url, nil)
		if err != nil {
			plog.Info("World Web Server : ", err)
		}
	}()*/
}
