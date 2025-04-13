package basic

import (
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func Filter(pattern string, vals ...string) (rets []string, err error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	for _, val := range vals {
		if re.MatchString(val) {
			continue
		}
		rets = append(rets, val)
	}
	return
}

// 阻塞接受信号
func SignalHandle(f func(os.Signal), sigs ...os.Signal) {
	ch := make(chan os.Signal, 0)
	args := append([]os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL}, sigs...)
	signal.Notify(ch, args...)
	for item := range ch {
		f(item)
		os.Exit(0)
	}
}
