package main

import (
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
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
)

var ServiceName string
var ServerIp string
var ServerPort int64
var NacosIp string
var NacosPort uint64

func main() {
	conf.InitConfig()
	ServiceName = conf.Viper.GetString("application.ServiceName")
	ServerIp = conf.Viper.GetString("application.Ip")
	ServerPort = conf.Viper.GetInt64("application.Port")
	NacosIp = conf.Viper.GetString("nacos.Ip")
	NacosPort = conf.Viper.GetUint64("nacos.Port")
	cache.Init()
	dao.Init()
	oss.Init()
	RegisterNacos()
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(52428800), //50Mb
		grpc.MaxSendMsgSize(52428800))
	pb.RegisterVideoCenterServer(server, &handler.VideoServer{})
	Sip := fmt.Sprintf("%s:%d", ServerIp, ServerPort)
	listen, err := net.Listen("tcp", Sip)
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
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: NacosIp,
			Port:   NacosPort,
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
		Ip:          ServerIp,
		Port:        uint64(ServerPort),
		ServiceName: ServiceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if !success {
		return
	} else {
		fmt.Println("namingClient.RegisterInstance Success")
	}
}
