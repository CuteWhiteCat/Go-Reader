import { useEffect, useMemo, useState, useLayoutEffect, useRef, useCallback } from 'react'
import { ChevronLeft, ChevronRight, X, List } from 'lucide-react'
import { useReaderStore } from '@/store/readerStore'
import { useSettingsStore } from '@/store/settingsStore'
import Button from '../common/Button'
import { useI18n } from '@/i18n/useI18n'
import { convertText } from '@/utils/textConversion'
import { Sun, Moon, Settings as SettingsIcon, Plus, Minus } from 'lucide-react'
import { progressService } from '@/services/progressService'

interface ReaderViewProps {
  onClose: () => void
}

export default function ReaderView({ onClose }: ReaderViewProps) {
  const {
    chapters,
    currentChapter,
    nextChapter,
    previousChapter,
    loadChapterContent,
    loadingChapters,
    setCurrentChapter,
    updateProgress,
    progress,
  } = useReaderStore()

  const {
    fontSize,
    lineSpacing,
    pageMargin,
    contentLanguage,
    toggleContentLanguage,
    readingTheme,
    setReadingTheme,
    fontFamily,
    setFontFamily,
    setFontSize,
  } = useSettingsStore()
  const { t } = useI18n()
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const [isSettingsOpen, setIsSettingsOpen] = useState(false)

  const chapter = chapters[currentChapter]
  const isChapterLoading = loadingChapters[currentChapter]
  const volumeLabel =
    chapter?.volume_number != null ? t('reader.volumeChapterLabel', {
      volume: chapter.volume_number,
      chapter: chapter?.volume_chapter_number ?? chapter?.chapter_number ?? currentChapter + 1,
    }) : null
  const chapterLabel =
    chapter?.volume_chapter_number != null || chapter?.chapter_number != null
      ? t('reader.chapterOnlyLabel', {
          chapter: chapter?.volume_chapter_number ?? chapter?.chapter_number ?? currentChapter + 1,
        })
      : t('reader.chapterOnlyLabel', { chapter: currentChapter + 1 })

  const convertedTitle = useMemo(
    () => (chapter ? convertText(chapter.title, contentLanguage) : ''),
    [chapter, contentLanguage]
  )

  const convertedLabel = useMemo(
    () => convertText(volumeLabel ?? chapterLabel, contentLanguage),
    [volumeLabel, chapterLabel, contentLanguage]
  )

  const convertedParagraphs = useMemo(() => {
    if (!chapter?.content) return []
    const converted = convertText(chapter.content, contentLanguage)
    return converted.split('\n')
  }, [chapter?.content, contentLanguage])

  const isVolumePage =
    (chapter?.volume_chapter_number ?? -1) === 0 ||
    (!!chapter?.volume_number && chapter?.word_count === 0 && (chapter?.content ?? '') === '')

  // Scroll management for progress and chapter transitions
  const scrollContainerRef = useRef<HTMLDivElement>(null)
  const saveTimerRef = useRef<number | null>(null)
  const restoredChaptersRef = useRef<Set<number>>(new Set())

  // Ensure the current chapter content is loaded when needed
  useEffect(() => {
    if (chapter && !chapter.content && !isVolumePage) {
      loadChapterContent(currentChapter)
    }
  }, [chapter, currentChapter, loadChapterContent, isVolumePage])

  const volumeTitles = useMemo(() => {
    const map = new Map<number, string>()
    chapters.forEach((ch) => {
      if ((ch.volume_chapter_number ?? -1) === 0) {
        map.set(ch.volume_number ?? 1, ch.title)
      }
    })
    return map
  }, [chapters])

  const volumeGroups = useMemo(() => {
    const groups = new Map<number, { volume: number; chapters: { index: number; title: string; number: number }[] }>()
    chapters.forEach((ch, idx) => {
      if ((ch.volume_chapter_number ?? -1) === 0) {
        return
      }
      const vol = ch.volume_number ?? 1
      if (!groups.has(vol)) {
        groups.set(vol, { volume: vol, chapters: [] })
      }
      groups.get(vol)?.chapters.push({
        index: idx,
        title: ch.title,
        number: ch.volume_chapter_number ?? ch.chapter_number,
      })
    })

    return Array.from(groups.values()).sort((a, b) => a.volume - b.volume)
  }, [chapters])

  // Content chapters only (跳過卷首頁)
  const contentChapters = useMemo(
    () =>
      chapters.reduce<
        { chapter: typeof chapters[number]; originalIndex: number; number: number }[]
      >((arr, ch, idx) => {
        if ((ch.volume_chapter_number ?? -1) === 0) return arr
        const number = ch.volume_chapter_number ?? ch.chapter_number ?? arr.length + 1
        arr.push({ chapter: ch, originalIndex: idx, number })
        return arr
      }, []),
    [chapters]
  )

  const lastReadIndex = useMemo(() => {
    if (!progress) return -1
    return contentChapters.findIndex((item) => item.originalIndex === progress.current_chapter)
  }, [contentChapters, progress])

  const readCount = lastReadIndex >= 0 ? lastReadIndex + 1 : 0
  const totalContent = contentChapters.length
  const readPercent = totalContent > 0 ? Math.round((readCount / totalContent) * 100) : 0

  const handleJump = (index: number) => {
    setCurrentChapter(index)
    setIsSidebarOpen(false)
  }

  const clearSaveTimer = () => {
    if (saveTimerRef.current) {
      window.clearTimeout(saveTimerRef.current)
      saveTimerRef.current = null
    }
  }

  const computeScrollProgress = useCallback(() => {
    const el = scrollContainerRef.current
    if (!el) return { position: 0, percent: 0 }
    const scrollTop = el.scrollTop
    const total = el.scrollHeight || 1
    const viewport = el.clientHeight || 1
    const percent = Math.min(100, Math.max(0, ((scrollTop + viewport) / total) * 100))
    return { position: Math.floor(scrollTop), percent }
  }, [])

  const scheduleSaveProgress = useCallback(() => {
    clearSaveTimer()
    // Debounce saves to avoid flooding backend while scrolling
    saveTimerRef.current = window.setTimeout(async () => {
      if (!chapter || isVolumePage) return
      const currentBook = useReaderStore.getState().currentBook
      if (!currentBook) return

      const { position, percent } = computeScrollProgress()

      try {
        const currentIndexInContent = contentChapters.findIndex((item) => item.chapter.id === chapter.id)
        const displayCurrent = currentIndexInContent >= 0 ? currentIndexInContent : currentChapter
        const displayTotal = contentChapters.length || chapters.length || 1
        const fallbackPercent = Math.min(100, Math.max(0, ((displayCurrent + 1) / displayTotal) * 100))
        const progressPercent = Number.isFinite(percent) && percent > 0 ? percent : fallbackPercent

        const res = await progressService.updateProgress(currentBook, {
          current_chapter: currentChapter,
          current_position: position,
          progress_percentage: progressPercent,
        })
        updateProgress(res)
      } catch (err) {
        console.error('Failed to save progress', err)
      }
    }, 400)
  }, [chapter, chapters, computeScrollProgress, currentChapter, isVolumePage, updateProgress])

  const themeStyles = {
    day: {
      container: 'bg-[#f6f4ef] text-gray-900',
      panel: 'bg-white/90 text-gray-900',
      glass: 'glass',
    },
    night: {
      container: 'bg-[#10131a] text-gray-200',
      panel: 'bg-[#161b24]/90 text-gray-100',
      glass: 'glass dark',
    },
    sepia: {
      container: 'bg-[#f2e8d5] text-[#3e2f1c]',
      panel: 'bg-[#f7efdf]/90 text-[#3e2f1c]',
      glass: 'glass',
    },
  } as const

  const themeClass = themeStyles[readingTheme] || themeStyles.day
  const isNight = readingTheme === 'night'
  const textMain = isNight ? 'text-gray-200' : 'text-gray-900'
  const textMuted = isNight ? 'text-gray-300' : 'text-gray-500'
  const [isContentReady, setIsContentReady] = useState(true)

  useLayoutEffect(() => {
    setIsContentReady(false)
    const id = requestAnimationFrame(() => setIsContentReady(true))
    return () => cancelAnimationFrame(id)
  }, [contentLanguage, chapter?.id])

  useEffect(() => {
    // Handle keyboard navigation
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') previousChapter()
      if (e.key === 'ArrowRight') nextChapter()
      if (e.key === 'Escape') onClose()
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [nextChapter, previousChapter, onClose])

  // Persist reading progress when chapter changes (skip volume pages)
  useEffect(() => {
    if (!chapter || !chapter.id || isVolumePage) return
    scheduleSaveProgress()
  }, [chapter, currentChapter, isVolumePage, scheduleSaveProgress])

  // On chapter change, jump to top
  useEffect(() => {
    const el = scrollContainerRef.current
    if (el) {
      el.scrollTo({ top: 0 })
    }
    restoredChaptersRef.current.delete(currentChapter)
    scheduleSaveProgress()
  }, [currentChapter, scheduleSaveProgress])

  // Restore saved scroll position when opening the book/chapter
  useEffect(() => {
    if (!progress) return
    if (progress.current_chapter !== currentChapter) return
    if (progress.current_position <= 0) return
    if (isVolumePage) return
    if (restoredChaptersRef.current.has(currentChapter)) return

    const el = scrollContainerRef.current
    if (!el) return

    const restore = () => {
      const maxScroll = Math.max(0, el.scrollHeight - el.clientHeight)
      const target = Math.min(progress.current_position, maxScroll)
      el.scrollTo({ top: target })
      restoredChaptersRef.current.add(currentChapter)
    }

    // Wait for layout/content ready before restoring position
    const id = window.requestAnimationFrame(restore)
    return () => window.cancelAnimationFrame(id)
  }, [currentChapter, progress, isVolumePage, convertedParagraphs.length])

  useEffect(() => () => clearSaveTimer(), [])

  if (!chapter) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p>Chapter not found</p>
      </div>
    )
  }

  return (
    <div className={`fixed inset-0 z-50 overflow-hidden ${themeClass.container}`}>
      {/* Sidebar */}
      <div
        className={`fixed top-0 bottom-0 left-0 w-72 ${themeClass.panel} backdrop-blur-lg border-r ${isNight ? 'border-[#1d2332]' : 'border-gray-200'} z-40 transition-transform duration-300 overflow-y-auto ${
          isSidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="px-4 py-5 border-b border-gray-200 dark:border-gray-800">
          <div className={`flex items-center gap-2 ${textMain}`}>
            <List className="w-4 h-4" />
            <span className="font-semibold">目錄</span>
          </div>
        </div>
        <div className="p-4 space-y-4">
          {volumeGroups.map((vol) => {
            const volumeTitle = convertText(
              volumeTitles.get(vol.volume) ?? `第${vol.volume}卷`,
              contentLanguage
            )
            return (
              <div key={vol.volume}>
                <p className={`text-sm font-semibold ${textMain} mb-2`}>
                  {volumeTitle}
                </p>
                <div className="space-y-1">
                  {vol.chapters.map((ch) => {
                    const label = convertText(ch.title, contentLanguage)
                    const isActive = ch.index === currentChapter
                    const readChapterIdx = progress?.current_chapter ?? -1
                    const isRead =
                      ch.index < readChapterIdx ||
                      (ch.index === readChapterIdx && (progress?.current_position ?? 0) > 0)

                    return (
                      <button
                        key={ch.index}
                        onClick={() => handleJump(ch.index)}
                        className={`w-full text-left px-3 py-2 rounded-lg transition border ${
                          isActive
                            ? 'bg-primary-100/20 text-primary-100 border-primary-400/40'
                            : isRead
                              ? 'bg-white/5 text-gray-400 border-white/5'
                              : 'hover:bg-white/10 text-gray-200 border-transparent'
                        }`}
                      >
                        <span className={`text-xs mr-2 ${textMuted}`}>
                          第{ch.number}章
                        </span>
                        <span className="text-sm truncate">{label}</span>
                      </button>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>
      </div>

      <div
        className={`h-full transition-transform duration-300 ${isSidebarOpen ? 'translate-x-72' : ''}`}
      >
        {/* Header */}
        <div className={`absolute top-0 left-0 right-0 h-16 ${themeClass.panel} backdrop-blur-md flex items-center justify-between px-6 z-10 border-b ${isNight ? 'border-[#1d2332]' : 'border-gray-200'}`}>
          <div className="flex items-center gap-3">
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setIsSidebarOpen((v) => !v)}
            >
              <List className="w-4 h-4" />
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={onClose}
            >
              <ChevronLeft className="w-4 h-4" />
              返回
            </Button>
            <div className="flex flex-col">
              {!isVolumePage && (
                <>
                  <h2 className={`text-lg font-semibold ${textMain}`}>
                    {convertedTitle || chapter.title}
                  </h2>
                  {chapter && (
                    <p className={`text-xs ${textMuted}`}>
                      {convertedLabel}
                    </p>
                  )}
                </>
              )}
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="secondary"
              size="sm"
              onClick={toggleContentLanguage}
              className="w-10 justify-center"
            >
              {contentLanguage === 'zh-Hant' ? '繁' : '簡'}
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setIsSettingsOpen((v) => !v)}
              className="w-10 justify-center"
            >
              <SettingsIcon className="w-4 h-4" />
            </Button>
            <Button variant="ghost" size="sm" onClick={onClose}>
              <X className="w-5 h-5" />
            </Button>
          </div>
        </div>

      {/* Settings Panel */}
      <div
        className={`fixed top-16 right-0 w-72 ${themeClass.panel} backdrop-blur-lg border-l border-gray-200 dark:border-gray-800 z-40 transition-transform duration-300 ${
          isSettingsOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        <div className="p-4 space-y-4">
          <div>
            <p className="text-sm font-semibold mb-2">閱讀主題</p>
            <div className="flex gap-2">
              {([
                { key: 'day', label: '日間', icon: <Sun className="w-4 h-4" /> },
                { key: 'night', label: '夜間', icon: <Moon className="w-4 h-4" /> },
                { key: 'sepia', label: '米黃', icon: <Sun className="w-4 h-4" /> },
              ] as const).map((item) => (
                <Button
                  key={item.key}
                  variant={readingTheme === item.key ? 'primary' : 'secondary'}
                  size="sm"
                  onClick={() => setReadingTheme(item.key)}
                >
                  {item.icon}
                  {item.label}
                </Button>
              ))}
            </div>
          </div>

          <div>
            <p className="text-sm font-semibold mb-2">字體</p>
            <div className="flex gap-2">
              <Button
                variant={fontFamily === 'serif' ? 'primary' : 'secondary'}
                size="sm"
                onClick={() => setFontFamily('serif')}
              >
                明體
              </Button>
              <Button
                variant={fontFamily === 'sans' ? 'primary' : 'secondary'}
                size="sm"
                onClick={() => setFontFamily('sans')}
              >
                黑體
              </Button>
            </div>
          </div>

          <div>
            <p className="text-sm font-semibold mb-2">字體大小</p>
            <div className="flex items-center gap-2">
              <Button variant="secondary" size="sm" onClick={() => setFontSize(Math.max(12, fontSize - 1))}>
                <Minus className="w-4 h-4" />
              </Button>
              <span className="text-sm w-10 text-center">{fontSize}px</span>
              <Button variant="secondary" size="sm" onClick={() => setFontSize(Math.min(40, fontSize + 1))}>
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div
        className="h-full pt-16 pb-20 overflow-y-auto"
        ref={scrollContainerRef}
        onScroll={scheduleSaveProgress}
      >
        <div
          className="max-w-4xl mx-auto py-8"
          style={{
            paddingLeft: `${pageMargin}px`,
            paddingRight: `${pageMargin}px`,
            fontFamily: fontFamily === 'serif'
              ? '"Noto Serif TC", "Noto Serif SC", "Songti TC", "PMingLiU", serif'
              : '"Noto Sans TC", "Noto Sans SC", "PingFang TC", "Microsoft JhengHei", sans-serif',
          }}
        >
          {isVolumePage ? (
            <div className="min-h-[70vh] flex items-center justify-center text-center">
              <div className="space-y-3">
                <p className={`text-sm uppercase tracking-wide ${textMuted}`}>
                  {chapter?.volume_number ? t('reader.volumeChapterLabel', {
                    volume: chapter.volume_number,
                    chapter: '',
                  }).replace(/·.*$/, '') : ''}
                </p>
                <h1 className={`text-5xl font-bold ${textMain}`}>
                  {convertedTitle || chapter.title}
                </h1>
              </div>
            </div>
          ) : (
            <>
              <h1 className={`text-3xl font-bold mb-6 transition-opacity duration-75 ${isContentReady ? 'opacity-100' : 'opacity-0'} ${textMain}`}>
                {convertedTitle || chapter.title}
              </h1>
              <div
                className={`prose dark:prose-invert max-w-none transition-opacity duration-75 ${isContentReady ? 'opacity-100' : 'opacity-0'}`}
                style={{
                  fontSize: `${fontSize}px`,
                  lineHeight: lineSpacing,
                }}
              >
                {isChapterLoading && (
                  <p className={textMuted}>{t('reader.loading')}</p>
                )}

                {!isChapterLoading && chapter.content &&
                  convertedParagraphs.map((paragraph, index) => (
                    <p key={index} className={`mb-4 ${textMain}`}>
                      {paragraph}
                    </p>
                  ))}

                {!isChapterLoading && !chapter.content && (
                  <p className={textMuted}>
                    {t('reader.preparing')}
                  </p>
                )}
              </div>
            </>
          )}
        </div>
      </div>

      {/* Navigation */}
      <div className={`absolute bottom-0 left-0 right-0 h-20 ${themeClass.panel} backdrop-blur-md flex items-center justify-between px-6 border-t ${isNight ? 'border-[#1d2332]' : 'border-gray-200'}`}>
        <Button
          variant="secondary"
          onClick={previousChapter}
          disabled={currentChapter === 0}
        >
          <ChevronLeft className="w-5 h-5 mr-2" />
          {t('actions.previous')}
        </Button>

        <div className={`text-sm ${textMuted}`}>
          {totalContent > 0
            ? `已讀：${readCount} / ${totalContent} 章（${readPercent}%）`
            : t('library.chapterProgress', { current: 0, total: 0 })}
        </div>

        <Button
          variant="secondary"
          onClick={nextChapter}
          disabled={currentChapter === chapters.length - 1}
        >
          {t('actions.next')}
          <ChevronRight className="w-5 h-5 ml-2" />
        </Button>
      </div>
    </div>
    </div>
  )
}
