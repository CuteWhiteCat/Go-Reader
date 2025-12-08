import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './styles/globals.css'

function preloadFonts() {
  if (typeof document === 'undefined' || !('fonts' in document)) return
  const fontFaces = [
    '1rem "Noto Sans TC"',
    '1rem "Noto Sans SC"',
    '1rem "Noto Serif TC"',
    '1rem "Noto Serif SC"',
  ]
  fontFaces.forEach((face) => {
    // Fire-and-forget; ensures fonts are requested early
    ;(document as any).fonts.load(face).catch(() => {})
  })
}

preloadFonts()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
