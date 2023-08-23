package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
	"os"
	"os/exec"
	"video-center/pkg/util"
)

var FormUploader *storage.FormUploader
var Policy *storage.PutPolicy
var Mac *qbox.Mac
var ret *storage.PutRet
var upToken *string
var ffmpetPath string

// UploadVideo 视频上传，需要传入userID、data字节数组、视频title,返回存放的视频路径和封面路径
func UploadVideo(ctx context.Context, authorId int64, data []byte, title string) (playURL string, CoverUrl string, err error) {
	storageURL := Viper.GetString("qiniu.url")
	digest := md5digest(title)

	// 生成视频和图片的Key
	videoKey := fmt.Sprintf("public/%d/%s_%s.mp4", authorId, title, digest)
	pictureKey := fmt.Sprintf("public/%d/%s_%s.jpg", authorId, title, digest)

	// 缓存视频数据到本地文件
	videoPath := "temp_video.mp4"
	if err := ioutil.WriteFile(videoPath, data, 0644); err != nil {
		glog.Fatalf("Failed to write video data: %v", err)
		return "", "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(videoPath) // 删除本地缓存的视频文件

	// 使用 FFmpeg 提取第一帧图像到本地缓存的图像文件
	picturePath := "temp_frame.jpg"
	ffmpegCmd := exec.Command(ffmpetPath, "-i", videoPath, "-vframes", "1", picturePath)
	if err := ffmpegCmd.Run(); err != nil {
		util.LogrusObj.Errorf("OSS ERROR Failed to run FFmpeg: %v", err)
		return "", "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(picturePath) // 删除本地缓存的图像文件

	// 读取第一帧图片数据
	frameImageData, err := ioutil.ReadFile(picturePath)
	if err != nil {
		util.LogrusObj.Errorf("OSS ERROR Failed to read frame image data: %v", err)
		return "", "", err
	}

	// 上传视频和第一帧图片到七牛云
	err = FormUploader.Put(ctx, ret, *upToken, videoKey, bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		util.LogrusObj.Errorf("OSS ERROR Failed to upload video: %v", err)
		return "", "", err
	}

	err = FormUploader.Put(ctx, ret, *upToken, pictureKey, bytes.NewReader(frameImageData), int64(len(frameImageData)), nil)
	if err != nil {
		util.LogrusObj.Errorf("OSS ERROR Failed to upload image: %v", err)
		return "", "", err
	}

	return storageURL + videoKey, storageURL + pictureKey, nil
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
	ffmpetPath = Viper.GetString("ffmpet.Path")
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
