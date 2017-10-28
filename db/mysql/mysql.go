package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	Address     string `flag:"|192.168.2.188:3306|mysql host:port"`
	User        string `flag:"|root|mysql user"`
	Password    string
	DataBase    string `flag:"|test|DataBase name"`
	TablePrefix string `flag:"||table name  Prefix"`
}

func Start(debug bool, conf Config, tables ...interface{}) *gorm.DB {
	mysqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Password, conf.Address, conf.DataBase)
	DB, err := gorm.Open("mysql", mysqlStr)

	if err != nil {
		panic("failed to connect database")
	}

	if debug {
		DB = DB.Debug()
	}
	// //	defer db.Close()

	// // Migrate the schema
	DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(tables)
	return DB
}
