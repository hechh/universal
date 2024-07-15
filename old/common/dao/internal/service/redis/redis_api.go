package redis

import (
	"universal/common/dao/internal/manager"

	"github.com/gomodule/redigo/redis"
)

func Do(dbname, cmd string, args ...interface{}) (reply interface{}, err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	reply, err = cli.Do(cmd, args...)
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func Expire(dbname, key string, ttl int) (err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	_, err = cli.Do("EXPIRE", key, ttl)
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func Exists(dbname, key string) (ok bool, err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	ok, err = redis.Bool(cli.Do("EXISTS", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func TTL(dbname, key string) (expire int, err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	expire, err = redis.Int(cli.Do("TTL", key))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func IncrBy(dbname, key string, incr int) (count int, err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	count, err = redis.Int(cli.Do("INCRBY", key, incr))
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func Del(dbname string, keys ...string) (err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	args := redis.Args{}.AddFlat(keys)
	_, err = cli.Do("DEL", args...)
	if err == redis.ErrNil {
		err = nil
	}
	return
}

func Get(dbname, key string) (data []byte, err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	data, err = redis.Bytes(cli.Do("GET", key))
	if err == redis.ErrNil {
		data, err = nil, nil
	}
	return
}

func Set(dbname, key string, value interface{}) (err error) {
	cli, err := manager.GetRedis(dbname)
	defer manager.PutRedis(cli)
	if err != nil {
		return
	}
	_, err = cli.Do("SET", key, value)
	if err == redis.ErrNil {
		err = nil
	}
	return
}
