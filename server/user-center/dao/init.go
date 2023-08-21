package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB

func Database(conn string) error {
	var ormLogger logger.Interface

	// 根据运行模式设置日志级别
	switch gin.Mode() {
	case "debug":
		ormLogger = logger.Default.LogMode(logger.Info)
	case "release":
		ormLogger = logger.Default.LogMode(logger.Warn)
	default:
		ormLogger = logger.Default
	}
	_db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conn,  // 数据源名称，包含用户名，密码，主机，端口，数据库等信息
		DefaultStringSize:         256,   // 设置默认的字符串长度为256
		DisableDatetimePrecision:  true,  // 禁用日期时间的精度，例如不使用毫秒
		DontSupportRenameIndex:    true,  // 不支持重命名索引，如果需要修改索引，需要先删除再创建
		DontSupportRenameColumn:   true,  // 不支持重命名列，如果需要修改列，需要先删除再创建
		SkipInitializeWithVersion: false, // 不跳过根据版本初始化数据库
	}), &gorm.Config{
		Logger: ormLogger, // 设置gorm的日志记录器，用于打印SQL语句等信息
		NamingStrategy: schema.NamingStrategy{ // 设置gorm的命名策略，用于映射模型和数据库表的名称
			SingularTable: true, // 使用单数形式的表名，例如user而不是users
		},
	})
	if err != nil {
		return err
	}
	sqlDB, _ := _db.DB()                       // 获取底层的sql.DB对象
	sqlDB.SetConnMaxLifetime(20)               //配置连接池
	sqlDB.SetConnMaxIdleTime(time.Second * 30) // 设置连接的最大存活时

	db = _db
	//自动迁移数据表结构
	err = migration()
	if err != nil {
		return err
	}
	return nil
}

func NewDBClient(ctx context.Context) *gorm.DB {
	_db := db
	return _db.WithContext(ctx)
}
