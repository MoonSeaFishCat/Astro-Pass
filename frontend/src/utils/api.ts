import axios from 'axios'

const API_BASE_URL = '/api'

// 创建axios实例
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('astro-pass-auth')
    if (token) {
      try {
        const authData = JSON.parse(token)
        if (authData.accessToken) {
          config.headers.Authorization = `Bearer ${authData.accessToken}`
        }
      } catch (e) {
        // 忽略解析错误
      }
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config

    // 如果是401错误且未重试过
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        // 尝试刷新token
        const token = localStorage.getItem('astro-pass-auth')
        if (token) {
          const authData = JSON.parse(token)
          if (authData.refreshToken) {
            const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
              refresh_token: authData.refreshToken,
            })

            const { data } = response.data
            const newAuthData = {
              ...authData,
              accessToken: data.access_token,
              refreshToken: data.refresh_token,
            }
            localStorage.setItem('astro-pass-auth', JSON.stringify(newAuthData))

            // 重试原请求
            originalRequest.headers.Authorization = `Bearer ${data.access_token}`
            return apiClient(originalRequest)
          }
        }
      } catch (refreshError) {
        // 刷新失败，清除认证信息并跳转到登录页
        localStorage.removeItem('astro-pass-auth')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default apiClient


