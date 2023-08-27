package grpcClient

import (
	"context"
	"gateway-center/pkg/e"
	"gateway-center/util"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// ResetVideoStreamClient 重置VideoStreamClient
func ResetVideoStreamClient() {
	VideoStreamClient = nil
	client, err := NewVideoStreamClient(VideoConn)
	if err != nil {
		panic(err)
	}
	VideoStreamClient = client
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

// Feed Feed流
func Feed(ctx context.Context, req *pb.DouyinFeedRequest) ([]*pb.Video, int64, error) {
	videos := make([]*pb.Video, 0)
	r, err := VideoClient.Feed(ctx, req)
	if err != nil {
		return videos, 0, err
	}
	if r.StatusCode != 0 {
		return videos, 0, e.NewCustomError(int64(r.StatusCode), *r.StatusMsg)
	}
	videos = r.VideoList
	return videos, *r.NextTime, nil
}

// GetPublishListCount 通过Token获取userId对应的VideoCount
func GetPublishListCount(ctx context.Context, token string) (count int64, err error) {
	userId := ValidateToken(token)
	req := &pb.DouyinPublishListRequest{UserId: userId, Token: token}
	r, err := VideoClient.PublishList(ctx, req)
	if err != nil {
		return 0, err
	}
	if r.StatusCode != 0 {
		return 0, e.NewCustomError(int64(r.StatusCode), *r.StatusMsg)
	}
	return int64(len(r.VideoList)), nil
}

// PublishAction 投稿
func PublishAction(ctx context.Context, req *pb.DouyinPublishActionRequest) error {
	err := VideoStreamClient.Send(req)
	if err != nil {
		return err
	}
	recv, err := VideoStreamClient.CloseAndRecv()
	ResetVideoStreamClient() //重置流式连接
	if err != nil {
		return err
	}
	if recv.StatusCode == 0 {
		return nil
	} else {
		return e.NewCustomError(int64(recv.StatusCode), *recv.StatusMsg)
	}
}

// PublishList 获取发布列表
func PublishList(ctx context.Context, req *pb.DouyinPublishListRequest) ([]*pb.Video, error) {
	videos := make([]*pb.Video, 0)
	r, err := VideoClient.PublishList(ctx, req)
	if err != nil {
		return videos, err
	}
	if r.StatusCode != 0 {
		return videos, e.NewCustomError(int64(r.StatusCode), *r.StatusMsg)
	}
	videos = r.VideoList
	return videos, nil
}

// todo Favorite 迁移

func ActionFavorite(ctx context.Context, selfUserId int64, videoId int64, actionType int32) (*pb.DouyinFavoriteActionResponse, error) {
	return VideoInteractionClient.ActionFavorite(ctx, &pb.DouyinFavoriteActionRequest{
		SelfUserId: selfUserId,
		VideoId:    videoId,
		ActionType: actionType,
	})
}

func GetFavoriteCount(ctx context.Context, userId int64) (*pb.DouyinFavoriteCountResponse, error) {
	return VideoInteractionClient.CountFavorite(ctx, &pb.DouyinFavoriteCountRequest{
		UserId: userId,
	})
}

func GetFavoriteList(ctx context.Context, request *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
	return VideoInteractionClient.ListFavorite(ctx, request)
}

// todo Comment 迁移

func ActionComment(ctx context.Context, req *pb.DouyinCommentActionRequest) (*pb.DouyinCommentActionResponse, error) {
	return VideoInteractionClient.ActionComment(ctx, req)
}

func ListComment(ctx context.Context, req *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	return VideoInteractionClient.ListComment(ctx, req)
}

// ValidateToken 校验Token
func ValidateToken(token string) int64 {
	if len(token) == 0 {
		return -1
	}
	parseToken, err := util.ParseToken(token)
	if err != nil {
		return -1
	}
	// 判断 token 是否过期
	if parseToken.ExpiresAt < time.Now().Unix() {
		return -1
	}
	return int64(parseToken.ID)
}
