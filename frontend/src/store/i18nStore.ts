import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Language = 'zh-Hant' | 'zh-Hans'

interface I18nStore {
  language: Language
  setLanguage: (lang: Language) => void
  toggleLanguage: () => void
}

export const useI18nStore = create<I18nStore>()(
  persist(
    (set, get) => ({
      language: 'zh-Hant',
      setLanguage: (language) => set({ language }),
      toggleLanguage: () => {
        const next = get().language === 'zh-Hant' ? 'zh-Hans' : 'zh-Hant'
        set({ language: next })
      },
    }),
    {
      name: 'go-reader-language',
    }
  )
)
