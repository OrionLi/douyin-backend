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
	dbUrl := conf.Viper.GetString("db.mysql.url")
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbUrl,
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
