<template>
  <div class="warning-data-container">
    <h2>高温或高湿度数据</h2>
    <div class="button-container">
      <button @click="fetchTemWarningData" class="fetch-Tem-button">获取最近十条高温数据</button>
      <button @click="fetchHumWarningData" class="fetch-Hum-button">获取最近十条高湿度数据</button>
      <button @click="clearWarningData" class="clear-button">清空数据</button>
    </div>
    <ul v-if="warningData.length" class="data-list">
      <li v-for="(data, index) in warningData" :key="index" class="data-item">
        <p :class="{ 'high-temp': data.temperature > 30 }">温度: {{ data.temperature }}°C</p>
        <p :class="{ 'high-humidity': data.humidity > 80 }">湿度: {{ data.humidity }}%</p>
        <p>时间: {{ data.timestamp }}</p>
      </li>
    </ul>
    <div v-else>
      <p>没有高温或高湿度数据。</p>
    </div>

    <h2>查询预警数据</h2>
    <form @submit.prevent="queryWarningData">
      <label for="date">日期:</label>
      <input type="date" v-model="queryDate" required>

      <label for="hour">小时:</label>
      <input type="number" v-model="queryHour" min="0" max="23" required>

      <button type="submit" class="query-button">查询</button>
    </form>
    <ul v-if="queriedData.length" class="data-list">
      <li v-for="(data, index) in queriedData" :key="index" class="data-item">
        <p :class="{ 'high-temp': data.temperature > 30 }">温度: {{ data.temperature }}°C</p>
        <p :class="{ 'high-humidity': data.humidity > 80 }">湿度: {{ data.humidity }}%</p>
        <p>时间: {{ data.timestamp }}</p>
      </li>
    </ul>
    <div v-else-if="queriedData.length === 0 && queryDate && queryHour !== null">
      <p>没有找到匹配的预警数据。</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      warningData: [],
      queryDate: '',
      queryHour: '',
      queriedData: [],
    };
  },
  methods: {
    fetchTemWarningData() {
      this.warningData = [];
      axios.get('http://localhost:8080/Tem-warning-data')
          .then(response => {
            this.warningData = response.data;
          })
          .catch(error => {
            console.error('Error fetching temperature warning data:', error);
          });
    },
    fetchHumWarningData() {
      this.warningData = [];
      axios.get('http://localhost:8080/Hum-warning-data')
          .then(response => {
            this.warningData = response.data;
          })
          .catch(error => {
            console.error('Error fetching Humidity warning data:', error);
          });
    },
    clearWarningData() {
      this.warningData = [];
      this.queriedData = [];
    },
    queryWarningData() {
      this.queriedData = []; // 清空之前的查询数据
      axios.get('http://localhost:8080/query-warning-data', {
        params: {
          date: this.queryDate,
          hour: this.queryHour,
        }
      })
          .then(response => {
            this.queriedData = response.data;
          })
          .catch(error => {
            console.error('Error querying warning data:', error);
          });
    },
  },
};
</script>

<style>
.warning-data-container {
  padding: 20px;
  font-family: Arial, sans-serif;
}

h2 {
  color: #333;
}

.button-container {
  margin-bottom: 20px;
}

.fetch-Tem-button, .fetch-Hum-button, .query-button, .clear-button {
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-left: 10px;
}

.fetch-Tem-button:hover, .fetch-Hum-button:hover, .query-button:hover, .clear-button:hover {
  background-color: #45a049;
}

.data-list {
  list-style-type: none;
  padding: 0;
}

.data-item {
  background-color: #f9f9f9;
  border: 1px solid #ddd;
  padding: 15px;
  margin-bottom: 10px;
  border-radius: 5px;
}

.data-item p {
  margin: 5px 0;
}

.high-temp {
  color: red;
  font-weight: bold;
}

.high-humidity {
  color: red;
  font-weight: bold;
}
</style>
