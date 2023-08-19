package dao

import "fmt"

func Migrate() {
	Init()
	m := DB.Migrator()
	if m.HasTable(&Fav{}) {
		return
	}
	fmt.Println("创建Fav")
	if err := m.CreateTable(&Fav{}); err != nil {
		panic(err)
	}
	if m.HasTable(&Video{}) {
		return
	}
	if err := m.CreateTable(&Video{}); err != nil {
		panic(err)
	}

}
