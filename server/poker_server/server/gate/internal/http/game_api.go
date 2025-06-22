package http

import (
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/mlog"
	"strconv"
)

// todo
func gameReconnect(w http.ResponseWriter, r *http.Request) {
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

	dst := framework.NewGameRouter(userid, "PlayerMgr", "QueryPlayerData")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeUid, userid)

	req := &pb.QueryPlayerDataReq{
		Uid: userid,
	}
	if err := framework.Request(head, req, rsp); err != nil {
		mlog.Infof("%v", err)
		rsp.Head = &pb.RspHead{
			Code: int32(pb.ErrorCode_PARAM_INVALID),
			Msg:  "请求失败",
		}
		response(w, http.StatusInternalServerError, rsp)
		return
	}
	response(w, http.StatusOK, rsp)
}
