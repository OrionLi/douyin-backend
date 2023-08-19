package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var FormUploader *storage.FormUploader
var Key *string
var Policy *storage.PutPolicy
var Mac *qbox.Mac
var ret *storage.PutRet
var upToken *string

// UploadVideo 视频上传，需要传入userID、data字节数组、视频title,返回存放的视频路径和封面路径
func UploadVideo(ctx context.Context, authorId int64, data []byte, title string) (playURL string, CoverUrl string, err error) {
	storageURL := Viper.GetString("qiniu.url")
	digest := md5digest(title)
	//获取Key
	url := fmt.Sprintf("public/%d/%s_%s.png",
		authorId,
		title,
		digest,
	)
	Key = &url
	dataLen := int64(len(data))
	err = FormUploader.Put(ctx, ret, *upToken, *Key, bytes.NewReader(data), dataLen, nil)
	if err != nil {
		return "", "", err
	}
	return "http://" + storageURL + *Key, "http://rzhmys0lm.hn-bkt.clouddn.com/psc.jpg", nil
}
func md5digest(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash)
	return md5str
}

// InitFormUploader 构建formUploader
func InitFormUploader() {
	accessKey := Viper.GetString("qiniu.accessKey")
	secretKey := Viper.GetString("qiniu.secretKey")

	bucket := Viper.GetString("qiniu.bucket")
	zone := Viper.Get("qiniu.zone")
	z, ok := zone.(storage.Region)
	if !ok {
		panic(errors.New("区域转化错误"))
	}
	//PutPolicy
	Policy = &storage.PutPolicy{
		Scope: bucket,
	}
	//获取Mac
	Mac = qbox.NewMac(accessKey, secretKey)
	//获取uploadToken
	token := Policy.UploadToken(Mac)
	upToken = &token
	cfg := storage.Config{}
	cfg.Zone = &z
	cfg.UseCdnDomains = false
	FormUploader = storage.NewFormUploader(&cfg)
	ret = &storage.PutRet{}
}
