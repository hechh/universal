package player

import (
	"poker_server/common/dao"
	"poker_server/common/dao/domain"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"reflect"
	"sync/atomic"
	"time"
)

func init() {
	dao.RegisterMysqlTable(domain.MYSQL_DB_PLAYER_DATA, &pb.PlayerData{})
}

type PlayerData struct {
	atomic.Pointer[pb.PlayerData]
	updateTime int64
}

type PlayerDataMgr struct {
	actor.Actor
	pool  *PlayerDataPool
	datas map[uint64]*PlayerData
}

func NewPlayerDataMgr() *PlayerDataMgr {
	ret := &PlayerDataMgr{
		pool:  NewPlayerDataPool(100),
		datas: make(map[uint64]*PlayerData),
	}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.SetId(uint64(pb.DataType_DataTypePlayerData))
	ret.Start()
	actor.Register(ret)
	return ret
}

func (p *PlayerDataMgr) Init() error {
	return p.RegisterTimer(&pb.Head{ActorName: "PlayerDataMgr", FuncName: "OnTick"}, 15*time.Minute, -1)
}

func (p *PlayerDataMgr) Close() {
	p.Actor.Stop()
	p.pool.Stop()
}

func (p *PlayerDataMgr) OnTick() {
	now := time.Now().Unix()
	for key, data := range p.datas {
		if atomic.LoadInt64(&data.updateTime)+30*60 <= now {
			delete(p.datas, key)
		}
	}
}

func (p *PlayerDataMgr) Login(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	if data, ok := p.datas[head.Uid]; ok {
		playerData := data.Load()
		if playerData != nil {
			req.PlayerData = playerData
			return framework.Send(framework.SwapToGame(head, head.Uid, "PlayerMgr", "Login"), req)
		}
	} else {
		p.datas[head.Uid] = new(PlayerData)
	}

	framework.StopAutoSendToClient(head)
	return p.pool.SendMsg(head, p.datas[head.Uid], req)
}

func (p *PlayerDataMgr) Update(head *pb.Head, data *pb.UpdatePlayerDataNotify) error {
	if val, ok := p.datas[head.Uid]; ok {
		return p.pool.SendMsg(head, val, data.Data)
	}

	val := new(PlayerData)
	val.Store(data.Data)
	atomic.StoreInt64(&val.updateTime, time.Now().Unix())
	p.datas[head.Uid] = val
	return p.pool.SendMsg(head, val, data.Data)
}

func (p *PlayerDataMgr) Remove(head *pb.Head) {
	delete(p.datas, head.Uid)
}
