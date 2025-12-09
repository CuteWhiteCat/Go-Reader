# Go-Reader

Go-Reader is a cross-platform desktop e-book reader built with Go, React, and Electron. It's designed for an excellent reading experience with both local files and online content.

## ✨ Features

- **Local Book Management**: Supports various formats like EPUB, TXT, and Markdown.
- **Online Content Crawler**: Fetch and read novels from supported online sources.
- **Cross-Platform**: Works on Windows, macOS, and Linux.
- **Modern UI**: Clean and intuitive user interface built with React and Tailwind CSS.
- **Reading Progression**: Keeps track of your reading progress.

## 🛠️ Tech Stack

- **Backend**: Go
- **Frontend**: React, TypeScript, Vite, Tailwind CSS
- **Desktop Framework**: Electron
- **UI Components**: Shadcn/ui (implied by usage)

## 📂 Project Structure

```
.
├── backend/         # Go backend server
├── electron/        # Electron main process and preload scripts
├── frontend/        # React frontend application
├── scripts/         # Development and build scripts
└── ...
```

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.20 or newer)
- [Node.js](https://nodejs.org/) (version 18 or newer)
- `make` (optional, for convenience scripts)

### Installation & Development

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/Go-Reader.git
    cd Go-Reader
    ```

2.  **Install dependencies:**
    This project uses `npm` for package management. Install dependencies for the root, frontend, and backend.
    ```bash
    npm install
    cd frontend
    npm install
    cd ../backend
    go mod tidy
    cd ..
    ```

3.  **Run the development server:**
    This command starts the Go backend, the React frontend, and the Electron app concurrently.
    ```bash
    npm run dev
    ```
    The app will launch in a new window. Hot-reloading is enabled for the frontend.

### Building for Production

To build the application for your current platform:

1.  **Build the frontend and backend:**
    ```bash
    npm run build
    ```
    This command creates a production build of the React app in `frontend/dist` and compiles the Go backend into the `bin` directory.

2.  **Package the Electron app:**
    ```bash
    npm run electron:build
    ```
    The packaged application will be available in the `dist` directory.

## 📜 Available Scripts

- `npm run dev`: Starts the application in development mode.
- `npm run build`: Builds the frontend and backend for production.
- `npm run electron:dev`: Starts the Electron app using pre-built resources (requires running `build` first).
- `npm run electron:build`: Packages the application for distribution.
- `npm run stop`: Forcefully stops all related services (backend, frontend, Electron).
