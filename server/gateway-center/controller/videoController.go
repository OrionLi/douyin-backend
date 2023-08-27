package controller

import (
	"context"
	"fmt"
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	baseResponse "gateway-center/response"
	"gateway-center/util"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

// PublishAction 视频投稿
func PublishAction(c *gin.Context) {
	var params baseResponse.PublishActionParam
	//参数绑定
	if err := c.ShouldBind(&params); err != nil {
		convertErr := e.ConvertErr(err)
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.PostForm, convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	fmt.Println(params.Token)
	//绑定二进制数据
	file, err2 := c.FormFile("data")
	if err2 != nil {
		convertErr := e.ConvertErr(err2)
		util.LogrusObj.Errorf("获取二进制流错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.PostForm, convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	open, err2 := file.Open()
	if err2 != nil {
		convertErr := e.ConvertErr(err2)
		util.LogrusObj.Errorf("OpenFile Error 错误原因:%s", convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	fileData := make([]byte, file.Size)
	_, err2 = open.Read(fileData)
	defer func(open multipart.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)
	if err2 != nil {
		convertErr := e.ConvertErr(err2)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	//调用rpcClient
	err := grpcClient.PublishAction(context.Background(), &pb.DouyinPublishActionRequest{
		Token: params.Token,
		Data:  fileData,
		Title: params.Title,
	})
	if err != nil {
		convertErr := e.ConvertErr(err)
		util.LogrusObj.Errorf("Upload Error ErrorMSG:%s", convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	c.JSON(http.StatusOK, baseResponse.PublishListResponse{
		VBResponse: baseResponse.VBResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)},
	})
}

// PublishList 获取用户投稿列
func PublishList(c *gin.Context) {
	var params baseResponse.PublishListParam
	if err := c.ShouldBindQuery(&params); err != nil {
		convertErr := e.ConvertErr(err)
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.PostForm, convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	if params.UserId <= 0 || len(params.Token) == 0 {
		util.LogrusObj.Errorf("Token格式错误 URL:%s Token:%s UserId:%d", c.Request.RequestURI, params.Token, params.UserId)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.ParamErr, StatusMsg: e.GetMsg(e.ParamErr)},
		})
		return
	}
	_, err2 := util.ParseToken(params.Token)
	if err2 != nil {
		util.LogrusObj.Errorf("Token验证失败 URL:%s Token:%s UserId:%d", c.Request.RequestURI, params.Token, params.UserId)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.TokenErr, StatusMsg: e.GetMsg(e.TokenErr)},
		})
		return
	}
	videos, err := grpcClient.PublishList(context.Background(), &pb.DouyinPublishListRequest{
		UserId: params.UserId,
		Token:  params.Token,
	})
	if err != nil {
		convertErr := e.ConvertErr(err)
		util.LogrusObj.Errorf("RPC Error ErrorMSG:%s", convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	videoList := baseResponse.VideoArray{}
	for _, video := range videos {
		info, err2 := grpcClient.GetUserById(context.Background(), uint(video.Author.Id))
		if err2 != nil {
			util.LogrusObj.Errorf("gRPC getUserInfo Error userId:%d", video.Author.Id)
			continue
		}
		user := baseResponse.Vuser{
			Id:            info.User.Id,
			Name:          info.User.Name,
			FollowerCount: info.User.FollowerCount,
			FollowCount:   info.User.FollowCount,
			IsFollow:      false,
		}
		v := baseResponse.Video{
			Id:            video.Id,
			User:          user,
			CoverUrl:      video.CoverUrl,
			PlayUrl:       video.PlayUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		videoList = append(videoList, v)
	}
	c.JSON(http.StatusOK, baseResponse.PublishListResponse{
		VBResponse: baseResponse.VBResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)},
		VideoList:  videoList,
	})
}

// Feed Feed流
func Feed(c *gin.Context) {
	var params baseResponse.FeedParam
	if err := c.ShouldBindQuery(&params); err != nil {
		convertErr := e.ConvertErr(err)
		//记录日志
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.RequestURI, convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.FeedResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	isLogin := false
	isFollow := false
	userId := grpcClient.ValidateToken(params.Token)
	if userId != -1 {
		isLogin = true
	}
	fmt.Printf("LatestTime: %d, Token: %s\n", params.LatestTime, params.Token)
	videos, nextTime, err := grpcClient.Feed(context.Background(), &pb.DouyinFeedRequest{
		Token:      &params.Token,
		LatestTime: &params.LatestTime,
	})
	if err != nil {
		convertErr := e.ConvertErr(err)
		util.LogrusObj.Errorf("rpc调用错误 URL:%s 错误原因:%s", c.Request.URL, convertErr.Msg)
		c.JSON(http.StatusOK, baseResponse.FeedResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.Code), StatusMsg: convertErr.Msg},
		})
		return
	}
	VideoList := baseResponse.VideoArray{}
	for _, video := range videos {
		info, err := grpcClient.GetUserById(context.Background(), uint(video.Author.Id))
		if err != nil {
			util.LogrusObj.Errorf("获取User失败 UserId:%d UserToken:%d", video.Author.Id, &params.Token)
			continue
		}
		if isLogin {
			follow, err := grpcClient.IsFollow(context.Background(), uint(video.Author.Id), uint(userId))
			if err != nil {
				isFollow = false
			}
			isFollow = follow.IsFollow
		}
		user := baseResponse.Vuser{
			Id:            info.User.Id,
			Name:          info.User.Name,
			FollowerCount: info.User.FollowerCount,
			FollowCount:   info.User.FollowCount,
			IsFollow:      isFollow,
		}
		fmt.Println(isLogin)
		fmt.Println(isFollow)
		v := &baseResponse.Video{
			Id:            video.Id,
			User:          user,
			CoverUrl:      video.CoverUrl,
			PlayUrl:       video.PlayUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		VideoList = append(VideoList, *v)
	}
	c.JSON(http.StatusOK, baseResponse.FeedResponse{
		VBResponse: baseResponse.VBResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)},
		VideoList:  VideoList,
		NextTime:   nextTime,
	})
}

