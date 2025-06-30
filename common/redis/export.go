package redis

import (
	"time"
	"universal/common/pb"
	"universal/common/redis/internal/manager"
	"universal/library/uerror"

	"github.com/go-redis/redis/v8"
)

func IncrBy(dbid int32, key string, val int64) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.IncrBy(key, val)
}

func Get(dbid int32, key string) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.Get(key)
}

func Set(dbid int32, key string, val interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.Set(key, val)
}

// 不存在key时，设置该key的值未val
func SetNX(dbid int32, key string, val interface{}) (bool, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return false, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.SetNX(key, val)
}

// 设置key的过期时间
func SetEX(dbid int32, key string, val interface{}, ttl time.Duration) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.SetEX(key, val, ttl)
}

// 批量读取操作
func MGet(dbid int32, keys ...string) ([]interface{}, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	args := []string{}
	for i := 0; i < len(keys); i++ {
		args = append(args, (keys[i]))
	}
	return cli.MGet(args...)
}

// 批量设置键值对
func MSet(dbid int32, args ...interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.MSet(args...)
}

func HGetAll(dbid int32, key string) (map[string]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HGetAll(key)
}

func HGet(dbid int32, key string, field string) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HGet(key, field)
}

func HDel(dbid int32, key string, fields ...string) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HDel(key, fields...)
}

func HKeys(dbid int32, key string) ([]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HKeys(key)
}

func HIncrBy(dbid int32, key string, field string, incr int64) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HIncrBy(key, field, incr)
}

func HSet(dbid int32, key string, field string, val interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HSet(key, field, val)
}

func HMSet(dbid int32, key string, vals ...interface{}) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.HMSet(key, vals...)
}

func ZAdd(dbid int32, key string, members ...*redis.Z) error {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.ZAdd(key, members...)
}

func ZCard(dbid int32, key string) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.ZCard(key)
}

// 返回有序集合中指定成员的排名，有序集合成员按分数值递减排序
func ZRevRank(dbid int32, key string, member string) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.ZRevRank(key, member)
}

// 返回有序集合中指定区间内的成员
func ZRevRange(dbid int32, key string, start, stop int64) ([]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.ZRevRange(key, start, stop)
}

func ZScore(dbid int32, key string, member string) (float64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.ZScore(key, member)
}

func ZRevRangeWithScores(dbid int32, key string, start, stop int64) ([]redis.Z, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.ZRevRangeWithScores(key, start, stop)
}

func RPush(dbid int32, key string, values ...interface{}) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.RPush(key, values...)
}

func RPop(dbid int32, key string) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.RPop(key)
}

func LLen(dbid int32, key string) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.LLen(key)
}

func LRange(dbid int32, key string, start, stop int64) ([]string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return nil, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.LRange(key, start, stop)
}

func LTrim(dbid int32, key string, start, stop int64) (string, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return "", uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.LTrim(key, start, stop)
}

func LRem(dbid int32, key string, val interface{}) (int64, error) {
	cli := manager.GetRedis(dbid)
	if cli == nil {
		return 0, uerror.N(1, int32(pb.ErrorCode_RedisClientNotFound), "DB(%d)", dbid)
	}
	return cli.LRem(key, val)
}
