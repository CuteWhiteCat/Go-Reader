import api from './api'
import type { Book, Chapter, ChapterSummary, CreateBookRequest, UpdateBookRequest } from '../types'

export const bookService = {
  // Get all books
  async getAllBooks(): Promise<Book[]> {
    const response = await api.get('/books')
    return response.data || []
  },

  // Get book by ID
  async getBook(id: string): Promise<Book> {
    const response = await api.get(`/books/${id}`)
    return response.data
  },

  // Create a new book
  async createBook(data: CreateBookRequest): Promise<Book> {
    const response = await api.post('/books', data)
    return response.data
  },

  // Update a book
  async updateBook(id: string, data: UpdateBookRequest): Promise<Book> {
    const response = await api.put(`/books/${id}`, data)
    return response.data
  },

  // Delete a book
  async deleteBook(id: string): Promise<void> {
    await api.delete(`/books/${id}`)
  },

  // Get book content (all chapters with content)
  async getBookContent(id: string): Promise<Chapter[]> {
    const response = await api.get(`/books/${id}/content`)
    return response.data || []
  },

  // Get book chapters (summaries only)
  async getBookChapters(id: string): Promise<ChapterSummary[]> {
    const response = await api.get(`/books/${id}/chapters`)
    return response.data || []
  },

  // Get a specific chapter
  async getChapter(bookId: string, chapterNumber: number): Promise<Chapter> {
    const response = await api.get(`/books/${bookId}/chapters/${chapterNumber}`)
    return response.data
  },
}
