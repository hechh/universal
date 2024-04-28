package game

import "flag"

func Run() {
	var addr, logPath string
	flag.StringVar(&addr, "addr", "", "ip:port地址")
	flag.StringVar(&logPath, "log", "", "日志输出目录")
	flag.Parse()

}
