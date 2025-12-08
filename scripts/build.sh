#!/bin/bash

# Go-Reader Build Script
# This script builds the application for production

set -e

echo "🏗️  Building Go-Reader..."

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf dist bin frontend/dist

# Create bin directory
mkdir -p bin

# Build backend
echo "🔨 Building backend..."
cd backend
go build -o ../bin/go-reader-server cmd/server/main.go
cd ..
echo "✅ Backend built successfully"

# Build frontend
echo "🔨 Building frontend..."
cd frontend
npm run build
cd ..
echo "✅ Frontend built successfully"

# Build Electron app
echo "🔨 Building Electron app..."
npm run electron:build
echo "✅ Electron app built successfully"

echo ""
echo "🎉 Build complete!"
echo "📦 Distribution files are in the 'dist' directory"
