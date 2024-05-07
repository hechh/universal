package domain

import "universal/common/pb"

const (
	NotifyTypeNats = 1
)

type NotifyHandle func(*pb.Packet)

type INotify interface {
	Subscribe(string, NotifyHandle) error
	Publish(string, *pb.Packet) error
}
