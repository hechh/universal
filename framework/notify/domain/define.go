package domain

import "universal/common/pb"

const (
	NotifyTypeNats = 1
)

type INotify interface {
	Subscribe(string, func(*pb.Packet)) error
	Publish(string, *pb.Packet) error
}
