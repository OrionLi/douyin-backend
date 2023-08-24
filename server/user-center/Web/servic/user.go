package servic

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc/status"
	"web/grpcClient"
	"web/pkg/util"
	"web/serializer"
)

func UserRegister(ctx context.Context, username, password string) interface{} {
	resp, err := grpcClient.UserRegister(ctx, username, password)
	if err != nil {
		// 该错误如果是status类型就解析错误，不是则另外处理
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息
			code := int32(st.Code())
			msg := st.Message()
			return serializer.DouyinUserRegisterResponse{
				StatusCode: code,
				StatusMsg:  msg,
			}
		} else {
			util.LogrusObj.Error("service GetUserById:", err)
			return serializer.ErrResponse{
				StatusCode: 1,
				StatusMsg:  "fail",
			}
		}
	}
	return serializer.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "用户注册成功",
		UserID:     resp.GetUserId(),
		Token:      resp.GetToken(),
	}
}

func UserLogin(ctx context.Context, username, password string) interface{} {
	resp, err := grpcClient.UserLogin(ctx, username, password)
	if err != nil {
		// 该错误如果是status类型就解析错误，不是则另外处理
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息
			code := int32(st.Code())
			msg := st.Message()
			return serializer.DouyinUserLoginResponse{
				StatusCode: code,
				StatusMsg:  msg,
			}
		} else {
			util.LogrusObj.Error("service GetUserById:", err)
			return serializer.ErrResponse{
				StatusCode: 1,
				StatusMsg:  "fail",
			}
		}
	}
	return serializer.DouyinUserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserID:     resp.GetUserId(),
		Token:      resp.GetToken(),
	}
}

func GetUserById(ctx context.Context, myId, uId uint, token string) interface{} {
	user, err := GetUser(ctx, myId, uId, token)
	if err != nil {
		// 该错误如果是status类型就解析错误，不是则另外处理
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息
			code := int32(st.Code())
			msg := st.Message()
			fmt.Println("c:", code)
			return serializer.ErrResponse{
				StatusCode: code,
				StatusMsg:  msg,
			}
		} else {
			util.LogrusObj.Error("service GetUserById:", err)
			return serializer.ErrResponse{
				StatusCode: 1,
				StatusMsg:  "fail",
			}
		}
	}
	return serializer.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		User:       *user,
	}
}
func GetUser(ctx context.Context, myId, uId uint, token string) (*serializer.User, error) {

	// respUserInfo 有id，name，关注总数和粉丝总数字段
	respUserInfo, err := grpcClient.GetUserById(ctx, uId)
	if err != nil {
		return nil, err
	}
	fmt.Println(`1`)
	respIsFollow, err := grpcClient.IsFollow(ctx, myId, uId)
	if err != nil {
		return nil, err
	}
	fmt.Println(`2`)
	//todo: 剩余信息需从其他模块获取
	respFavCount, err := grpcClient.GetFavCount(ctx, uId)
	if err != nil || respFavCount.GetStatusCode() != 0 {
		fmt.Println("err", err)
		return nil, err
	}

	respVideo, err := grpcClient.GetPublishList(ctx, uId, token)
	if err != nil || respVideo.GetStatusCode() != 0 {
		return nil, err
	}

	user := serializer.User{
		ID:              respUserInfo.User.GetId(),
		Name:            respUserInfo.User.Name,
		FollowCount:     respUserInfo.User.FollowCount,
		FollowerCount:   respUserInfo.User.FollowerCount,
		IsFollow:        respIsFollow.IsFollow,
		Avatar:          "https://th.bing.com/th/id/OIP.7puQ571IXynjU6anJWm_XAHaHa?w=214&h=214&c=7&r=0&o=5&dpr=1.1&pid=1.7",
		BackgroundImage: "https://th.bing.com/th/id/R.476b455c002094fac528b20cf23db88c?rik=iEHmrlVbrFcATw&pid=ImgRaw&r=0",
		Signature:       "test",
		TotalFavorited:  favorSum(respVideo.VideoList),
		WorkCount:       int64(len(respVideo.VideoList)),
		FavoriteCount:   int64(respFavCount.GetFavCount()),
	}
	return &user, nil
}

func favorSum(vs []*pb.Video) int64 {
	var sum int64
	for _, v := range vs {
		sum += v.FavoriteCount
	}
	return sum
}
