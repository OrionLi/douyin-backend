package dao

func Migrate() {
	Init()
	m := DB.Migrator()
	if m.HasTable(&Video{}) {
		return
	}
	if err := m.CreateTable(&Video{}); err != nil {
		panic(err)
	}
}
