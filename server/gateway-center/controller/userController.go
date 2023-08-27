package controller

import (
	"fmt"
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	"gateway-center/response"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

// todo: 封装json接口
// 用户登录
func UserLogin(ctx *gin.Context) {
	code := e.Success
	var err error
	userName := ctx.Query("username")
	password := ctx.Query("password")
	if len(userName) > 32 || len(password) > 32 { // 最长32位字符
		// todo: 添加返回
		code = e.InvalidParams
		ctx.JSON(http.StatusOK, response.DouyinUserLoginResponse{
			StatusCode: int32(code),
			StatusMsg:  e.GetMsg(code),
		})
		return
	}
	// todo: 添加grpc调用
	resp, err := grpcClient.UserLogin(ctx, userName, password)

	if err != nil {
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息

			ctx.JSON(http.StatusOK, response.DouyinUserLoginResponse{
				StatusCode: int32(st.Code()),
				StatusMsg:  st.Message(),
			})
			return
		}
		// todo: 添加返回
		code = e.Error
		ctx.JSON(http.StatusOK, response.DouyinUserLoginResponse{
			StatusCode: int32(code),
			StatusMsg:  e.GetMsg(code),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.DouyinUserLoginResponse{
		StatusCode: int32(code),
		StatusMsg:  e.GetMsg(code),
		UserID:     resp.UserId,
		Token:      resp.Token,
	})

}

func UserRegister(ctx *gin.Context) {
	code := e.Success
	var err error
	userName := ctx.Query("username")
	password := ctx.Query("password")
	if len(userName) > 32 || len(password) > 32 { // 最长32位字符
		// todo: 添加返回
		code = e.InvalidParams
		ctx.JSON(http.StatusOK, response.DouyinUserRegisterResponse{
			StatusCode: int32(code),
			StatusMsg:  e.GetMsg(code),
		})
		return
	}

	// todo: 添加grpc调用
	resp, err := grpcClient.UserRegister(ctx, userName, password)

	if err != nil {
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息
			ctx.JSON(http.StatusOK, response.DouyinUserRegisterResponse{
				StatusCode: int32(st.Code()),
				StatusMsg:  st.Message(),
			})
			return
		}
		// todo: 添加返回
		code = e.Error
		ctx.JSON(http.StatusOK, response.DouyinUserRegisterResponse{
			StatusCode: int32(code),
			StatusMsg:  e.GetMsg(code),
		})
		return
	}
	// todo: 添加返回
	ctx.JSON(http.StatusOK, response.DouyinUserRegisterResponse{
		StatusCode: int32(code),
		StatusMsg:  e.GetMsg(code),
		UserID:     resp.UserId,
		Token:      resp.Token,
	})
	return
}

// 获取用户信息
func GetUser(ctx *gin.Context) {
	code := e.Success
	var err error
	userId := ctx.Query("user_id")
	myId, _ := ctx.Get("UserId")
	mId := myId.(uint)
	uId, err := strconv.Atoi(userId)
	if err != nil {
		// todo: 添加返回
		fmt.Println(err)
		code = e.InvalidParams
		ctx.JSON(http.StatusOK, response.DouyinUserResponse{
			StatusCode: int32(code),
			StatusMsg:  e.GetMsg(code),
		})
		return
	}
	user, err := GetUserInfo(ctx, mId, uint(uId), " ")
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息

			ctx.JSON(http.StatusOK, response.DouyinUserResponse{
				StatusCode: int32(st.Code()),
				StatusMsg:  st.Message(),
			})
			return
		}
		// todo: 添加返回
		code = e.Error
		ctx.JSON(http.StatusOK, response.DouyinUserResponse{
			StatusCode: int32(code),
			StatusMsg:  e.GetMsg(code),
		})
		return
	}
	// todo: 添加返回
	ctx.JSON(http.StatusOK, response.DouyinUserResponse{
		StatusCode: int32(code),
		StatusMsg:  e.GetMsg(code),
		User:       user,
	})
	return
}

func GetUserInfo(ctx *gin.Context, myId, uId uint, token string) (*response.UserInfo, error) {
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

	user := response.UserInfo{
		ID:              respUserInfo.User.GetId(),
		Name:            respUserInfo.User.Name,
		FollowCount:     respUserInfo.User.FollowCount,
		FollowerCount:   respUserInfo.User.FollowerCount,
		IsFollow:        respIsFollow.IsFollow,
		Avatar:          "https://th.bing.com/th/id/OIP.7puQ571IXynjU6anJWm_XAHaHa?w=214&h=214&c=7&r=0&o=5&dpr=1.1&pid=1.7",
		BackgroundImage: "https://th.bing.com/th/id/R.476b455c002094fac528b20cf23db88c?rik=iEHmrlVbrFcATw&pid=ImgRaw&r=0",
		Signature:       "test",
		//需从video模块获取
		TotalFavorited: "",
		WorkCount:      0,
		FavoriteCount:  0,
	}
	return &user, nil
}
