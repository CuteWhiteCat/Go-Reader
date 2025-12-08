import axios from 'axios'
import type { ApiResponse } from '../types'

const api = axios.create({
  baseURL: '/api',
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
