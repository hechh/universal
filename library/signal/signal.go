package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalNotify(ff func(), sigs ...os.Signal) {
	defaults := []os.Signal{syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	if len(sigs) > 0 {
		defaults = append(defaults, sigs...)
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, defaults...)
	<-sig
	ff()
}
