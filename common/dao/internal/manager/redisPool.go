package manager

import (
	"time"
	"universal/common/pb"
	"universal/framework/fbasic"

	"github.com/gomodule/redigo/redis"
)

var (
	redisPool = make(map[string]*redis.Pool)
)

func InitRedis(name, user, passwd, addr string) error {
	if _, ok := redisPool[name]; ok {
		return fbasic.NewUError(1, pb.ErrorCode_HasRegistered, name, addr)
	}
	// 连接参数
	opts := []redis.DialOption{}
	if len(user) > 0 {
		opts = append(opts, redis.DialUsername(user), redis.DialPassword(passwd))
	}
	// 建立连接池
	pool := redis.Pool{
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr, opts...) },
		IdleTimeout: 3 * time.Second, // 连接超时时间
		MaxIdle:     100,             // 最大空闲数
		MaxActive:   800,             // 限制最大连接数量
		Wait:        true,            // 最大连接超出限制时，需要等待
	}
	// 测试联通性
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_RedisPing, err)
	}
	redisPool[name] = &pool
	return nil
}

func GetRedis(name string) (conn redis.Conn, err error) {
	pool, ok := redisPool[name]
	if !ok {
		return nil, fbasic.NewUError(1, pb.ErrorCode_NotSupported, name)
	}
	conn = pool.Get()
	err = conn.Err()
	return
}

func PutRedis(conn redis.Conn) {
	if conn != nil {
		conn.Close()
	}
}
