package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	preKey string
}

func NewRedisClient(cli *redis.Client, key string) *RedisClient {
	return &RedisClient{client: cli, preKey: key}
}

func (d *RedisClient) getKey(key string) string {
	return fmt.Sprintf("%s_%s", d.preKey, key)
}

func (d *RedisClient) IncrBy(dbid uint32, key string, val int64) (ret int64, err error) {
	ret, err = d.client.IncrBy(context.Background(), d.getKey(key), val).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) Get(dbid uint32, key string) (str string, err error) {
	str, err = d.client.Get(context.Background(), d.getKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) Set(dbid uint32, key string, val interface{}) (err error) {
	_, err = d.client.Set(context.Background(), d.getKey(key), val, 0).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 不存在key时，设置该key的值未val
func (d *RedisClient) SetNX(dbid uint32, key string, val interface{}) (exist bool, err error) {
	exist, err = d.client.SetNX(context.Background(), d.getKey(key), val, 0).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 设置key的过期时间
func (d *RedisClient) SetEX(dbid uint32, key string, val interface{}, ttl time.Duration) (err error) {
	_, err = d.client.SetEX(context.Background(), d.getKey(key), val, ttl).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 批量读取操作
func (d *RedisClient) MGet(dbid uint32, keys ...string) (rets []interface{}, err error) {
	args := []string{}
	for i := 0; i < len(keys); i++ {
		args = append(args, d.getKey(keys[i]))
	}
	rets, err = d.client.MGet(context.Background(), args...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 批量设置键值对
func (d *RedisClient) MSet(dbid uint32, args ...interface{}) (err error) {
	_, err = d.client.MSet(context.Background(), args...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HGetAll(dbid uint32, key string) (ret map[string]string, err error) {
	ret, err = d.client.HGetAll(context.Background(), d.getKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HGet(dbid uint32, key string, field string) (ret string, err error) {
	ret, err = d.client.HGet(context.Background(), d.getKey(key), field).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HDel(dbid uint32, key string, fields ...string) (err error) {
	_, err = d.client.HDel(context.Background(), d.getKey(key), fields...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HKeys(dbid uint32, key string) (rets []string, err error) {
	rets, err = d.client.HKeys(context.Background(), d.getKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HIncrBy(dbid uint32, key string, field string, incr int64) (err error) {
	_, err = d.client.HIncrBy(context.Background(), d.getKey(key), field, incr).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HSet(dbid uint32, key string, field string, val interface{}) (err error) {
	_, err = d.client.HSet(context.Background(), d.getKey(key), field, val).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) HMSet(dbid uint32, key string, vals ...interface{}) (err error) {
	_, err = d.client.HMSet(context.Background(), d.getKey(key), vals...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) ZAdd(dbid uint32, key string, members ...*redis.Z) (err error) {
	_, err = d.client.ZAdd(context.Background(), d.getKey(key), members...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) ZCard(dbid uint32, key string) (count int64, err error) {
	count, err = d.client.ZCard(context.Background(), d.getKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 返回有序集合中指定成员的排名，有序集合成员按分数值递减排序
func (d *RedisClient) ZRevRank(dbid uint32, key string, member string) (count int64, err error) {
	count, err = d.client.ZRevRank(context.Background(), d.getKey(key), member).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 返回有序集合中指定区间内的成员
func (d *RedisClient) ZRevRange(dbid uint32, key string, start, stop int64) (members []string, err error) {
	members, err = d.client.ZRevRange(context.Background(), d.getKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) ZScore(dbid uint32, key string, member string) (score float64, err error) {
	score, err = d.client.ZScore(context.Background(), d.getKey(key), member).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) ZRevRangeWithScores(dbid uint32, key string, start, stop int64) (members []redis.Z, err error) {
	members, err = d.client.ZRevRangeWithScores(context.Background(), d.getKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) RPush(dbid uint32, key string, values ...interface{}) (ret int64, err error) {
	ret, err = d.client.RPush(context.Background(), d.getKey(key), values...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) RPop(dbid uint32, key string) (ret string, err error) {
	ret, err = d.client.RPop(context.Background(), d.getKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) LLen(dbid uint32, key string) (ret int64, err error) {
	ret, err = d.client.LLen(context.Background(), d.getKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) LRange(dbid uint32, key string, start, stop int64) (ret []string, err error) {
	ret, err = d.client.LRange(context.Background(), d.getKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) LTrim(dbid uint32, key string, start, stop int64) (ret string, err error) {
	ret, err = d.client.LTrim(context.Background(), d.getKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *RedisClient) LRem(dbid uint32, key string, val interface{}) (ret int64, err error) {
	ret, err = d.client.LRem(context.Background(), d.getKey(key), 0, val).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}
