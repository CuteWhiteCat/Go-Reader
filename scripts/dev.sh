#!/bin/bash

# Go-Reader Development Script
# This script starts the development environment

set -e

echo "🚀 Starting Go-Reader Development Environment..."

# Add Go to PATH if it exists but not in PATH
if [ -d "/usr/local/go/bin" ] && [[ ":$PATH:" != *":/usr/local/go/bin:"* ]]; then
    export PATH=$PATH:/usr/local/go/bin
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo ""
    echo "Installation instructions:"
    echo "  wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz"
    echo "  sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz"
    echo "  echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc"
    echo "  source ~/.bashrc"
    echo ""
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or higher."
    exit 1
fi

# Install backend dependencies
echo "📦 Installing backend dependencies..."
cd backend
if [ ! -f "go.sum" ]; then
    go mod download
fi
cd ..

# Install frontend dependencies
echo "📦 Installing frontend dependencies..."
cd frontend
if [ ! -d "node_modules" ]; then
    npm install
fi
cd ..

# Install root dependencies (for Electron)
echo "📦 Installing Electron dependencies..."
if [ ! -d "node_modules" ]; then
    npm install
fi

# Start development servers
echo "✅ Setup complete! Starting development servers..."
echo ""
echo "Backend API: http://localhost:8080"
echo "Frontend: http://localhost:5173"
echo ""

npm run dev
