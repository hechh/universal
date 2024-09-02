package nodes

import (
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"universal/common/pb"
	"universal/framework/basic/util"
	"universal/framework/plog"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
)

var (
	root  string                    // etcd的跟路径
	self  *pb.ClusterInfo           // 当前服务节点的clusterID
	infos sync.Map                  // 所有节点
	stubs = make(map[string]uint64) // stub配置数据
)

func Init(selfNode *pb.ClusterInfo, path string, st map[string]uint64) {
	root = path
	self = selfNode
	stubs = st
}

func GetSelfTopicChannel() string {
	return filepath.Clean(filepath.Join(root, self.Type.String()))
}

func GetSelfChannel() string {
	return filepath.Clean(filepath.Join(root, self.Type.String(), cast.ToString(self.ClusterID)))
}

// 服务广播消息
func GetTopicChannel(typ pb.SERVICE) string {
	return filepath.Clean(filepath.Join(root, typ.String()))
}

// 服务间单播消息
func GetChannel(typ pb.SERVICE, clusterID uint32) string {
	return filepath.Clean(filepath.Join(root, typ.String(), cast.ToString(clusterID)))
}

// 删除节点通知
func DeleteNotify(key, value string) {
	strs := strings.Split(key, "/")
	if ll := len(strs); ll >= 2 {
		if val, ok := pb.SERVICE_value[strs[ll-2]]; ok && val != 0 {
			infos.Delete(cast.ToUint32(strs[ll-1]))
		}
	}
}

// 添加节点通知
func AddNotify(key, value string) {
	data := &pb.ClusterInfo{}
	proto.Unmarshal(util.StringToBytes(value), data)
	Insert(data)
}

// 添加节点
func Insert(item *pb.ClusterInfo) {
	infos.Store(item.ClusterID, item)
	plog.Info("新增服务节点: %v", item)
}

// 查询节点
func Get(clusterID uint32) *pb.ClusterInfo {
	if val, ok := infos.Load(clusterID); ok && val != nil {
		return val.(*pb.ClusterInfo)
	}
	return nil
}

func GetSelf() *pb.ClusterInfo {
	return self
}

// 查询指定类型的所有节点
func Gets(typ pb.SERVICE) (rets []*pb.ClusterInfo) {
	infos.Range(func(key, val interface{}) bool {
		item, ok := val.(*pb.ClusterInfo)
		if !ok || item == nil {
			return true
		}
		if item.Type == typ {
			rets = append(rets, item)
		}
		return true
	})
	sort.Slice(rets, func(i, j int) bool {
		return rets[i].CreateTime < rets[j].CreateTime
	})
	return
}

// 随机一个指定类型的服务节点
func Random(typ pb.SERVICE, id uint64, actorName string) (info *pb.ClusterInfo) {
	rets := Gets(typ)
	ll := len(rets)
	if ll <= 0 {
		return
	}
	// 随机一个节点
	if id <= 0 {
		info = rets[rand.Intn(ll)]
	} else {
		// 根据id或者actorName随机
		if val, ok := stubs[actorName]; ok && val > 0 {
			info = rets[int(id/val)%ll]
		} else {
			info = rets[int(id)%ll]
		}
	}
	return
}
