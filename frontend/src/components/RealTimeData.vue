<template>
  <div>
    <div ref="lineChart" style="width: 100%; height: 400px;"></div>
    <div style="display: flex; justify-content: space-around; margin-top: 20px;">
      <div ref="temperatureGauge" style="width: 45%; height: 400px;"></div>
      <div ref="humidityGauge" style="width: 45%; height: 400px;"></div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import axios from 'axios';

export default {
  data() {
    return {
      lineChart: null,
      temperatureGauge: null,
      humidityGauge: null,
      intervalId: null,
    };
  },
  mounted() {
    this.initCharts();
    this.fetchData();
    this.intervalId = setInterval(this.fetchData, 1000); // Fetch data every second
  },
  beforeDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  },
  methods: {
    initCharts() {
      this.lineChart = echarts.init(this.$refs.lineChart);
      this.temperatureGauge = echarts.init(this.$refs.temperatureGauge);
      this.humidityGauge = echarts.init(this.$refs.humidityGauge);

      //折线图
      this.lineChart.setOption({
        title: {
          text: 'Real-time Sensor Data',
        },
        tooltip: {
          trigger: 'axis',
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: [],
        },
        yAxis: {
          type: 'value',
        },
        series: [
          {
            name: 'Temperature',
            type: 'line',
            data: [],
          },
          {
            name: 'Humidity',
            type: 'line',
            data: [],
          },
        ],
      });

      //温度仪表盘
      this.temperatureGauge.setOption({
        title: {
          text: 'Temperature',
        },
        series: [
          {
            name: 'Temperature',
            type: 'gauge',
            min: 0,
            max: 40,
            splitNumber: 10,
            axisLine: {
              lineStyle: {
                color: [[0.25, '#48b'], [0.875, '#f80'], [1, '#c23531']],
                width: 8,
              },
            },
            pointer: {
              width: 5,
            },
            detail: {
              formatter: '{value}°C',
            },
            data: [{ value: 0, name: 'Temperature' }],
          },
        ],
      });

      //湿度仪表盘
      this.humidityGauge.setOption({
        title: {
          text: 'Humidity',
        },
        series: [
          {
            name: 'Humidity',
            type: 'gauge',
            min: 0,
            max: 100,
            splitNumber: 10,
            axisLine: {
              lineStyle: {
                color: [[0.8, '#48b'], [1, '#c23531']],
                width: 8,
              },
            },
            pointer: {
              width: 5,
            },
            detail: {
              formatter: '{value}%',
            },
            data: [{ value: 0, name: 'Humidity' }],
          },
        ],
      });
    },
    async fetchData() {
      try {
        const response = await axios.get('http://localhost:8080/realtime-data');
        const data = response.data;

        if (data) {
          const timestamp = new Date().toLocaleTimeString();

          //更新折线图
          const lineChartOption = this.lineChart.getOption();
          lineChartOption.xAxis[0].data.push(timestamp);
          lineChartOption.series[0].data.push(data.temperature);
          lineChartOption.series[1].data.push(data.humidity);
          if (lineChartOption.xAxis[0].data.length > 20) {
            lineChartOption.xAxis[0].data.shift();
            lineChartOption.series[0].data.shift();
            lineChartOption.series[1].data.shift();
          }
          this.lineChart.setOption(lineChartOption);

          //更新仪表盘
          this.temperatureGauge.setOption({
            series: [
              {
                data: [{ value: data.temperature, name: 'Temperature' }],
              },
            ],
          });
          this.humidityGauge.setOption({
            series: [
              {
                data: [{ value: data.humidity, name: 'Humidity' }],
              },
            ],
          });
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    },
  },
};
</script>

<style>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  background-color: #b8ecd5;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.charts {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.chart {
  width: 100%;
  height: 400px;
  margin-bottom: 20px;
}

.gauges {
  display: flex;
  justify-content: space-around;
  width: 100%;
}

.gauge {
  width: 45%;
  height: 300px;
}

.chart, .gauge {
  background-color: #9fd8e2;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
</style>