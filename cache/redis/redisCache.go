package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	defaultExpiration time.Duration
	client            *redis.Client
	prefix            string
}

func NewRedisCache(options redis.Options, prefix string, defaultExpiration time.Duration, cleanupInterval ...time.Duration) *RedisCache {
	client := redis.NewClient(&options)

	redisCache := &RedisCache{
		defaultExpiration: defaultExpiration,
		client:            client,
	}

	return redisCache
}

func (rc *RedisCache) formatKey(key string) string {
	if len(rc.prefix) > 0 {
		return rc.prefix + "_" + key
	}
	return key
}

func (rc *RedisCache) Set(key string, value interface{}, expiration ...time.Duration) error {
	d := rc.defaultExpiration
	if len(expiration) > 0 {
		d = expiration[0]
	}
	return rc.client.Set(rc.formatKey(key), value, d).Err()
}

func (rc *RedisCache) Get(key string) interface{} {
	return rc.client.Get(rc.formatKey(key)).Val
}

func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(rc.formatKey(key)).Err()
}

func (rc *RedisCache) Count() int64 {
	return rc.client.DBSize().Val()
}

func (rc *RedisCache) Keys(prefix ...string) []string {
	if len(prefix) == 0 {
		prefix = append(prefix, "*")
	}
	return rc.client.Keys(prefix[0]).Val()
}

func (rc *RedisCache) Clear() error {
	return rc.client.FlushAllAsync().Err()
}

func (mc *RedisCache) LoadFromFile(filename string) error {
	return nil
}

func (mc *RedisCache) SaveToFile(filename string) error {
	return nil
}

func (rc *RedisCache) gc(gcInterval time.Duration) {
	go func() {
	}()
}
