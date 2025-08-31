package manager

import (
	"context"
	"time"
	"universal/common/redis/internal/client"
	"universal/common/yaml"

	goredis "github.com/go-redis/redis/v8"
)

var (
	redisPool = make(map[string]*client.RedisClient)
)

func InitRedis(cfgs map[string]*yaml.DbConfig) error {
	for _, cfg := range cfgs {
		// 建立redis连接
		cli := goredis.NewClient(&goredis.Options{
			IdleTimeout:  1 * time.Minute,
			MinIdleConns: 100,
			DB:           int(cfg.Db),
			ReadTimeout:  -1,
			WriteTimeout: -1,
			Addr:         cfg.Host,
			Username:     cfg.User,
			Password:     cfg.Password,
			OnConnect:    func(ctx context.Context, cn *goredis.Conn) error { return nil },
		})
		// 连接到redis服务器，测试连通性
		if _, err := cli.Ping(context.Background()).Result(); err != nil {
			return err
		}
		redisPool[cfg.DbName] = client.NewRedisClient(cli, cfg.Prefix)
	}
	return nil
}

func GetRedis(db string) *client.RedisClient {
	return redisPool[db]
}
