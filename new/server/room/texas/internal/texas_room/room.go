package texas_room

import (
	"poker_server/common/dao/repository/redis/texas_room_data"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/internal/machine"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

type TexasRoom struct {
	framework.Actor
	texasCfg   *pb.TexasConfig       // Texas配置
	machineCfg *pb.MachineConfig     // 机器配置
	data       *pb.TexasRoomData     // 游戏数据
	record     *pb.TexasGameRecord   // 游戏记录
	machine    *machine.TexasMachine // 机器
	isChange   int32                 // 是否变更
	isFinished int32                 // 是否结束
	isStart    int32                 // 是否开始
}

func NewTexasRoom() *TexasRoom {
	ret := &TexasRoom{}
	ret.Actor.Register(ret)
	return ret
}

func (d *TexasRoom) GetRecord() *pb.TexasGameRecord {
	return d.record
}

func (d *TexasRoom) SetRecord(rr *pb.TexasGameRecord) {
	d.record = rr
}

func (d *TexasRoom) GetTexasRoomData() *pb.TexasRoomData {
	return d.data
}

// 查询或者创建玩家
func (d *TexasRoom) GetOrNewPlayer(uid uint64) *pb.TexasPlayerData {
	usr := d.data.Table.Players[uid]
	if usr == nil {
		usr = &pb.TexasPlayerData{
			Uid:      uid,
			GameInfo: &pb.TexasPlayerGameInfo{},
		}
		d.data.Table.Players[uid] = usr
	}
	return usr
}

// 是否变更数据
func (d *TexasRoom) IsChange() bool {
	return atomic.LoadInt32(&d.isChange) > 0
}

// 变更数据
func (d *TexasRoom) Change() {
	atomic.AddInt32(&d.isChange, 1)
}

// 保存数据
func (d *TexasRoom) Save() error {
	if d.data == nil {
		return nil
	}
	ttl := time.Duration(d.texasCfg.RoomKeepLive+15) * time.Minute
	return texas_room_data.Set(d.data.RoomId, d.data, ttl)
}

// 定时触发
func (d *TexasRoom) OnTick(nowMs int64) {
	if d.isStart <= 0 {
		return
	}
	if d.isFinished > 0 {
		head := &pb.Head{ActorName: "TexasRoomMgr", FuncName: "RemoveRoom"}
		framework.SendMsg(head, d.data.RoomId)
		return
	}
	if d.data == nil || d.data.Table == nil {
		return
	}
	if d.machine == nil {
		d.machine = machine.NewTexasMachine(nowMs, d.data.Table.CurState, d)
	}
	d.machine.Handle(nowMs, d)
}

func (d *TexasRoom) Report(proto.Message) error {
	// todo
	return nil
}

func (d *TexasRoom) SendToClient(uid uint64, msgId pb.CMD, msg proto.Message) error {
	mlog.Infof("SendToClient uid:%d, eventType:%s, msg:%v", uid, msgId.String(), msg)
	return SendPackageByPlayerIDRummy(uid, msgId, msg)
}

func (d *TexasRoom) NotifyToClient(notifyId pb.TexasEventType, msg proto.Message) error {
	uids := []uint64{}
	for uid := range d.data.Table.Players {
		uids = append(uids, uid)
	}

	mlog.Infof("NotifyToClient eventType:%s, msg:%v", notifyId.String(), msg)

	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	event := &pb.TexasEventNotify{
		Event:   notifyId,
		Content: buf,
	}
	return BroadcastPackageByPlayerIDRummy(uids, pb.CMD_TEXAS_EVENT_NOTIFY, event)
}

func (d *TexasRoom) NotifyToPlayerClient(uid uint64, notifyId pb.TexasEventType, msg proto.Message) error {
	mlog.Infof("NotifyToClient uid:%d, eventType:%s, msg:%v", uid, notifyId.String(), msg)

	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	event := &pb.TexasEventNotify{
		Event:   notifyId,
		Content: buf,
	}
	return SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_EVENT_NOTIFY, event)
}
