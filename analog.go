package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// SensorData 模拟传感器的数据
type SensorData struct {
	Temperature int    `json:"temperature"`
	Humidity    int    `json:"humidity"`
	Timestamp   string `json:"timestamp"`
}

// randDate 生成随机的传感器温度和湿度数据，并且将它转换成json格式，以便传输和接受处理
func (t *SensorData) randData() []byte {
	t.Temperature = rand.Intn(31) + 10
	t.Humidity = rand.Intn(101)
	t.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	jStr, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return jStr
}

func Start(serverip string, port int) {
	//建立udp链接
	conn, err := net.Dial("udp", serverip+":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	//关闭链接
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	fmt.Println("running......")

	s := SensorData{}
	for {
		//发送传感器数据
		_, err := conn.Write(s.randData())
		if err != nil {
			panic(err)
		}

		//间隔一秒
		time.Sleep(time.Second)
	}

}
