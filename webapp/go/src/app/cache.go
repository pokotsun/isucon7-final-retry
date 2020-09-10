package main

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gomodule/redigo/redis"
)

type redisClient struct {
	pool *redis.Pool
}

// protocol is unix or tcp
func NewRedis(protocol, server string) *redisClient {
	pool := &redis.Pool{
		IdleTimeout: 30 * time.Second,
		Wait:        false,
		Dial:        func() (redis.Conn, error) { return redis.Dial(protocol, server) },
	}

	return &redisClient{
		pool: pool,
	}
}

func (rc *redisClient) FetchConn() redis.Conn {
	return rc.pool.Get()
}

func (rc *redisClient) SingleSet(key string, value []byte) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (rc *redisClient) SingleGet(key string) ([]byte, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", key))
}

func (rc *redisClient) MultiSet(set map[string][]byte) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("MSET", redis.Args{}.AddFlat(set)...)
	return err
}

func (rc *redisClient) MultiGet(keys []string) ([][]byte, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("MGET", redis.Args{}.AddFlat(keys)...))
}

func (rc *redisClient) SingleDelete(key string) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (rc *redisClient) MultiDelete(keys []string) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
	return err
}

func (rc *redisClient) SingleSetNX(key string, value []byte) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	ok, err := redis.Int(conn.Do("SETNX", key, value))
	if err != nil {
		return 0, err
	}
	return ok, nil
}

func (rc *redisClient) Increment(key string, delta uint64) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	length, _ := redis.Int(conn.Do("STRLEN", key))
	if length <= 0 {
		return 0, redis.ErrNil
	}

	ok, err := redis.Int(conn.Do("INCRBY", key, delta))
	if err != nil {
		return 0, err
	}
	return ok, nil
}

func (rc *redisClient) Flush() error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}

type memcacheClient struct {
	client *memcache.Client
}

// protocol unused.
// if host contains "/", it will be used unix socket.
func NewMemcache(protocol, host string) *memcacheClient {
	client := memcache.New(fmt.Sprintf("%s", host))
	err := client.Ping()
	fmt.Println(err)

	return &memcacheClient{
		client: client,
	}
}

func (mc *memcacheClient) FetchConn() *memcache.Client {
	return mc.client
}

func (mc *memcacheClient) SingleSet(key string, value []byte) error {
	return mc.client.Set(&memcache.Item{Key: key, Value: value})
}

func (mc *memcacheClient) SingleGet(key string) ([]byte, error) {
	item, err := mc.client.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (mc *memcacheClient) MultiSet(set map[string][]byte) error {
	var err error
	for k, v := range set {
		err = mc.SingleSet(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mc *memcacheClient) MultiGet(keys []string) ([][]byte, error) {
	res := make([][]byte, len(keys))
	itemMap, err := mc.client.GetMulti(keys)
	if err != nil {
		return nil, err
	}
	for key, _ := range itemMap {
		res = append(res, itemMap[key].Value)
	}
	return res, nil
}

func (mc *memcacheClient) SingleDelete(key string) error {
	mc.client.Delete(key)
	return nil
}

func (mc *memcacheClient) MultiDelete(keys []string) error {
	for _, key := range keys {
		mc.SingleDelete(key)
	}

	return nil
}

func (mc *memcacheClient) SingleSetNX(key string, value []byte) (int, error) {
	err := mc.client.Add(&memcache.Item{Key: key, Value: value})
	if err == memcache.ErrNotStored {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (mc *memcacheClient) Increment(key string, delta uint64) (int, error) {
	v, err := mc.client.Increment(key, delta)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (mc *memcacheClient) Flush() error {
	return mc.client.FlushAll()
}
