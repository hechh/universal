/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package user_data

import (
	"fmt"
	"universal/common/pb"
	"universal/common/redis/internal/manager"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

const (
	DBNAME = "universal"
)

func GetKey() string {
	return fmt.Sprintf("uesr_data")
}

func GetField(uid uint64) string {
	return fmt.Sprintf("%d", uid)
}

func HGetAll() (ret map[string]*pb.UserData, err error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		err = uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
		return
	}
	key := GetKey()

	// 加载数据
	kvs, err := client.HGetAll(key)
	if err != nil {
		return
	}

	// 解析数据
	ret = make(map[string]*pb.UserData)
	for k, item := range kvs {
		if len(item) <= 0 {
			continue
		}

		data := &pb.UserData{}
		if err := proto.Unmarshal([]byte(item), data); err != nil {
			return nil, err
		}
		ret[k] = data
	}
	return
}

func HMSet(data map[string]*pb.UserData) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()

	// 设置数据
	vals := []interface{}{}
	for k, v := range data {
		buf, err := proto.Marshal(v)
		if err != nil {
			return err
		}
		vals = append(vals, k, buf)
	}
	return client.HMSet(key, vals...)
}

func HGet(uid uint64) (*pb.UserData, bool, error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return nil, false, uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()
	field := GetField(uid)

	// 加载数据
	str, err := client.HGet(key, field)
	if err != nil {
		return nil, false, err
	}

	// 解析数据
	data := &pb.UserData{}
	if err := proto.Unmarshal([]byte(str), data); err != nil {
		return nil, len(str) > 0, err
	}
	return data, len(str) > 0, nil
}

func HSet(uid uint64, data *pb.UserData) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()
	field := GetField(uid)

	// 序列化数据
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// 设置数据
	return client.HSet(key, field, buf)
}

func HDel(fields ...string) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, -1, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()
	return client.HDel(key, fields...)
}
