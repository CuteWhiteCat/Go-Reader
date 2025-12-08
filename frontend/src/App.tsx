import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import { useSettingsStore } from './store/settingsStore'
import LibraryPage from './pages/LibraryPage'
import SettingsPage from './pages/SettingsPage'
import OnlineSearchPage from './pages/OnlineSearchPage'

function App() {
  const { theme } = useSettingsStore()

  useEffect(() => {
    // Apply theme on mount
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme])

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/library" replace />} />
        <Route path="/library" element={<LibraryPage />} />
        <Route path="/search" element={<OnlineSearchPage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Routes>
    </Router>
  )
}

export default App
