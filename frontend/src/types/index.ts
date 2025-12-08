// Book types
export interface Book {
  id: string
  title: string
  author: string
  description: string
  cover_path: string
  file_path: string
  file_format: 'txt' | 'md' | 'epub' | 'web'
  file_size: number
  created_at: string
  updated_at: string
  tags?: Tag[]
}

export interface CreateBookRequest {
  title: string
  author?: string
  description?: string
  file_path: string
  file_format: 'txt' | 'md' | 'epub' | 'web'
  tag_ids?: string[]
}

export interface UpdateBookRequest {
  title?: string
  author?: string
  description?: string
  cover_path?: string
  tag_ids?: string[]
}

// Chapter types
export interface Chapter {
  id: string
  book_id: string
  chapter_number: number
  volume_number?: number
  volume_chapter_number?: number
  title: string
  content?: string
  word_count: number
  created_at: string
}

export interface ChapterSummary {
  id: string
  book_id: string
  chapter_number: number
  volume_number?: number
  volume_chapter_number?: number
  title: string
  word_count: number
  created_at: string
}

// Tag types
export interface Tag {
  id: string
  name: string
  color: string
  created_at: string
}

export interface CreateTagRequest {
  name: string
  color?: string
}

export interface UpdateTagRequest {
  name?: string
  color?: string
}

// Progress types
export interface ReadingProgress {
  book_id: string
  current_chapter: number
  current_position: number
  progress_percentage: number
  last_read_at: string
}

export interface UpdateProgressRequest {
  current_chapter: number
  current_position: number
  progress_percentage: number
}

// Bookmark types
export interface Bookmark {
  id: string
  book_id: string
  chapter_id: string
  position: number
  note: string
  created_at: string
}

export interface CreateBookmarkRequest {
  book_id: string
  chapter_id: string
  position: number
  note?: string
}

// Source types
export interface BookSource {
  id: string
  name: string
  url: string
  type: string
  rules: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface CreateSourceRequest {
  name: string
  url: string
  type?: string
  rules?: string
  enabled?: boolean
}

export interface UpdateSourceRequest {
  name?: string
  url?: string
  type?: string
  rules?: string
  enabled?: boolean
}

// Settings types
export interface AppSettings {
  theme: 'light' | 'dark'
  font_size: number
  line_spacing: number
  page_margin: number
  default_library_path: string
}

// API Response types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}
