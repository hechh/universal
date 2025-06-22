/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package rummy_room_data

import (
	"fmt"
	"poker_server/common/dao/internal/manager"
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/library/uerror"

	"github.com/golang/protobuf/proto"
)

const (
	DBNAME = "poker"
)

func GetKey() string {
	return fmt.Sprintf("room")
}

func GetField(RoomId uint64) string {
	return fmt.Sprintf("%d", RoomId)
}

func HGetAll() (ret map[string]*pb.RummyRoomData, err error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		err = uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
		return
	}
	key := GetKey()

	// 加载数据
	kvs, err := client.HGetAll(key)
	if err != nil {
		return
	}

	// 解析数据
	ret = make(map[string]*pb.RummyRoomData)
	for k, item := range kvs {
		if len(item) <= 0 {
			continue
		}

		data := &pb.RummyRoomData{}
		if err := proto.Unmarshal([]byte(item), data); err != nil {
			return nil, err
		}
		ret[k] = data
	}
	return
}

func HMGet(fields ...string) (map[string]*pb.RummyRoomData, error) {
	if len(fields) <= 0 {
		return nil, nil
	}
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return nil, uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()

	results, err := client.HMGet(key, fields...)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]*pb.RummyRoomData)
	for i, field := range fields {
		if results[i] == nil {
			continue
		}
		buf, ok := results[i].([]byte)
		if !ok {
			mlog.Errorf("数据类型不支持: %s", field)
			continue
		}
		item := &pb.RummyRoomData{}
		if err := proto.Unmarshal(buf, item); err != nil {
			mlog.Errorf("数据序列化失败：%s: %v", field, err)
		}
		ret[field] = item
	}
	return ret, nil
}

func HMSet(data map[string]*pb.RummyRoomData) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
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

func HGet(RoomId uint64) (*pb.RummyRoomData, bool, error) {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return nil, false, uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()
	field := GetField(RoomId)

	// 加载数据
	str, err := client.HGet(key, field)
	if err != nil {
		return nil, false, err
	}

	// 解析数据
	data := &pb.RummyRoomData{}
	if err := proto.Unmarshal([]byte(str), data); err != nil {
		return nil, len(str) > 0, err
	}
	return data, len(str) > 0, nil
}

func HSet(RoomId uint64, data *pb.RummyRoomData) error {
	// 获取redis连接
	client := manager.GetRedis(DBNAME)
	if client == nil {
		return uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()
	field := GetField(RoomId)

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
		return uerror.New(1, pb.ErrorCode_CLIENT_NOT_FOUND, "redis数据库不存在: %s", DBNAME)
	}
	key := GetKey()
	return client.HDel(key, fields...)
}
