import { create } from 'zustand'
import type { Book, Tag } from '../types'
import { bookService } from '../services/bookService'
import { tagService } from '../services/tagService'

interface BookStore {
  books: Book[]
  tags: Tag[]
  selectedBook: Book | null
  isLoading: boolean
  error: string | null

  // Actions
  fetchBooks: () => Promise<void>
  fetchTags: () => Promise<void>
  selectBook: (book: Book | null) => void
  addBook: (book: Book) => void
  updateBook: (id: string, book: Partial<Book>) => void
  removeBook: (id: string) => void
  clearError: () => void
}

export const useBookStore = create<BookStore>((set) => ({
  books: [],
  tags: [],
  selectedBook: null,
  isLoading: false,
  error: null,

  fetchBooks: async () => {
    set({ isLoading: true, error: null })
    try {
      const books = await bookService.getAllBooks()
      set({ books, isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },

  fetchTags: async () => {
    try {
      const tags = await tagService.getAllTags()
      set({ tags })
    } catch (error) {
      set({ error: (error as Error).message })
    }
  },

  selectBook: (book) => {
    set({ selectedBook: book })
  },

  addBook: (book) => {
    set((state) => ({ books: [...state.books, book] }))
  },

  updateBook: (id, updatedBook) => {
    set((state) => ({
      books: state.books.map((book) =>
        book.id === id ? { ...book, ...updatedBook } : book
      ),
    }))
  },

  removeBook: (id) => {
    set((state) => ({
      books: state.books.filter((book) => book.id !== id),
      selectedBook: state.selectedBook?.id === id ? null : state.selectedBook,
    }))
  },

  clearError: () => {
    set({ error: null })
  },
}))
