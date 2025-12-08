# Go-Reader

A modern, cross-platform desktop novel reader application built with Go, React, and Electron.

## Features

- 📚 **Multi-format Support**: Read `.txt`, `.md`, and `.epub` files
- 🎨 **Modern UI**: Liquid design with blur effects and dark mode
- 📖 **Reading Progress**: Track your reading progress across all books
- 🏷️ **Tag System**: Organize your library with custom tags
- 🔖 **Bookmarks**: Save your favorite passages
- 🌐 **Book Sources**: Download books from external sources
- ⚙️ **Customizable**: Adjust font size, line spacing, and page margins

## Tech Stack

- **Backend**: Go 1.21+ with Chi router and SQLite
- **Frontend**: React 18+ with TypeScript and TailwindCSS
- **Desktop**: Electron 28+
- **State Management**: Zustand
- **Build Tool**: Vite

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Go-Reader
```

2. Run the development setup script:
```bash
./scripts/dev.sh
```

Or manually install dependencies:

```bash
# Install backend dependencies
cd backend
go mod download

# Install frontend dependencies
cd ../frontend
npm install

# Install Electron dependencies
cd ..
npm install
```

## Development

Start the development environment:

```bash
npm run dev
```

This will:
1. Start the Go backend server on `http://localhost:8080`
2. Start the React development server on `http://localhost:5173`
3. Launch the Electron application

### Running Individual Components

Backend only:
```bash
cd backend
go run cmd/server/main.go
```

Frontend only:
```bash
cd frontend
npm run dev
```

Electron only (requires backend and frontend running):
```bash
npm run electron:dev
```

## Building for Production

Run the build script:

```bash
./scripts/build.sh
```

Or manually build:

```bash
# Build backend
cd backend
go build -o ../bin/go-reader-server cmd/server/main.go

# Build frontend
cd ../frontend
npm run build

# Build Electron app
cd ..
npm run electron:build
```

The packaged application will be in the `dist` directory.

## Project Structure

```
Go-Reader/
├── backend/                 # Go backend
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── api/           # HTTP handlers and routes
│   │   ├── service/       # Business logic
│   │   ├── repository/    # Data access layer
│   │   ├── models/        # Data structures
│   │   ├── parser/        # File parsers
│   │   └── config/        # Configuration
│   └── migrations/        # Database migrations
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── pages/        # Page components
│   │   ├── services/     # API services
│   │   ├── store/        # State management
│   │   └── types/        # TypeScript types
│   └── public/           # Static assets
├── electron/             # Electron main process
│   ├── main.js          # Main process
│   └── preload.js       # Preload script
└── scripts/             # Build scripts
```

## API Endpoints

### Books
- `GET /api/books` - List all books
- `GET /api/books/:id` - Get book details
- `POST /api/books` - Add a new book
- `PUT /api/books/:id` - Update book metadata
- `DELETE /api/books/:id` - Delete a book
- `GET /api/books/:id/content` - Get book content
- `GET /api/books/:id/chapters` - Get chapter list

### Tags
- `GET /api/tags` - List all tags
- `POST /api/tags` - Create a tag
- `PUT /api/tags/:id` - Update a tag
- `DELETE /api/tags/:id` - Delete a tag

### Progress
- `GET /api/progress/:bookId` - Get reading progress
- `PUT /api/progress/:bookId` - Update reading progress

### Bookmarks
- `GET /api/bookmarks/:bookId` - Get bookmarks
- `POST /api/bookmarks` - Create a bookmark
- `DELETE /api/bookmarks/:id` - Delete a bookmark

### Sources
- `GET /api/sources` - List book sources
- `POST /api/sources` - Add a book source
- `PUT /api/sources/:id` - Update a source
- `DELETE /api/sources/:id` - Delete a source
- `POST /api/sources/:id/search` - Search books
- `POST /api/sources/:id/download` - Download a book

## Configuration

The application can be configured through environment variables or a configuration file:

```yaml
server:
  port: 8080
  host: localhost

database:
  path: ./data/database.db

storage:
  books_dir: ./data/books
  covers_dir: ./data/covers
```

Environment variables use the prefix `GOREADER_`, e.g., `GOREADER_SERVER_PORT=8080`.

## Contributing

Contributions are welcome! Please read the contributing guidelines before submitting a pull request.

## License

This project is licensed under the MIT License.

## Acknowledgments

Built with ❤️ using Go, React, and Electron
