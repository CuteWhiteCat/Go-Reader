# Go-Reader: Desktop Novel Reader Application

## Overview
A cross-platform desktop novel reader application built with:
- **Backend**: Go (REST API server running locally)
- **Frontend**: React (chosen for better Electron integration and ecosystem)
- **Desktop Packaging**: Electron (wraps the React app and manages Go backend process)

The application will run a local Go HTTP server that the React frontend communicates with via REST API.

## Core Features

### 1. Book Reading & Management
- Support multiple text formats: `.txt`, `.md`, `.epub`
- Local library with bookshelf view
- Book metadata management (title, author, cover, description)
- Reading progress tracking
- Bookmarks and annotations
- Tag-based organization system

### 2. Book Source Integration (書源)
- Parse and manage external book sources (web scrapers)
- Auto-download books from configured sources
- Chapter-by-chapter synchronization
- Source management interface

### 3. Modern UI/UX
- **Liquid Design**: Blur effects, transparency, frosted glass aesthetics
- **Dark Mode**: Full dark theme support
- **Responsive Reading Experience**: Adjustable font size, line spacing, margins
- **Smooth Animations**: Page transitions and UI interactions

## Technical Requirements

### Backend (Go 1.21+)
- Local HTTP REST API server
- SQLite database for metadata storage
- File system management for book storage
- Web scraping for book sources

### Frontend (React 18+)
- Modern React with Hooks
- TypeScript for type safety
- TailwindCSS for styling
- Component-based architecture

### Desktop (Electron)
- Package React app
- Manage Go backend lifecycle
- Handle system tray and native menus

## Implementation Notes
1. This project is focused on improving Go skills, so the backend architecture should follow Go best practices including clean architecture, proper error handling, and idiomatic Go code patterns.
2. For every function we made, we also need to write unit tests to ensure code quality and reliability. And we can also implement Swagger API documentation for the REST API.

---

## Project Structure

```
Go-Reader/
├── backend/                      # Go backend application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Application entry point
│   ├── internal/                # Private application code
│   │   ├── api/                 # HTTP handlers and routes
│   │   │   ├── handlers/        # Request handlers
│   │   │   ├── middleware/      # HTTP middleware
│   │   │   └── router.go        # Route definitions
│   │   ├── service/             # Business logic layer
│   │   │   ├── book.go
│   │   │   ├── library.go
│   │   │   ├── source.go
│   │   │   └── reader.go
│   │   ├── repository/          # Data access layer
│   │   │   ├── book.go
│   │   │   ├── tag.go
│   │   │   └── progress.go
│   │   ├── models/              # Data structures and entities
│   │   │   ├── book.go
│   │   │   ├── chapter.go
│   │   │   ├── tag.go
│   │   │   └── source.go
│   │   ├── parser/              # File format parsers
│   │   │   ├── txt.go
│   │   │   ├── markdown.go
│   │   │   └── epub.go
│   │   ├── scraper/             # Book source scrapers
│   │   │   ├── scraper.go
│   │   │   └── rules.go
│   │   └── config/              # Configuration management
│   │       └── config.go
│   ├── pkg/                     # Public libraries (reusable)
│   │   └── utils/
│   ├── migrations/              # Database migration files
│   │   └── 001_initial.sql
│   ├── go.mod
│   └── go.sum
│
├── frontend/                    # React frontend application
│   ├── src/
│   │   ├── components/          # React components
│   │   │   ├── common/          # Reusable components
│   │   │   ├── reader/          # Reading view components
│   │   │   ├── library/         # Library/bookshelf components
│   │   │   └── settings/        # Settings components
│   │   ├── pages/               # Page-level components
│   │   │   ├── LibraryPage.tsx
│   │   │   ├── ReaderPage.tsx
│   │   │   └── SettingsPage.tsx
│   │   ├── hooks/               # Custom React hooks
│   │   │   ├── useBooks.ts
│   │   │   └── useReader.ts
│   │   ├── services/            # API client services
│   │   │   ├── api.ts           # Axios instance
│   │   │   ├── bookService.ts
│   │   │   └── sourceService.ts
│   │   ├── store/               # State management
│   │   │   ├── bookStore.ts
│   │   │   └── settingsStore.ts
│   │   ├── styles/              # Global styles
│   │   │   ├── globals.css
│   │   │   └── theme.css
│   │   ├── types/               # TypeScript type definitions
│   │   │   └── index.ts
│   │   ├── utils/               # Utility functions
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── public/
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   └── vite.config.ts
│
├── electron/                    # Electron main process
│   ├── main.js                  # Electron entry point
│   ├── preload.js              # Preload script for IPC
│   └── backend-manager.js      # Go backend process manager
│
├── scripts/                     # Build and utility scripts
│   ├── build.sh
│   └── dev.sh
│
├── data/                        # Application data (created at runtime)
│   ├── books/                   # Stored book files
│   ├── covers/                  # Book cover images
│   └── database.db              # SQLite database
│
├── CLAUDE.md                    # This file
├── README.md
├── .gitignore
└── package.json                 # Root package.json for Electron
```

