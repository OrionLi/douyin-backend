package main

import (
	"user-center/conf"
	"user-center/pkg/util"
)

func main() {

	//初始化配置文件
	err := conf.Init()
	if err != nil {
		util.LogrusObj.Error("<Main> : ", err)
	}

}
