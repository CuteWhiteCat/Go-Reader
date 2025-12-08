import api from './api'
import type { Tag, CreateTagRequest, UpdateTagRequest } from '../types'

export const tagService = {
  // Get all tags
  async getAllTags(): Promise<Tag[]> {
    const response = await api.get('/tags')
    return response.data || []
  },

  // Create a new tag
  async createTag(data: CreateTagRequest): Promise<Tag> {
    const response = await api.post('/tags', data)
    return response.data
  },

  // Update a tag
  async updateTag(id: string, data: UpdateTagRequest): Promise<Tag> {
    const response = await api.put(`/tags/${id}`, data)
    return response.data
  },

  // Delete a tag
  async deleteTag(id: string): Promise<void> {
    await api.delete(`/tags/${id}`)
  },
}
