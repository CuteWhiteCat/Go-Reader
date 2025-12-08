import { Book, ReadingProgress } from '@/types'
import BookCard from '../common/BookCard'
import { useI18n } from '@/i18n/useI18n'

interface BookListProps {
  books: Book[]
  onBookClick: (book: Book) => void
  onDeleteBook?: (book: Book) => void
  deletingBookId?: string | null
  progressMap?: Record<string, ReadingProgress | null>
  chapterCountMap?: Record<string, number>
}

export default function BookList({
  books,
  onBookClick,
  onDeleteBook,
  deletingBookId,
  progressMap = {},
  chapterCountMap = {},
}: BookListProps) {
  const { t } = useI18n()

  if (books.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-gray-500">
        <p className="text-lg">{t('library.noBooks')}</p>
        <p className="text-sm mt-2">{t('library.noBooksDesc')}</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
      {books.map((book) => (
        <BookCard
          key={book.id}
          book={book}
          onClick={() => onBookClick(book)}
          onDelete={onDeleteBook ? () => onDeleteBook(book) : undefined}
          isDeleting={deletingBookId === book.id}
          progress={progressMap[book.id] || null}
          totalChapters={chapterCountMap[book.id] || 0}
        />
      ))}
    </div>
  )
}
