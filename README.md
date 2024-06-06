# GoSensor

## 项目简介
这是一个包含Golang后端和Vue3前端的传感器模拟，学习并构建和组织一个全栈应用。
## 步骤
### 后端
1. 安装Go语言：https://golang.org/doc/install
2. 在`backend`目录下运行：
   ```bash
   go mod tidy
   go build main.go server.go analog.go router.go rabbit.go
   go run main.go
   ```
### 前端
1. 安装Node.js和npm：https://nodejs.org/
2. 在`frontend`目录下运行：
    ```bash
    npm install
    npm run dev
    ```
## 使用说明
访问` http://localhost:5174/ `就可以看前端应用，包含了传感器数据的实时折线图和仪表盘，获取预警数据，查询某时间内的预警数据

## 具体实现
### 传感器模拟端
在后端模拟进行，每一秒都能模拟出0-40的温度数据和0-100的湿度数据，并通过udp进行数据传输
### 后端
后端使用了golang语言来处理来自传感器模拟端的数据，运用了gin的web框架来为前端提供数据接口，实现了实时数据的传输、预警数据的查询传输
使用了本地的mysql数据库来进行数据的存储和更新
通过RabbitMQ来构建消息队列服务的架构
### 前端
前端使用vue3实现，使用了eCharts来实现了折线图、仪表盘的可视化