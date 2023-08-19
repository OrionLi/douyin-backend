package conf

import (
	"github.com/spf13/viper"
	"log"
)

var Viper *viper.Viper

func InitConfig() {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./conf")
	v.SetConfigType("yaml")
	v.SetConfigName("DalConf.yml")

	if err := v.ReadInConfig(); err == nil {
		log.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		panic(err)
	}
	Viper = v
}
