package oss

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"log"
)

var Viper *viper.Viper

// Init 传入包含AK和SK的文件路径,必须为yml文件 格式为qiniu.accessKey和qiniu.secretKey
func Init() {
	v := viper.New()
	v.AddConfigPath("./oss")
	v.SetConfigType("yaml")
	v.SetConfigName("OssConf.yaml")
	//设置bucket和地区
	v.SetDefault("qiniu.bucket", "bytedance-bucket")
	v.SetDefault("qiniu.zone", storage.ZoneHuanan)
	v.SetDefault("qiniu.url", "http://cnd0.raqtpie.xyz")
	if err := v.ReadInConfig(); err == nil {
		log.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		panic(err)
	}
	Viper = v
	InitFormUploader()
}
