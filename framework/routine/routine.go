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
	tables sync.Map
)

type RoutineInfo struct {
	ClusterType pb.ClusterType
	ClusterID   uint32
}

type RoutineTable struct {
	uid        uint64         // 玩家uid
	status     int32          // 在线状态
	updateTime int64          // 更新时间
	list       []*RoutineInfo // 该玩家的路由信息
}

func Print() {
	tables.Range(func(key, value interface{}) bool {
		val, ok := value.(*RoutineTable)
		if !ok || val == nil {
			return true
		}
		log.Println(key, "-----routine----->", val.String())
		return true
	})
}

// 获取在线uid列表
func GetOnlines() (list []uint64) {
	tables.Range(func(key, value interface{}) bool {
		val, ok := value.(*RoutineTable)
		if !ok || val == nil {
			return true
		}

		if val.GetStatus() > 0 {
			list = append(list, val.uid)
		}
		return true
	})
	return
}

// 获取玩家路由信息
func GetRoutine(uid uint64) *RoutineTable {
	if val, ok := tables.Load(uid); ok && val != nil {
		return val.(*RoutineTable)
	}
	// 新建路由表
	rlist := &RoutineTable{
		updateTime: fbasic.GetNow(),
		status:     1,
		uid:        uid,
	}
	// 存储玩家路由表
	tables.Store(uid, rlist)
	return rlist
}

func (d *RoutineTable) SetStatus(online int32) {
	atomic.StoreInt32(&d.status, online)
}

func (d *RoutineTable) GetStatus() int32 {
	return atomic.LoadInt32(&d.status)
}

func (d *RoutineTable) String() string {
	buf, _ := json.Marshal(d.list)
	return fmt.Sprintf("%d: %s", d.updateTime, string(buf))
}

func (d *RoutineTable) SetUpdateTime(now int64) {
	atomic.StoreInt64(&d.updateTime, now)
}

func (d *RoutineTable) GetUpdateTime() int64 {
	return atomic.LoadInt64(&d.updateTime)
}

// 查询路由节点
func (d *RoutineTable) Get(typ pb.ClusterType) (dst *RoutineInfo) {
	for _, item := range d.list {
		if item.ClusterType == typ {
			dst = item
			break
		}
	}
	return
}

// 创建路由节点
func (d *RoutineTable) New(typ pb.ClusterType) (dst *RoutineInfo) {
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
func (d *RoutineTable) UpdateRoutine(head *pb.PacketHead, node *pb.ClusterNode) error {
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

func init() {
	go clearExpire()
}

func clearExpire() {
	timer := time.NewTicker(10 * time.Second)
	for {
		<-timer.C
		tables.Range(func(key, value interface{}) bool {
			val, ok := value.(*RoutineTable)
			if !ok || val == nil {
				return true
			}

			// 判断路由信息是否过期
			now := fbasic.GetNow()
			if val.GetUpdateTime()+30*60 <= now {
				tables.Delete(key)
			}
			return true
		})
	}
}
