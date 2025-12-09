import axios from 'axios'
import type { ApiResponse } from '../types'

// In production, the frontend is served from file://, so we need the full URL.
// In development, we use a relative path to utilize the Vite proxy.
const isProduction = import.meta.env.PROD

const api = axios.create({
  baseURL: isProduction ? 'http://localhost:8080/api' : '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor to extract data from ApiResponse
api.interceptors.response.use(
  (response) => {
    const data: ApiResponse<any> = response.data
    if (data.success) {
      return { ...response, data: data.data }
    }
    return Promise.reject(new Error(data.error || 'Unknown error'))
  },
  (error) => {
    const message = error.response?.data?.error || error.message || 'An error occurred'
    return Promise.reject(new Error(message))
  }
)

export default api
