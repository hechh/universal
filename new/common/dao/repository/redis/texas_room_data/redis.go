package texas_room_data

import (
	"context"
	"fmt"
	"poker_server/common/dao"
	"poker_server/common/pb"
	"time"

	"github.com/golang/protobuf/proto"
)

func GetKey(uniqueId uint64) string {
	return fmt.Sprintf("texas_room_data_%d", uniqueId)
}

// 房间数据操作
func Get(roomId uint64) (*pb.TexasRoomData, error) {
	val, err := dao.Get("poker", GetKey(roomId))
	if err != nil {
		return nil, err
	}

	item := &pb.TexasRoomData{}
	if err = proto.Unmarshal([]byte(val), item); err != nil {
		return nil, err
	}
	return item, nil
}

func Set(roomId uint64, data *pb.TexasRoomData, ttl time.Duration) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return dao.SetEX("poker", GetKey(roomId), buf, ttl)
}

func MGet(keys ...string) (rets map[uint64]*pb.TexasRoomData, err error) {
	rets = make(map[uint64]*pb.TexasRoomData)
	if len(keys) <= 0 {
		return
	}
	// 批量加载数据
	results, err := dao.MGet("poker", keys...)
	if err != nil {
		return nil, err
	}
	// 解析数据
	for _, vv := range results {
		if vv == nil {
			continue
		}
		item := &pb.TexasRoomData{}
		if err = proto.Unmarshal([]byte(vv.(string)), item); err != nil {
			return
		}
		rets[item.RoomId] = item
	}
	return
}

type DataInfo struct {
	ttl   int64
	value string
}

func MSet(datas map[uint64]*pb.TexasRoomData) error {
	if len(datas) <= 0 {
		return nil
	}

	// 获取redis客户端
	rets := map[string]*DataInfo{}
	for key, data := range datas {
		buf, err := proto.Marshal(data)
		if err != nil {
			return err
		}
		rets[GetKey(key)] = &DataInfo{
			ttl:   data.BaseInfo.FinishTime - data.BaseInfo.CreateTime + 15*60,
			value: string(buf),
		}
	}

	// 批量加载数据
	client := dao.GetRedisClient("poker").GetClient()
	pipe := client.Pipeline()
	defer pipe.Close()

	for key, data := range rets {
		pipe.Set(context.Background(), key, data.value, time.Duration(data.ttl)*time.Second)
	}
	_, err := pipe.Exec(context.Background())
	return err
}
