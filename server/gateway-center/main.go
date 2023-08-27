package main

import (
	"gateway-center/conf"
	"gateway-center/grpcClient"
)

func main() {
	grpcClient.UserClientInit()
	conf.Init()
	NewRouter().Run(":3010")

}
