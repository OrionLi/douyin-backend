package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web/pkg/util"
	"web/serializer"
	"web/servic"
)

func ErrorResponse(err error) serializer.Response {
	return serializer.Response{
		Code: 4,
		Msg:  "错误请求",
	}
}
func UserRegister(ctx *gin.Context) {

	username := ctx.Query("username")
	password := ctx.Query("password")
	res := servic.UserRegister(ctx, username, password)
	ctx.JSON(http.StatusOK, res)

}

func UserLogin(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	fmt.Println(username)
	res := servic.UserLogin(ctx, username, password)
	ctx.JSON(http.StatusOK, res)

}
func GetUserById(ctx *gin.Context) {
	token := ctx.Query("token")
	userId := ctx.Query("user_id")
	uId, _ := strconv.Atoi(userId)
	//获取token中的id
	claims, _ := util.ParseToken(token)

	res := servic.GetUserById(ctx, claims.ID, uint(uId))
	ctx.JSON(http.StatusOK, res)

}
