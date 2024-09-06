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
	root  string         // etcd的跟路径
	self  *pb.ServerInfo // 当前服务节点的clusterID
	infos sync.Map       // 所有节点
)

func Init(path string, selfNode *pb.ServerInfo) {
	root = path
	self = selfNode
}

func GetSelfTopicChannel() string {
	return filepath.Clean(filepath.Join(root, self.Type.String()))
}

func GetSelfChannel() string {
	return filepath.Clean(filepath.Join(root, self.Type.String(), cast.ToString(self.ServerID)))
}

// 服务广播消息
func GetTopicChannel(typ pb.SERVER) string {
	return filepath.Clean(filepath.Join(root, typ.String()))
}

// 服务间单播消息
func GetChannel(typ pb.SERVER, serverID uint32) string {
	return filepath.Clean(filepath.Join(root, typ.String(), cast.ToString(serverID)))
}

// 删除节点通知
func DeleteNotify(key, value string) {
	strs := strings.Split(key, "/")
	if ll := len(strs); ll >= 2 {
		if val, ok := pb.SERVER_value[strs[ll-2]]; ok && val != 0 {
			infos.Delete(cast.ToUint32(strs[ll-1]))
		}
	}
}

// 添加节点通知
func AddNotify(key, value string) {
	data := &pb.ServerInfo{}
	proto.Unmarshal(util.StringToBytes(value), data)
	Insert(data)
}

// 添加节点
func Insert(item *pb.ServerInfo) {
	infos.Store(item.ServerID, item)
	plog.Info("新增服务节点: %v", item)
}

// 查询节点
func Get(serverID uint32) *pb.ServerInfo {
	if val, ok := infos.Load(serverID); ok && val != nil {
		return val.(*pb.ServerInfo)
	}
	return nil
}

func GetSelf() *pb.ServerInfo {
	return self
}

// 查询指定类型的所有节点
func Gets(typ pb.SERVER) (rets []*pb.ServerInfo) {
	infos.Range(func(key, val interface{}) bool {
		item, ok := val.(*pb.ServerInfo)
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
func Random(typ pb.SERVER, id uint64) (info *pb.ServerInfo) {
	rets := Gets(typ)
	ll := len(rets)
	if ll <= 0 {
		return
	}
	// 随机一个节点
	if id <= 0 {
		info = rets[rand.Intn(ll)]
	} else {
		// 根据id随机
		info = rets[int(id)%ll]
	}
	return
}
