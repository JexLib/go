package levelDB

import (
	"github.com/golang/leveldb"
	"github.com/golang/leveldb/db"
	"github.com/golang/leveldb/memfs"
)

type JexLevelDB struct {
	db   *leveldb.DB
	Bath jexLevelBath
}

type jexLevelBath struct {
	bath leveldb.Batch
	db   *leveldb.DB
}

func NewLevelDBFromFile(file string) *JexLevelDB {
	db, err := leveldb.Open(file, &db.Options{})
	if err != nil {
		panic(err)
	}
	return &JexLevelDB{
		db:   db,
		Bath: jexLevelBath{db: db},
	}

}

func NewLevelDBFormMem() *JexLevelDB {
	db, err := leveldb.Open("", &db.Options{
		FileSystem: memfs.New(),
	})
	if err != nil {
		panic(err)
	}
	return &JexLevelDB{
		db:   db,
		Bath: jexLevelBath{db: db},
	}

}

func (ldb *JexLevelDB) Set(key string, val string) error {
	// fmt.Println("===>", key, val)
	return ldb.db.Set([]byte(key), []byte(val), nil)
}

func (ldb *JexLevelDB) Get(key string) (string, error) {
	bytes, err := ldb.db.Get([]byte(key), nil)
	// fmt.Println("<<<<<", key, string(bytes))
	return string(bytes), err
}

func (ldb *JexLevelDB) Delete(key string) error {
	return ldb.db.Delete([]byte(key), nil)
}

func (bt jexLevelBath) Set(key, val string) {
	bt.bath.Set([]byte(key), []byte(val))
}

func (bt jexLevelBath) Delete(key string) {
	bt.bath.Delete([]byte(key))
}

func (bt jexLevelBath) Apply() error {
	return bt.db.Apply(bt.bath, nil)
}
