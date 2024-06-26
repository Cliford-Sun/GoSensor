package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"

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
func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:123SZCszc@tcp(127.0.0.1:3306)/sensordata")
	if err != nil {
		log.Fatalln("open Sqlserver err:", err)
	}
	stmt, err = db.Prepare("INSERT INTO sdata(Temperature, Humidity, time, highTemWarning, highHumWarning, lowTemWarning, lowHumWarning) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalln("prepare Sqlserver err:", err)
	}
}

// closeDB 关闭数据库连接
func closeDB(db *sql.DB, stmt *sql.Stmt) {
	if db != nil {
		db.Close()
	}
	if stmt != nil {
		stmt.Close()
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

// Warning 判断是否存在极端数据情况
func Warning(sensor *SensorData) (highTemWarning bool,highHumWarning bool,lowTemWarning bool,lowHumWarning bool) {
	highTemWarning = sensor.Temperature > 35
	highHumWarning = sensor.Humidity > 80
	lowTemWarning = sensor.Temperature < 5
	lowHumWarning =	sensor.Humidity < 20
	return
}

// InsertSensorData 将传感器数据存入本地数据库
func InsertSensorData(sensor *SensorData) {
	highTemWarn, highHumWarn,lowTemWarning,lowHumWarning:= Warning(sensor)
	_, err := stmt.Exec(sensor.Temperature, sensor.Humidity, sensor.Timestamp, highTemWarn, highHumWarn, lowTemWarning, lowHumWarning)
	if err != nil {
		fmt.Println("Data insert err:", err)
	}
}

// ReceiveJsonData 接收UDP传输的数据并插入数据库
func ReceiveJsonData(conn *net.UDPConn, mq *Mq) {
	sensor := &SensorData{}
	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("conn ReadFromUDP err:", err)
			continue
		}
		if n > 0 {
			err := json.Unmarshal(buf[:n], sensor)
			if err != nil {
				fmt.Println("json Unmarshal err:", err)
				continue
			}
			InsertSensorData(sensor)
			mq.PublishMq(sensor)
		}
	}
}
