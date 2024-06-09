<template>
  <div class="real-time-data-container">
    <h2>实时数据</h2>
    <div ref="temperatureChart" style="width: 100%; height: 400px;"></div>
    <div ref="humidityChart" style="width: 100%; height: 400px; margin-top: 20px;"></div>
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
  name: "RealTimeData",
  data() {
    return {
      temperatureChart: null,
      humidityChart: null,
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
      this.temperatureChart = echarts.init(this.$refs.temperatureChart);
      this.humidityChart = echarts.init(this.$refs.humidityChart);
      this.temperatureGauge = echarts.init(this.$refs.temperatureGauge);
      this.humidityGauge = echarts.init(this.$refs.humidityGauge);
      // 初始化温度折线图
      this.temperatureChart.setOption({
        title: {
          text: '温度变化折线图',
          left: 'center',
          textStyle: {
            color: '#ff7f50',
            fontSize: 20,
          },
        },
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(50, 50, 50, 0.7)',
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: [],
          axisLine: {
            lineStyle: {
              color: '#ff7f50',
            },
          },
          axisLabel: {
            color: '#ff7f50',
          },
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value} °C',
            color: '#ff7f50',
          },
          axisLine: {
            lineStyle: {
              color: '#ff7f50',
            },
          },
        },
        series: [
          {
            name: '温度',
            type: 'line',
            data: [],
            lineStyle: {
              color: '#ff7f50',
              width: 2,
            },
            itemStyle: {
              color: '#ff7f50',
            },
            areaStyle: {
              color: 'rgba(255, 127, 80, 0.3)',
            },
            markLine: {
              silent: true,
              data: [
                { yAxis: 35, name: '高温警戒线', lineStyle: { color: 'red', type: 'dashed' } },
                { yAxis: 5, name: '低温警戒线', lineStyle: { color: 'blue', type: 'dashed' } }
              ]
            }
          },
        ],
      });

      // 初始化湿度折线图
      this.humidityChart.setOption({
        title: {
          text: '湿度变化折线图',
          left: 'center',
          textStyle: {
            color: '#87cefa',
            fontSize: 20,
          },
        },
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(50, 50, 50, 0.7)',
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: [],
          axisLine: {
            lineStyle: {
              color: '#87cefa',
            },
          },
          axisLabel: {
            color: '#87cefa',
          },
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value} %',
            color: '#87cefa',
          },
          axisLine: {
            lineStyle: {
              color: '#87cefa',
            },
          },
        },
        series: [
          {
            name: '湿度',
            type: 'line',
            data: [],
            lineStyle: {
              color: '#87cefa',
              width: 2,
            },
            itemStyle: {
              color: '#87cefa',
            },
            areaStyle: {
              color: 'rgba(135, 206, 250, 0.3)',
            },
            markLine: {
              silent: true,
              data: [
                { yAxis: 80, name: '高湿警戒线', lineStyle: { color: 'red', type: 'dashed' } },
                { yAxis: 20, name: '低湿警戒线', lineStyle: { color: 'blue', type: 'dashed' } }
              ]
            }
          },
        ],
      });


      //温度仪表盘
      this.temperatureGauge.setOption({
      title: {
        text: '温度仪表盘',
        left: 'center',
        textStyle: {
          color: '#ff3b10',
          fontSize: 20,
        },
      },
      series: [
        {
          name: '温度',
          type: 'gauge',
          min: 0,
          max: 40,
          splitNumber: 10,
          axisLine: {
            lineStyle: {
              color: [[0.125, '#48b'], [0.875, '#f80'], [1, '#c23531']],
              width: 15,
            },
          },
          pointer: {
            width: 5,
            color: '#ff3b10',
          },
          title: {
            offsetCenter: [0, '-30%'],
            color: '#ff3b10',
            fontSize: 16,
          },
          detail: {
            formatter: '{value}°C',
            textStyle: {
              color: '#ff3b10',
              fontSize: 20,
            },
          },
          data: [{ value: 0, name: '温度' }],
        },
      ],
    });

      //湿度仪表盘
      this.humidityGauge.setOption({
        title: {
          text: '湿度仪表盘',
          left: 'center',
          textStyle: {
            color: '#35a5f1',
            fontSize: 20,
          },
        },
        series: [
          {
            name: '湿度',
            type: 'gauge',
            min: 0,
            max: 100,
            splitNumber: 10,
            axisLine: {
              lineStyle: {
                color: [[0.2, '#48b'], [0.8, '#f80'], [1, '#c23531']],
                width: 15,
              },
            },
            pointer: {
              width: 5,
              color: '#35a5f1',
            },
            title: {
              offsetCenter: [0, '-30%'],
              color: '#35a5f1',
              fontSize: 16,
            },
            detail: {
              formatter: '{value}%',
              textStyle: {
                color: '#35a5f1',
                fontSize: 20,
              },
            },
            data: [{ value: 0, name: '湿度' }],
          },
        ],
      });
    },
    async fetchData() {
      try {
        const response = await axios.get('http://localhost:8080/realtime-data');
        const data = response.data;
        const timestamp = new Date().toLocaleTimeString();

        if (data) {
          // 更新温度折线图
          const temperatureOption = this.temperatureChart.getOption();
          temperatureOption.xAxis[0].data.push(timestamp);
          temperatureOption.series[0].data.push(data.temperature);
          if (temperatureOption.xAxis[0].data.length > 20) {
            temperatureOption.xAxis[0].data.shift();
            temperatureOption.series[0].data.shift();
          }
          this.temperatureChart.setOption(temperatureOption);

          // 更新湿度折线图
          const humidityOption = this.humidityChart.getOption();
          humidityOption.xAxis[0].data.push(timestamp);
          humidityOption.series[0].data.push(data.humidity);
          if (humidityOption.xAxis[0].data.length > 20) {
            humidityOption.xAxis[0].data.shift();
            humidityOption.series[0].data.shift();
          }
          this.humidityChart.setOption(humidityOption);

          //更新仪表盘
          this.temperatureGauge.setOption({
            series: [
              {
                data: [{ value: data.temperature, name: '温度' }],
              },
            ],
          });
          this.humidityGauge.setOption({
            series: [
              {
                data: [{ value: data.humidity, name: '湿度' }],
              },
            ],
          });
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    }
  },
};
</script>

<style>
.real-time-data-container {
  padding: 5px;
  font-family: Arial, sans-serif;
  background-color: #fff9f1;
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
  background-color: #f5f5f5; /* 图表背景颜色 */
}

.gauges {
  display: flex;
  justify-content: space-around;
  width: 100%;
}

.gauge {
  width: 45%;
  height: 300px;
  background-color: #f5f5f5; /* 仪表背景颜色 */
}

.chart, .gauge {
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
</style>