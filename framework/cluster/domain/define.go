package domain

import "hego/common/pb"

const (
	ROOT_DIR = "server"
)

type WatchFunc func(key string, val string) // 监控服务节点变更
type HandleFunc func(*pb.Head, []byte) bool // 内网统一使用RpcPacket
