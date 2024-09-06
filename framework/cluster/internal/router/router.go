package router

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

var (
	mutex  = sync.Mutex{}
	routes = new(sync.Map)
)

// 玩家路由
type RouteInfo struct {
	uid        uint64
	updateTime int64
	table      *pb.RouteTable
}

// 设置过期清理机制
func Init(expire int64) {
	tt := time.NewTicker(5 * time.Second)
	util.SafeGo(nil, func() {
		for {
			<-tt.C
			routes.Range(func(key, value interface{}) bool {
				val, ok := value.(*RouteInfo)
				if !ok || val == nil {
					return true
				}
				// 判断路由信息是否过期
				if atomic.LoadInt64(&val.updateTime)+expire <= util.GetNowUnixSecond() {
					routes.Delete(key)
					plog.Info("clear route success: %v", val)
				}
				return true
			})
		}
	})
}

func GetOrNew(uid uint64) *RouteInfo {
	if vv, ok := routes.Load(uid); ok {
		return vv.(*RouteInfo)
	}
	// 枷锁创建路由
	mutex.Lock()
	defer mutex.Unlock()
	// 判断是否已经创建
	if vv, ok := routes.Load(uid); ok {
		return vv.(*RouteInfo)
	}
	// 新建
	item := &RouteInfo{
		uid:        uid,
		updateTime: util.GetNowUnixSecond(),
	}
	routes.Store(uid, item)
	return item
}

func Get(uid uint64) *RouteInfo {
	if vv, ok := routes.Load(uid); ok {
		rr := vv.(*RouteInfo)
		atomic.StoreInt64(&rr.updateTime, util.GetNowUnixSecond())
		return rr
	}
	return nil
}

func (d *RouteInfo) String() string {
	buf, _ := json.Marshal(d.Get())
	return fmt.Sprintf("路由信息 uid: %d, update: %d, route: %s\n", d.uid, d.updateTime, string(buf))
}

func (d *RouteInfo) Get() *pb.RouteTable {
	return &pb.RouteTable{
		Gate: atomic.LoadUint32(&d.table.Gate),
		Game: atomic.LoadUint32(&d.table.Game),
		Gm:   atomic.LoadUint32(&d.table.Gm),
		Db:   atomic.LoadUint32(&d.table.Db),
	}
}

func (d *RouteInfo) Refresh(rr *pb.RouteTable) {
	if rr != nil {
		atomic.CompareAndSwapUint32(&d.table.Gate, 0, rr.Gate)
		atomic.CompareAndSwapUint32(&d.table.Game, 0, rr.Game)
		atomic.CompareAndSwapUint32(&d.table.Gm, 0, rr.Gm)
		atomic.CompareAndSwapUint32(&d.table.Db, 0, rr.Db)
	}
}

func (d *RouteInfo) Update(typ pb.SERVER, clusterId uint32) {
	switch typ {
	case pb.SERVER_Gate:
		atomic.SwapUint32(&d.table.Gate, clusterId)
	case pb.SERVER_Game:
		atomic.SwapUint32(&d.table.Game, clusterId)
	case pb.SERVER_Gm:
		atomic.SwapUint32(&d.table.Gm, clusterId)
	case pb.SERVER_Db:
		atomic.SwapUint32(&d.table.Db, clusterId)
	}
}

func (d *RouteInfo) GetServerID(typ pb.SERVER) uint32 {
	switch typ {
	case pb.SERVER_Gate:
		return atomic.LoadUint32(&d.table.Gate)
	case pb.SERVER_Game:
		return atomic.LoadUint32(&d.table.Game)
	case pb.SERVER_Gm:
		return atomic.LoadUint32(&d.table.Gm)
	case pb.SERVER_Db:
		return atomic.LoadUint32(&d.table.Db)
	}
	return 0
}
