package user_info

import (
	"poker_server/common/dao/repository/redis/user_info"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"reflect"
	"time"
)

type UserInfo struct {
	*pb.UserInfo
	updateTime int64
}

type UserInfoMgr struct {
	actor.Actor
	datas map[uint64]*UserInfo
	uids  []uint64
}

func NewUserInfoMgr() *UserInfoMgr {
	ret := &UserInfoMgr{datas: make(map[uint64]*UserInfo)}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeUserInfo))
	ret.Actor.Start()
	actor.Register(ret)
	ret.RegisterTimer(&pb.Head{FuncName: "OnTick"}, 5*time.Second, -1)
	return ret
}

func (d *UserInfoMgr) Init() error {
	return nil
}

func (d *UserInfoMgr) Close() {
	d.Save()
	d.Stop()
}

func (d *UserInfoMgr) OnTick() {
	if len(d.uids) > 0 {
		d.SendMsg(&pb.Head{FuncName: "Save"})
	}
}

func (d *UserInfoMgr) Save() error {
	if len(d.uids) <= 0 {
		return nil
	}
	tmps := map[string]*pb.UserInfo{}
	for _, uid := range d.uids {
		tmps[user_info.GetField(uid)] = d.datas[uid].UserInfo
	}
	if err := user_info.HMSet(tmps); err != nil {
		return err
	}
	d.uids = d.uids[:0]
	return nil
}

func (m *UserInfoMgr) Update(head *pb.Head, req *pb.UpdateUserInfoNotify) error {
	if req.Data == nil {
		return nil
	}
	if vals, ok := m.datas[req.Data.Uid]; ok {
		vals.UserInfo = req.Data
		vals.updateTime = time.Now().Unix()
	} else {
		m.datas[req.Data.Uid] = &UserInfo{UserInfo: req.Data, updateTime: time.Now().Unix()}
	}
	for _, uid := range m.uids {
		if uid == req.Data.Uid {
			return nil
		}
	}
	m.uids = append(m.uids, req.Data.Uid)
	return nil
}

func (m *UserInfoMgr) Query(head *pb.Head, req *pb.GetUserInfoReq, rsp *pb.GetUserInfoRsp) error {
	fields := []string{}
	for _, uid := range req.UidList {
		if _, ok := m.datas[uid]; !ok {
			fields = append(fields, user_info.GetField(uid))
			continue
		}
		rsp.UserList = append(rsp.UserList, m.datas[uid].UserInfo)
	}
	rets, err := user_info.HMGet(fields...)
	if err != nil {
		return err
	}
	for _, item := range rets {
		m.datas[item.Uid] = &UserInfo{updateTime: time.Now().Unix(), UserInfo: item}
		rsp.UserList = append(rsp.UserList, item)
	}
	return nil
}
