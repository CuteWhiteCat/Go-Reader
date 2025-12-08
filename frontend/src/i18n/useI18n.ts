import { useMemo } from 'react'
import { useI18nStore, type Language } from '@/store/i18nStore'

type TranslationKey =
  | 'library.title'
  | 'library.addBook'
  | 'library.searchPlaceholder'
  | 'library.noBooks'
  | 'library.noBooksDesc'
  | 'library.loading'
  | 'library.error'
  | 'library.chapterProgress'
  | 'actions.cancel'
  | 'actions.add'
  | 'actions.adding'
  | 'actions.previous'
  | 'actions.next'
  | 'reader.loading'
  | 'reader.preparing'
  | 'reader.close'
  | 'reader.volumeChapterLabel'
  | 'reader.chapterOnlyLabel'
  | 'addBook.title'
  | 'addBook.field.title'
  | 'addBook.field.author'
  | 'addBook.field.description'
  | 'addBook.field.bookFile'
  | 'addBook.field.format'
  | 'addBook.chooseFile'
  | 'addBook.manualPathPlaceholder'
  | 'addBook.fileHint'
  | 'addBook.error.required'
  | 'addBook.error.selectFileFailed'
  | 'addBook.error.fileUnavailable'
  | 'bookCard.unknownAuthor'
  | 'bookCard.deleting'

type Vars = Record<string, string | number>

const translations: Record<Language, Record<TranslationKey, string>> = {
  'zh-Hant': {
    'library.title': '我的書庫',
    'library.addBook': '新增書籍',
    'library.searchPlaceholder': '搜尋書籍...',
    'library.noBooks': '沒有找到書籍',
    'library.noBooksDesc': '新增你的第一本書開始閱讀吧',
    'library.loading': '載入書籍中...',
    'library.error': '載入時發生錯誤',
    'library.chapterProgress': '第{{current}}章 / 共{{total}}章',
    'actions.cancel': '取消',
    'actions.add': '新增',
    'actions.adding': '新增中...',
    'actions.previous': '上一章',
    'actions.next': '下一章',
    'reader.loading': '章節載入中...',
    'reader.preparing': '正在準備此章節...',
    'reader.close': '關閉',
    'reader.volumeChapterLabel': '第{{volume}}卷 · 第{{chapter}}章',
    'reader.chapterOnlyLabel': '第{{chapter}}章',
    'addBook.title': '新增書籍',
    'addBook.field.title': '書名 *',
    'addBook.field.author': '作者',
    'addBook.field.description': '簡介',
    'addBook.field.bookFile': '書籍檔案 *',
    'addBook.field.format': '格式 *',
    'addBook.chooseFile': '選擇檔案',
    'addBook.manualPathPlaceholder': '或貼上檔案路徑',
    'addBook.fileHint': '點擊「選擇檔案」或貼上完整路徑（支援 .txt、.md、.epub）',
    'addBook.error.required': '書名與檔案路徑為必填',
    'addBook.error.selectFileFailed': '選擇檔案失敗',
    'addBook.error.fileUnavailable': '目前環境無法選擇檔案',
    'bookCard.unknownAuthor': '未知作者',
    'bookCard.deleting': '刪除中...',
  },
  'zh-Hans': {
    'library.title': '我的书库',
    'library.addBook': '新增书籍',
    'library.searchPlaceholder': '搜索书籍...',
    'library.noBooks': '没有找到书籍',
    'library.noBooksDesc': '新增你的第一本书开始阅读吧',
    'library.loading': '加载书籍中...',
    'library.error': '加载时发生错误',
    'library.chapterProgress': '第{{current}}章 / 共{{total}}章',
    'actions.cancel': '取消',
    'actions.add': '新增',
    'actions.adding': '新增中...',
    'actions.previous': '上一章',
    'actions.next': '下一章',
    'reader.loading': '章节加载中...',
    'reader.preparing': '正在准备此章节...',
    'reader.close': '关闭',
    'reader.volumeChapterLabel': '第{{volume}}卷 · 第{{chapter}}章',
    'reader.chapterOnlyLabel': '第{{chapter}}章',
    'addBook.title': '新增书籍',
    'addBook.field.title': '书名 *',
    'addBook.field.author': '作者',
    'addBook.field.description': '简介',
    'addBook.field.bookFile': '书籍文件 *',
    'addBook.field.format': '格式 *',
    'addBook.chooseFile': '选择文件',
    'addBook.manualPathPlaceholder': '或粘贴文件路径',
    'addBook.fileHint': '点击「选择文件」或粘贴完整路径（支持 .txt、.md、.epub）',
    'addBook.error.required': '书名与文件路径为必填',
    'addBook.error.selectFileFailed': '选择文件失败',
    'addBook.error.fileUnavailable': '当前环境无法选择文件',
    'bookCard.unknownAuthor': '未知作者',
    'bookCard.deleting': '删除中...',
  },
}

function interpolate(template: string, vars?: Vars) {
  if (!vars) return template
  return template.replace(/{{(.*?)}}/g, (_, key) => String(vars[key.trim()] ?? ''))
}

export function useI18n() {
  const language = useI18nStore((state) => state.language)
  const setLanguage = useI18nStore((state) => state.setLanguage)
  const toggleLanguage = useI18nStore((state) => state.toggleLanguage)

  const t = useMemo(
    () =>
      (key: TranslationKey, vars?: Vars) => {
        const dict = translations[language] ?? translations['zh-Hant']
        const template = dict[key] ?? key
        return interpolate(template, vars)
      },
    [language]
  )

  return { t, language, setLanguage, toggleLanguage }
}
