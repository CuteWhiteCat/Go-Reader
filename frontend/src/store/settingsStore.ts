import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface SettingsStore {
  theme: 'light' | 'dark'
  fontSize: number
  lineSpacing: number
  pageMargin: number
  contentLanguage: 'zh-Hant' | 'zh-Hans'
  readingTheme: 'day' | 'night' | 'sepia'
  fontFamily: 'serif' | 'sans'

  // Actions
  setTheme: (theme: 'light' | 'dark') => void
  toggleTheme: () => void
  setFontSize: (size: number) => void
  setLineSpacing: (spacing: number) => void
  setPageMargin: (margin: number) => void
  setContentLanguage: (lang: 'zh-Hant' | 'zh-Hans') => void
  toggleContentLanguage: () => void
  setReadingTheme: (theme: 'day' | 'night' | 'sepia') => void
  setFontFamily: (font: 'serif' | 'sans') => void
}

export const useSettingsStore = create<SettingsStore>()(
  persist(
    (set) => ({
      theme: 'dark',
      fontSize: 17,
      lineSpacing: 1.7,
      pageMargin: 20,
      contentLanguage: 'zh-Hant',
      readingTheme: 'night',
      fontFamily: 'serif',

      setTheme: (theme) => {
        set({ theme })
        document.documentElement.classList.toggle('dark', theme === 'dark')
      },

      toggleTheme: () => {
        set((state) => {
          const newTheme = state.theme === 'light' ? 'dark' : 'light'
          document.documentElement.classList.toggle('dark', newTheme === 'dark')
          return { theme: newTheme }
        })
      },

      setFontSize: (fontSize) => set({ fontSize }),
      setLineSpacing: (lineSpacing) => set({ lineSpacing }),
      setPageMargin: (pageMargin) => set({ pageMargin }),
      setContentLanguage: (contentLanguage) => set({ contentLanguage }),
      toggleContentLanguage: () =>
        set((state) => ({
          contentLanguage: state.contentLanguage === 'zh-Hant' ? 'zh-Hans' : 'zh-Hant',
        })),
      setReadingTheme: (readingTheme) => set({ readingTheme }),
      setFontFamily: (fontFamily) => set({ fontFamily }),
    }),
    {
      name: 'go-reader-settings',
    }
  )
)
