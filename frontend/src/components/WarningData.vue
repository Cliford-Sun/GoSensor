<template>
  <div class="warning-data-container">
    <h2>警告数据</h2>
    <div class="button-container">
      <button @click="fetchTemWarningData" class="fetch-Tem-button">获取最近十条极端温度数据</button>
      <button @click="fetchHumWarningData" class="fetch-Hum-button">获取最近十条极端湿度数据</button>
      <button @click="clearWarningData" class="clear-button">清空数据</button>
    </div>
    <ul v-if="warningData.length" class="data-list">
      <li v-for="(data, index) in warningData" :key="index" class="data-item">
        <p :class="{
          'high-temp': data.temperature > 35,
          'low-temp': data.temperature < 5
        }">温度: {{ data.temperature }}°C</p>
        <p :class="{
          'high-humidity': data.humidity > 80,
          'low-humidity': data.humidity < 20
        }">湿度: {{ data.humidity }}%</p>
        <p>时间: {{ data.timestamp }}</p>
      </li>
    </ul>
    <div v-else>
      <p>没有极端数据。</p>
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
        <p :class="{
          'high-temp': data.temperature > 35,
          'low-temp': data.temperature < 5
        }">温度: {{ data.temperature }}°C</p>
        <p :class="{
          'high-humidity': data.humidity > 80,
          'low-humidity': data.humidity < 20
        }">湿度: {{ data.humidity }}%</p>
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
  font-family: 'Roboto', Arial, sans-serif;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin: 10px;
}

h2 {
  color: #333;
  text-align: center;
  margin-bottom: 20px;
}

.button-container {
  margin-bottom: 20px;
  text-align: center;
}

.fetch-Tem-button, .fetch-Hum-button, .query-button, .clear-button {
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin: 0 10px;
  transition: background-color 0.3s ease;
}

.fetch-Tem-button:hover, .fetch-Hum-button:hover, .query-button:hover, .clear-button:hover {
  background-color: #45a049;
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
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
  transition: all 0.2s ease-in-out;
}

.data-item:hover {
  background-color: #f0f0f0;
  transform: translateY(-2px);
}

.data-item p {
  margin: 5px 0;
}

.high-temp, .high-humidity {
  color: #ff3b10;
  font-weight: bold;
}

.low-temp, .low-humidity {
  color: #35a5f1;
  font-weight: bold;
}

form {
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
}

form label {
  margin-bottom: 5px;
  font-weight: bold;
}

form input[type="date"],
form input[type="number"] {
  padding: 10px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
}

form input[type="date"]:focus,
form input[type="number"]:focus {
  outline: none;
  border-color: #4CAF50;
  box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
}

form button.query-button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

form button.query-button:hover {
  background-color: #45a049;
}

@media (max-width: 768px) {
  form {
    flex-direction: column;
  }
}
</style>
