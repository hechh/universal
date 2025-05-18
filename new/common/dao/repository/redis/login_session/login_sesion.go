package login_session

import (
	"fmt"
	"poker_server/common/dao/internal/manager"
	"poker_server/common/pb"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
)

func GetKey(id string) string {
	return fmt.Sprintf("login_session:%s", id)
}

func Get(id string) (*pb.LoginSession, error) {
	// 获取redis连接
	client := manager.GetRedis("poker")

	// 加载数据
	str, err := client.Get(GetKey(id))
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	// 解析数据
	data := &pb.LoginSession{}
	if err := proto.Unmarshal([]byte(str), data); err != nil {
		return nil, err
	}
	return data, nil
}

func Set(id string, data *pb.LoginSession, ttl time.Duration) error {
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	// 获取redis连接
	client := manager.GetRedis("poker")
	return client.SetEX(GetKey(id), buf, ttl)
}
