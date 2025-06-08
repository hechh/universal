package domain

import (
	"time"
	"universal/common/pb"
)

// Actor接口定义
type IActor interface {
	GetId() uint64                                      // 获取Actor ID
	SetId(uint64)                                       // 设置Actor ID
	Start()                                             // 启动Actor
	Stop()                                              // 停止Actor
	GetActorName() string                               // 获取Actor名称
	Register(IActor, ...int)                            // 注册Actor
	ParseFunc(interface{})                              // 解析方法列表
	SendMsg(*pb.Head, ...interface{}) error             // 发送消息
	Send(*pb.Head, []byte) error                        // 发送远程调用
	RegisterTimer(*pb.Head, time.Duration, int32) error // 注册定时器
}
