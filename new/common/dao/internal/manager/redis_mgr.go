package manager

import (
	"context"
	"fmt"
	"poker_server/common/dao/internal/redis"
	"poker_server/common/yaml"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

var (
	redisPool = make(map[string]*redis.RedisClient)
)

func InitRedis(cfgs map[int32]*yaml.RedisConfig) error {
	if len(cfgs) <= 0 {
		return fmt.Errorf("redis配置为空")
	}
	for _, cfg := range cfgs {
		// 建立redis连接
		cli := goredis.NewClient(&goredis.Options{
			IdleTimeout:  1 * time.Minute,
			MinIdleConns: 100,
			DB:           cfg.Db,
			ReadTimeout:  -1,
			WriteTimeout: -1,
			Addr:         cfg.Host,
			Username:     cfg.User,
			Password:     cfg.Password,
			OnConnect:    func(ctx context.Context, cn *goredis.Conn) error { return nil },
		})
		// 连接到redis服务器，测试连通性
		if _, err := cli.Ping(context.Background()).Result(); err != nil {
			return fmt.Errorf("Redis connecting is failed, error: %v, cfg: %v", err, cfg)
		}
		redisPool[cfg.DbName] = redis.NewRedisClient(cli, cfg.DbName)
	}
	return nil
}

func GetRedis(dbid string) *redis.RedisClient {
	return redisPool[dbid]
}
