package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// MQURL RabbitMQ服务器地址
const MQURL = "amqp://guest:guest@localhost:5672/"

// Mq RabbitMQ结构体
type Mq struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel

	//队列名称
	QueueName string
	//交换机
	Exchange string
	//routing Key
	RoutingKey string
	//MQl链接字符串
	MQurl string
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName, exchange, routineKey string) *Mq {
	rabbitMQ := Mq{
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routineKey,
		MQurl:      MQURL,
	}
	var err error
	//创建rabbitmq链接
	rabbitMQ.Conn, err = amqp.Dial(rabbitMQ.MQurl)
	if err != nil {
		log.Fatalln("RabbitMQ connection error:", err)
	}
	//创建Channel
	rabbitMQ.Channel, err = rabbitMQ.Conn.Channel()
	if err != nil {
		log.Fatalln("RabbitMQ channel error:", err)
	}
	return &rabbitMQ
}

// ReleaseMq 释放资源
func (mq *Mq) ReleaseMq() {
	mq.Conn.Close()
	mq.Channel.Close()
}

// PublishMq 将获取的模拟数据传输到RabbitMQ的服务器中
func (mq *Mq) PublishMq(sensor *SensorData) {
	//声明队列
	_, err := mq.Channel.QueueDeclare(
		mq.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Queue Declare error:", err)
	}

	//声明交换器
	err = mq.Channel.ExchangeDeclare(
		mq.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Exchange Declare error:", err)
	}

	//建立关系
	err = mq.Channel.QueueBind(
		mq.QueueName,
		mq.RoutingKey,
		mq.Exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Queue Bind error:", err)
	}

	//将数据转化成json格式
	jStr, err := json.Marshal(sensor)
	if err != nil {
		fmt.Println("Json Marshal error:", err)
	}

	//发送json格式消息
	err = mq.Channel.Publish(
		mq.Exchange,
		mq.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jStr,
		},
	)
	if err != nil {
		fmt.Println("Publish error:", err)
	} else {
		fmt.Println(sensor)
	}
}
