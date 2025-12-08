import { useState } from 'react'
import { X, Upload, Loader2, FileText } from 'lucide-react'
import Button from '../common/Button'
import { bookService } from '@/services/bookService'
import type { CreateBookRequest } from '@/types'
import { useI18n } from '@/i18n/useI18n'

interface AddBookModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export default function AddBookModal({ isOpen, onClose, onSuccess }: AddBookModalProps) {
  const [formData, setFormData] = useState<CreateBookRequest>({
    title: '',
    author: '',
    description: '',
    file_path: '',
    file_format: 'txt',
    tag_ids: [],
  })
  const [selectedFileName, setSelectedFileName] = useState<string>('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { t } = useI18n()

  const detectFormat = (fileName: string): 'txt' | 'md' | 'epub' => {
    const ext = fileName.split('.').pop()?.toLowerCase()
    if (ext === 'md' || ext === 'markdown') return 'md'
    if (ext === 'epub') return 'epub'
    return 'txt'
  }

  if (!isOpen) return null

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!formData.title || !formData.file_path) {
      setError(t('addBook.error.required'))
      return
    }

    setIsLoading(true)

    try {
      await bookService.createBook(formData)
      onSuccess()
      onClose()
      // Reset form
      setFormData({
        title: '',
        author: '',
        description: '',
        file_path: '',
        file_format: 'txt',
        tag_ids: [],
      })
      setSelectedFileName('')
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-md">
      <div className="rounded-2xl w-full max-w-md mx-4 p-6 bg-[#121622] border border-white/10 shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-100">
            {t('addBook.title')}
          </h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Title */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('addBook.field.title')}
            </label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-100"
              required
            />
          </div>

          {/* Author */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('addBook.field.author')}
            </label>
            <input
              type="text"
              value={formData.author}
              onChange={(e) => setFormData({ ...formData, author: e.target.value })}
              className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-100"
            />
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('addBook.field.description')}
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows={3}
              className="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-100 resize-none"
            />
          </div>

          {/* File Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('addBook.field.bookFile')}
            </label>
            <div className="flex items-center gap-3 mb-2 min-w-0">
            <button
              type="button"
              onClick={async () => {
                  console.log('Choose File clicked')
                  console.log('window.electron:', window.electron)

                  if (window.electron?.selectFile) {
                    try {
                      const filePath = await window.electron.selectFile()
                      console.log('Selected file:', filePath)

                      if (filePath) {
                        // Extract just the filename
                        const fileName = filePath.split('/').pop() || filePath.split('\\').pop() || filePath
                        setSelectedFileName(fileName)
                        setFormData((prev) => ({
                          ...prev,
                          file_path: filePath,
                          file_format: detectFormat(fileName),
                          title: prev.title || fileName.replace(/\.[^.]+$/, ''),
                        }))
                      }
                    } catch (err) {
                  console.error('Error selecting file:', err)
                  setError(t('addBook.error.selectFileFailed'))
                }
              } else {
                console.error('window.electron.selectFile not available')
                setError(t('addBook.error.fileUnavailable'))
              }
                }}
            className="px-4 py-2 rounded-lg bg-[#2a3552] hover:bg-[#1f2a45] border border-white/10 text-gray-100 transition-colors flex items-center gap-2 shrink-0 shadow-lg"
          >
            <Upload className="w-4 h-4" />
            {t('addBook.chooseFile')}
          </button>

              {selectedFileName && (
                <div className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 min-w-0 flex-1">
                  <FileText className="w-4 h-4" />
                  <span className="truncate">{selectedFileName}</span>
                </div>
              )}
            </div>
          </div>

          {/* Error Message */}
          {error && (
            <div className="p-3 rounded-lg bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 text-sm">
              {error}
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-3 pt-2">
            <Button
              type="button"
              variant="secondary"
              onClick={onClose}
              className="flex-1 justify-center"
              disabled={isLoading}
            >
              {t('actions.cancel')}
            </Button>
            <Button
              type="submit"
              variant="secondary"
              className="flex-1 justify-center bg-gradient-to-r from-[#1f2a45] via-[#1a2237] to-[#121a2c] text-gray-100 border border-white/15 hover:shadow-lg hover:-translate-y-0.5"
              disabled={isLoading}
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  {t('actions.adding')}
                </>
              ) : (
                <>
                  <Upload className="w-4 h-4 mr-2" />
                  {t('actions.add')}
                </>
              )}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
