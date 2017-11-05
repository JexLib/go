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
	Count() int
	SaveToFile(filename string) error
	LoadFromFile(filename string) error
	gc(gcInterval time.Duration)
}

type _CacheItem struct {
	Key        string
	Val        interface{}
	Expiration int64 `json:"-"`
}

func new_CacheItem(key string, val interface{}, d time.Duration) *_CacheItem {
	return &_CacheItem{
		Key:        key,
		Val:        val,
		Expiration: time.Now().Add(d).UnixNano(),
	}
}

//是否过期
func (ci _CacheItem) Expired() bool {
	return ci.Expiration > 0 && ci.Expiration < time.Now().UnixNano()
}
