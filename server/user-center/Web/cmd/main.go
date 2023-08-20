package main

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	_ "github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc/status"
	"web/common"
)

func main() {
	//grpc连接
	ctx := context.Background()
	common.Grpc_conn()
	conn := common.GetConn()
	defer conn.Close()
	clientUser := pb.NewUserServiceClient(conn)
	pb.NewRelationServiceClient(conn)

	resp, err := clientUser.Register(ctx, &pb.DouyinUserRegisterRequest{
		Username: "唐家三少",
		Password: "114514",
	})
	if err != nil {
		// 将错误转换为status.Status
		st, _ := status.FromError(err)
		// 获取错误码和错误信息
		code := st.Code()
		msg := st.Message()
		fmt.Println("code:", code)
		fmt.Println("msg:", msg)
	}
	fmt.Println("token", resp.GetToken(), " userId: ", resp.GetUserId())
}
