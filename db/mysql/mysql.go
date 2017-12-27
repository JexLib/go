package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

type Config struct {
	Address     string `flag:"|192.168.2.188:3306|mysql host:port"`
	User        string `flag:"|root|mysql user"`
	Password    string
	DataBase    string `flag:"|test|DataBase name"`
	TablePrefix string `flag:"||table name  Prefix"`
	TableSuffix string `flag:"||table name  Suffix"`
}

var tables []interface{}

func RegisterTables(table ...interface{}) {
	tables = append(tables, table...)
}

func autoCreateDataBase(mDB *gorm.DB, dbname string) {
	mDB.Exec("Create Database If Not Exists " + dbname + " Character Set UTF8")
	mDB.Exec("Use " + dbname)
}

func NewMySQL(debug bool, conf Config) *gorm.DB {
	mysqlStr := fmt.Sprintf("%s:%s@tcp(%s)/mysql?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Password, conf.Address)

	//	mysqlStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Password, conf.Address, conf.DataBase)
	mDB, err := gorm.Open("mysql", mysqlStr)

	if err != nil {
		panic("failed to connect databas e")
	}

	if debug {
		mDB = mDB.Debug()
	}

	autoCreateDataBase(mDB, conf.DataBase)
	// //	defer db.Close()
	// DB.SetMaxIdleConns(10)
	mDB.DB().SetMaxOpenConns(100)

	if len(conf.TablePrefix) > 0 || len(conf.TableSuffix) > 0 {
		// 增加表名前缀/后缀
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return conf.TablePrefix + defaultTableName + conf.TableSuffix
		}
	}

	// // Migrate the schema
	mDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(tables...)
	DB = mDB
	return mDB
}
