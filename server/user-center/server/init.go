package server

import (
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"net"
	"user-center/pkg/util"
)

// Grpc 启动grpc服务的
func Grpc(addr string) error {
	addr = ":" + addr
	listen, _ := net.Listen("tcp", addr)
	// 创建一个grpc服务，并设置不安全的证书 todo: 后期改用STL
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &UserRPCServer{})      //注册用户服务
	pb.RegisterRelationServiceServer(grpcServer, &RelationServer{}) //注册关注服务

	//启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve%v", err)
		util.LogrusObj.Error("Service startup error ", err)
		return err
	}
	fmt.Println("启动成功")
	return nil
}

func RegisterNacos(
	serverIp,
	serviceName,
	nacosIp string,
	nacosPort,
	serverPort uint64) {
	// 创建clientConfig，用于配置客户端行为
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,    // 设置超时时间为5秒
		NotLoadCacheAtStart: true,    // 在启动时不加载缓存
		LogLevel:            "debug", // 设置日志级别为调试模式
	}

	// 至少一个ServerConfig，指定Nacos服务器的地址和端口
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosIp,   // Nacos服务器的IP地址
			Port:   nacosPort, // Nacos服务器的端口
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
		Ip:          serverIp,    // 注册实例的IP地址
		Port:        serverPort,  // 注册实例的端口
		ServiceName: serviceName, // 服务的名称
		Weight:      10,          // 权重为10
		Enable:      true,        // 设置实例为可用状态
		Healthy:     true,        // 设置实例为健康状态
		Ephemeral:   true,        // 设置实例为临时实例
	})
	if !success {
		return // 注册失败则退出函数
	} else {
		fmt.Println("namingClient.RegisterInstance Success") // 输出注册成功信息
	}
}
