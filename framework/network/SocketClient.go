package network

import (
	"net"
	"universal/common/pb"
)

type SocketClient struct {
	conn     net.Conn
	receives []byte
}

func (d *SocketClient) Read() (*pb.Packet, error) {
	for {

	}

	return nil, nil
}
