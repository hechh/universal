package routine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/fbasic"
)

var (
	_routines sync.Map
)

type RoutineInfo struct {
	clusterType pb.ClusterType
	clusterID   uint32
}

type RoutineInfoList struct {
	updateTime int64
	list       []*RoutineInfo
}

// 获取玩家路由信息
func GetRoutine(head *pb.PacketHead) *RoutineInfoList {
	if val, ok := _routines.Load(head.UID); ok && val != nil {
		return val.(*RoutineInfoList)
	}
	// 新建路由表
	rlist := &RoutineInfoList{updateTime: fbasic.GetNow()}
	// 存储玩家路由表
	_routines.Store(head.UID, rlist)

	return rlist
}

func (d *RoutineInfo) GetType() pb.ClusterType { return d.clusterType }
func (d *RoutineInfo) GetClusterID() uint32    { return d.clusterID }

func (d *RoutineInfoList) SetUpdateTime(now int64) {
	atomic.StoreInt64(&d.updateTime, now)
}

func (d *RoutineInfoList) GetUpdateTime() int64 {
	return atomic.LoadInt64(&d.updateTime)
}

// 查询路由节点
func (d *RoutineInfoList) Get(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.clusterType == typ {
			dst = item
			break
		}
	}
	return
}

// 创建路由节点
func (d *RoutineInfoList) New(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.clusterType == typ {
			dst = item
			break
		}
	}
	if dst == nil {
		dst = &RoutineInfo{clusterType: typ}
		d.list = append(d.list, dst)
	}
	return
}

// 更新玩家路由信息
func (d *RoutineInfoList) UpdateRoutine(head *pb.PacketHead, node *pb.ClusterNode) error {
	if node == nil {
		return fbasic.NewUError(1, pb.ErrorCode_ClusterNodeNotFound, fmt.Sprint(head))
	}
	// 设置路由关系
	head.DstClusterID = node.ClusterID
	// 更新路由信息
	dst := d.New(head.DstClusterType)
	dst.clusterID = node.ClusterID
	// 更新时间
	atomic.StoreInt64(&d.updateTime, time.Now().Unix())
	return nil
}

func init() {
	timer := time.NewTicker(10 * time.Second)
	for {
		<-timer.C
		_routines.Range(func(key, value interface{}) bool {
			val, ok := value.(*RoutineInfoList)
			if !ok || val == nil {
				return true
			}

			// 判断路由信息是否过期
			now := fbasic.GetNow()
			if val.GetUpdateTime()+30*60 <= now {
				_routines.Delete(key)
			}
			return true
		})
	}
}
