package test

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

// HACK 临时测试函数
func TestGRPC(t *testing.T) {
	// 连接 gRPC 服务
	conn, err := grpc.Dial("localhost:9422", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewDouyinMessageServiceClient(conn)

	// 调用 gRPC 方法
	getMessageRequest := &pb.DouyinMessageChatRequest{
		SelfUserId: 1,
		ToUserId:   123,
		PreMsgTime: 0, // 设置合适的时间戳
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	getMessageResponse, err := client.GetMessage(ctx, getMessageRequest)
	if err != nil {
		log.Fatalf("Failed to call GetMessage: %v", err)
	}

	// 打印响应
	fmt.Printf("GetMessage Response: %+v\n", getMessageResponse)

	fmt.Println("=====================================")

	// 调用 gRPC 方法
	sendMessageRequest := &pb.DouyinMessageActionRequest{
		SelfUserId: 1,
		ToUserId:   123,
		Content:    "Hello World!",
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sendMessageResponse, err := client.SendMessage(ctx, sendMessageRequest)
	if err != nil {
		log.Fatalf("Failed to call SendMessage: %v", err)
	}

	// 打印响应
	fmt.Printf("SendMessage Response: %+v\n", sendMessageResponse)
}
