package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	stmt *sql.Stmt
)

// Server 服务器对象
type Server struct {
	serverIP   string
	serverPort int
}

// init 初始化数据库并且预编译SQL语句
func init() {
	// 初始化数据库连接池
	var err error
	db, err = sql.Open("mysql", "root:123SZCszc@tcp(127.0.0.1:3306)/sensordata")
	if err != nil {
		fmt.Println("open Sqlserver err:", err)
	}
	// 预编译SQL语句
	stmt, err = db.Prepare("INSERT INTO sdata(Temperature, Humidity, time, TemWarning, HumWarning) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("prepare Sqlserver err:", err)
	}
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
	temWarn, humWarn := Warning(sensor)
	_, err := stmt.Exec(sensor.Temperature, sensor.Humidity, sensor.Timestamp, temWarn, humWarn)
	if err != nil {
		fmt.Println("Data insert err:", err)
	}

	//fmt.Println(sensor.Timestamp, " : temperature = ", sensor.Temperature, ";Humidity = ", sensor.Humidity, "high Temperature?", temWarn, "high Humidity?", humWarn)
}

// ReceiveJsonData 接收UDP传输的数据并插入数据库
func ReceiveJsonData(conn *net.UDPConn, mq *Mq) {
	sensor := &SensorData{}
	buf := make([]byte, 1024) //缓存区重用
	for {
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
			mq.PublishMq(sensor)
		}
	}
}

func main() {
	mq := NewRabbitMQ("SensorQueue", "SensorExchange", "SensorExchange")
	defer mq.ReleaseMq()

	server := newServer("127.0.0.1", 9000)
	// 开启模拟传感器的数据发送
	go Start(server.serverIP, server.serverPort)

	// 建立udp链接
	addr, err := net.ResolveUDPAddr("udp", server.serverIP+":"+strconv.Itoa(server.serverPort))
	if err != nil {
		fmt.Println("Resolve udp err:", err)
	}

	//监听udp链接
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("ListenUDP err:", err)
	}

	//关闭监听
	defer conn.Close()

	//建立go程接收传感器数据
	go ReceiveJsonData(conn, mq)

	//阻塞主程序
	select {}
}
