package dao

import (
	conf "douyin-backend/video-center/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

var DB *gorm.DB

func Init() {

	//通过viper读取配置
	conf.InitConfig()
	dbUrl := conf.Viper.GetString("db.mysql.url")
	var err error
	DB, err = gorm.Open(mysql.Open(dbUrl),
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

	//m := DB.Migrator()
	//if m.HasTable(&Video{}) {
	//	return
	//}
	//
	//if err = m.CreateTable(&Video{}); err != nil {
	//	panic(err)
	//}
}
