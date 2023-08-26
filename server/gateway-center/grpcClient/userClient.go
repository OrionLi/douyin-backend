package grpcClient

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	UserConn *grpc.ClientConn
)

// UserClientInit grpc初始化
func UserClientInit() {
	addr := "127.0.0.1"
	port := "3001"
	addr = addr + ":" + port
	fmt.Println(addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: #{err}")
	}
	UserConn = conn

}

func GetUserConn() *grpc.ClientConn {
	return UserConn
}

func UserRegister(ctx context.Context, username, password string) (*pb.DouyinUserRegisterResponse, error) {

	conn := GetUserConn()

	clientUser := pb.NewUserServiceClient(conn)
	pb.NewRelationServiceClient(conn)

	resp, err := clientUser.Register(ctx, &pb.DouyinUserRegisterRequest{
		Username: username,
		Password: password,
	})
	return resp, err
}

func UserLogin(ctx context.Context, username, password string) (*pb.DouyinUserLoginResponse, error) {

	conn := GetUserConn()

	clientUser := pb.NewUserServiceClient(conn)
	resp, err := clientUser.Login(ctx, &pb.DouyinUserLoginRequest{
		Username: username,
		Password: password,
	})
	return resp, err
}

func GetUserById(ctx context.Context, uId uint) (*pb.DouyinUserResponse, error) {

	conn := GetUserConn()
	clientUser := pb.NewUserServiceClient(conn)
	resp, err := clientUser.GetUserById(ctx, &pb.DouyinUserRequest{UserId: int64(uId)})
	return resp, err
}

func IsFollow(ctx context.Context, uId, followId uint) (*pb.IsFollowResponse, error) {

	conn := GetUserConn()

	clientUser := pb.NewUserServiceClient(conn)
	resp, err := clientUser.IsFollow(ctx, &pb.IsFollowRequest{
		UserId:       int64(uId),
		FollowUserId: int64(followId),
	})
	return resp, err
}
