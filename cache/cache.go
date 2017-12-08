package cache

import "time"

var (
	//默认清除间隔时间
	defCleanupInterval = time.Second * 5
)

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{}, expiration ...time.Duration) error
	Delete(key string) error
	Clear() error
	Count() int64
	Keys(prefix ...string) []string
	SaveToFile(filename string) error
	LoadFromFile(filename string) error
}

type CacheItem struct {
	Key        string
	Val        interface{}
	Expiration int64 `json:"-"`
}

func NewCacheItem(key string, val interface{}, d time.Duration) *CacheItem {
	return &CacheItem{
		Key:        key,
		Val:        val,
		Expiration: time.Now().Add(d).UnixNano(),
	}
}

//是否过期
func (ci CacheItem) Expired() bool {
	return ci.Expiration > 0 && ci.Expiration < time.Now().UnixNano()
}
