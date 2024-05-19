package broadcast

import (
	"sync"
	"universal/common/pb"
)

type Broadcast struct {
	sync.RWMutex
	funcs map[uint64]func(*pb.Packet)
}

func NewBroadcast() *Broadcast {
	return &Broadcast{
		funcs: make(map[uint64]func(*pb.Packet)),
	}
}

func (d *Broadcast) Add(uid uint64, f func(*pb.Packet)) {
	d.Lock()
	defer d.Unlock()
	d.funcs[uid] = f
}

func (d *Broadcast) Delete(uid uint64) {
	d.Lock()
	defer d.Unlock()
	delete(d.funcs, uid)
}

func (d *Broadcast) Send(pac *pb.Packet) {
	d.RLock()
	defer d.RUnlock()
	for uid, f := range d.funcs {
		pac.Head.UID = uid
		f(pac)
	}
}
