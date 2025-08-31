package redis

import (
	"time"
	"universal/common/redis/internal/manager"
	"universal/library/uerror"

	"github.com/go-redis/redis/v8"
)

func IncrBy(dbid string, key string, val int64) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.IncrBy(key, val)
}

func Get(dbid string, key string) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.Get(key)
}

func Set(dbid string, key string, val interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.Set(key, val)
}

// 不存在key时，设置该key的值未val
func SetNX(dbid string, key string, val interface{}) (bool, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return false, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.SetNX(key, val)
}

// 设置key的过期时间
func SetEX(dbid string, key string, val interface{}, ttl time.Duration) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.SetEX(key, val, ttl)
}

// 批量读取操作
func MGet(dbid string, keys ...string) ([]interface{}, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.New(1, -1, "DB(%s)", dbid)
	}
	args := []string{}
	for i := 0; i < len(keys); i++ {
		args = append(args, (keys[i]))
	}
	return cli.MGet(args...)
}

// 批量设置键值对
func MSet(dbid string, args ...interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.MSet(args...)
}

func HGetAll(dbid string, key string) (map[string]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HGetAll(key)
}

func HGet(dbid string, key string, field string) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HGet(key, field)
}

func HDel(dbid string, key string, fields ...string) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HDel(key, fields...)
}

func HKeys(dbid string, key string) ([]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HKeys(key)
}

func HIncrBy(dbid string, key string, field string, incr int64) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HIncrBy(key, field, incr)
}

func HSet(dbid string, key string, field string, val interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HSet(key, field, val)
}

func HMSet(dbid string, key string, vals ...interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.HMSet(key, vals...)
}

func ZAdd(dbid string, key string, members ...*redis.Z) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.ZAdd(key, members...)
}

func ZCard(dbid string, key string) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.ZCard(key)
}

// 返回有序集合中指定成员的排名，有序集合成员按分数值递减排序
func ZRevRank(dbid string, key string, member string) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.ZRevRank(key, member)
}

// 返回有序集合中指定区间内的成员
func ZRevRange(dbid string, key string, start, stop int64) ([]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.ZRevRange(key, start, stop)
}

func ZScore(dbid string, key string, member string) (float64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.ZScore(key, member)
}

func ZRevRangeWithScores(dbid string, key string, start, stop int64) ([]redis.Z, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.ZRevRangeWithScores(key, start, stop)
}

func RPush(dbid string, key string, values ...interface{}) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.RPush(key, values...)
}

func RPop(dbid string, key string) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.RPop(key)
}

func LLen(dbid string, key string) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.LLen(key)
}

func LRange(dbid string, key string, start, stop int64) ([]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.LRange(key, start, stop)
}

func LTrim(dbid string, key string, start, stop int64) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.LTrim(key, start, stop)
}

func LRem(dbid string, key string, val interface{}) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.New(1, -1, "DB(%s)", dbid)
	}
	return cli.LRem(key, val)
}
