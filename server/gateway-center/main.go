package main

import (
	"gateway-center/conf"
	"gateway-center/grpcClient"
)

func main() {
	conf.Init()
	grpcClient.Init()
	NewRouter().Run(":3010")

}
