package manager

import (
	"sync"
	"time"
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

func LoadActor(uuid string) *base.Actor {
	if val, ok := srvPool.Load(uuid); ok && val != nil {
		return val.(*base.Actor)
	}
	return nil
}

func StoreActor(aa *base.Actor) {
	srvPool.Store(aa.GetUUID(), aa)
}

func GetIActor(uuid string) domain.IActor {
	if val, ok := srvPool.Load(uuid); ok && val != nil {
		return val.(*base.Actor)
	}
	user := base.NewActor(uuid, srvFunc)
	// 启动协程
	user.Start()
	// 存储
	srvPool.Store(uuid, user)
	return user
}

// 清理过期玩家
func init() {
	timer := time.NewTicker(5 * time.Second)
	for {
		<-timer.C
		srvPool.Range(func(key, val interface{}) bool {
			vv, ok := val.(*base.Actor)
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
