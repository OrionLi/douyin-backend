package rpc

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"video-center/conf"
)

var ServiceName string
var NacosIp string
var NacosPort uint64

func Init() {
	ServiceName = conf.Viper.GetString("application.ServiceName")
	NacosIp = conf.Viper.GetString("nacos.Ip")
	NacosPort = conf.Viper.GetUint64("nacos.Port")
	initVideoRpc()
	InitUserClient()
}

// VideoClient 非流式
var VideoClient pb.VideoCenterClient

// VideoStreamClient 流式
var VideoStreamClient pb.VideoCenter_PublishActionClient

// VideoInteractionClient 视频互动rpc端口
var VideoInteractionClient pb.DouyinVideoInteractionServiceClient

var UserClient pb.UserServiceClient

// Conn 共有连接
var Conn *grpc.ClientConn

func initVideoRpc() {
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: NacosIp,
			Port:   NacosPort,
		},
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		panic(err)
	}
	instances, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: ServiceName,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instances.Ip, instances.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	Conn = conn
	//创建非流式client
	client := pb.NewVideoCenterClient(conn)
	//创建流式client
	streamClient, err := NewVideoStreamClient(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err != nil {
		panic(err)
	}
	VideoClient = client
	VideoStreamClient = streamClient

	// VideoInteractionClient
	VideoInteractionClient = pb.NewDouyinVideoInteractionServiceClient(Conn)
}
func InitUserClient() {
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: NacosIp,
			Port:   NacosPort,
		},
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		panic(err)
	}
	instances, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "UserServer",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instances.Ip, instances.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	Conn = conn
	//创建非流式client
	client := pb.NewUserServiceClient(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err != nil {
		panic(err)
	}
	UserClient = client
}

// StreamClient 流式client
type StreamClient struct {
	client pb.VideoCenterClient
	stream grpc.ClientStream
}

func (c *StreamClient) Send(request *pb.DouyinPublishActionRequest) error {
	if err := c.stream.SendMsg(request); err != nil {
		return err
	}
	return nil
}

func (c *StreamClient) CloseAndRecv() (*pb.DouyinPublishActionResponse, error) {
	if err := c.stream.CloseSend(); err != nil {
		return nil, err
	}

	response := new(pb.DouyinPublishActionResponse)
	if err := c.stream.RecvMsg(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *StreamClient) Header() (metadata.MD, error) {
	return nil, nil
}

func (c *StreamClient) Trailer() metadata.MD {
	return nil
}

func (c *StreamClient) CloseSend() error {
	return nil
}

func (c *StreamClient) Context() context.Context {
	return context.Background()
}

func (c *StreamClient) SendMsg(m interface{}) error {
	return nil
}

func (c *StreamClient) RecvMsg(m interface{}) error {
	return nil
}

// NewVideoStreamClient 新建video流式client
func NewVideoStreamClient(conn *grpc.ClientConn) (*StreamClient, error) {
	client := pb.NewVideoCenterClient(conn)
	stream, err := client.PublishAction(context.Background())
	if err != nil {
		return nil, err
	}
	return &StreamClient{
		client: client,
		stream: stream,
	}, nil
}
func (c *StreamClient) SendData(data *pb.DouyinPublishActionRequest) error {
	return c.stream.SendMsg(data)
}
func (c *StreamClient) CloseAndReceive() (*pb.DouyinPublishActionResponse, error) {
	if err := c.stream.CloseSend(); err != nil {
		return nil, err
	}

	response := new(pb.DouyinPublishActionResponse)
	if err := c.stream.RecvMsg(response); err != nil {
		return nil, err
	}

	return response, nil
}

// ResetVideoStreamClient 重置VideoStreamClient
func ResetVideoStreamClient() {
	VideoStreamClient = nil
	client, err := NewVideoStreamClient(Conn)
	if err != nil {
		panic(err)
	}
	VideoStreamClient = client
}
