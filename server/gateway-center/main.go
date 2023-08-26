package main

import (
	"gateway/conf"
	"gateway/grpcClient"
)

func main() {
	grpcClient.UserClientInit()
	conf.Init()
	NewRouter().Run(":3010")

}
