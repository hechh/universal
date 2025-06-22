package util

import "poker_server/common/pb"

func DestructRoomId(roomId uint64) *pb.DefaultRoomId {
	return &pb.DefaultRoomId{
		GameType:  pb.GameType(roomId >> 32 & 0xFF),
		CoinType:  pb.CoinType(roomId >> 24 & 0xFF),
		MatchType: pb.MatchType(roomId >> 40 & 0xFF),
		Incr:      uint32(roomId & 0xFFFFFF),
	}
}

func GenMatchId(types *pb.DefaultRoomId) uint64 {
	switch types.GetGameType() {
	default: // rummy texas类型
		return uint64(types.GameType)<<32 | uint64(types.CoinType)
	}
}
