package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func ActionComment(ctx context.Context, req *pb.DouyinCommentActionRequest) (*pb.DouyinCommentActionResponse, error) {
	return VideoInteractionClient.ActionComment(ctx, req)
}

func ListComment(ctx context.Context, req *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	return VideoInteractionClient.ListComment(ctx, req)
}
