package middleware

import (
	"github.com/gin-gonic/gin"
	"time"

	"web/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		code = 0
		token := ctx.Query("token")

		if token == "" {
			code = 404
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {

				code = 1
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = 2
			}
		}
		if code != 0 {
			ctx.JSON(200, gin.H{
				"code": code,
				"msg":  "请登录",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
