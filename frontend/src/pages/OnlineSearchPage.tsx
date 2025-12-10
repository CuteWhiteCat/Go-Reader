import { useState } from 'react'
import { Search, Download, Loader2, ChevronLeft, ChevronRight, ArrowLeft } from 'lucide-react'
import { crawlerService, type CrawlSearchResult } from '@/services/crawlerService'
import { useNavigate } from 'react-router-dom'
import { convertText } from '@/utils/textConversion'

const ITEMS_PER_PAGE = 10

export default function OnlineSearchPage() {
  const navigate = useNavigate()
  const [keyword, setKeyword] = useState('')
  const [results, setResults] = useState<CrawlSearchResult[]>([])
  const [isSearching, setIsSearching] = useState(false)
  const [currentPage, setCurrentPage] = useState(1)
  const [downloadingItems, setDownloadingItems] = useState<Set<string>>(new Set())
  const [downloadProgress, setDownloadProgress] = useState<Map<string, number>>(new Map())
  const [searchError, setSearchError] = useState('')
  const [hasSearched, setHasSearched] = useState(false)


  const totalPages = Math.ceil(results.length / ITEMS_PER_PAGE)
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
  const endIndex = startIndex + ITEMS_PER_PAGE
  const currentResults = results.slice(startIndex, endIndex)

  const handleSearch = async () => {
    if (!keyword.trim()) return

    setIsSearching(true)
    setCurrentPage(1)
    setSearchError('')
    setHasSearched(true)
    try {
      // 將搜尋詞轉為簡體以符合網站語言
      const query = convertText(keyword.trim(), 'zh-Hans')
      const data = await crawlerService.search(query)
      setResults(data || [])
      // 搜索成功但没有结果 - 这是正常情况
      if (!data || data.length === 0) {
        setSearchError('')
      } else {
        setSearchError('')
      }
    } catch (error) {
      console.error('Search failed:', error)
      const msg = error instanceof Error ? error.message : '未知錯誤'
      const normalizedMsg = msg.toLowerCase()
      const rateLimitSignals = [
        'rate_limit',
        '搜索次数已耗尽',
        '搜索过于频繁',
        '搜索次數已耗盡',
        '搜索過於頻繁',
        '提供10次搜索机会',
        '提供10次搜索機會',
        '一分钟只提供10次搜索机会',
        '一分鐘只提供10次搜索機會',
        '為防止惡意搜索',
        '为防止恶意搜索',
        '為防止惡意搜尋',
        '为防止恶意搜尋',
        '429',
      ]

      // 检查是否是频率限制错误
      const isRateLimited =
        rateLimitSignals.some((signal) => msg.includes(signal)) ||
        normalizedMsg.includes('too many')

      if (isRateLimited) {
        setSearchError('⚠️ 搜尋過於頻繁，已超過網站限制，請稍後再試。')
      } else {
        setSearchError(`搜尋失敗：${msg}`)
      }
    } finally {
      setIsSearching(false)
    }
  }

  const handleDownload = async (novel: CrawlSearchResult) => {
    const itemKey = novel.url
    if (downloadingItems.has(itemKey)) return

    setDownloadingItems(prev => new Set(prev).add(itemKey))
    setDownloadProgress(prev => new Map(prev).set(itemKey, 0))

    try {
      const jobId = await crawlerService.startImport({
        title: novel.title,
        author: novel.author,
        latest: novel.latest,
        url: novel.url,
      })


      const intervalId = window.setInterval(async () => {
        const currentJobId = jobId
        if (!currentJobId) return
        try {
          const status = await crawlerService.getImportStatus(currentJobId)
          const percent = status.total > 0 ? Math.round((status.done / status.total) * 100) : 0
          setDownloadProgress(prev => {
            const next = new Map(prev)
            next.set(itemKey, percent)
            return next
          })

          if (status.status === 'success') {
            clearInterval(intervalId)
            setDownloadingItems(prev => {
              const next = new Set(prev)
              next.delete(itemKey)
              return next
            })
            setDownloadProgress(prev => {
              const next = new Map(prev)
              next.set(itemKey, 100)
              return next
            })
            setTimeout(() => {
              setDownloadProgress(prev => {
                const next = new Map(prev)
                next.delete(itemKey)
                return next
              })
            }, 1000)
            alert(`《${novel.title}》下載成功！`)
          } else if (status.status === 'error') {
            clearInterval(intervalId)
            setDownloadingItems(prev => {
              const next = new Set(prev)
              next.delete(itemKey)
              return next
            })
            setDownloadProgress(prev => {
              const next = new Map(prev)
              next.delete(itemKey)
              return next
            })
            alert(`下載失敗：${status.error || '未知錯誤'}`)
          }
        } catch (error) {
          console.error('Poll failed', error)
        }
      }, 800)
    } catch (error) {
      console.error('Download failed:', error)
      setDownloadingItems(prev => {
        const next = new Set(prev)
        next.delete(itemKey)
        return next
      })
      setDownloadProgress(prev => {
        const next = new Map(prev)
        next.delete(itemKey)
        return next
      })
      alert(`下載失敗：${error instanceof Error ? error.message : '未知錯誤'}`)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch()
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0d1018] via-[#101523] to-[#0a0c14] text-gray-100">
      <header className="sticky top-0 z-10 bg-[#161b24]/90 border-b border-white/10 backdrop-blur">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center gap-4">
          <button
            onClick={() => navigate('/library')}
            className="p-2 rounded-lg hover:bg-white/5 transition-colors border border-white/10"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <h1 className="text-2xl font-bold text-gray-100">線上搜尋</h1>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-6 py-10">
        <div className="bg-[#161b24]/90 border border-white/10 rounded-2xl shadow-xl p-6 mb-8 backdrop-blur">
          {searchError && (
            <div className="mb-4 rounded-lg border border-amber-400/50 bg-amber-500/10 text-amber-100 px-4 py-3 text-sm">
              {searchError}
            </div>
          )}
          <div className="flex gap-4">
            <div className="flex-1 relative">
              <input
                type="text"
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="輸入書名或作者名稱..."
                className="w-full px-4 py-3 pl-12 bg-[#1d2332] border border-white/10 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#5f7fff] text-gray-100 placeholder:text-gray-500"
              />
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
            </div>
            <button
              onClick={handleSearch}
              disabled={isSearching || !keyword.trim()}
              className="px-8 py-3 bg-gradient-to-r from-[#4f7bfd] via-[#5c6ffb] to-[#6f63f9] hover:brightness-110 disabled:opacity-50 text-white rounded-xl transition-all font-medium disabled:cursor-not-allowed flex items-center gap-2 shadow-lg shadow-blue-900/40"
            >
              {isSearching ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  搜尋中...
                </>
              ) : (
                '搜尋'
              )}
            </button>
          </div>
        </div>

        {results.length > 0 && (
          <div className="space-y-4">
            {currentResults.map((novel) => {
              const isDownloading = downloadingItems.has(novel.url)
              const progress = downloadProgress.get(novel.url) || 0

              return (
                <div
                  key={novel.url}
                  className="bg-[#161b24]/90 border border-white/10 backdrop-blur rounded-xl shadow-md p-6 hover:-translate-y-0.5 hover:shadow-2xl transition-all"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <h3 className="text-xl font-bold text-gray-100 mb-2">
                        {novel.title}
                      </h3>
                      <div className="space-y-1 text-sm text-gray-400">
                        {novel.author && (
                          <p>
                            <span className="font-medium text-gray-300">作者：</span>
                            {novel.author}
                          </p>
                        )}
                        {novel.latest && (
                          <p>
                            <span className="font-medium text-gray-300">最新：</span>
                            {novel.latest}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="flex flex-col items-end gap-2">
                      <button
                        onClick={() => handleDownload(novel)}
                        disabled={isDownloading}
                        className="px-6 py-2 bg-gradient-to-r from-[#34d399] to-[#10b981] hover:brightness-110 disabled:opacity-50 text-white rounded-lg transition-all font-medium disabled:cursor-not-allowed flex items-center gap-2 min-w-[120px] justify-center shadow-lg shadow-emerald-900/40"
                      >
                        {isDownloading ? (
                          <>
                            <Loader2 className="w-4 h-4 animate-spin" />
                            下載中
                          </>
                        ) : (
                          <>
                            <Download className="w-4 h-4" />
                            下載
                          </>
                        )}
                      </button>

                      {isDownloading && (
                        <div className="w-full">
                          <div className="h-2 bg-white/5 rounded-full overflow-hidden">
                            <div
                              className="h-full bg-gradient-to-r from-[#34d399] to-[#10b981] transition-all duration-300"
                              style={{ width: `${progress}%` }}
                            />
                          </div>
                          <p className="text-xs text-gray-400 mt-1 text-right">
                            {progress}%
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              )
            })}

            {totalPages > 1 && (
              <div className="flex flex-col items-center gap-4 mt-8">
                <div className="flex items-center justify-center gap-2 bg-[#161b24]/90 border border-white/10 backdrop-blur rounded-xl shadow-md p-4">
                  <button
                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                    className="px-4 py-2 bg-gradient-to-r from-[#4f7bfd] to-[#6f63f9] hover:brightness-110 disabled:opacity-50 text-white rounded-lg disabled:cursor-not-allowed transition-all flex items-center gap-2 font-medium"
                  >
                    <ChevronLeft className="w-5 h-5" />
                    上一頁
                  </button>

                  <div className="flex items-center gap-2">
                    {currentPage > 3 && (
                      <>
                        <button
                          onClick={() => setCurrentPage(1)}
                          className="w-10 h-10 rounded-lg bg-[#1d2332] hover:bg-[#4f7bfd] hover:text-white transition-colors border border-white/10"
                        >
                          1
                        </button>
                        {currentPage > 4 && <span className="px-2">...</span>}
                      </>
                    )}

                    {Array.from({ length: totalPages }, (_, i) => i + 1)
                      .filter(page => {
                        return page === currentPage ||
                               page === currentPage - 1 ||
                               page === currentPage + 1 ||
                               (page < 4 && currentPage < 4) ||
                               (page > totalPages - 3 && currentPage > totalPages - 3)
                      })
                      .map(page => (
                        <button
                          key={page}
                          onClick={() => setCurrentPage(page)}
                          className={`w-10 h-10 rounded-lg font-medium transition-colors ${
                            page === currentPage
                              ? 'bg-gradient-to-r from-[#4f7bfd] to-[#6f63f9] text-white border border-transparent'
                              : 'bg-[#1d2332] border border-white/10 hover:bg-[#4f7bfd] hover:text-white'
                          }`}
                        >
                          {page}
                        </button>
                      ))}

                    {currentPage < totalPages - 2 && (
                      <>
                        {currentPage < totalPages - 3 && <span className="px-2">...</span>}
                        <button
                          onClick={() => setCurrentPage(totalPages)}
                          className="w-10 h-10 rounded-lg bg-[#1d2332] border border-white/10 hover:bg-[#4f7bfd] hover:text-white transition-colors"
                        >
                          {totalPages}
                        </button>
                      </>
                    )}
                  </div>

                  <button
                    onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                    className="px-4 py-2 bg-gradient-to-r from-[#4f7bfd] to-[#6f63f9] hover:brightness-110 disabled:opacity-50 text-white rounded-lg disabled:cursor-not-allowed transition-all flex items-center gap-2 font-medium"
                  >
                    下一頁
                    <ChevronRight className="w-5 h-5" />
                  </button>
                </div>

                <div className="text-sm text-gray-400">
                  顯示第 {startIndex + 1} - {Math.min(endIndex, results.length)} 筆，共 {results.length} 筆結果
                </div>
              </div>
            )}
          </div>
        )}

        {!isSearching && results.length === 0 && (
          <div className="text-center py-16 text-gray-400 space-y-2">
            <Search className="w-16 h-16 text-gray-600 mx-auto mb-4" />
            {searchError ? (
              <p className="text-amber-100">{searchError}</p>
            ) : hasSearched ? (
              <p className="text-gray-300">沒有找到相關結果，請嘗試其他關鍵字。</p>
            ) : (
              <p className="text-gray-400">輸入關鍵字搜尋線上小說</p>
            )}
          </div>
        )}
      </main>
    </div>
  )
}
