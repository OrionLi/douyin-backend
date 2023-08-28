package controller

import (
	"encoding/json"
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	"gateway-center/response"
	"gateway-center/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

// 用户登录
func UserLogin(ctx *gin.Context) {
	code := e.Success
	var err error
	userName := ctx.Query("username")
	password := ctx.Query("password")
	if len(userName) > 32 || len(password) > 32 { // 最长32位字符

		code = e.InvalidParams
		response.ErrorJSON(ctx, int64(code), e.GetMsg(code))
		return
	}

	resp, err := grpcClient.UserLogin(ctx, userName, password)

	if err != nil {
		util.LogrusObj.Error("<Grpc>  ", err)
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息
			response.ErrorJSON(ctx, int64(st.Code()), st.Message())
			return
		}

		code = e.Error
		response.ErrorJSON(ctx, int64(code), e.GetMsg(code))
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
		response.ErrorJSON(ctx, int64(code), e.GetMsg(code))

		return
	}

	resp, err := grpcClient.UserRegister(ctx, userName, password)
	if err != nil {
		util.LogrusObj.Error("<Grpc>  ", err)
		if st, ok := status.FromError(err); ok {
			// 获取错误码和错误信息
			response.ErrorJSON(ctx, int64(st.Code()), st.Message())
			return
		}

		code = e.Error
		response.ErrorJSON(ctx, int64(code), e.GetMsg(code))
		return
	}

	res := response.DouyinUserRegisterResponse{
		StatusCode: int32(code),
		StatusMsg:  e.GetMsg(code),
		UserID:     resp.UserId,
		Token:      resp.Token,
	}
	response.SuccessJSON(ctx, res)

	return
}

// 获取用户信息
func GetUser(ctx *gin.Context) {
	code := e.Success
	var err error
	userId := ctx.Query("user_id")
	token := ctx.Query("token")
	myId, _ := ctx.Get("UserId")
	mId := myId.(uint)
	uId, err := strconv.Atoi(userId)
	if err != nil {

		code = e.InvalidParams
		response.ErrorJSON(ctx, int64(code), e.GetMsg(code))
		return
	}
	user, err := GetUserInfo(ctx, mId, uint(uId), token)
	if err != nil {

		//日志
		util.LogrusObj.Error("<Grpc>  ", err)
		if st, ok := status.FromError(err); ok {
			// 获取封装的错误码和错误信息
			response.ErrorJSON(ctx, int64(st.Code()), st.Message())
			return
		}

		code = e.Error
		response.ErrorJSON(ctx, int64(code), e.GetMsg(code))
		return
	}

	res := response.DouyinUserResponse{
		StatusCode: int32(code),
		StatusMsg:  e.GetMsg(code),
		User:       user,
	}
	response.SuccessJSON(ctx, res)
	return
}

func GetUserInfo(ctx *gin.Context, myId, uId uint, token string) (user *response.UserInfo, err error) {
	// respUserInfo 有id，name，关注总数和粉丝总数字段
	respUserInfo, err := grpcClient.GetUserById(ctx, myId, uId, token)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(respUserInfo)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &user)
	return
}
