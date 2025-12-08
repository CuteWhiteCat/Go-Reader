import { useEffect, useState } from 'react'
import { Plus, Search, Globe2 } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useBookStore } from '@/store/bookStore'
import { useReaderStore } from '@/store/readerStore'
import BookList from '@/components/library/BookList'
import ReaderView from '@/components/reader/ReaderView'
import AddBookModal from '@/components/library/AddBookModal'
import Button from '@/components/common/Button'
import { bookService } from '@/services/bookService'
import { progressService } from '@/services/progressService'
import { useI18n } from '@/i18n/useI18n'
import type { Book, Chapter, ReadingProgress } from '@/types'

export default function LibraryPage() {
  const navigate = useNavigate()
  const { books, fetchBooks, selectBook, isLoading, error, removeBook } = useBookStore()
  const { startReading, isReading, stopReading, currentBook } = useReaderStore()
  const [searchQuery, setSearchQuery] = useState('')
  const [isAddModalOpen, setIsAddModalOpen] = useState(false)
  const [deletingBookId, setDeletingBookId] = useState<string | null>(null)
  const [lastReadMap, setLastReadMap] = useState<Record<string, string>>({})
  const [progressMap, setProgressMap] = useState<Record<string, ReadingProgress | null>>({})
  const [chapterCountMap, setChapterCountMap] = useState<Record<string, number>>({})
  const [sortKey, setSortKey] = useState<'recent' | 'lastRead' | 'title'>('recent')
  const { t } = useI18n()
  const loadProgress = async (bookList: Book[]) => {
    if (!bookList.length) return
    try {
      const results = await Promise.all(
        bookList.map(async (b) => {
          try {
            const [p, chapters] = await Promise.all([
              progressService.getProgress(b.id),
              bookService.getBookChapters(b.id),
            ])
            return { id: b.id, progress: p, count: chapters.length }
          } catch (err) {
            console.warn('loadProgress failed for book', b.id, err)
            return { id: b.id, progress: null, count: 0 }
          }
        })
      )
      const lastMap: Record<string, string> = {}
      const progMap: Record<string, ReadingProgress | null> = {}
      const countMap: Record<string, number> = {}
      results.forEach(({ id, progress, count }) => {
        if (progress?.last_read_at) lastMap[id] = progress.last_read_at
        progMap[id] = progress
        countMap[id] = count
      })
      setLastReadMap(lastMap)
      setProgressMap(progMap)
      setChapterCountMap(countMap)
    } catch (err) {
      console.error('Failed to load progress map', err)
    }
  }

  useEffect(() => {
    fetchBooks()
  }, [fetchBooks])

  useEffect(() => {
    loadProgress(books)
  }, [books])

  useEffect(() => {
    if (!isReading) {
      loadProgress(books)
    }
  }, [isReading, books])

  const handleBookClick = async (book: Book) => {
    selectBook(book)

    try {
      // Fetch chapter summaries and progress first to avoid loading the entire book
      const [chapterSummaries, progress] = await Promise.all([
        bookService.getBookChapters(book.id),
        progressService.getProgress(book.id),
      ])

      // Fallback for cases where chapter summaries are unavailable
      if (chapterSummaries.length === 0) {
        const chapters = await bookService.getBookContent(book.id)
        startReading(book.id, chapters, progress)
        return
      }

      // Prepare placeholders for all chapters and preload the first batch
      const chapters: Chapter[] = chapterSummaries.map((summary) => ({
        ...summary,
        content: undefined,
      }))

      const INITIAL_CHAPTERS_TO_LOAD = 10
      const chaptersToLoad = new Set<number>()

      // Always load the first chapters
      chapterSummaries
        .slice(0, INITIAL_CHAPTERS_TO_LOAD)
        .forEach((chapter) => {
          if (chapter.volume_chapter_number !== 0) {
            chaptersToLoad.add(chapter.chapter_number)
          }
        })

      // Ensure the user's current chapter has content even if it's beyond the first batch
      const progressChapterNumber = progress.current_chapter + 1
      if (progressChapterNumber > 0 && progressChapterNumber <= chapterSummaries.length) {
        const target = chapterSummaries[progressChapterNumber - 1]
        if (target?.volume_chapter_number !== 0) {
          chaptersToLoad.add(progressChapterNumber)
        }
      }

      const loadedChapters = await Promise.all(
        Array.from(chaptersToLoad).map((chapterNumber) =>
          bookService.getChapter(book.id, chapterNumber)
        )
      )

      // Merge loaded content back into the placeholder list
      loadedChapters.forEach((fullChapter) => {
        const chapterIndex = fullChapter.chapter_number - 1
        if (chapterIndex >= 0 && chapterIndex < chapters.length) {
          chapters[chapterIndex] = fullChapter
        }
      })

      startReading(book.id, chapters, progress)
    } catch (err) {
      console.error('Failed to start reading:', err)
    }
  }

  const filteredBooks = books.filter((book) =>
    book.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    book.author?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const sortedBooks = [...filteredBooks].sort((a, b) => {
    if (sortKey === 'title') {
      return a.title.localeCompare(b.title, 'zh-Hant')
    }
    if (sortKey === 'lastRead') {
      const aTime = lastReadMap[a.id] ? new Date(lastReadMap[a.id]).getTime() : 0
      const bTime = lastReadMap[b.id] ? new Date(lastReadMap[b.id]).getTime() : 0
      if (aTime !== bTime) return bTime - aTime
    }
    // default recent added
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })

  if (isReading) {
    return <ReaderView onClose={stopReading} />
  }

  const handleDeleteBook = async (book: Book) => {
    const confirmed = window.confirm(`確定要刪除「${book.title}」嗎？`)
    if (!confirmed) return

    setDeletingBookId(book.id)
    try {
      await bookService.deleteBook(book.id)
      removeBook(book.id)
      setProgressMap((prev) => {
        const next = { ...prev }
        delete next[book.id]
        return next
      })
      setChapterCountMap((prev) => {
        const next = { ...prev }
        delete next[book.id]
        return next
      })

      if (currentBook === book.id) {
        stopReading()
      }
    } catch (err) {
      console.error('Failed to delete book:', err)
      alert('刪除書籍時發生錯誤，請稍後再試。')
    } finally {
      setDeletingBookId(null)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0d1018] via-[#101523] to-[#0a0c14] text-gray-100">
      {/* Header */}
      <header className="sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-6 py-4 rounded-b-2xl bg-[#161b24] border border-white/10 shadow-lg">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-100">
              {t('library.title')}
            </h1>
            <div className="flex items-center gap-3">
              <Button
                variant="secondary"
                className="bg-gradient-to-r from-[#1f2a45] via-[#1a2237] to-[#121a2c] text-gray-100 border border-white/15 hover:-translate-y-0.5 hover:shadow-lg"
                onClick={() => setIsAddModalOpen(true)}
              >
                <Plus className="w-5 h-5 mr-2" />
                {t('library.addBook')}
              </Button>
              <Button
                variant="secondary"
                className="border border-white/15"
                onClick={() => navigate('/search')}
              >
                <Globe2 className="w-4 h-4 mr-2" />
                線上搜尋
              </Button>
              <select
                value={sortKey}
                onChange={(e) => setSortKey(e.target.value as any)}
                className="bg-[#1d2332] border border-white/10 text-gray-100 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                <option value="recent">依加入時間</option>
                <option value="lastRead">依上次閱讀</option>
                <option value="title">依字典序</option>
              </select>
            </div>
          </div>

          {/* Search Bar */}
          <div className="mt-4 relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder={t('library.searchPlaceholder')}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 rounded-xl bg-[#1d2332] border border-white/10 focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-100 placeholder:text-gray-500 shadow-inner"
            />
          </div>
        </div>
      </header>

      {/* Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <p className="text-gray-400">{t('library.loading')}</p>
          </div>
        ) : error ? (
          <div className="flex items-center justify-center h-64">
            <p className="text-red-400">{t('library.error')}：{error}</p>
          </div>
        ) : (
          <BookList
            books={sortedBooks}
            onBookClick={handleBookClick}
            onDeleteBook={handleDeleteBook}
            deletingBookId={deletingBookId}
            progressMap={progressMap}
            chapterCountMap={chapterCountMap}
          />
        )}
      </main>

      {/* Add Book Modal */}
      <AddBookModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onSuccess={() => {
          fetchBooks()
          setIsAddModalOpen(false)
        }}
      />
    </div>
  )
}
