package middleware

import (
	"chat-center/pkg/util"
	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		util.LogrusObj.Infof("URL:%s host:%s method:%s remoteIp:%s", request.URL, request.Host, request.Method, request.RemoteAddr)
		c.Next()
	}
}
