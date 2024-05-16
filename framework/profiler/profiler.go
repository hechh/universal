package profiler

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"universal/common/pb"
	"universal/framework/fbasic"

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

// 本地服务端口
func Init(port int) error {
	addr := fmt.Sprintf("localhost:%d", port)
	server := &http.Server{Addr: addr, Handler: local}
	if err := server.ListenAndServe(); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_TcpListen, err)
	}
	if err := agent.Listen(agent.Options{Addr: addr}); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_TcpListen, err)
	}
	return nil
}

// 默认开启pprof工具
func init() {
	local.HandleFunc("/debug/pprof/", pprof.Index)
	local.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	local.HandleFunc("/debug/pprof/profile", pprof.Profile)
	local.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	local.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
