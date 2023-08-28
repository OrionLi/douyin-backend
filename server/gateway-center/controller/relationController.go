package controller

import (
	"gateway-center/grpcClient"
	"gateway-center/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

//type RelationController struct {
//	relationServer *grpcClient.RelationServer
//}
//
//func NewRelationController(server *server.RelationServer) *RelationController {
//	return &RelationController{relationServer: server}
//}

func RelationAction(ctx *gin.Context) {
	// 获取请求参数
	token := ctx.Query("token")
	toUserId := util.StringToInt64(ctx.Query("to_user_id"))
	actionType := util.StringToInt64(ctx.Query("action_type"))

	// 调用rpc
	resp, err := grpcClient.RelationAction(ctx, token, toUserId, actionType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 返回响应
	ctx.JSON(http.StatusOK, gin.H{"message": resp.StatusMsg})
}

func GetFollowList(ctx *gin.Context) {
	userId := util.StringToInt64(ctx.Query("user_id"))
	token := ctx.Query("token")

	resp, err := grpcClient.GetFollowList(ctx, userId, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_list": resp.UserList})
}

func GetFollowerList(ctx *gin.Context) {
	userId := util.StringToInt64(ctx.Query("user_id"))
	token := ctx.Query("token")

	resp, err := grpcClient.GetFollowerList(ctx, userId, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_list": resp.UserList})
}

func GetFriendList(ctx *gin.Context) {
	userId := util.StringToInt64(ctx.Query("user_id"))
	token := ctx.Query("token")

	resp, err := grpcClient.GetFriendList(ctx, userId, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_list": resp.UserList})
}
