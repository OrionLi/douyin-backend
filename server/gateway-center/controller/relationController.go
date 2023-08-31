package controller

import (
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	"gateway-center/response"
	"gateway-center/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RelationAction(ctx *gin.Context) {
	// 获取请求参数
	token := ctx.Query("token")
	toUserId := util.StringToInt64(ctx.Query("to_user_id"))
	actionType := util.StringToInt64(ctx.Query("action_type"))

	// 调用rpc
	resp, err := grpcClient.RelationAction(ctx, token, toUserId, actionType)
	if err != nil || resp.StatusCode != e.Success {
		ctx.JSON(http.StatusInternalServerError, response.RelationActionResponse{
			StatusCode: e.Error,
			StatusMsg:  e.GetMsg(e.Error),
		})
		return
	}

	// 返回响应
	ctx.JSON(http.StatusOK, response.RelationActionResponse{
		StatusCode: e.Success,
		StatusMsg:  e.GetMsg(e.Success),
	})
}

func GetFollowList(ctx *gin.Context) {
	userId := util.StringToInt64(ctx.Query("user_id"))
	token := ctx.Query("token")

	println(userId)
	resp, err := grpcClient.GetFollowList(ctx, userId, token)
	if err != nil || resp.StatusCode != e.Success {
		ctx.JSON(http.StatusInternalServerError, response.GetFollowListResponse{
			StatusCode: e.Error,
			StatusMsg:  e.GetMsg(e.Error),
			UserList:   nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.GetFollowListResponse{
		StatusCode: e.Success,
		StatusMsg:  e.GetMsg(e.Success),
		UserList:   resp.GetUserList(),
	})
}

func GetFollowerList(ctx *gin.Context) {
	userId := util.StringToInt64(ctx.Query("user_id"))
	token := ctx.Query("token")

	resp, err := grpcClient.GetFollowerList(ctx, userId, token)
	if err != nil || resp.StatusCode != e.Success {
		ctx.JSON(http.StatusInternalServerError, response.GetFollowerListResponse{
			StatusCode: e.Error,
			StatusMsg:  e.GetMsg(e.Error),
			UserList:   nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.GetFollowerListResponse{
		StatusCode: e.Success,
		StatusMsg:  e.GetMsg(e.Success),
		UserList:   resp.GetUserList(),
	})
}

func GetFriendList(ctx *gin.Context) {
	userId := util.StringToInt64(ctx.Query("user_id"))
	token := ctx.Query("token")

	resp, err := grpcClient.GetFriendList(ctx, userId, token)
	if err != nil || resp.StatusCode != e.Success {
		ctx.JSON(http.StatusInternalServerError, response.GetFriendListResponse{
			StatusCode: e.Error,
			StatusMsg:  e.GetMsg(e.Error),
			UserList:   nil,
		})
		return
	}
	// 更新分割线-----------------------
	var friendUserList []response.FriendUser
	for _, friendUserTemp := range resp.GetUserList() {
		interActiveId := friendUserTemp.GetUser().GetId()
		interActiveName := friendUserTemp.GetUser().GetName()
		messageResponse, _ := grpcClient.GetMessage(userId, interActiveId, 0)
		list := messageResponse.GetMessageList()
		var content string
		if len(list) == 0 {
			content = ""
		} else {
			content = list[len(list)-1].GetContent()
		}
		//message := list[len(list)-1]
		friendUser := response.FriendUser{
			ID:              interActiveId,
			Name:            interActiveName,
			FollowCount:     0,
			FollowerCount:   0,
			IsFollow:        true,
			FavoriteCount:   0,
			WorkCount:       0,
			TotalFavorited:  0,
			BackgroundImage: "https://ts2.cn.mm.bing.net/th?id=OIP-C.HfZqICAPqMQslH0cMrIDFQHaKe&w=210&h=297&c=8&rs=1&qlt=90&o=6&dpr=1.1&pid=3.1&rm=2",
			Signature:       "测试用户",
			Avatar:          "https://ts2.cn.mm.bing.net/th?id=OIP-C.druUEHdZrBEuZPn2w80Y1QHaNK&w=187&h=333&c=8&rs=1&qlt=90&o=6&dpr=1.1&pid=3.1&rm=2",
			//Message: message.GetContent(),
			Message: content,
			MsgType: 0,
		}
		friendUserList = append(friendUserList, friendUser)
	}

	// 更新分割线-----------------------

	ctx.JSON(http.StatusOK, response.GetFriendListResponse{
		StatusCode: e.Success,
		StatusMsg:  e.GetMsg(e.Success),
		UserList:   friendUserList,
	})
}
