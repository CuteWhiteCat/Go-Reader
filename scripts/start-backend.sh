#!/bin/bash

# Start Go Backend Server

set -e

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo ""
    echo "Please install Go 1.21 or higher:"
    echo "  wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz"
    echo "  sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz"
    echo "  export PATH=\$PATH:/usr/local/go/bin"
    echo ""
    exit 1
fi

echo "🚀 Starting Go backend server..."
echo "Go version: $(go version)"
echo ""

cd backend

# Download dependencies if needed
if [ ! -f "go.sum" ]; then
    echo "📦 Downloading Go dependencies..."
    go mod download
    go mod tidy
fi

# Start the server
echo "✅ Starting server on http://localhost:8080"
go run cmd/server/main.go
