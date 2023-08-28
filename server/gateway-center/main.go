package main

import (
	"gateway-center/cache"
	"gateway-center/conf"
	"gateway-center/grpcClient"
)

func main() {
	conf.Init()
	grpcClient.Init()
	cache.Init()
	NewRouter().Run(":3010")
}
