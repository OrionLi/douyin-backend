package main

import (
	_ "github.com/OrionLi/douyin-backend/pkg/pb"
)

func main() {

	r := NewRouter()

	r.Run(":3301")
}
