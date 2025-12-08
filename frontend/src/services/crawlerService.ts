import api from './api'
import type { Book } from '@/types'

export interface CrawlSearchResult {
  title: string
  author: string
  latest: string
  url: string
}

export const crawlerService = {
  async search(query: string): Promise<CrawlSearchResult[]> {
    const res = await api.post('/crawler/search', { query })
    return res.data || []
  },
  async importBook(payload: { title: string; author?: string; latest?: string; url: string }): Promise<Book> {
    const res = await api.post('/crawler/import', payload)
    return res.data
  },
  async startImport(payload: { title: string; author?: string; latest?: string; url: string }): Promise<string> {
    const res = await api.post('/crawler/import/start', payload)
    return res.data?.job_id
  },
  async getImportStatus(jobId: string): Promise<{
    id: string
    status: string
    error?: string
    total: number
    done: number
    book_id?: string
  }> {
    const res = await api.get('/crawler/import/status', { params: { id: jobId } })
    return res.data
  },
}
