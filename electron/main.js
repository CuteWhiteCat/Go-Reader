const { app, BrowserWindow, Menu, dialog, ipcMain } = require('electron')
const path = require('path')
const { spawn } = require('child_process')
const net = require('net')

let mainWindow
let backendProcess

// Backend server manager
class BackendManager {
  constructor() {
    this.process = null
    this.port = Number(process.env.GOREADER_SERVER_PORT) || 8080
  }

  async isPortInUse() {
    return new Promise((resolve) => {
      const socket = net.createConnection({ port: this.port, host: '127.0.0.1' })
      socket.once('connect', () => {
        socket.destroy()
        resolve(true)
      })
      socket.once('error', () => {
        socket.destroy()
        resolve(false)
      })
    })
  }

  async isPortAvailable() {
    return !(await this.isPortInUse())
  }

  async waitForExistingBackend(timeoutMs = 5000, intervalMs = 200) {
    const start = Date.now()

    while (Date.now() - start < timeoutMs) {
      if (await this.isPortInUse()) {
        return true
      }

      await new Promise((resolve) => setTimeout(resolve, intervalMs))
    }

    return false
  }

  async start() {
    const isDev = !app.isPackaged

    // In dev, wait briefly in case the separate dev server is starting on this port
    if (isDev) {
      const existingBackendDetected = await this.waitForExistingBackend()
      if (existingBackendDetected) {
        console.log(
          `Backend already running on port ${this.port}, skipping spawn to avoid conflicts.`,
        )
        return
      }
    }

    if (!(await this.isPortAvailable())) {
      console.log(
        `Backend already running on port ${this.port}, skipping spawn to avoid conflicts.`,
      )
      return
    }

    const backendPath = isDev
      ? path.join(__dirname, '../backend/cmd/server/main.go')
      : path.join(process.resourcesPath, 'bin', 'go-reader-server')

    console.log('Starting backend server...')

    if (isDev) {
      // Development mode: run with go run
      this.process = spawn('go', ['run', backendPath], {
        cwd: path.join(__dirname, '../backend'),
        env: {
          ...process.env,
          GOREADER_SERVER_PORT: this.port.toString(),
        },
      })
    } else {
      // Production mode: run compiled binary
      const userDataPath = app.getPath('userData')
      const envForBackend = {
        ...process.env,
        GOREADER_SERVER_PORT: this.port.toString(),
        GOREADER_DATABASE_PATH: path.join(userDataPath, 'database.db'),
        GOREADER_STORAGE_BOOKS_DIR: path.join(userDataPath, 'books'),
        GOREADER_STORAGE_COVERS_DIR: path.join(userDataPath, 'covers'),
      }

      console.log('--- GO-READER DEBUG INFO ---')
      console.log('Running in production mode:', app.isPackaged)
      console.log('Backend executable path:', backendPath)
      console.log('Backend working directory (cwd):', process.resourcesPath)
      console.log('User data path for backend:', userDataPath)
      console.log('Final ENV for backend:', JSON.stringify(envForBackend, null, 2))
      console.log('--- END GO-READER DEBUG INFO ---')

      this.process = spawn(backendPath, [], {
        cwd: process.resourcesPath,
        env: envForBackend,
      })
    }

    this.process.stdout.on('data', (data) => {
      console.log(`Backend: ${data}`)
    })

    this.process.stderr.on('data', (data) => {
      console.error(`Backend Error: ${data}`)
    })

    this.process.on('close', (code) => {
      console.log(`Backend process exited with code ${code}`)
    })
  }

  stop() {
    if (this.process) {
      console.log('Stopping backend server...')
      this.process.kill()
      this.process = null
    }
  }
}

const backendManager = new BackendManager()

// IPC Handlers
ipcMain.handle('select-file', async (event, options) => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openFile'],
    filters: [
      { name: 'Book Files', extensions: ['txt', 'md', 'epub'] },
      { name: 'Text Files', extensions: ['txt'] },
      { name: 'Markdown Files', extensions: ['md'] },
      { name: 'EPUB Files', extensions: ['epub'] },
      { name: 'All Files', extensions: ['*'] },
    ],
    ...options,
  })

  if (result.canceled) {
    return null
  }

  return result.filePaths[0]
})

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
    },
    titleBarStyle: 'hiddenInset',
    trafficLightPosition: { x: 15, y: 15 },
  })

  // Create application menu
  const template = [
    {
      label: 'File',
      submenu: [
        {
          label: 'Add Book',
          accelerator: 'CmdOrCtrl+O',
          click: () => {
            mainWindow.webContents.send('add-book')
          },
        },
        { type: 'separator' },
        { role: 'quit' },
      ],
    },
    {
      label: 'Edit',
      submenu: [
        { role: 'undo' },
        { role: 'redo' },
        { type: 'separator' },
        { role: 'cut' },
        { role: 'copy' },
        { role: 'paste' },
      ],
    },
    {
      label: 'View',
      submenu: [
        { role: 'reload' },
        { role: 'forceReload' },
        { role: 'toggleDevTools' },
        { type: 'separator' },
        { role: 'resetZoom' },
        { role: 'zoomIn' },
        { role: 'zoomOut' },
        { type: 'separator' },
        { role: 'togglefullscreen' },
      ],
    },
    {
      label: 'Window',
      submenu: [
        { role: 'minimize' },
        { role: 'zoom' },
        { type: 'separator' },
        { role: 'front' },
      ],
    },
  ]

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)

  // Load the app
  const isDev = !app.isPackaged
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(path.join(__dirname, '../frontend/dist/index.html'))
    mainWindow.webContents.openDevTools()
  }

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

// Enable GPU hardware acceleration explicitly (default is on, but ensure not disabled elsewhere)
app.commandLine.appendSwitch('ignore-gpu-blacklist')
app.commandLine.appendSwitch('enable-gpu-rasterization')
app.commandLine.appendSwitch('enable-zero-copy')

// App lifecycle
app.whenReady().then(async () => {
  // Start backend server (skip if something is already listening)
  await backendManager.start()

  // Wait a bit for backend to start, then create window
  setTimeout(() => {
    createWindow()
  }, 2000)

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    backendManager.stop()
    app.quit()
  }
})

app.on('before-quit', () => {
  backendManager.stop()
})

app.on('will-quit', () => {
  backendManager.stop()
})
