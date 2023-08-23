package rpc

//
//import (
//	"context"
//	"github.com/OrionLi/douyin-backend/pkg/pb"
//	"video-center/dao"
//	"video-center/pkg/errno"
//)
//
//type FavoriteRPCService struct {
//	pb.UnimplementedDouyinMessageServiceServer
//}
//
//func NewFavoriteRPCService() *FavoriteRPCService {
//	return &FavoriteRPCService{}
//}
//
//func (s *FavoriteRPCService) ActionFavorite(ctx context.Context, request *pb.DouyinFavoriteActionRequest) (*pb.DouyinFavoriteActionResponse, error) {
//	selfUserId := request.GetSelfUserId()
//	videoId := request.GetVideoId()
//	actionType := request.GetActionType()
//	switch actionType {
//	case 1:
//		err := dao.CreateFav(ctx, selfUserId, videoId)
//		if err != nil {
//			return &pb.DouyinFavoriteActionResponse{
//				StatusCode: errno.FavActionErrCode,
//				StatusMsg:  errno.FavActionErr.ErrMsg,
//			}, err
//		}
//		return &pb.DouyinFavoriteActionResponse{
//			StatusCode: errno.SuccessCode,
//			StatusMsg:  errno.Success.ErrMsg,
//		}, nil
//	case 2:
//		err := dao.DeleteFav(ctx, selfUserId, videoId)
//		if err != nil {
//			return &pb.DouyinFavoriteActionResponse{
//				StatusCode: errno.FavActionErrCode,
//				StatusMsg:  errno.FavActionErr.ErrMsg,
//			}, err
//		}
//		return &pb.DouyinFavoriteActionResponse{
//			StatusCode: errno.SuccessCode,
//			StatusMsg:  errno.Success.ErrMsg,
//		}, nil
//	default:
//		return &pb.DouyinFavoriteActionResponse{
//			StatusCode: errno.ParamErrCode,
//			StatusMsg:  errno.ParamErr.ErrMsg,
//		}, nil
//	}
//}
//
//func (s *FavoriteRPCService) GetFavoriteCount(ctx context.Context, request *pb.DouyinFavoriteCountRequest) (*pb.DouyinFavoriteCountResponse, error) {
//	// userId: 要查询赞数量和被赞数量的用户ID
//	userId := request.GetUserId()
//	favCount, getFavCount, err := dao.GetFavoriteCount(ctx, userId)
//	if err != nil {
//		return &pb.DouyinFavoriteCountResponse{
//			StatusCode: errno.FavCountErrCode,
//			StatusMsg:  errno.FavCountErr.ErrMsg,
//		}, err
//	}
//	return &pb.DouyinFavoriteCountResponse{
//		StatusCode:   errno.SuccessCode,
//		StatusMsg:    errno.Success.ErrMsg,
//		FavCount:     favCount,
//		GetFavCount_: getFavCount,
//	}, nil
//}
//
//func (s *FavoriteRPCService) GetFavoriteList(ctx context.Context, request *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
//	userId := request.GetUserId()
//	favs := dao.ListFav(ctx, userId)
//	if len(favs) == 0 {
//		return &pb.DouyinFavoriteListResponse{
//			StatusCode: errno.FavListEmptyCode,
//			StatusMsg:  errno.FavListEmptyErr.ErrMsg,
//			VideoList:  nil,
//		}, errno.NewErrno(errno.FavListEmptyCode, errno.FavListEmptyErr.ErrMsg)
//	}
//	//pb
//	var favVideoList []*pb.Video
//	for _, v := range favs {
//		//todo 得到用户ID，然后调用rpc查询用户信息
//		//var user := xxx(v.AuthorID) //然后修改Author
//		video := &pb.Video{
//			Id:            v.Id,
//			Author:        nil,
//			PlayUrl:       v.PlayUrl,
//			CoverUrl:      v.CoverUrl,
//			FavoriteCount: v.FavoriteCount,
//			CommentCount:  v.CommentCount,
//			IsFavorite:    true,
//			Title:         v.Title,
//		}
//		favVideoList = append(favVideoList, video)
//	}
//	return &pb.DouyinFavoriteListResponse{
//		StatusCode: 0,
//		StatusMsg:  errno.Success.ErrMsg,
//		VideoList:  favVideoList,
//	}, nil
//}
