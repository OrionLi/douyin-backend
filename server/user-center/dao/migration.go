package dao

import (
	"user-center/model"
)

// migration 自动迁移数据库表结构
func migration() error {
	err := db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Relation{},
		) //自动创建或更新数据库表结构
	if err != nil {
		return err
	}
	return nil
}
