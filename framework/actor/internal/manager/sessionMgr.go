package manager

import (
	"sync"
	"time"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/session"
)

var (
	srvPool = sync.Map{}
	srvFunc domain.PacketHandle
)

func SetPacketHandle(h domain.PacketHandle) {
	srvFunc = h
}

func GetSession(id string) domain.ISession {
	if val, ok := srvPool.Load(id); ok || val != nil {
		return val.(*session.Session)
	}
	user := session.NewSession(id, srvFunc)
	// 启动协程
	user.Start()
	// 存储
	srvPool.Store(id, user)
	return user
}

// 清理过期玩家
func init() {
	timer := time.NewTicker(5 * time.Second)
	for {
		<-timer.C
		srvPool.Range(func(key, val interface{}) bool {
			vv, ok := val.(*session.Session)
			if !ok || vv == nil {
				return true
			}
			if vv.IsExpired() {
				// 停止协程
				vv.Stop()
				// 删除缓存
				srvPool.Delete(key)
			}
			return true
		})
	}
}
