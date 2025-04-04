package manager

import (
	"context"
	"fmt"
	"hego/common/dao/domain"
	"hego/common/dao/internal/redis"
	"hego/common/global"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

var (
	redisPool = make(map[uint32]*redis.RedisClient)
)

func InitRedis(cfgs map[uint32]*global.DbConfig) error {
	if len(cfgs) <= 0 {
		return fmt.Errorf("redis配置为空")
	}
	for dbid, cfg := range cfgs {
		// 建立redis连接
		cli := goredis.NewClient(&goredis.Options{
			IdleTimeout:  1 * time.Minute,
			MinIdleConns: 100,
			ReadTimeout:  -1,
			WriteTimeout: -1,
			Addr:         cfg.Host,
			Password:     cfg.Password,
			OnConnect:    func(ctx context.Context, cn *goredis.Conn) error { return nil },
		})
		// 连接到redis服务器，测试连通性
		if _, err := cli.Ping(context.Background()).Result(); err != nil {
			return fmt.Errorf("Redis connecting is failed, error: %v, cfg: %v", err, cfg)
		}
		redisPool[dbid] = redis.NewRedisClient(cli, cfg.DbName)
	}
	return nil
}

func GetRedis(dbid uint32) *redis.RedisClient {
	return redisPool[dbid]
}

func GetRedisByUID(uid uint64) *redis.RedisClient {
	return GetRedis(uint32(uid&domain.MIN_ACCOUNT_ID/domain.MAX_GROUP_ACCOUNT) + 1)
}
