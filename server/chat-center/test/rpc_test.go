package test

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

// TestGRPC 测试 gRPC 服务
func TestGRPC(t *testing.T) {
	t.Parallel() // 可选：允许并行运行测试

	// 连接 gRPC 服务
	conn, err := grpc.Dial("localhost:9422", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Fatalf("Failed to close gRPC connection: %v", err)
		}
	}(conn)

	// 创建 gRPC 客户端
	client := pb.NewDouyinMessageServiceClient(conn)

	t.Run("GetMessage", func(t *testing.T) {
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
			t.Fatalf("Failed to call GetMessage: %v", err)
		}

		// 打印响应
		t.Logf("GetMessage Response: %+v", getMessageResponse)
	})

	t.Run("SendMessage", func(t *testing.T) {
		// 调用 gRPC 方法
		sendMessageRequest := &pb.DouyinMessageActionRequest{
			SelfUserId: 1,
			ToUserId:   123,
			Content:    "Hello World!",
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		sendMessageResponse, err := client.SendMessage(ctx, sendMessageRequest)
		if err != nil {
			t.Fatalf("Failed to call SendMessage: %v", err)
		}

		// 打印响应
		t.Logf("SendMessage Response: %+v", sendMessageResponse)
	})
}
