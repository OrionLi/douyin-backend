package dao

import (
	conf "douyin-backend/video-center/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

var DB *gorm.DB
var isInit = false

func Init() {
	if isInit != false { //确保只执行一次Init
		return
	}
	//通过viper读取配置
	conf.InitConfig()
	dsn := conf.Viper.GetString("db.mysql.dsn")
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
		},
	)
	if err != nil {
		panic(err)
	}

	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}
	isInit = true
}

// Init() 使用以下格式
//application:
//  name: video-center
//
//db:
//  mysql:
//    host: your URL
//    port: 3306
//    username: your Name
//    password: your Password
//    database: DouyinDB
//    charset: utf8mb4
//  redis:
//    host: 127.0.0.1
//    port:  6067
//    db: 0
//    passwd: 123456
//func Init() {
//	if isInit != false { //确保只执行一次Init
//		return
//	}
//	//通过viper读取配置
//	conf.InitConfig()
//	username := conf.Viper.GetString("db.mysql.username")
//	password := conf.Viper.GetString("db.mysql.password")
//	port := conf.Viper.GetString("db.mysql.port")
//	host := conf.Viper.GetString("db.mysql.host")
//	database := conf.Viper.GetString("db.mysql.database")
//	charset := conf.Viper.GetString("db.mysql.charset")
//	fmt.Println(host)
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
//		username,
//		password,
//		host,
//		port,
//		database,
//		charset)
//	var err error
//	DB, err = gorm.Open(mysql.New(mysql.Config{
//		DSN:                       dsn,
//		DefaultStringSize:         256,
//		DisableDatetimePrecision:  true,
//		DontSupportRenameIndex:    true,
//		DontSupportRenameColumn:   true,
//		SkipInitializeWithVersion: false,
//	}),
//		&gorm.Config{
//			PrepareStmt:            true,
//			SkipDefaultTransaction: true,
//		},
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	if err = DB.Use(gormopentracing.New()); err != nil {
//		panic(err)
//	}
//	isInit = true
//}
