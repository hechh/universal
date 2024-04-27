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

func Load(uuid string) domain.ICustom {
	if val, ok := srvPool.Load(uuid); ok && val != nil {
		return val.(domain.ICustom)
	}
	return nil
}

func Store(aa domain.ICustom) {
	srvPool.Store(aa.UUID(), aa)
}

func Send(key string, pa *pb.Packet) {
	// 获取actor
	var act domain.ICustom
	if val, ok := srvPool.Load(key); ok && val != nil {
		act = val.(domain.ICustom)
	} else {
		act = base.NewActor(key, srvFunc)
		// 启动协程
		act.Start()
		// 存储
		srvPool.Store(key, act)
	}
	// 刷新时间
	act.SetUpdateTime(time.Now().Unix())
	// 发送
	act.Send(pa)
}

// 清理过期玩家
func init() {
	timer := time.NewTicker(5 * time.Second)
	for {
		<-timer.C
		srvPool.Range(func(key, val interface{}) bool {
			vv, ok := val.(domain.ICustom)
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
