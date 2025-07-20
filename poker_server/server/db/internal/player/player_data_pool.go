package player

import (
	"poker_server/common/mysql"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
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

func (d *PlayerDataPool) Relogin(head *pb.Head, dd *PlayerData, req *pb.GateLoginRequest) error {
	playerInfoRsp := &pb.HttpPlayerInfoRsp{RespData: &pb.PlayerInfo{Uid: head.Uid}}
	if err := PlayerInfoRequest(head.Uid, playerInfoRsp); err != nil {
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "获取玩家信息失败: %v", err)
	}
	mlog.Debug(head, "PlayerDataPool.Login: rsp:%v", playerInfoRsp)

	actor.SendMsg(&pb.Head{ActorName: "UserInfoMgr", FuncName: "Update"}, &pb.UpdateUserInfoNotify{
		DataType: pb.DataType_DataTypeUserInfo,
		Data: &pb.UserInfo{
			Uid:    head.Uid,
			Name:   playerInfoRsp.RespData.NickName,
			Avatar: playerInfoRsp.RespData.Avatar,
		},
	})

	playerData := dd.Load()
	playerData.Base.PlayerInfo = playerInfoRsp.RespData
	req.PlayerData = playerData
	return cluster.Send(framework.SwapToGame(head, head.Uid, "PlayerMgr", "Login"), req)
}

func (d *PlayerDataPool) Login(head *pb.Head, dd *PlayerData, req *pb.GateLoginRequest) error {
	playerInfoRsp := &pb.HttpPlayerInfoRsp{RespData: &pb.PlayerInfo{Uid: head.Uid}}
	/*
		if err := PlayerInfoRequest(head.Uid, playerInfoRsp); err != nil {
			return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "获取玩家信息失败: %v", err)
		}
	*/
	mlog.Debug(head, "PlayerDataPool.Login: rsp:%v", playerInfoRsp)

	actor.SendMsg(&pb.Head{ActorName: "UserInfoMgr", FuncName: "Update"}, &pb.UpdateUserInfoNotify{
		DataType: pb.DataType_DataTypeUserInfo,
		Data: &pb.UserInfo{
			Uid:    head.Uid,
			Name:   playerInfoRsp.RespData.NickName,
			Avatar: playerInfoRsp.RespData.Avatar,
		},
	})

	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	// 从mysql中加载玩家数据
	userData := &pb.PlayerData{Uid: head.Uid}
	ok, err := session.Get(userData)
	if err != nil {
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "mysql查询失败: %v", err)
	}
	if ok {
		userData.Base.PlayerInfo = playerInfoRsp.RespData
	} else {
		userData.Base = &pb.PlayerDataBase{
			CreateTime: time.Now().Unix(),
			PlayerInfo: playerInfoRsp.RespData,
		}
	}
	atomic.StoreInt64(&dd.updateTime, time.Now().Unix())
	dd.Store(userData)

	req.PlayerData = userData
	return cluster.Send(framework.SwapToGame(head, head.Uid, "PlayerMgr", "Login"), req)
}

func (d *PlayerDataPool) Update(head *pb.Head, data *PlayerData, newData *pb.PlayerData) error {
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	playerData := data.Load()
	newData.Id = playerData.Id
	newData.Version = playerData.Version

	// 保存玩家断线重连信息
	roomInfo := &pb.RoomInfo{}
	if newData.Base != nil && newData.Base.RoomInfo != nil {
		roomInfo.Uid = head.Uid
		roomInfo.GameType = newData.Base.RoomInfo.GameType
		roomInfo.RoomId = newData.Base.RoomInfo.RoomId
	}
	actor.SendMsg(&pb.Head{ActorName: "RoomInfoMgr", FuncName: "Update"}, &pb.UpdateRoomInfoNotify{
		DataType: pb.DataType_DataTypeRoomInfo,
		Data:     roomInfo,
	})

	if newData.Version <= 0 {
		// 新玩家的版本号一定是从0开始
		if _, err := session.Insert(newData); err != nil {
			return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "mysql插入数据失败: data:%v, error:%v", data, err)
		}
	} else {
		affected, err := session.Where("uid = ?", newData.Uid).Update(newData)
		if err != nil {
			return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "mysql更新数据失败: data:%v, error: %v", data, err)
		}
		if affected == 0 {
			actor.SendMsg(&pb.Head{ActorName: "PlayerDataMgr", FuncName: "Remove"}, head.Uid)
			return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "数据版本不一致，丢弃更新: data:%v", data)
		}
	}
	data.Store(newData)
	atomic.StoreInt64(&data.updateTime, time.Now().Unix())
	return nil
}
