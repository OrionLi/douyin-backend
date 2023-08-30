package grpc

import (
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"user-center/conf"
)

var (
	VideoClient            pb.VideoCenterClient
	VideoInteractionClient pb.DouyinVideoInteractionServiceClient
	VideoStreamClient      pb.VideoCenter_PublishActionClient
)
var VideoConn *grpc.ClientConn

func Init(nacosAddress string,
	nacosPort uint64) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	// 创建ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosAddress,
			Port:   nacosPort,
		},
	}

	// 创建服务发现客户端
	nacosClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		})

	if err != nil {
		log.Fatalf("Failed to create Nacos client: %v", err)
	}

	// 获取 gRPC 服务实例信息

	instances, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: conf.VideoCenterServiceName,
	})
	if err != nil {
		log.Fatalf("Failed to get video-service instances: %v", err)
	}
	VideoConn, err = grpc.Dial(fmt.Sprintf("%s:%d", instances.Ip, instances.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	VideoClient = pb.NewVideoCenterClient(VideoConn)
	VideoInteractionClient = pb.NewDouyinVideoInteractionServiceClient(VideoConn)

}
