package dao

import (
	conf "douyin-backend/video-center/conf"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormopentracing "gorm.io/plugin/opentracing"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var isInit = false

func Init() {
	if isInit != false { //确保只执行一次Init
		return
	}
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	username := conf.Viper.GetString("db.mysql.username")
	password := conf.Viper.GetString("db.mysql.password")
	port := conf.Viper.GetString("db.mysql.port")
	host := conf.Viper.GetString("db.mysql.host")
	database := conf.Viper.GetString("db.mysql.database")
	charset := conf.Viper.GetString("db.mysql.charset")
	fmt.Println(host)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger:                 newLogger},
	)
	if err != nil {
		panic(err)
	}

	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}
	isInit = true
}
