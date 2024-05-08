package base

import (
	"sync"
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/notify/domain"
)

type Broadcast struct {
	sync.RWMutex
	key   string
	funcs map[uint64]domain.NotifyHandle
}

func NewBroadcast(key string) *Broadcast {
	return &Broadcast{
		key:   key,
		funcs: make(map[uint64]domain.NotifyHandle),
	}
}

func (d *Broadcast) Register(uid uint64, f domain.NotifyHandle) error {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.funcs[uid]; ok {
		return fbasic.NewUError(1, pb.ErrorCode_NotSupported, uid, "Broadcast")
	}
	d.funcs[uid] = f
	return nil
}

func (d *Broadcast) Delete(uid uint64) {
	d.Lock()
	defer d.Unlock()
	delete(d.funcs, uid)
}

func (d *Broadcast) GetKey() string {
	return d.key
}

func (d *Broadcast) Handle(pac *pb.Packet) {
	d.RLock()
	defer d.RUnlock()
	for uid, f := range d.funcs {
		pac.Head.UID = uid
		f(pac)
	}
}
