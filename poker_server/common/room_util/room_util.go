package room_util

import (
	"poker_server/common/pb"

	"github.com/golang/protobuf/proto"
)

func TexasRoomIdTo(roomId uint64) (pb.MatchType, pb.GameType, pb.CoinType) {
	return pb.MatchType((roomId >> 40) & 0xFF), pb.GameType((roomId >> 32) & 0xFF), pb.CoinType((roomId >> 24) & 0xFF)
}

func ToTexasRoomId(m pb.MatchType, g pb.GameType, c pb.CoinType) uint64 {
	return uint64(m&0xFF)<<40 | uint64(g&0xFF)<<32 | uint64(c&0xFF)<<24
}

func ToMatchGameId(m pb.MatchType, g pb.GameType, c pb.CoinType) uint64 {
	return uint64(m&0xFFFF)<<32 | uint64(g&0xFFFF)<<16 | uint64(c&0xFFFF)
}

func DeepCopy(oldData proto.Message, newData proto.Message) error {
	buf, err := proto.Marshal(oldData)
	if err != nil {
		return err
	}
	return proto.Unmarshal(buf, newData)
}
