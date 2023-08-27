package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"user-center/conf"
	"user-center/pkg/util"
)

func main() {
	RegisterNacos()
	//初始化配置文件
	err := conf.Init()
	if err != nil {
		util.LogrusObj.Error("<Main> : ", err)
	}

}
func RegisterNacos() {
	// 创建clientConfig，用于配置客户端行为
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,    // 设置超时时间为5秒
		NotLoadCacheAtStart: true,    // 在启动时不加载缓存
		LogLevel:            "debug", // 设置日志级别为调试模式
	}

	// 至少一个ServerConfig，指定Nacos服务器的地址和端口
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: conf.NacosIp,   // Nacos服务器的IP地址
			Port:   conf.NacosPort, // Nacos服务器的端口
		},
	}

	// 创建服务发现客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig, // 设置客户端配置
			ServerConfigs: serverConfigs, // 设置Nacos服务器配置
		},
	)
	if err != nil {
		fmt.Println("clients.NewNamingClient err,", err) // 输出错误信息
	}

	// 注册实例到Nacos服务中
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          conf.ServerIp,           // 注册实例的IP地址
		Port:        uint64(conf.ServerPort), // 注册实例的端口
		ServiceName: conf.ServiceName,        // 服务的名称
		Weight:      10,                      // 权重为10
		Enable:      true,                    // 设置实例为可用状态
		Healthy:     true,                    // 设置实例为健康状态
		Ephemeral:   true,                    // 设置实例为临时实例
	})
	if !success {
		return // 注册失败则退出函数
	} else {
		fmt.Println("namingClient.RegisterInstance Success") // 输出注册成功信息
	}
}
