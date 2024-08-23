package handler

import (
	"universal/common/pb"
	"universal/framework/basic/util"
	"universal/framework/cluster/domain"
)

var (
	topics = []domain.HandleFunc{}
	points = []domain.HandleFunc{}
)

func BindPoint(fs ...domain.HandleFunc) {
	points = append(points, fs...)
}

func BindTopic(fs ...domain.HandleFunc) {
	topics = append(topics, fs...)
}

func HandlePoint(head *pb.RpcHead, buf []byte) {
	for _, f := range points {
		if f(util.CopyRpcHead(head), buf) {
			return
		}
	}
}

func HandleTopic(head *pb.RpcHead, buf []byte) {
	for _, f := range topics {
		if f(util.CopyRpcHead(head), buf) {
			return
		}
	}
}
