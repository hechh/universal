package http_api

import (
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/uerror"
	"strconv"
)

func TexasRoomList(w http.ResponseWriter, r *http.Request) {
	// 1. 只允许GET方法
	rsp := &pb.TexasRoomListRsp{}
	if r.Method != http.MethodGet {
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "参数错误",
		}
		response(w, http.StatusMethodNotAllowed, rsp)
		return
	}
	// 2. 解析URL参数
	query := r.URL.Query()
	gameType, err := strconv.ParseInt(query.Get("GameType"), 10, 32)
	if err != nil {
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "参数错误",
		}
		response(w, http.StatusBadRequest, rsp)
		return
	}
	coinType, err := strconv.ParseInt(query.Get("CoinType"), 10, 32)
	if err != nil {
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "参数错误",
		}
		response(w, http.StatusBadRequest, rsp)
		return
	}

	// 发送请求
	head := &pb.Head{
		Src: framework.NewSrcRouter(0, "Player"),
		Dst: framework.NewMatchRouter(uint64(gameType)<<16|uint64(coinType), "MatchTexasRoom", "RoomListReq"),
	}
	req := &pb.TexasRoomListReq{
		GameType: int32(gameType),
		CoinType: int32(coinType),
	}
	if err := cluster.Request(head, req, rsp); err != nil {
		rsp.Head = uerror.ToRspHead(err)
		response(w, http.StatusInternalServerError, rsp)
		return
	}
	response(w, http.StatusOK, rsp)
}
