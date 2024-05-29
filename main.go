package main

import (
	"log"
	"net"
	"strconv"
)

func main() {
	initDB()
	defer closeDB(db, stmt)

	mq := NewRabbitMQ("SensorQueue", "SensorExchange", "SensorExchange")
	defer mq.ReleaseMq()

	server := newServer("127.0.0.1", 9000)
	// 开启模拟传感器的数据发送
	go Start(server.serverIP, server.serverPort)

	// 建立udp链接
	addr, err := net.ResolveUDPAddr("udp", server.serverIP+":"+strconv.Itoa(server.serverPort))
	if err != nil {
		log.Fatalln("Resolve udp err:", err)
	}

	//监听udp链接
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln("ListenUDP err:", err)
	}

	//关闭监听
	defer conn.Close()

	//建立go程接收传感器数据
	go ReceiveJsonData(conn, mq)

	//阻塞主程序
	select {}
}
