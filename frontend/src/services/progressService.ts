import api from './api'
import type { ReadingProgress, UpdateProgressRequest, Bookmark, CreateBookmarkRequest } from '../types'

export const progressService = {
  // Get reading progress
  async getProgress(bookId: string): Promise<ReadingProgress> {
    const response = await api.get(`/progress/${bookId}`)
    return response.data
  },

  // Update reading progress
  async updateProgress(bookId: string, data: UpdateProgressRequest): Promise<ReadingProgress> {
    const response = await api.put(`/progress/${bookId}`, data)
    return response.data
  },

  // Get bookmarks for a book
  async getBookmarks(bookId: string): Promise<Bookmark[]> {
    const response = await api.get(`/bookmarks/${bookId}`)
    return response.data || []
  },

  // Create a bookmark
  async createBookmark(data: CreateBookmarkRequest): Promise<Bookmark> {
    const response = await api.post('/bookmarks', data)
    return response.data
  },

  // Delete a bookmark
  async deleteBookmark(id: string): Promise<void> {
    await api.delete(`/bookmarks/${id}`)
  },
}
