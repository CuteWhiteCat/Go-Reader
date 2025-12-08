import { Moon, Sun, Type, AlignLeft, Maximize2 } from 'lucide-react'
import { useSettingsStore } from '@/store/settingsStore'
import Button from '@/components/common/Button'

export default function SettingsPage() {
  const {
    theme,
    fontSize,
    lineSpacing,
    pageMargin,
    toggleTheme,
    setFontSize,
    setLineSpacing,
    setPageMargin,
  } = useSettingsStore()

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0d1018] via-[#101523] to-[#0a0c14] text-gray-100">
      <div className="max-w-4xl mx-auto px-6 py-8">
        <h1 className="text-3xl font-bold mb-8 text-gray-100">
          Settings
        </h1>

        <div className="space-y-6">
          {/* Theme */}
          <div className="rounded-xl p-6 bg-[#161b24]/90 border border-white/10 shadow-lg">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                {theme === 'dark' ? (
                  <Moon className="w-6 h-6 text-primary-600" />
                ) : (
                  <Sun className="w-6 h-6 text-primary-600" />
                )}
                <div>
                  <h3 className="font-semibold text-gray-100">
                    Theme
                  </h3>
                  <p className="text-sm text-gray-400">
                    Current: {theme === 'dark' ? 'Dark' : 'Light'}
                  </p>
                </div>
              </div>
              <Button onClick={toggleTheme}>Toggle Theme</Button>
            </div>
          </div>

          {/* Font Size */}
          <div className="rounded-xl p-6 bg-[#161b24]/90 border border-white/10 shadow-lg">
            <div className="flex items-center gap-3 mb-4">
              <Type className="w-6 h-6 text-primary-600" />
              <div>
                <h3 className="font-semibold text-gray-100">
                  Font Size
                </h3>
                <p className="text-sm text-gray-400">
                  Current: {fontSize}px
                </p>
              </div>
            </div>
            <input
              type="range"
              min="12"
              max="24"
              value={fontSize}
              onChange={(e) => setFontSize(Number(e.target.value))}
              className="w-full"
            />
          </div>

          {/* Line Spacing */}
          <div className="rounded-xl p-6 bg-[#161b24]/90 border border-white/10 shadow-lg">
            <div className="flex items-center gap-3 mb-4">
              <AlignLeft className="w-6 h-6 text-primary-600" />
              <div>
                <h3 className="font-semibold text-gray-100">
                  Line Spacing
                </h3>
                <p className="text-sm text-gray-400">
                  Current: {lineSpacing}
                </p>
              </div>
            </div>
            <input
              type="range"
              min="1"
              max="2.5"
              step="0.1"
              value={lineSpacing}
              onChange={(e) => setLineSpacing(Number(e.target.value))}
              className="w-full"
            />
          </div>

          {/* Page Margin */}
          <div className="rounded-xl p-6 bg-[#161b24]/90 border border-white/10 shadow-lg">
            <div className="flex items-center gap-3 mb-4">
              <Maximize2 className="w-6 h-6 text-primary-600" />
              <div>
                <h3 className="font-semibold text-gray-100">
                  Page Margin
                </h3>
                <p className="text-sm text-gray-400">
                  Current: {pageMargin}px
                </p>
              </div>
            </div>
            <input
              type="range"
              min="0"
              max="100"
              value={pageMargin}
              onChange={(e) => setPageMargin(Number(e.target.value))}
              className="w-full"
            />
          </div>
        </div>
      </div>
    </div>
  )
}
