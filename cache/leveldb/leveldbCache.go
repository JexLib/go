package leveldb

import (
	"time"

	"github.com/JexLib/golang/cache"
	"github.com/JexLib/golang/utils"
	"github.com/golang/leveldb"
	"github.com/golang/leveldb/db"
	"github.com/golang/leveldb/memfs"
)

type LevelDBCache struct {
	defaultExpiration time.Duration
	db                *leveldb.DB
}

func NewCache_LevelDBFile(file string, defaultExpiration time.Duration, cleanupInterval ...time.Duration) *LevelDBCache {
	db, err := leveldb.Open(file, &db.Options{})
	if err != nil {
		panic(err)
	}
	c := &LevelDBCache{
		defaultExpiration: defaultExpiration,
		db:                db,
	}

	defCleanupInterval := time.Second * 5
	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		defCleanupInterval = cleanupInterval[0]
	}
	c.gc(defCleanupInterval)

	return c
}

func NewCache_LevelDBMem(defaultExpiration time.Duration, cleanupInterval ...time.Duration) *LevelDBCache {
	db, err := leveldb.Open("", &db.Options{
		FileSystem: memfs.New(),
	})
	if err != nil {
		panic(err)
	}
	c := &LevelDBCache{
		defaultExpiration: defaultExpiration,
		db:                db,
	}
	defCleanupInterval := time.Second * 5
	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		defCleanupInterval = cleanupInterval[0]
	}
	c.gc(defCleanupInterval)

	return c
}

func (rc *LevelDBCache) Set(key string, value interface{}, expiration ...time.Duration) error {
	d := rc.defaultExpiration
	if len(expiration) > 0 {
		d = expiration[0]
	}
	item := cache.NewCacheItem(key, value, d)
	bs, err := utils.Encode(item)
	if err != nil {
		return err
	}
	return rc.db.Set([]byte(key), bs, nil)
}

func (rc *LevelDBCache) Get(key string) interface{} {
	if bs, err := rc.db.Get([]byte(key), nil); err == nil {
		item := cache.CacheItem{}
		if utils.Decode(bs, &item) == nil {
			return item.Val
		}
	}
	return nil
}

func (rc *LevelDBCache) Delete(key string) error {
	return rc.db.Delete([]byte(key), nil)
}

func (rc *LevelDBCache) Keys(prefix ...string) []string {
	if len(prefix) == 0 {
		prefix = append(prefix, "*")
	}
	return nil

}

func (rc *LevelDBCache) Count() int64 {
	return 0
}

func (rc *LevelDBCache) Clear() error {
	return nil
}

func (mc *LevelDBCache) LoadFromFile(filename string) error {
	return nil
}

func (mc *LevelDBCache) SaveToFile(filename string) error {
	return nil
}

func (rc *LevelDBCache) gc(gcInterval time.Duration) {
	go func() {
	}()
}
