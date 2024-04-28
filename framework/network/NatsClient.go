package network

import (
	"log"
	"universal/common/pb"
	"universal/framework/fbasic"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_NatsBuildClient, err)
	}
	return &NatsClient{conn: conn}, nil
}

func (d *NatsClient) Subscribe(key string, f func(*pb.Packet)) error {
	_, err := d.conn.Subscribe(key, func(msg *nats.Msg) {
		pac := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pac); err != nil {
			log.Printf("Subscribe error: %v", err)
		} else {
			f(pac)
		}
	})
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_NatsSubscribe, err)
	}
	return nil
}

func (d *NatsClient) Publish(key string, pac *pb.Packet) error {
	buf, err := proto.Marshal(pac)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}

	if err := d.conn.Publish(key, buf); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_NatsPublish, err)
	}
	return nil
}
