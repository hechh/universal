package router

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/common/fbasic"
	"universal/framework/common/plog"
)

var (
	routings sync.Map
)

// 删除路由表
func DeleteRouteList(uid uint64) {
	routings.Delete(uid)
}

// 获取玩家路由信息
func GetRouteList(uid uint64) *RouteList {
	if val, ok := routings.Load(uid); ok {
		vv := val.(*RouteList)
		atomic.StoreInt64(&vv.updateTime, fbasic.GetNow())
		return vv
	}
	// 新建路由表
	rlist := &RouteList{
		updateTime: fbasic.GetNow(),
		uid:        uid,
	}
	// 存储玩家路由表
	routings.Store(uid, rlist)
	return rlist
}

func SetClearExpire(expire int64) {
	go func() {
		timer := time.NewTicker(5 * time.Second)
		for {
			<-timer.C
			routings.Range(func(key, value interface{}) bool {
				val, ok := value.(*RouteList)
				if !ok || val == nil {
					return true
				}

				// 判断路由信息是否过期
				now := fbasic.GetNow()
				if atomic.LoadInt64(&val.updateTime)+expire <= now {
					routings.Delete(key)
				}
				return true
			})
		}
	}()
}

func Print() {
	routings.Range(func(key, value interface{}) bool {
		val, ok := value.(*RouteList)
		if !ok || val == nil {
			return true
		}
		plog.InfoSkip(1, "router: %s", val.String())
		return true
	})
}