## Required Packages & Tools

### Go Backend Dependencies

```go
// go.mod dependencies
require (
    // HTTP Router - lightweight and fast
    github.com/go-chi/chi/v5 v5.0.10

    // CORS middleware
    github.com/go-chi/cors v1.2.1

    // Database
    github.com/mattn/go-sqlite3 v1.14.18
    github.com/jmoiron/sqlx v1.3.5

    // EPUB parser
    github.com/taylorskalyo/goreader v0.0.0-20230626095242-3b9c48f0e8e4

    // Web scraping
    github.com/PuerkitoBio/goquery v1.8.1
    github.com/gocolly/colly/v2 v2.1.0

    // Markdown parser (for .md files)
    github.com/gomarkdown/markdown v0.0.0-20231115200524-a660076da3fd

    // Configuration
    github.com/spf13/viper v1.18.2

    // Logging
    github.com/sirupsen/logrus v1.9.3

    // UUID generation
    github.com/google/uuid v1.5.0

    // Validation
    github.com/go-playground/validator/v10 v10.16.0
)
```

### Frontend Dependencies

```json
{
  "dependencies": {
    // Core
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.21.0",

    // State Management
    "zustand": "^4.4.7",

    // HTTP Client
    "axios": "^1.6.2",

    // UI & Styling
    "tailwindcss": "^3.4.0",
    "framer-motion": "^10.16.16",
    "lucide-react": "^0.303.0",

    // Utilities
    "clsx": "^2.0.0",
    "date-fns": "^3.0.6"
  },
  "devDependencies": {
    // TypeScript
    "@types/react": "^18.2.45",
    "@types/react-dom": "^18.2.18",
    "typescript": "^5.3.3",

    // Build Tools (Vite)
    "vite": "^5.0.8",
    "@vitejs/plugin-react": "^4.2.1",

    // Linting
    "eslint": "^8.56.0",
    "@typescript-eslint/eslint-plugin": "^6.15.0",
    "@typescript-eslint/parser": "^6.15.0",

    // PostCSS (for Tailwind)
    "autoprefixer": "^10.4.16",
    "postcss": "^8.4.32"
  }
}
```

### Electron Dependencies

```json
{
  "dependencies": {
    "electron": "^28.1.0",
    "electron-builder": "^24.9.1"
  }
}
```

## Development Tools

### Required Software
- **Go** 1.21 or higher
- **Node.js** 18 or higher
- **npm** or **yarn**
- **Git**
- **SQLite3** (for database inspection)

### Recommended VS Code Extensions
- Go (golang.go)
- ESLint
- Prettier
- Tailwind CSS IntelliSense
- SQLite Viewer

### Build Tools
- **Air** (Go hot reload): `go install github.com/cosmtrek/air@latest`
- **Vite** (Frontend build tool): Included in frontend dependencies
- **Electron Builder** (Package desktop app): Included in Electron dependencies

## API Endpoints (Planned)

### Books
- `GET /api/books` - List all books
- `GET /api/books/:id` - Get book details
- `POST /api/books` - Add book manually
- `PUT /api/books/:id` - Update book metadata
- `DELETE /api/books/:id` - Delete book
- `GET /api/books/:id/content` - Get book content
- `GET /api/books/:id/chapters` - Get book chapters

### Reading Progress
- `GET /api/progress/:bookId` - Get reading progress
- `PUT /api/progress/:bookId` - Update reading progress

### Tags
- `GET /api/tags` - List all tags
- `POST /api/tags` - Create tag
- `PUT /api/books/:id/tags` - Update book tags

### Sources (書源)
- `GET /api/sources` - List all book sources
- `POST /api/sources` - Add new source
- `PUT /api/sources/:id` - Update source
- `DELETE /api/sources/:id` - Delete source
- `POST /api/sources/:id/download` - Download book from source

### Settings
- `GET /api/settings` - Get application settings
- `PUT /api/settings` - Update settings
