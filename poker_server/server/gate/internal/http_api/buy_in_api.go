package http_api

import (
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/uerror"
	"poker_server/server/game/module/reward"
	"strconv"
)

func BuyInApi(w http.ResponseWriter, r *http.Request) {
	// 1. 只允许GET方法
	rsp := &pb.RewardRsp{}
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
	uid, err := strconv.ParseInt(query.Get("Uid"), 10, 64)
	if err != nil {
		rsp.Head = &pb.RspHead{Code: int32(pb.ErrorCode_PARAM_INVALID), Msg: "参数错误"}
		response(w, http.StatusBadRequest, rsp)
		return
	}
	coinType, err := strconv.ParseInt(query.Get("CoinType"), 10, 32)
	if err != nil {
		rsp.Head = &pb.RspHead{Code: int32(pb.ErrorCode_PARAM_INVALID), Msg: "参数错误"}
		response(w, http.StatusBadRequest, rsp)
		return
	}
	Incr, err := strconv.ParseInt(query.Get("Incr"), 10, 64)
	if err != nil {
		rsp.Head = &pb.RspHead{Code: int32(pb.ErrorCode_PARAM_INVALID), Msg: "参数错误"}
		response(w, http.StatusBadRequest, rsp)
		return
	}
	// 发送请求
	req := &pb.RewardReq{
		RewardList: []*pb.Reward{
			reward.ToReward(pb.CoinType(coinType), Incr),
		},
	}
	head := &pb.Head{
		Uid: uint64(uid),
		Dst: framework.NewGameRouter(uint64(uid), "Player", "RewardReq"),
		Src: framework.NewSrcRouter(uint64(uid), "Player"),
	}
	if err := cluster.Request(head, req, rsp); err != nil {
		rsp.Head = uerror.ToRspHead(err)
		response(w, http.StatusInternalServerError, rsp)
		return
	}
	response(w, http.StatusOK, rsp)
}
