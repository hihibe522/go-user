package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func init() {
	const (
		UserName     string = "bebe"
		Password     string = "qwe123"
		Addr         string = "127.0.0.1"
		Port         int    = 3306
		Database     string = "test"
		MaxLifetime  int    = 10
		MaxOpenConns int    = 10
		MaxIdleConns int    = 10
	)
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", UserName, Password, Addr, Port, Database)
	//連接MySQL
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("db connect error" + err.Error())
	}
	fmt.Println("mysql connect success")

	dbSQL, errdb := db.DB()

	if errdb != nil {
		panic("get db failed" + errdb.Error())
	}
	//設置最大連接數
	dbSQL.SetMaxOpenConns(MaxOpenConns)
	//設置最大閒置數
	dbSQL.SetMaxIdleConns(MaxIdleConns)
	//設置每個連接的過期時間
	dbSQL.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
}
