package router

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/async"
	"universal/library/mlog"
)

type Table struct {
	ttl   int64
	mutex sync.RWMutex
	exit  chan struct{}
	infos map[string]domain.IRouter
}

func New() *Table {
	return &Table{
		exit:  make(chan struct{}),
		infos: make(map[string]domain.IRouter),
	}
}

func getKey(idType pb.IdType, id uint64) string {
	return fmt.Sprintf("%s:%d", strings.ToLower(idType.String()), id)
}

func (d *Table) Get(idType pb.IdType, id uint64) domain.IRouter {
	// 读取路由信息
	key := getKey(idType, id)
	if info := d.get(key); info != nil {
		return info
	}
	// 创建路由信息
	val := &Router{updateTime: time.Now().Unix(), Router: &pb.Router{}}
	d.mutex.Lock()
	d.infos[key] = val
	d.mutex.Unlock()
	return val
}

func (d *Table) Add(idType pb.IdType, id uint64, info *pb.Router) {
	if info == nil {
		return
	}
	// 读取路由信息
	key := getKey(idType, id)
	oldInfo := d.get(key)
	if oldInfo == nil {
		d.set(key, info)
		return
	}
	// 更新路由信息
	newInfo := &Router{Router: info}
	for i := pb.NodeType_Begin + 1; i < pb.NodeType_End; i++ {
		oldInfo.Set(i, newInfo.Get(i))
	}
}

func (d *Table) SetExpire(ttl int64) {
	d.ttl = ttl
	async.SafeGo(mlog.Fatalf, d.run)
}

func (r *Table) Close() {
	r.exit <- struct{}{}
}

func (d *Table) set(key string, rr *pb.Router) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.infos[key] = &Router{Router: rr, updateTime: time.Now().Unix()}
}

func (d *Table) get(key string) domain.IRouter {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if router, ok := d.infos[key]; ok {
		return router
	}
	return nil
}

func (r *Table) run() {
	tt := time.NewTicker(time.Duration(r.ttl) * time.Second)
	defer tt.Stop()

	// 定时清理
	for {
		select {
		case <-tt.C:
			// 获取过期节点
			now := time.Now().Unix()
			if rets := r.getExpires(now); len(rets) > 0 {
				// 删除节点
				r.mutex.Lock()
				for _, val := range rets {
					delete(r.infos, val)
				}
				r.mutex.Unlock()
			}
		case <-r.exit:
			return
		}
	}
}

func (r *Table) getExpires(now int64) (rets []string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for routeId, val := range r.infos {
		if val.IsExpire(now, r.ttl) {
			rets = append(rets, routeId)
		}
	}
	return
}
