package dao

func Migrate() {
	Init()
	m := DB.Migrator()
	if !m.HasTable(&Fav{}) {
		if err := m.CreateTable(&Fav{}); err != nil {
			panic(err)
		}
	}
	if !m.HasTable(&Video{}) {
		if err := m.CreateTable(&Video{}); err != nil {
			panic(err)
		}
	}
}
