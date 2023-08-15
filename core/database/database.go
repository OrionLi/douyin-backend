package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	//"user-center-go/models"
)

// 如果你要使用 MySQL 数据库，请修改下面的配置信息，并在main.go中改为调用InitMysqlDB
const (
	USER   = "root"
	PASS   = "123456"
	HOST   = "127.0.0.1"
	PORT   = "3306"
	DBNAME = "db1"
)

var DB *gorm.DB

func InitMysqlDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	migrateTables()
}

func InitSqliteDB() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	migrateTables()
	//insertInitialData()
}

func migrateTables() {
	//err := DB.AutoMigrate(&models.User{})
	//if err != nil {
	//	log.Fatalf("Failed to migrate database tables: %v", err)
	//}
}

//func insertInitialData() {
//	initialData := []models.User{
//		{ID: 1, Username: "Alice", Account: "10001"},
//		{ID: 2, Username: "Bob", Account: "10002"},
//		{ID: 3, Username: "Cathy", Account: "10003"},
//		{ID: 4, Username: "Dave", Account: "10004"},
//		{ID: 5, Username: "Eric", Account: "10005"},
//		{ID: 6, Username: "Frank", Account: "10006"},
//		{ID: 7, Username: "Gary", Account: "10007"},
//		{ID: 8, Username: "Helen", Account: "10008"},
//		{ID: 9, Username: "Irene", Account: "10009"},
//		{ID: 10, Username: "Jack", Account: "10010"},
//		{ID: 11, Username: "Kate", Account: "10011"},
//		{ID: 12, Username: "Lily", Account: "10012"},
//		{ID: 13, Username: "Mike", Account: "10013"},
//		{ID: 14, Username: "Nancy", Account: "10014"},
//		{ID: 15, Username: "Olivia", Account: "10015"},
//		{ID: 16, Username: "Penny", Account: "10016"},
//		{ID: 17, Username: "Qun", Account: "10017"},
//		// 添加更多初始数据...
//	}
//
//	for _, data := range initialData {
//		result := DB.Create(&data)
//		if result.Error != nil {
//			log.Fatal(result.Error)
//		}
//	}
//}
