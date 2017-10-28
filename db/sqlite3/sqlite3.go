package sqlite3

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Config struct {
	DBFile string `flag:"|sqlite.db|sqlite3 db filename"`
}

func Start(conf SQLite3Config, tables ...interface{}) *gorm.DB {
	DB, err := gorm.Open("sqlite3", conf.DBFile)
	if err != nil {
		panic("failed to connect database")
	}

	// //	defer db.Close()

	// // Migrate the schema
	DB.AutoMigrate(tables)
	return DB
}
