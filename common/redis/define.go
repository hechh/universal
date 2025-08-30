package redis

import (
	"universal/common/redis/internal/client"
	"universal/common/redis/internal/manager"
	"universal/common/yaml"
)

// 定义redis数据库枚举
const (
	REDIS_DB_DEFAULT = "universal"
)

func Init(cfgs map[string]*yaml.DbConfig) error {
	return manager.InitRedis(cfgs)
}

func GetClient(db string) *client.RedisClient {
	return manager.GetRedis(db)
}