// FavoriteParam 点赞请求参数
type FavoriteParam struct {
	Token      string `json:"token"`
	VideoId    string `json:"video_id"`
	ActionType string `json:"action_type"`
}

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	userIdAny, _ := c.Get("UserId")
	userId := userIdAny.(int64)
	var param baseResponse.CommentActionParam
	if err := c.ShouldBind(&param); err != err {
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			Comment:    &pb.Comment{},
		})
		return
	}
	if param.ActionType == "" || param.VideoID == "" {
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			Comment:    &pb.Comment{},
		})
		return
	}
	//判断是哪种操作？
	if param.ActionType == "1" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
				Comment:    &pb.Comment{},
			})
			return
		}
		request := pb.DouyinCommentActionRequest{
			SelfUserId:  userId,
			VideoId:     videoId,
			ActionType:  0, //保存
			CommentText: param.CommentText,
		}
		response, err := grpcClient.ActionComment(c, &request)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: e.FailedToCallRpc, StatusMsg: e.GetMsg(e.FailedToCallRpc)},
				Comment:    &pb.Comment{},
			})
		}
		c.JSON(http.StatusOK, response)
		return
	}
	if param.ActionType == "2" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		commentID, err := strconv.ParseInt(param.CommentID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
				Comment:    &pb.Comment{},
			})
			return
		}
		request := pb.DouyinCommentActionRequest{
			SelfUserId: userId,
			VideoId:    videoId,
			ActionType: 1, //删除
			CommentId:  commentID,
		}
		response, err2 := grpcClient.ActionComment(c, &request)
		if err2 != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: e.FailedToCallRpc, StatusMsg: e.GetMsg(e.FailedToCallRpc)},
				Comment:    &pb.Comment{},
			})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
		VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
		Comment:    &pb.Comment{},
	})
}

// CommentList 评论列表
func CommentList(c *gin.Context) {
	userIdAny, _ := c.Get("UserId")
	userId := userIdAny.(int64)
	videoId := c.Query("video_id")
	if videoId == "" {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			Comment:    []*pb.Comment{},
		})
		return
	}
	videoID, err1 := strconv.ParseInt(videoId, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			Comment:    []*pb.Comment{},
		})
		return
	}

	request := pb.DouyinCommentListRequest{
		SelfUserId: userId,
		VideoId:    videoID,
	}
	response, _ := grpcClient.ListComment(c, &request)
	c.JSON(http.StatusOK, response)
}

// ActionFav 点赞操作
func ActionFav(c *gin.Context) {
	userIdAny, _ := c.Get("UserId")
	userId := userIdAny.(int64)
	videoId := util.StringToInt64(c.Query("video_id"))
	if videoId == -1 {
		c.JSON(http.StatusOK, baseResponse.DouyinFavoriteActionResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}
	actionType := util.StringToInt64(c.Query("action_type"))
	resp, err := grpcClient.ActionFavorite(context.Background(), userId, videoId, int32(actionType))
	if err != nil || resp.StatusCode != e.Success {
		c.JSON(http.StatusOK, baseResponse.DouyinFavoriteActionResponse{StatusCode: e.Error, StatusMsg: e.GetMsg(e.Error)})
		return
	}
	c.JSON(http.StatusOK, baseResponse.DouyinFavoriteActionResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)})
}

// ListFav 获取喜欢列表
func ListFav(c *gin.Context) {
	userId := c.Query("user_id")
	userIdAny, _ := c.Get("UserId")
	userIdToken := userIdAny.(int64)
	if userId == "" {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			FavList:    []*pb.Video{},
		})
		return
	}

	UserIdParseInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			FavList:    []*pb.Video{},
		})
		return
	}
	if userIdToken != UserIdParseInt {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)},
			FavList:    []*pb.Video{},
		})
		return
	}
	request := pb.DouyinFavoriteListRequest{UserId: userIdToken}
	response, _ := grpcClient.GetFavoriteList(c, &request)
	if response == nil {
		c.JSON(http.StatusOK, &pb.DouyinFavoriteListResponse{
			StatusCode: e.FavListEmpty,
			StatusMsg:  e.GetMsg(e.FavListEmpty),
			VideoList:  []*pb.Video{},
		})
		return
	}
	println(response)
	c.JSON(http.StatusOK, response)
}
