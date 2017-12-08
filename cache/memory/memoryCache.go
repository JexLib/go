package memory

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JexLib/golang/cache"
)

type MemoryCache struct {
	defaultExpiration time.Duration
	items             *list.List
	mu                *sync.Mutex
}

func NewMemoryCache(defaultExpiration time.Duration, cleanupInterval ...time.Duration) *MemoryCache {
	c := &MemoryCache{
		defaultExpiration: defaultExpiration,
		items:             list.New(),
		mu:                new(sync.Mutex),
	}

	defCleanupInterval := time.Second * 5
	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		defCleanupInterval = cleanupInterval[0]
	}
	c.gc(defCleanupInterval)

	return c
}

func (mc *MemoryCache) Set(key string, value interface{}, expiration ...time.Duration) error {
	d := mc.defaultExpiration
	if len(expiration) > 0 {
		d = expiration[0]
	}
	item := cache.NewCacheItem(key, value, d)

	defer mc.mu.Unlock()
	mc.mu.Lock()
	if e := mc.find(key); e != nil {
		e.Value = *item
	} else {
		mc.items.PushBack(*item)
	}
	fmt.Println("set:", *item)
	return nil
}

func (mc *MemoryCache) Get(key string) interface{} {
	if e := mc.find(key); e != nil {
		return e.Value.(cache.CacheItem).Val
	}
	return nil
}

func (mc *MemoryCache) Delete(key string) error {
	if e := mc.find(key); e != nil {
		defer mc.mu.Unlock()
		mc.mu.Lock()
		mc.items.Remove(e)
	}
	return nil
}

func (mc *MemoryCache) Keys(prefix ...string) []string {
	slice := []string{}
	for e := mc.items.Front(); e != nil; e = e.Next() {
		key := e.Value.(cache.CacheItem).Key
		if len(prefix) == 0 || strings.HasPrefix(key, prefix[0]) {
			slice = append(slice, key)
		}
	}
	return slice
}

func (mc *MemoryCache) find(key string) *list.Element {
	for e := mc.items.Front(); e != nil; e = e.Next() {
		if e.Value.(cache.CacheItem).Key == key && !e.Value.(cache.CacheItem).Expired() {
			return e
		}
	}
	return nil
}

func (mc *MemoryCache) Count() int64 {
	return int64(mc.items.Len())
}

//清空数据
func (mc *MemoryCache) Clear() error {
	defer mc.mu.Unlock()
	mc.mu.Lock()
	mc.items.Init()
	return nil
}

func (mc *MemoryCache) LoadFromFile(filename string) error {
	return nil
}

func (mc *MemoryCache) SaveToFile(filename string) error {

	// var d1 = []byte(wireteString);
	// err2 := ioutil.WriteFile("./output2.txt", d1, 0666)  //写入文件(字节数组)
	return nil
}

func (mc *MemoryCache) gc(gcInterval time.Duration) {
	go func() {
		for {
			for e := mc.items.Front(); e != nil; e = e.Next() {
				if e.Value.(cache.CacheItem).Expired() {
					mc.mu.Lock()
					mc.items.Remove(e)
					mc.mu.Unlock()
				}
			}
			time.Sleep(gcInterval)
		}
	}()
}
