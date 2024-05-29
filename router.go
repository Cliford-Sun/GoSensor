package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func setupRouter(mq *Mq, db *sql.DB) *gin.Engine {
	r := gin.Default()

	// 创建一个通道用于传递从RabbitMQ读取的数据
	sensorDataChannel := make(chan SensorData)

	// 启动RabbitMQ消费者
	go consumeRabbitMQ(mq, sensorDataChannel)

	// 实时数据接口
	r.GET("/realtime-data", func(c *gin.Context) {
		select {
		case sensorData := <-sensorDataChannel:
			c.JSON(http.StatusOK, sensorData)
		default:
			c.JSON(http.StatusOK, gin.H{"message": "No data available"})
		}
	})

	// 获取数据库中高温或高湿度数据接口
	r.GET("/warning-data", func(c *gin.Context) {
		rows, err := db.Query("SELECT Temperature, Humidity, time FROM sdata WHERE TemWarning = true OR HumWarning = true")
		if err != nil {
			log.Printf("Query error: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database query error"})
			return
		}
		defer rows.Close()

		var warningData []SensorData
		for rows.Next() {
			var sensorData SensorData
			err := rows.Scan(&sensorData.Temperature, &sensorData.Humidity, &sensorData.Timestamp)
			if err != nil {
				log.Printf("Row scan error: %s", err)
				continue
			}
			warningData = append(warningData, sensorData)
		}

		c.JSON(http.StatusOK, warningData)
	})

	return r
}
