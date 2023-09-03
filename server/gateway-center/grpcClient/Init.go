package grpcClient

import (
	"fmt"
	"gateway-center/conf"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

var (
	ChatClient             pb.DouyinMessageServiceClient
	UserClient             pb.UserServiceClient
	RelationClient         pb.RelationServiceClient
	VideoClient            pb.VideoCenterClient
	VideoInteractionClient pb.DouyinVideoInteractionServiceClient
	VideoStreamClient      pb.VideoCenter_PublishActionClient
)
var VideoConn *grpc.ClientConn

func Init() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	// 创建ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: conf.NacosAddress,
			Port:   uint64(conf.NacosPort),
		},
	}
	// 保存公钥
	creds, _ := credentials.NewClientTLSFromFile("../../server/key/test.pem", "*.ygxiaobai111.com")
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
		ServiceName: conf.ChatCenterServiceName,
	})
	if err != nil {
		log.Fatalf("Failed to get chat-service instances: %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instances.Ip, instances.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	ChatClient = pb.NewDouyinMessageServiceClient(conn)

	inst, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: conf.UserCenterServiceName,
	})
	if err != nil {
		log.Fatalf("Failed to get user-service instances: %v", err)
	}

	UserConn, err := grpc.Dial(fmt.Sprintf("%s:%d", inst.Ip, inst.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	UserClient = pb.NewUserServiceClient(UserConn)
	RelationClient = pb.NewRelationServiceClient(UserConn)

	instances, err = nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: conf.VideoCenterServiceName,
	})
	if err != nil {
		log.Fatalf("Failed to get video-service instances: %v", err)
	}
	VideoConn, err = grpc.Dial(fmt.Sprintf("%s:%d", instances.Ip, instances.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	VideoClient = pb.NewVideoCenterClient(VideoConn)
	VideoInteractionClient = pb.NewDouyinVideoInteractionServiceClient(VideoConn)
	//初始化VideoStreamClient
	client, err := NewVideoStreamClient(VideoConn)
	if err != nil {
		log.Fatalf("Failed to get video-stream-service instances: %v", err)
	}
	VideoStreamClient = client
}
