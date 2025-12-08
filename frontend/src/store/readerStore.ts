import { create } from 'zustand'
import { bookService } from '@/services/bookService'
import type { Chapter, ReadingProgress } from '../types'

interface ReaderStore {
  currentBook: string | null
  chapters: Chapter[]
  currentChapter: number
  progress: ReadingProgress | null
  isReading: boolean
  loadingChapters: Record<number, boolean>

  // Actions
  startReading: (bookId: string, chapters: Chapter[], progress: ReadingProgress) => void
  stopReading: () => void
  setCurrentChapter: (chapterNumber: number) => void
  updateProgress: (progress: ReadingProgress) => void
  nextChapter: () => void
  previousChapter: () => void
  loadChapterContent: (chapterIndex: number) => Promise<void>
}

export const useReaderStore = create<ReaderStore>((set, get) => ({
  currentBook: null,
  chapters: [],
  currentChapter: 0,
  progress: null,
  isReading: false,
  loadingChapters: {},

  startReading: (bookId, chapters, progress) => {
    const safeChapterIndex = Math.min(
      Math.max(progress.current_chapter, 0),
      Math.max(chapters.length - 1, 0)
    )

    set({
      currentBook: bookId,
      chapters,
      currentChapter: safeChapterIndex,
      progress,
      isReading: true,
      loadingChapters: {},
    })
  },

  stopReading: () => {
    set({
      currentBook: null,
      chapters: [],
      currentChapter: 0,
      progress: null,
      isReading: false,
      loadingChapters: {},
    })
  },

  setCurrentChapter: (chapterNumber) => {
    set({ currentChapter: chapterNumber })
  },

  updateProgress: (progress) => {
    set({ progress, currentChapter: progress.current_chapter })
  },

  nextChapter: () => {
    const { currentChapter, chapters } = get()
    if (currentChapter < chapters.length - 1) {
      set({ currentChapter: currentChapter + 1 })
    }
  },

  previousChapter: () => {
    const { currentChapter } = get()
    if (currentChapter > 0) {
      set({ currentChapter: currentChapter - 1 })
    }
  },

  loadChapterContent: async (chapterIndex) => {
    const { chapters, currentBook, loadingChapters } = get()
    const targetChapter = chapters[chapterIndex]

    // Skip loading if no book, no chapter, already has content, or is a volume page without content.
    if (
      !currentBook ||
      !targetChapter ||
      targetChapter.content ||
      targetChapter.volume_chapter_number === 0 ||
      (targetChapter.volume_number && targetChapter.word_count === 0 && targetChapter.content === '')
    ) {
      return
    }

    if (loadingChapters[chapterIndex]) {
      return
    }

    // Mark chapter as loading to prevent duplicate requests
    set((state) => ({
      loadingChapters: { ...state.loadingChapters, [chapterIndex]: true },
    }))

    try {
      const fullChapter = await bookService.getChapter(currentBook, targetChapter.chapter_number)

      set((state) => {
        const updatedChapters = [...state.chapters]
        updatedChapters[chapterIndex] = fullChapter

        const updatedLoading = { ...state.loadingChapters }
        delete updatedLoading[chapterIndex]

        return {
          chapters: updatedChapters,
          loadingChapters: updatedLoading,
        }
      })
    } catch (error) {
      console.error('Failed to load chapter content', error)

      set((state) => {
        const updatedLoading = { ...state.loadingChapters }
        delete updatedLoading[chapterIndex]
        return { loadingChapters: updatedLoading }
      })
    }
  },
}))
