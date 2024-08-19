package manager

import (
	"context"
	"fmt"
	"time"
	"universal/common/config"
	"universal/common/dao/domain"
	"universal/common/dao/internal/db"

	"github.com/go-redis/redis/v8"
)

var (
	redisPool = make(map[uint32]*db.RedisClient)
)

func InitRedis(cfgs map[uint32]*config.DbConfig) error {
	for dbid, cfg := range cfgs {
		// 建立redis连接
		cli := redis.NewClient(&redis.Options{
			IdleTimeout:  1 * time.Minute,
			MinIdleConns: 100,
			ReadTimeout:  -1,
			WriteTimeout: -1,
			Addr:         cfg.Host,
			Password:     cfg.Password,
			OnConnect:    func(ctx context.Context, cn *redis.Conn) error { return nil },
		})
		// 连接到redis服务器，测试连通性
		if _, err := cli.Ping(context.Background()).Result(); err != nil {
			return fmt.Errorf("Redis connecting is failed, error: %v, cfg: %v", err, cfg)
		}
		redisPool[dbid] = db.NewRedisClient(cli, cfg.DbName)
	}
	return nil
}

func GetRedis(dbid uint32) *db.RedisClient {
	return redisPool[dbid]
}

func GetRedisByUID(uid uint64) *db.RedisClient {
	return GetRedis(uint32(uid&domain.MIN_ACCOUNT_ID/domain.MAX_GROUP_ACCOUNT) + 1)
}
