package conf

import "github.com/spf13/viper"

// Init 初始化配置文件与引擎
func Init() error {

	// 设置配置文件的名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// todo: 获取配置项的值并赋值给变量

	return nil
}
