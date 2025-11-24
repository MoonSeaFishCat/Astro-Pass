import { create } from 'zustand'
import axios from 'axios'

const API_BASE_URL = '/api'

interface User {
  id: number
  uuid: string
  username: string
  email: string
  nickname?: string
  mfa_enabled?: boolean
  email_verified?: boolean
  roles?: Array<{ id: number; name: string; display_name: string }>
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  register: (username: string, email: string, password: string, nickname?: string) => Promise<void>
  logout: () => void
  refreshAccessToken: () => Promise<void>
}

// 简单的持久化实现
const STORAGE_KEY = 'astro-pass-auth'

const loadFromStorage = (): Partial<AuthState> => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    return stored ? JSON.parse(stored) : {}
  } catch {
    return {}
  }
}

const saveToStorage = (state: Partial<AuthState>) => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({
      user: state.user,
      accessToken: state.accessToken,
      refreshToken: state.refreshToken,
      isAuthenticated: state.isAuthenticated,
    }))
  } catch {
    // 忽略存储错误
  }
}

export const useAuthStore = create<AuthState>()(
  (set, get) => {
    // 初始化时从localStorage加载
    const stored = loadFromStorage()
    const initialState = {
      user: (stored.user as User) || null,
      accessToken: stored.accessToken || null,
      refreshToken: stored.refreshToken || null,
      isAuthenticated: stored.isAuthenticated || false,
    }
    
    if (initialState.accessToken) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${initialState.accessToken}`
    }

    return ({
      ...initialState,

      login: async (username: string, password: string) => {
        try {
          const response = await axios.post(`${API_BASE_URL}/auth/login`, {
            username,
            password,
          })

          const { data } = response.data
          const newState = {
            user: data.user,
            accessToken: data.access_token,
            refreshToken: data.refresh_token,
            isAuthenticated: true,
          }
          set(newState)
          saveToStorage(newState)

          // 设置axios默认header
          axios.defaults.headers.common['Authorization'] = `Bearer ${data.access_token}`
        } catch (error: any) {
          throw new Error(error.response?.data?.message || '登录失败')
        }
      },

      register: async (username: string, email: string, password: string, nickname?: string) => {
        try {
          const response = await axios.post(`${API_BASE_URL}/auth/register`, {
            username,
            email,
            password,
            nickname,
          })

          const { data } = response.data
          // 注册成功后自动登录
          await get().login(username, password)
        } catch (error: any) {
          throw new Error(error.response?.data?.message || '注册失败')
        }
      },

      logout: () => {
        const newState = {
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        }
        set(newState)
        saveToStorage(newState)
        localStorage.removeItem(STORAGE_KEY)
        delete axios.defaults.headers.common['Authorization']
      },

      refreshAccessToken: async () => {
        const { refreshToken } = get()
        if (!refreshToken) {
          throw new Error('没有刷新令牌')
        }

        try {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })

          const { data } = response.data
          const newState = {
            accessToken: data.access_token,
            refreshToken: data.refresh_token,
          }
          set(newState)
          saveToStorage({ ...get(), ...newState })

          axios.defaults.headers.common['Authorization'] = `Bearer ${data.access_token}`
        } catch (error: any) {
          // 刷新失败，清除认证状态
          get().logout()
          throw new Error(error.response?.data?.message || '刷新令牌失败')
        }
      },
    })
  }
)

// 初始化axios拦截器
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        await useAuthStore.getState().refreshAccessToken()
        originalRequest.headers['Authorization'] = `Bearer ${useAuthStore.getState().accessToken}`
        return axios(originalRequest)
      } catch (refreshError) {
        useAuthStore.getState().logout()
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

