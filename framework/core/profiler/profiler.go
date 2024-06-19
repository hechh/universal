package profiler

import (
	"net/http"
	"net/http/pprof"
	"universal/framework/library/plog"
	"universal/framework/util"

	"github.com/google/gops/agent"
)

// 本地http服务
var (
	local = http.NewServeMux()
)

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	local.HandleFunc(pattern, handler)
}

func Handle(pattern string, handler http.Handler) {
	local.Handle(pattern, handler)
}

func InitGops(addr string) error {
	return agent.Listen(agent.Options{Addr: addr})
}

// 本地服务端口
func InitPprof(addr string) {
	util.SafeGo(plog.Fatal, func() {
		server := &http.Server{Addr: addr, Handler: local}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	})
}

// 默认开启pprof工具
func init() {
	local.HandleFunc("/debug/pprof/", pprof.Index)
	local.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	local.HandleFunc("/debug/pprof/profile", pprof.Profile)
	local.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	local.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
