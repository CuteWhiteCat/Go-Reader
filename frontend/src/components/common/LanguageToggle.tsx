import { useI18n } from '@/i18n/useI18n'

export default function LanguageToggle() {
  const { language, toggleLanguage } = useI18n()
  const label = language === 'zh-Hant' ? '繁' : '简'
  const nextLabel = language === 'zh-Hant' ? '简' : '繁'

  return (
    <button
      type="button"
      onClick={toggleLanguage}
      className="px-3 py-2 rounded-lg glass text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-200/60 dark:hover:bg-gray-800/60 transition"
      aria-label="Switch language"
      title={`切換到${nextLabel}`}
    >
      {label}
    </button>
  )
}
