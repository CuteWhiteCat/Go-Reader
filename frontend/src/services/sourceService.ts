import api from './api'
import type { BookSource, CreateSourceRequest, UpdateSourceRequest } from '../types'

export const sourceService = {
  // Get all sources
  async getAllSources(): Promise<BookSource[]> {
    const response = await api.get('/sources')
    return response.data || []
  },

  // Create a new source
  async createSource(data: CreateSourceRequest): Promise<BookSource> {
    const response = await api.post('/sources', data)
    return response.data
  },

  // Update a source
  async updateSource(id: string, data: UpdateSourceRequest): Promise<BookSource> {
    const response = await api.put(`/sources/${id}`, data)
    return response.data
  },

  // Delete a source
  async deleteSource(id: string): Promise<void> {
    await api.delete(`/sources/${id}`)
  },

  // Search books from a source
  async searchBooks(sourceId: string, query: string): Promise<any[]> {
    const response = await api.post(`/sources/${sourceId}/search`, { query })
    return response.data || []
  },

  // Import a book (crawl chapters and save)
  async importBook(sourceId: string, payload: { book_url: string; title?: string; author?: string; description?: string }) {
    const response = await api.post(`/sources/${sourceId}/import`, payload)
    return response.data
  },

  // Download book from a source
  async downloadBook(sourceId: string, bookUrl: string): Promise<string> {
    const response = await api.post(`/sources/${sourceId}/download`, { book_url: bookUrl })
    return response.data?.content || ''
  },
}
