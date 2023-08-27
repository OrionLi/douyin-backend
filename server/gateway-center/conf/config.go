package conf

import (
	"github.com/spf13/viper"
)

// nacos配置项
var (
	NacosAddress string
	NacosPort    int
)

// 各个微服务的服务名
var (
	ChatCenterServiceName  string
	UserCenterServiceName  string
	VideoCenterServiceName string
)

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
	// 解析 Nacos 配置
	NacosAddress = viper.GetString("nacos.Ip")
	NacosPort = viper.GetInt("nacos.Port")

	// 解析各个微服务的服务名
	ChatCenterServiceName = viper.GetString("application.chat-center.ServiceName")
	UserCenterServiceName = viper.GetString("application.user-center.ServiceName")
	VideoCenterServiceName = viper.GetString("application.video-center.ServiceName")

	return nil
}
