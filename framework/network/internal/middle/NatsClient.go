package middle

import (
	"log"
	"universal/common/pb"
	"universal/framework/common/uerror"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type NatsClient struct {
	conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, uerror.NewUError(1, -1, err)
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
		return uerror.NewUError(1, -1, err)
	}
	return nil
}

func (d *NatsClient) Publish(key string, pac *pb.Packet) error {
	buf, err := proto.Marshal(pac)
	if err != nil {
		return uerror.NewUError(1, -1, err)
	}

	if err := d.conn.Publish(key, buf); err != nil {
		return uerror.NewUError(1, -1, err)
	}
	return nil
}

func (d *NatsClient) Close() {
	d.conn.Close()
}
