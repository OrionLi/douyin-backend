package main

import (
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"os/signal"
	"syscall"
	"video-center/cache"
	"video-center/conf"
	"video-center/dao"
	"video-center/handler"
	"video-center/mq"
	"video-center/oss"
	"video-center/service"
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
	mq.Init()
	mq.ConsumerFavorite()
	RegisterNacos()
	BeforeExit()
	service.UpdateFavoriteCacheToMySQL()
	go service.UpdateFavoriteCacheToMySQLAtRegularTime()
	//引入证书
	cert, err := credentials.NewServerTLSFromFile("../../server/key/test.pem", "../../server/key/test.key")
	if err != nil {
		fmt.Printf("credentials.NewServerTLSFromFile Err :%s\n", err.Error())
	}
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(52428800), //50Mb
		grpc.MaxSendMsgSize(52428800),
		grpc.Creds(cert)) //加密
	pb.RegisterVideoCenterServer(server, &handler.VideoServer{})
	pb.RegisterDouyinVideoInteractionServiceServer(server, &handler.VideoInteractionServer{})
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
	// 创建clientConfig，用于配置客户端行为
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,    // 设置超时时间为5秒
		NotLoadCacheAtStart: true,    // 在启动时不加载缓存
		LogLevel:            "debug", // 设置日志级别为调试模式
	}

	// 至少一个ServerConfig，指定Nacos服务器的地址和端口
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: NacosIp,   // Nacos服务器的IP地址
			Port:   NacosPort, // Nacos服务器的端口
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
		Ip:          ServerIp,           // 注册实例的IP地址
		Port:        uint64(ServerPort), // 注册实例的端口
		ServiceName: ServiceName,        // 服务的名称
		Weight:      10,                 // 权重为10
		Enable:      true,               // 设置实例为可用状态
		Healthy:     true,               // 设置实例为健康状态
		Ephemeral:   true,               // 设置实例为临时实例
	})
	if !success {
		return // 注册失败则退出函数
	} else {
		fmt.Println("namingClient.RegisterInstance Success") // 输出注册成功信息
	}
}

// BeforeExit 服务器关闭前的最后操作
func BeforeExit() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		fmt.Printf("Received signal %s, performing final tasks...\n", sig)
		service.UpdateFavoriteCacheToMySQL()
		os.Exit(0)
	}()
}
