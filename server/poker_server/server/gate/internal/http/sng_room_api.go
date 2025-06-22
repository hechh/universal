package http

import (
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework"
	"strconv"
)

func sngRoomList(w http.ResponseWriter, r *http.Request) {
	// 1. 只允许GET方法
	rsp := &pb.SngRoomListRsp{}
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
	matchType, err := strconv.ParseInt(query.Get("MatchType"), 10, 32)
	if err != nil {
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "参数错误",
		}
		response(w, http.StatusBadRequest, rsp)
		return
	}
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
	dst := framework.NewMatchRouter(uint64(pb.DataType_DataTypeSngRoom), "SngRoomMgr", "RoomListReq")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeUid, 0)
	req := &pb.SngRoomListReq{
		MatchType: pb.MatchType(matchType),
		GameType:  pb.GameType(gameType),
		CoinType:  pb.CoinType(coinType),
	}
	if err := framework.Request(head, req, rsp); err != nil {
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "请求失败",
		}
		response(w, http.StatusInternalServerError, rsp)
		return
	}
	response(w, http.StatusOK, rsp)
}
