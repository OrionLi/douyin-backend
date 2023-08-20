package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"net"
	"video-center/cache"
	"video-center/conf"
	"video-center/dao"
	"video-center/handler"
	"video-center/oss"
	"video-center/pkg/pb"
)

func main() {
	conf.InitConfig()
	cache.Init()
	dao.Init()
	oss.Init("D://d", "OssConf.yaml")
	RegisterNacos()
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(52428800), //50Mb
		grpc.MaxSendMsgSize(52428800))
	pb.RegisterVideoCenterServer(server, &handler.VideoServer{})
	listen, err := net.Listen("tcp", "127.0.0.1:8800")
	if err != nil {
		fmt.Println("端口监听有误")
	}
	fmt.Println("端口监听成功")
	err = server.Serve(listen)
	if err != nil {
		return
	}
}
func RegisterNacos() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		//NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}

	// 创建服务发现客户端 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		fmt.Println("clients.NewNamingClient err,", err)
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8800,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})
	if !success {
		return
	} else {
		fmt.Println("namingClient.RegisterInstance Success")
	}
}
