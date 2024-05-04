package routine

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/fbasic"
)

var (
	_routines sync.Map
)

func init() {
	go clearExpire()
}

func clearExpire() {
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

func Print() {
	_routines.Range(func(key, value interface{}) bool {
		val, ok := value.(*RoutineInfoList)
		if !ok || val == nil {
			return true
		}
		log.Println(key, "-----routine----->", val.String())
		return true
	})
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

type RoutineInfo struct {
	ClusterType pb.ClusterType
	ClusterID   uint32
}

type RoutineInfoList struct {
	updateTime int64
	list       []*RoutineInfo
}

func (d *RoutineInfoList) String() string {
	buf, _ := json.Marshal(d.list)
	return fmt.Sprintf("%d: %s", d.updateTime, string(buf))
}

func (d *RoutineInfoList) SetUpdateTime(now int64) {
	atomic.StoreInt64(&d.updateTime, now)
}

func (d *RoutineInfoList) GetUpdateTime() int64 {
	return atomic.LoadInt64(&d.updateTime)
}

// 查询路由节点
func (d *RoutineInfoList) Get(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.ClusterType == typ {
			dst = item
			break
		}
	}
	return
}

// 创建路由节点
func (d *RoutineInfoList) New(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.ClusterType == typ {
			dst = item
			break
		}
	}
	if dst == nil {
		dst = &RoutineInfo{ClusterType: typ}
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
	dst.ClusterID = node.ClusterID
	// 更新时间
	atomic.StoreInt64(&d.updateTime, time.Now().Unix())
	return nil
}
