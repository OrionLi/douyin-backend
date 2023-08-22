package servic

import (
	"context"
	"google.golang.org/grpc/status"
	"web/grpcClient"
	"web/serializer"
)

func UserRegister(ctx context.Context, username, password string) serializer.DouyinUserRegisterResponse {
	resp, err := grpcClient.UserRegister(ctx, username, password)
	if err != nil {
		// 将错误转换为status.Status
		st, _ := status.FromError(err)
		// 获取错误码和错误信息
		code := int32(st.Code())
		msg := st.Message()
		return serializer.DouyinUserRegisterResponse{
			StatusCode: code,
			StatusMsg:  msg,
			UserID:     0,
			Token:      "",
		}
	}
	return serializer.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "用户注册成功",
		UserID:     resp.GetUserId(),
		Token:      resp.GetToken(),
	}
}

func UserLogin(ctx context.Context, username, password string) serializer.DouyinUserLoginResponse {
	resp, err := grpcClient.UserLogin(ctx, username, password)
	if err != nil {
		// 将错误转换为status.Status
		st, _ := status.FromError(err)
		// 获取错误码和错误信息
		code := int32(st.Code())
		msg := st.Message()
		return serializer.DouyinUserLoginResponse{
			StatusCode: code,
			StatusMsg:  msg,
		}
	}
	return serializer.DouyinUserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserID:     resp.GetUserId(),
		Token:      resp.GetToken(),
	}
}
func GetUserById(ctx context.Context, myId, uId uint) serializer.DouyinUserResponse {
	user, err := GetUser(ctx, myId, uId)
	if err != nil {
		// 将错误转换为status.Status
		st, _ := status.FromError(err)
		// 获取错误码和错误信息
		code := int32(st.Code())
		msg := st.Message()
		return serializer.DouyinUserResponse{
			StatusCode: code,
			StatusMsg:  msg,
		}
	}
	return serializer.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		User:       *user,
	}
}
func GetUser(ctx context.Context, myId, uId uint) (*serializer.User, error) {

	// respUserInfo 有id，name，关注总数和粉丝总数字段
	respUserInfo, err := grpcClient.GetUserById(ctx, uId)
	if err != nil {
		return nil, err
	}
	respIsFollow, err := grpcClient.IsFollow(ctx, myId, uId)
	if err != nil {
		return nil, err
	}

	//todo: 剩余信息需从其他模块获取
	user := serializer.User{
		ID:              respUserInfo.User.GetId(),
		Name:            respUserInfo.User.Name,
		FollowCount:     respUserInfo.User.FollowCount,
		FollowerCount:   respUserInfo.User.FollowerCount,
		IsFollow:        respIsFollow.IsFollow,
		Avatar:          "https://th.bing.com/th/id/OIP.7puQ571IXynjU6anJWm_XAHaHa?w=214&h=214&c=7&r=0&o=5&dpr=1.1&pid=1.7",
		BackgroundImage: "https://th.bing.com/th/id/R.476b455c002094fac528b20cf23db88c?rik=iEHmrlVbrFcATw&pid=ImgRaw&r=0",
		Signature:       "test",
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
	return &user, nil
}
