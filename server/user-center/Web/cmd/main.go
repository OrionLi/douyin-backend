package main

import (
	_ "github.com/OrionLi/douyin-backend/pkg/pb"
	"web/common"
	"web/router"
)

func main() {
	//grpc连接
	common.Grpc_conn()
	r := router.NewRouter()

	r.Run(":3301")
}
