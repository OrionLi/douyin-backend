package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

var LogrusObj *logrus.Logger

type MyFormatter struct {
}

// Format 自定义日志模式
func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	level := entry.Level.String()
	message := entry.Message
	return []byte(fmt.Sprintf("[%s]: %s [%s]\n", level, timestamp, message)), nil
}
func init() {
	src, err := setOutPutFile()
	if err != nil {
		log.Println(err)
	}
	if LogrusObj != nil {
		LogrusObj.Out = src
		return
	}

	// 实例化
	logger := logrus.New()
	logger.Out = src                   // 设置输出
	logger.SetLevel(logrus.DebugLevel) // 设置日志级别
	logger.SetFormatter(&MyFormatter{})

	LogrusObj = logger
}

func setOutPutFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	// dir 当前工作目录
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return src, nil
}
