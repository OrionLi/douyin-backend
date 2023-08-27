package middleware

import (
	"gateway-center/pkg/e"
	"github.com/gin-gonic/gin"
	"time"

	"gateway-center/util"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := e.Success

		token := ctx.Query("token")

		if token == "" {
			code = e.ErrorAuthToken
		} else {
			claims, err := util.ParseToken(token)
			// 将id存入上下文
			ctx.Set("UserId", claims.ID)
			if err != nil {

				code = e.Error
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout
			}
		}
		if code != e.Success {
			ctx.JSON(200, gin.H{
				"status_code": code,
				"status_msg":  e.GetMsg(code),
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}