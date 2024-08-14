package ttt

type PlayerBaseDisplay struct {
	Name int64
}

// 用户标记
type PBPlayerBase struct {
	Display *PlayerBaseDisplay `protobuf:"bytes,1,opt,name=Display,proto3" json:"Display,omitempty"` // 展示数据
}
