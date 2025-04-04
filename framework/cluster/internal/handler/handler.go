package handler

import (
	"hego/common/pb"
	"hego/common/util"
	"hego/framework/cluster/domain"
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

func HandlePoint(head *pb.Head, buf []byte) {
	for _, f := range points {
		if f(util.CopyHead(head), buf) {
			return
		}
	}
}

func HandleTopic(head *pb.Head, buf []byte) {
	for _, f := range topics {
		if f(util.CopyHead(head), buf) {
			return
		}
	}
}
