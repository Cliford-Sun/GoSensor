package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	serverIP   string
	serverPort int
}

// newServer 建立新的server链接
func newServer(serverIP string, serverPort int) *Server {
	server := Server{
		serverIP:   serverIP,
		serverPort: serverPort,
	}
	return &server
}

// Warning 判断是否温度过高或者湿度过高
func Warning(sensor *SensorData) (temWarning bool, humWarning bool) {
	temWarning = sensor.Temperature > 35
	humWarning = sensor.Humidity > 80
	return
}

// InsertSensorData 将传感器数据存入本地数据库
func InsertSensorData(sensor *SensorData) {
	//链接本地数据库
	db, err := sql.Open("mysql", "root:123SZCszc@tcp(127.0.0.1:3306)/sensordata")
	if err != nil {
		panic(err)
	}

	//关闭链接
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// 准备插入语句
	stmt, err := db.Prepare("INSERT INTO sdata(Temperature, Humidity, time, TemWarning, HumWarning) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}

	//关闭插入语句
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	//分析数据是否温度过高或者湿度过高
	temWarn, humWarn := Warning(sensor)

	// 执行插入语句
	_, err = stmt.Exec(sensor.Temperature, sensor.Humidity, sensor.Timestamp, temWarn, humWarn)
	if err != nil {
		panic(err)
	}

	//在终端打印
	fmt.Println(sensor.Timestamp, " : temperature = ", sensor.Temperature, ";Humidity = ", sensor.Humidity,
		"high Temperature?", temWarn, "high Humidity?", humWarn)
}

// ReceiveJsonData 接收UDP传输的数据并插入数据库
func ReceiveJsonData(conn *net.UDPConn) {
	sensor := &SensorData{}
	for {
		//将数据写入buf
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("conn ReadFromUDP err:", err)
			continue
		}
		if n > 0 {
			// 将数据从json格式转换
			ero := json.Unmarshal(buf[:n], sensor)
			if ero != nil {
				fmt.Println("json Unmarshal err:", ero)
				continue
			}
			InsertSensorData(sensor)
		}
	}
}

func main() {
	server := newServer("127.0.0.1", 9000)
	// 开启模拟传感器的数据发送
	go Start(server.serverIP, server.serverPort)

	// 建立udp链接
	addr, err := net.ResolveUDPAddr("udp", server.serverIP+":"+strconv.Itoa(server.serverPort))
	if err != nil {
		panic(err)
	}

	//监听udp链接
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	//关闭监听
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	//建立go程接收传感器数据
	go ReceiveJsonData(conn)

	select {}
}
