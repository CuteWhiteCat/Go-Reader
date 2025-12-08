import { useEffect, useState } from 'react'
import { Book, ReadingProgress } from '@/types'
import { BookOpen, Calendar, Trash2 } from 'lucide-react'
import { format } from 'date-fns'
import { useI18n } from '@/i18n/useI18n'

interface BookCardProps {
  book: Book
  onClick?: () => void
  onDelete?: () => void
  isDeleting?: boolean
  progress?: ReadingProgress | null
  totalChapters?: number
}

export default function BookCard({
  book,
  onClick,
  onDelete,
  isDeleting = false,
  progress,
  totalChapters = 0,
}: BookCardProps) {
  const { t } = useI18n()

  const fallbackCover =
    'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="240" height="320" viewBox="0 0 240 320"><defs><linearGradient id="g" x1="0" x2="1" y1="0" y2="1"><stop stop-color="%231f2a45" offset="0"/><stop stop-color="%230d1018" offset="1"/></linearGradient></defs><rect width="240" height="320" fill="url(%23g)"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="%23ccd4e0" font-family="sans-serif" font-size="20">No Cover</text></svg>'

  const resolveCover = () => {
    if (!book.cover_path) return fallbackCover
    if (book.cover_path.startsWith('http')) return book.cover_path
    if (book.cover_path.startsWith('/')) return book.cover_path
    return `/${book.cover_path}`
  }

  const [coverSrc, setCoverSrc] = useState(resolveCover())

  useEffect(() => {
    setCoverSrc(resolveCover())
  }, [book.cover_path])

  const readChapterIdx = progress?.current_chapter ?? -1
  const readCount = readChapterIdx >= 0 ? readChapterIdx + 1 : 0
  const percent =
    totalChapters > 0
      ? Math.min(100, Math.round((readCount / totalChapters) * 100))
      : Math.round(progress?.progress_percentage ?? 0)

  return (
    <div
      onClick={onClick}
      className="group cursor-pointer rounded-xl overflow-hidden bg-[#161b24] border border-white/10 hover:border-primary-500/40 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 relative"
    >
      {onDelete && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation()
            if (!isDeleting) {
              onDelete()
            }
          }}
          className="absolute top-3 right-3 p-2 rounded-full bg-white/80 text-gray-700 hover:bg-red-50 hover:text-red-600 shadow-sm transition"
          aria-label="Delete book"
          disabled={isDeleting}
        >
          <Trash2 className="w-4 h-4" />
        </button>
      )}

      {isDeleting && (
        <div className="absolute inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center text-sm text-gray-200 z-10">
          {t('bookCard.deleting')}
        </div>
      )}

      {/* Cover Image */}
      <div className="h-48 bg-gradient-to-br from-[#1f2937] to-[#0f172a] flex items-center justify-center overflow-hidden">
        {coverSrc ? (
          <img
            src={coverSrc}
            alt={book.title}
            className="w-full h-full object-cover"
            onError={() => setCoverSrc(fallbackCover)}
          />
        ) : (
          <BookOpen className="w-16 h-16 text-white opacity-80" />
        )}
      </div>

      {/* Book Info */}
      <div className="p-4">
        <h3 className="font-semibold text-lg mb-1 text-gray-100 line-clamp-2">
          {book.title}
        </h3>
        <p className="text-sm text-gray-400 mb-2">{book.author || t('bookCard.unknownAuthor')}</p>

        {/* Tags */}
        {book.tags && book.tags.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-2">
            {book.tags.slice(0, 3).map((tag) => (
              <span
                key={tag.id}
                className="px-2 py-0.5 text-xs rounded-full"
                style={{
                  backgroundColor: tag.color || '#e5e7eb',
                  color: '#000',
                }}
              >
                {tag.name}
              </span>
            ))}
          </div>
        )}

        {/* Meta Info */}
        <div className="flex items-center gap-2 text-xs text-gray-500">
          <Calendar className="w-3 h-3" />
          <span>{format(new Date(book.updated_at), 'MMM d, yyyy')}</span>
        </div>

        {/* Progress */}
        <div className="mt-3 space-y-2">
          <div className="flex items-center justify-between text-xs text-gray-400">
            <span>閱讀進度</span>
            {totalChapters > 0 ? (
              <span className="text-gray-300">
                {readCount} / {totalChapters} 章
              </span>
            ) : (
              <span className="text-gray-300">{percent}%</span>
            )}
          </div>
          <div className="h-2 rounded-full bg-white/5 overflow-hidden">
            <div
              className="h-full rounded-full bg-gradient-to-r from-[#4f7bfd] to-[#6f63f9] transition-all duration-300"
              style={{ width: `${Math.min(100, Math.max(0, percent))}%` }}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
