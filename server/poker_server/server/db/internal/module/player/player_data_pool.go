package player

import (
	"poker_server/common/dao"
	"poker_server/common/dao/domain"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"
	"sync/atomic"
	"time"
)

type PlayerDataPool struct {
	actor.ActorPool
}

func NewPlayerDataPool(size int) *PlayerDataPool {
	ret := &PlayerDataPool{}
	ret.Register(ret, size)
	ret.ParseFunc(reflect.TypeOf(ret))
	ret.SetId(uint64(pb.DataType_DataTypePlayerData))
	ret.Start()
	actor.Register(ret)
	return ret
}

func (d *PlayerDataPool) Login(head *pb.Head, dd *PlayerData, req *pb.GateLoginRequest) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	playerInfoRsp := &pb.HttpPlayerInfoRsp{RespData: &pb.PlayerInfo{Uid: head.Uid}}
	if err := PlayerInfoRequest(head.Uid, playerInfoRsp); err != nil {
		//return uerror.NEW(pb.ErrorCode_REQUEST_FAIELD, head, "获取玩家信息失败: %v", err)
	}
	mlog.Debugf("PlayerDataPool.Login: rsp:%v", playerInfoRsp)

	// 从mysql中加载玩家数据
	userData := &pb.PlayerData{Uid: head.Uid}
	ok, err := session.Get(userData)
	if err != nil {
		return uerror.NEW(pb.ErrorCode_REQUEST_FAIELD, head, "mysql查询失败: %v", err)
	}
	if !ok {
		userData.Base = &pb.PlayerDataBase{PlayerInfo: playerInfoRsp.RespData}
	} else {
		userData.Base.PlayerInfo = playerInfoRsp.RespData
	}

	actor.SendMsg(&pb.Head{ActorName: "UserInfoMgr", FuncName: "Update"}, &pb.UpdateUserInfoNotify{
		DataType: pb.DataType_DataTypeUserInfo,
		Data: &pb.UserInfo{
			Uid:    userData.Uid,
			Name:   userData.Base.PlayerInfo.NickName,
			Avatar: userData.Base.PlayerInfo.Avatar,
		},
	})
	atomic.StoreInt64(&dd.updateTime, time.Now().Unix())
	dd.Store(userData)
	req.PlayerData = userData
	return framework.Send(framework.SwapToGame(head, head.Uid, "PlayerMgr", "Login"), req)
}

func (d *PlayerDataPool) Update(head *pb.Head, data *PlayerData, newData *pb.PlayerData) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	playerData := data.Load()
	newData.Id = playerData.Id
	newData.Version = playerData.Version

	if newData.Version <= 0 {
		// 新玩家的版本号一定是从0开始
		if _, err := session.Insert(newData); err != nil {
			return uerror.NEW(pb.ErrorCode_REQUEST_FAIELD, nil, "mysql插入数据失败: data:%v, error:%v", data, err)
		}
	} else {
		affected, err := session.ID(newData.Id).Update(newData)
		if err != nil {
			return uerror.NEW(pb.ErrorCode_REQUEST_FAIELD, head, "mysql更新数据失败: data:%v, error: %v", data, err)
		}
		if affected == 0 {
			actor.SendMsg(&pb.Head{ActorName: "PlayerDataMgr", FuncName: "Remove"}, head.Uid)
			return uerror.NEW(pb.ErrorCode_PARAM_INVALID, head, "数据版本不一致，丢弃更新: data:%v", data)
		}
	}
	data.Store(newData)
	atomic.StoreInt64(&data.updateTime, time.Now().Unix())
	return nil
}
