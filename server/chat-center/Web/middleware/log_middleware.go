package middleware

import (
	"chat-center/pkg/utils"
	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		utils.LogrusObj.Infof("URL:%s host:%s method:%s remoteIp:%s", request.URL, request.Host, request.Method, request.RemoteAddr)
		c.Next()
	}
}
