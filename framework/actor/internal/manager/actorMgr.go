package manager

import (
	"sync"
	"time"
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/base"
	"universal/framework/common/fbasic"
)

var (
	srvPool = sync.Map{}
)

func GetIActor(key string, ff domain.ActorHandle) domain.IActor {
	var act domain.IActor
	if val, ok := srvPool.Load(key); ok && val != nil {
		act = val.(domain.IActor)
	} else {
		act = base.NewActor(key, ff)
		// 启动协程
		act.Start()
		// 存储
		srvPool.Store(key, act)
	}
	return act
}

func Send(key string, ff domain.ActorHandle, pa *pb.Packet) {
	// 获取actor
	act := GetIActor(key, ff)
	// 刷新时间
	act.SetUpdateTime(time.Now().Unix())
	// 发送
	act.Send(pa.Head, pa.Buff)
}

func StopAll() {
	srvPool.Range(func(key, val interface{}) bool {
		vv, ok := val.(domain.IActor)
		if !ok || vv == nil {
			return true
		}
		vv.Stop()
		return true
	})
}

// 清理过期玩家
func SetClearExpire(expire int64) {
	go func() {
		timer := time.NewTicker(5 * time.Second)
		for {
			<-timer.C
			srvPool.Range(func(key, val interface{}) bool {
				vv, ok := val.(domain.IActor)
				if !ok || vv == nil {
					return true
				}
				if upTime := vv.GetUpdateTime(); upTime+expire <= fbasic.GetNow() {
					// 停止协程
					vv.Stop()
					// 删除缓存
					srvPool.Delete(key)
				}
				return true
			})
		}
	}()
}
