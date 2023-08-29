package middleware

import (
	"gateway-center/pkg/e"
	"gateway-center/util"
	"github.com/gin-gonic/gin"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := e.Success

		token := ctx.Query("token")

		if token == "" {
			code = e.ErrorAuthToken
		} else {
			claims, err := util.ParseToken(token)

			if err != nil {
				code = e.Error
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout
			}
			if claims != nil {
				// 将id存入上下文
				ctx.Set("UserId", claims.ID)
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
