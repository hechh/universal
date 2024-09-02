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
	pb.RouteInfo
	uid        uint64
	updateTime int64
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

func (d *RouteInfo) Get() *pb.RouteInfo {
	return &pb.RouteInfo{
		Gate:   atomic.LoadUint32(&d.Gate),
		Game:   atomic.LoadUint32(&d.Game),
		Gm:     atomic.LoadUint32(&d.Gm),
		Db:     atomic.LoadUint32(&d.Db),
		Dip:    atomic.LoadUint32(&d.Dip),
		Record: atomic.LoadUint32(&d.Record),
	}
}

func (d *RouteInfo) Refresh(rr *pb.RouteInfo) {
	if rr != nil {
		atomic.CompareAndSwapUint32(&d.Gate, 0, rr.Gate)
		atomic.CompareAndSwapUint32(&d.Game, 0, rr.Game)
		atomic.CompareAndSwapUint32(&d.Gm, 0, rr.Gm)
		atomic.CompareAndSwapUint32(&d.Db, 0, rr.Db)
		atomic.CompareAndSwapUint32(&d.Record, 0, rr.Record)
		atomic.CompareAndSwapUint32(&d.Dip, 0, rr.Dip)
	}
}

func (d *RouteInfo) Update(typ pb.SERVICE, clusterId uint32) {
	switch typ {
	case pb.SERVICE_GATE:
		atomic.SwapUint32(&d.Gate, clusterId)
	case pb.SERVICE_GAME:
		atomic.SwapUint32(&d.Game, clusterId)
	case pb.SERVICE_GM:
		atomic.SwapUint32(&d.Gm, clusterId)
	case pb.SERVICE_DB:
		atomic.SwapUint32(&d.Db, clusterId)
	case pb.SERVICE_Dip:
		atomic.SwapUint32(&d.Dip, clusterId)
	case pb.SERVICE_Record:
		atomic.SwapUint32(&d.Record, clusterId)
	}
}

func (d *RouteInfo) GetClusterID(typ pb.SERVICE) uint32 {
	switch typ {
	case pb.SERVICE_GATE:
		return atomic.LoadUint32(&d.Gate)
	case pb.SERVICE_GAME:
		return atomic.LoadUint32(&d.Game)
	case pb.SERVICE_GM:
		return atomic.LoadUint32(&d.Gm)
	case pb.SERVICE_DB:
		return atomic.LoadUint32(&d.Db)
	case pb.SERVICE_Dip:
		return atomic.LoadUint32(&d.Dip)
	case pb.SERVICE_Record:
		return atomic.LoadUint32(&d.Record)
	}
	return 0
}
