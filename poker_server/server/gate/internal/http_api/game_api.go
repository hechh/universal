package http_api

import (
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/uerror"
	"strconv"
)

// todo
func GameReconnect(w http.ResponseWriter, r *http.Request) {
	// 1. 只允许GET方法
	rsp := &pb.QueryPlayerDataRsp{}
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
	userid, err := strconv.ParseUint(query.Get("userid"), 10, 64)
	if err != nil {
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "参数错误",
		}
		response(w, http.StatusBadRequest, rsp)
		return
	}

	// 发送请求
	req := &pb.QueryPlayerDataReq{Uid: userid}
	head := &pb.Head{
		Uid: userid,
		Dst: framework.NewGameRouter(userid, "PlayerMgr", "QueryPlayerData"),
		Src: framework.NewSrcRouter(userid, "Player"),
	}
	if err := cluster.Request(head, req, rsp); err != nil {
		rsp.Head = uerror.ToRspHead(err)
		response(w, http.StatusInternalServerError, rsp)
		return
	}
	response(w, http.StatusOK, rsp)
}
