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

const (
	EXPIRE = 30 * 60
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

// 删除路由表
func DelRoutine(uid uint64) {
	tables.Delete(uid)
}

// 获取玩家路由信息
func GetRoutine(uid uint64) *RoutineTable {
	if val, ok := tables.Load(uid); ok {
		vv := val.(*RoutineTable)
		atomic.StoreInt64(&vv.updateTime, fbasic.GetNow())
		return vv
	}
	// 新建路由表
	rlist := &RoutineTable{
		updateTime: fbasic.GetNow(),
		uid:        uid,
	}
	// 存储玩家路由表
	tables.Store(uid, rlist)
	return rlist
}

func (d *RoutineTable) String() string {
	buf, _ := json.Marshal(d.list)
	return fmt.Sprintf("%d: %s", d.updateTime, string(buf))
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
			if atomic.LoadInt64(&val.updateTime)+EXPIRE <= now {
				tables.Delete(key)
			}
			return true
		})
	}
}
