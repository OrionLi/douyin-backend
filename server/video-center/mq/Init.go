package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"video-center/conf"
	"video-center/pkg/util"
)

var RabbitChannel *amqp.Channel

func Init() {
	host := conf.Viper.GetString("mq.rabbitmq.host")
	port := conf.Viper.GetInt("mq.rabbitmq.port")
	user := conf.Viper.GetString("mq.rabbitmq.user")
	passwd := conf.Viper.GetString("mq.rabbitmq.passwd")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/tiktokmsg", user, passwd, host, port))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	RabbitChannel = ch
	if err != nil {
		panic(err)
	}
	var exchangeName = "douyinFavMQ"
	var createQueueName = "create_fav_queue"
	var deleteQueueName = "delete_fav_queue"
	// 声明交换机
	err = RabbitChannel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.ExchangeDeclare error", err)
	}
	// 声明队列
	_, err = RabbitChannel.QueueDeclare(
		createQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error", err)
	}
	_, err = RabbitChannel.QueueDeclare(
		deleteQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueDeclare error", err)
	}

	// 绑定队列到交换机，同时设置路由键
	err = RabbitChannel.QueueBind(
		createQueueName,
		"create",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueBind error", err)
	}

	err = RabbitChannel.QueueBind(
		deleteQueueName,
		"delete",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		util.LogrusObj.Error("RabbitChannel.QueueBind error", err)
	}
}
