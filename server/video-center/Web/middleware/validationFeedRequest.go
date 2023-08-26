package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
	"video-center/Web/pkg/baseResponse"
	"video-center/pkg/errno"
)

// ValidationPublishActionRequest 简单的对请求数据进行检验
func ValidationPublishActionRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params baseResponse.PublishActionParam
		if err := c.ShouldBind(&params); err != nil {
			convertErr := errno.ConvertErr(err)
			c.JSON(http.StatusOK, baseResponse.FeedResponse{
				Response: baseResponse.Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
			})
			return
		}
		fmt.Println("数据绑定成功")
		//校验Token的格式
		parts := strings.Split(params.Token, ".")
		if len(parts) != 3 {
			newErrno := errno.NewErrno(errno.TokenErrCode, "无效的Token")
			c.JSON(http.StatusOK, baseResponse.FeedResponse{
				Response: baseResponse.Response{StatusCode: int32(newErrno.ErrCode), StatusMsg: newErrno.ErrMsg},
			})
			return
		}
		fmt.Println("Token格式正确")
		//校验Title
		if !isValidTitle(params.Title) {
			newErrno := errno.NewErrno(errno.ParamErrCode, errno.ParamErr.ErrMsg)
			c.JSON(http.StatusOK, baseResponse.FeedResponse{
				Response: baseResponse.Response{StatusCode: int32(newErrno.ErrCode), StatusMsg: newErrno.ErrMsg},
			})
			return
		}
		fmt.Println("title格式正确")
		c.Next()
	}
}

func isValidTitle(title string) bool {
	chineseCharPattern := regexp.MustCompile(`[\p{Han}]`)
	chineseCharCount := len(chineseCharPattern.FindAllString(title, -1))
	wordCount := len(strings.Fields(title))
	totalChars := utf8.RuneCountInString(title)
	if (chineseCharCount >= 1 && chineseCharCount <= 20) || (wordCount <= 30 && totalChars <= 30) {
		return true
	}
	return false
}
