package request

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/uerror"
)

func TexasMatchRoomReq(head *pb.Head, matchReq *pb.TexasMatchRoomReq, matchRsp *pb.TexasMatchRoomRsp) error {
	cfg := texas_config.MGetID(int32(matchReq.TableId))
	if cfg == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "游戏配置不存在")
	}
	id := room_util.ToMatchGameId(cfg.MatchType, cfg.GameType, cfg.CoinType)
	newHead := &pb.Head{
		Src: head.Dst,
		Dst: framework.NewMatchRouter(id, "MatchTexasRoom", "MatchRoomReq"),
		Uid: head.Uid,
	}
	if err := cluster.Request(newHead, matchReq, matchRsp); err != nil {
		return err
	}
	if matchRsp.Head != nil {
		return uerror.ToError(matchRsp.Head)
	}
	return nil
}

func TexasJoinReq(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) error {
	newHead := &pb.Head{
		Src: head.Dst,
		Uid: head.Uid,
		Dst: framework.NewRoomRouter(req.RoomId, "TexasGameMgr", "JoinRoomReq"),
	}
	if err := cluster.Request(newHead, req, rsp); err != nil {
		return err
	}
	if rsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}
	return nil
}

func TexasQuitReq(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	newHead := &pb.Head{
		Src: head.Dst,
		Dst: framework.NewRoomRouter(req.RoomId, "TexasGameMgr", "QuitRoomReq"),
		Uid: head.Uid,
	}
	if err := cluster.Request(newHead, req, rsp); err != nil {
		return err
	}
	if rsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}
	return nil
}
