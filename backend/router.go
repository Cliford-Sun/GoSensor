package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter(mq *Mq, db *sql.DB) *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

	// 获取数据库中高温数据接口
	r.GET("/Tem-warning-data", func(c *gin.Context) {
		rows, err := db.Query("SELECT Temperature, Humidity, time FROM sdata WHERE TemWarning = true ORDER BY time DESC LIMIT 10")
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

// 获取数据库中高湿度数据接口
	r.GET("/Hum-warning-data", func(c *gin.Context) {
		rows, err := db.Query("SELECT Temperature, Humidity, time FROM sdata WHERE HumWarning = true ORDER BY time DESC LIMIT 10")
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

	// 添加查询接口
	r.GET("/query-warning-data", func(c *gin.Context) {
		date := c.Query("date")
		hour := c.Query("hour")

		// 解析日期和小时
		startTime, err := time.Parse("2006-01-02 15", date+" "+hour)
		if err != nil {
			log.Printf("Time parse error: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date or hour format"})
			return
		}
		endTime := startTime.Add(time.Hour)

		rows, err := db.Query("SELECT Temperature, Humidity, time FROM sdata WHERE (TemWarning = true OR HumWarning = true) AND time BETWEEN ? AND ? ORDER BY time ", startTime, endTime)
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
