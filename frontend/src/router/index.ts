import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import RealTimeData from '../components/RealTimeData.vue'
import WarningData from '../components/WarningData.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/realtime-data' // 将根路径重定向到实时数据路由
  },
  {
    path: '/realtime-data',
    name: 'RealTimeData',
    component: RealTimeData
  },
  {
    path: '/warning-data',
    name: 'WarningData',
    component: WarningData
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 移除之前添加的全局前置守卫
export default router
