/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package generator_data

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/common/redis/internal/manager"
	"poker_server/library/uerror"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	DBNAME = "poker"
)

func GetKey() string {
	return fmt.Sprintf("generator")
}

func Get() (*pb.GeneratorData, bool, error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return nil, false, uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()

	// 加载数据
	str, err := client.Get(key)
	if err != nil {
		return nil, false, err
	}

	// 解析数据
	data := &pb.GeneratorData{}
	if err := proto.Unmarshal([]byte(str), data); err != nil {
		return nil, len(str) > 0, err
	}
	return data, len(str) > 0, nil
}

func Set(data *pb.GeneratorData) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()

	// 存储数据
	return client.Set(key, buf)
}

func SetEX(data *pb.GeneratorData, ttl time.Duration) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()

	// 存储数据
	return client.SetEX(key, buf, ttl)
}

func Del() error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()

	// 删除数据
	return client.Del(key)
}
