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

	QueueName  string
	Exchange   string
	RoutingKey string
	MQurl      string
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
	rabbitMQ.Conn, err = amqp.Dial(rabbitMQ.MQurl)
	if err != nil {
		log.Fatalln("RabbitMQ connection error:", err)
	}
	rabbitMQ.Channel, err = rabbitMQ.Conn.Channel()
	if err != nil {
		log.Fatalln("RabbitMQ channel error:", err)
	}
	return &rabbitMQ
}

// ReleaseMq 释放资源
func (mq *Mq) ReleaseMq() {
	if mq.Conn != nil {
		mq.Conn.Close()
	}
	if mq.Channel != nil {
		mq.Channel.Close()
	}
}

// PublishMq 将获取的模拟数据传输到RabbitMQ的服务器中
func (mq *Mq) PublishMq(sensor *SensorData) {
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

	jStr, err := json.Marshal(sensor)
	if err != nil {
		fmt.Println("Json Marshal error:", err)
	}

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
	}
}

func consumeRabbitMQ(mq *Mq, ch chan<- SensorData) {
	msgs, err := mq.Channel.Consume(
		mq.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	for msg := range msgs {
		var sensorData SensorData
		err := json.Unmarshal(msg.Body, &sensorData)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			continue
		}
		ch <- sensorData
	}
}
