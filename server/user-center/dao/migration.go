package dao

import (
	"fmt"
	"user-center/model"
)

func migration() {
	err := db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
		) //自动创建或更新数据库表结构
	if err != nil {
		fmt.Println("err:", err)
	}
	return
}
