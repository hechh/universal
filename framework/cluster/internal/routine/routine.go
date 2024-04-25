package routine

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/basic"
)

type Routine struct {
	routines *sync.Map // uid -> *RoutineInfoList
}

func NewRoutine() *Routine {
	return &Routine{routines: new(sync.Map)}
}

// 定时清理过期路由信息
func (d *Routine) Refresh() {
	timer := time.NewTicker(10 * time.Second)
	for {
		<-timer.C
		now := time.Now().Unix()
		d.routines.Range(func(key, val interface{}) bool {
			if value := val.(*RoutineInfoList); value.updateTime+30*60 <= now {
				d.routines.Delete(key)
			}
			return true
		})
	}
}

// 获取玩家路由信息
func (d *Routine) GetRoutine(head *pb.PacketHead) *RoutineInfoList {
	if val, ok := d.routines.Load(head.UID); ok {
		return val.(*RoutineInfoList)
	}
	rlist := &RoutineInfoList{updateTime: time.Now().Unix()}
	d.routines.Store(head.UID, rlist)
	return rlist
}

type RoutineInfo struct {
	service   pb.ClusterType
	clusterID uint32
}

type RoutineInfoList struct {
	updateTime int64
	list       []*RoutineInfo
}

func (d *RoutineInfo) GetType() pb.ClusterType { return d.service }
func (d *RoutineInfo) GetClusterID() uint32    { return d.clusterID }

func (d *RoutineInfoList) Update() {
	d.updateTime = time.Now().Unix()
}

func (d *RoutineInfoList) Get(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.service == typ {
			dst = item
		}
	}
	return
}

func (d *RoutineInfoList) GetAndNew(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.service == typ {
			dst = item
		}
	}
	if dst == nil {
		dst = &RoutineInfo{service: typ}
		d.list = append(d.list, dst)
	}
	return
}

// 对玩家路由
func (d *RoutineInfoList) UpdateRoutine(head *pb.PacketHead, node *pb.ClusterNode) error {
	if node == nil {
		return basic.NewUError(1, pb.ErrorCode_NotExist, "Node not found")
	}
	// 设置路由关系
	head.DstClusterID = node.ClusterID
	//head.SocketID = node.SocketId
	// 更新路由信息
	dst := d.GetAndNew(head.DstClusterType)
	dst.clusterID = node.ClusterID
	// 更新时间
	atomic.StoreInt64(&d.updateTime, time.Now().Unix())
	return nil
}
