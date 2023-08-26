package mq

import (
	"context"
	"strconv"
	"strings"
	"video-center/dao"
	"video-center/pkg/util"
)

func ConsumerFavorite() {
	go ConsumeCreateMessages()
	go ConsumeDeleteMessages()
}

// ConsumeCreateMessages 消费点赞消息
func ConsumeCreateMessages() {
	createMsgList, err := RabbitChannel.Consume(
		"create_fav_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msg := range createMsgList {
		// 处理接收到的消息
		handleCreateMsg(msg.Body)
		// 手动确认消息已处理
		err := msg.Ack(false)
		if err != nil {
			util.LogrusObj.Error("RabbitChannel.QueueDeclare ACK error ", msg, "error msg: ", err)
			return
		}
	}
}

func ConsumeDeleteMessages() {
	deleteMsgList, err := RabbitChannel.Consume(
		"delete_fav_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msg := range deleteMsgList {
		handleDeleteMessage(msg.Body)
		err := msg.Ack(false)
		if err != nil {
			util.LogrusObj.Error("RabbitChannel.QueueDeclare ACK error ", msg, "error msg: ", err)
			return
		}
	}
}

func handleCreateMsg(body []byte) {
	// 将消息内容进行解析，获取 videoId 和 userId
	parts := strings.Split(string(body), ":")
	if len(parts) != 2 {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error")
		return
	}
	videoId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error ", err)
		return
	}
	userId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error ", err)
		return
	}
	err = dao.CreateFav(context.Background(), videoId, userId)
	if err != nil {
		util.LogrusObj.Error("MySQL create Favorite record failed ", err)
		return
	}
}

func handleDeleteMessage(body []byte) {
	parts := strings.Split(string(body), ":")
	if len(parts) != 2 {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error")
		return
	}
	videoId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error ", err)
		return
	}
	userId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error ", err)
		return
	}
	err = dao.DeleteFav(context.Background(), videoId, userId)
	if err != nil {
		util.LogrusObj.Error("MySQL delete Favorite record failed ", err)
		return
	}
}
