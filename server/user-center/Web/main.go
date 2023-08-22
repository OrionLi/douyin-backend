package main

import (
	_ "github.com/OrionLi/douyin-backend/pkg/pb"
	"web/pkg/util"
)

func main() {

	r := NewRouter()

	err := r.Run(":3301")
	if err != nil {
		util.LogrusObj.Error("<Main> : ", err)
	}
}
