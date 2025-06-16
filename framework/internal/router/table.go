package router

import (
	"sync"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/safe"
)

type Table struct {
	mutex   sync.RWMutex
	routers map[uint64]domain.IRouter
	exit    chan struct{}
}

func NewTable() *Table {
	return &Table{
		routers: make(map[uint64]domain.IRouter),
		exit:    make(chan struct{}),
	}
}

func (d *Table) Get(id uint64) domain.IRouter {
	d.mutex.RLock()
	rr, ok := d.routers[id]
	d.mutex.RUnlock()
	if ok {
		return rr
	}

	val := &Router{Router: &pb.Router{}, updateTime: time.Now().Unix()}
	d.mutex.Lock()
	d.routers[id] = val
	d.mutex.Unlock()
	return val
}

func (d *Table) Close() error {
	close(d.exit)
	return nil
}

func (d *Table) Expire(ttl int64) {
	safe.Go(func() {
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()

		for {
			select {
			case <-tt.C:
				if keys := d.getExpires(ttl); len(keys) > 0 {
					d.mutex.Lock()
					for _, k := range keys {
						delete(d.routers, k)
					}
					d.mutex.Unlock()
				}
			case <-d.exit:
				return
			}
		}
	})
}

func (d *Table) getExpires(ttl int64) (keys []uint64) {
	now := time.Now().Unix()
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	for key, rr := range d.routers {
		if rr.GetUpdateTime()+ttl <= now {
			keys = append(keys, key)
		}
	}
	return
}
