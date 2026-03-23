import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// 1. 创建 axios 实例
const service = axios.create({
  // 这里留空，让它自动匹配当前域名
  // 也就是请求会发给 http://localhost:5173/api/v1/...
  // 然后被 Vite 代理捕获
  baseURL: '',
  timeout: 5000, // 请求超时时间：5秒
})

// 2. 请求拦截器 (Request Interceptor)
service.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 后端统一响应格式
interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

let isRefreshing = false
let requestsQueue: Array<(token: string) => void> = []

// 3. 响应拦截器 (Response Interceptor)
service.interceptors.response.use(
  (response) => {
    // 2xx 范围内的状态码都会触发这里
    const res = response.data as ApiResponse

    // 后端统一包装格式：{ code, msg, data }
    // code 为 200 表示成功，直接返回 data
    if (res.code === 200) {
      return res.data as unknown as AxiosResponse
    }

    // code 不为 200，说明业务逻辑错误
    ElMessage.error(res.msg || '请求失败')
    return Promise.reject(new Error(res.msg || '请求失败'))
  },
  (error) => {
    // 超出 2xx 范围的状态码都会触发这里（网络错误、404、500等）
    console.error('请求出错:', error)

    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      const originalRequest = error.config
      
      if (!authStore.refreshToken) {
        authStore.logout()
        if (window.location.pathname !== '/auth') window.location.href = '/auth'
        return Promise.reject(error)
      }

      if (!isRefreshing) {
        isRefreshing = true
        
        // 使用独立的基础 axios 实例重试刷新，避免循环依赖和重复拦截
        return axios.post('/api/v1/users/refresh', { refresh_token: authStore.refreshToken })
          .then(res => {
            const data = res.data.data
            authStore.updateTokens(data.token, data.refresh_token)
            
            // 修正原请求的 header
            originalRequest.headers.Authorization = `Bearer ${data.token}`
            
            // 执行队列中等待的所有请求
            requestsQueue.forEach(cb => cb(data.token))
            requestsQueue = []
            
            // 继续重新发送当前失败的请求
            return service(originalRequest)
          })
          .catch(refreshErr => {
            // 刷新也失败了，说明 refresh_token 过期或无效
            authStore.logout()
            ElMessage.error('登录状态已过期，请重新登录')
            if (window.location.pathname !== '/auth') window.location.href = '/auth'
            return Promise.reject(refreshErr)
          })
          .finally(() => {
            isRefreshing = false
          })
      } else {
        // 如果正在刷新，挂起当前的请求
        return new Promise((resolve) => {
          requestsQueue.push((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            resolve(service(originalRequest))
          })
        })
      }
    }

    // 统一弹出错误提示 (过滤掉刷新 token 时弹出的重复错误消息)
    if (error.response?.status !== 401) {
      const msg = error.response?.data?.msg || error.message || '网络请求失败，请稍后重试'
      ElMessage.error(msg)
    }

    return Promise.reject(error)
  },
)

export default service
