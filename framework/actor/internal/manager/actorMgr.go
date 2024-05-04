package manager

import (
	"sync"
	"time"
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/base"
	"universal/framework/fbasic"
)

var (
	srvPool = sync.Map{}
	srvFunc domain.ActorHandle
)

const (
	ActorExpireTime = 30 * 60 //单位：秒
)

func SetActorHandle(h domain.ActorHandle) {
	srvFunc = h
}

func GetIActor(key string) domain.IActor {
	var act domain.IActor
	if val, ok := srvPool.Load(key); ok && val != nil {
		act = val.(domain.IActor)
	} else {
		act = base.NewActor(key, srvFunc)
		// 启动协程
		act.Start()
		// 存储
		srvPool.Store(key, act)
	}
	return act
}

func Send(key string, pa *pb.Packet) {
	// 获取actor
	act := GetIActor(key)
	// 刷新时间
	act.SetUpdateTime(time.Now().Unix())
	// 发送
	act.Send(pa)
}

func cleanExpire() {
	timer := time.NewTicker(5 * time.Second)
	for {
		<-timer.C
		srvPool.Range(func(key, val interface{}) bool {
			vv, ok := val.(domain.IActor)
			if !ok || vv == nil {
				return true
			}
			now := fbasic.GetNow()
			if upTime := vv.GetUpdateTime(); upTime+ActorExpireTime <= now {
				// 停止协程
				vv.Stop()
				// 删除缓存
				srvPool.Delete(key)
			}
			return true
		})
	}
}

// 清理过期玩家
func init() {
	go cleanExpire()
}
