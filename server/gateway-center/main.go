package main

import "gateway/conf"

func main() {
	conf.Init()
	NewRouter()
}
