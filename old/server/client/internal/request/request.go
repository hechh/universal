package request

import (
	"poker_server/common/pb"

	"github.com/golang/protobuf/proto"
)

var (
	Cmds = make(map[uint32]func() proto.Message)
)

func init() {
	Cmds[16777216] = func() proto.Message { return &pb.GateLoginRequest{} }
	Cmds[16777217] = func() proto.Message { return &pb.GateLoginResponse{} }
	Cmds[16777218] = func() proto.Message { return &pb.KickPlayerNotify{} }
	Cmds[33554432] = func() proto.Message { return &pb.GateHeartRequest{} }
	Cmds[33554433] = func() proto.Message { return &pb.GateHeartResponse{} }
	Cmds[33554434] = func() proto.Message { return &pb.GetBagReq{} }
	Cmds[33554435] = func() proto.Message { return &pb.GetBagRsp{} }
	Cmds[50331648] = func() proto.Message { return &pb.TexasRoomListReq{} }
	Cmds[50331649] = func() proto.Message { return &pb.TexasRoomListRsp{} }
	Cmds[50331650] = func() proto.Message { return &pb.RummyRoomListReq{} }
	Cmds[50331651] = func() proto.Message { return &pb.RummyRoomListRsp{} }
	Cmds[50331652] = func() proto.Message { return &pb.SngJoinRoomReq{} }
	Cmds[50331653] = func() proto.Message { return &pb.SngJoinRoomRsp{} }
	Cmds[67108864] = func() proto.Message { return &pb.TexasEventNotify{} }
	Cmds[67108866] = func() proto.Message { return &pb.TexasJoinRoomReq{} }
	Cmds[67108867] = func() proto.Message { return &pb.TexasJoinRoomRsp{} }
	Cmds[67108868] = func() proto.Message { return &pb.TexasQuitRoomReq{} }
	Cmds[67108869] = func() proto.Message { return &pb.TexasQuitRoomRsp{} }
	Cmds[67108870] = func() proto.Message { return &pb.TexasSitDownReq{} }
	Cmds[67108871] = func() proto.Message { return &pb.TexasSitDownRsp{} }
	Cmds[67108872] = func() proto.Message { return &pb.TexasStandUpReq{} }
	Cmds[67108873] = func() proto.Message { return &pb.TexasStandUpRsp{} }
	Cmds[67108880] = func() proto.Message { return &pb.TexasBuyInReq{} }
	Cmds[67108881] = func() proto.Message { return &pb.TexasBuyInRsp{} }
	Cmds[67108882] = func() proto.Message { return &pb.TexasDoBetReq{} }
	Cmds[67108883] = func() proto.Message { return &pb.TexasDoBetRsp{} }
	Cmds[67108884] = func() proto.Message { return &pb.RummyJoinRoomReq{} }
	Cmds[67108885] = func() proto.Message { return &pb.RummyJoinRoomRsp{} }
	Cmds[67108886] = func() proto.Message { return &pb.RummyEventNotify{} }
	Cmds[67108888] = func() proto.Message { return &pb.RummyQuitRoomReq{} }
	Cmds[67108889] = func() proto.Message { return &pb.RummyQuitRoomRsp{} }
	Cmds[67108890] = func() proto.Message { return &pb.TexasChangeRoomReq{} }
	Cmds[67108891] = func() proto.Message { return &pb.TexasChangeRoomRsp{} }
	Cmds[67108896] = func() proto.Message { return &pb.RummySaveCardGroupReq{} }
	Cmds[67108897] = func() proto.Message { return &pb.RummySaveCardGroupRsp{} }
	Cmds[67108898] = func() proto.Message { return &pb.RummyOprCardReq{} }
	Cmds[67108899] = func() proto.Message { return &pb.RummyOprCardRsp{} }
	Cmds[67108900] = func() proto.Message { return &pb.RummyFixCardReq{} }
	Cmds[67108901] = func() proto.Message { return &pb.RummyFixCardRsp{} }
	Cmds[67108902] = func() proto.Message { return &pb.RummyReadyRoomReq{} }
	Cmds[67108903] = func() proto.Message { return &pb.RummyReadyRoomRsp{} }
	Cmds[67108904] = func() proto.Message { return &pb.RummyGetOutCardsReq{} }
	Cmds[67108905] = func() proto.Message { return &pb.RummyGetOutCardsRsp{} }
	Cmds[67108906] = func() proto.Message { return &pb.TexasStatisticsReq{} }
	Cmds[67108907] = func() proto.Message { return &pb.TexasStatisticsRsp{} }
	Cmds[67108908] = func() proto.Message { return &pb.SngRankReq{} }
	Cmds[67108909] = func() proto.Message { return &pb.SngRankRsp{} }
	Cmds[67108910] = func() proto.Message { return &pb.TexasQuitRoomNotify{} }
	Cmds[67108912] = func() proto.Message { return &pb.RummyChangeRoomReq{} }
	Cmds[67108913] = func() proto.Message { return &pb.RummyChangeRoomRsp{} }
	Cmds[67108914] = func() proto.Message { return &pb.RummyMatchReq{} }
	Cmds[67108915] = func() proto.Message { return &pb.RummyMatchRsp{} }
	Cmds[67108916] = func() proto.Message { return &pb.RummyCancelMatchReq{} }
	Cmds[67108917] = func() proto.Message { return &pb.RummyCancelMatchRsp{} }
	Cmds[67108918] = func() proto.Message { return &pb.RummyGiveUpReq{} }
	Cmds[67108919] = func() proto.Message { return &pb.RummyGiveUpRsp{} }
	Cmds[83886080] = func() proto.Message { return &pb.TexasGameReportReq{} }
	Cmds[83886081] = func() proto.Message { return &pb.TexasGameReportRsp{} }
	Cmds[83886082] = func() proto.Message { return &pb.RummyMatchSelectReq{} }
	Cmds[83886083] = func() proto.Message { return &pb.RummyMatchSelectRsp{} }
	Cmds[83886084] = func() proto.Message { return &pb.GetUserInfoReq{} }
	Cmds[83886085] = func() proto.Message { return &pb.GetUserInfoRsp{} }
	Cmds[83886086] = func() proto.Message { return &pb.GetTexasGameReportReq{} }
	Cmds[83886087] = func() proto.Message { return &pb.GetTexasGameReportRsp{} }
}
