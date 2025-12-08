#!/bin/bash

# Clean Install Script - 清除舊依賴並重新安裝

set -e

echo "🧹 Cleaning old dependencies..."

# Clean root node_modules
if [ -d "node_modules" ]; then
    echo "Removing root node_modules..."
    rm -rf node_modules
fi

if [ -f "package-lock.json" ]; then
    echo "Removing root package-lock.json..."
    rm -f package-lock.json
fi

# Clean frontend node_modules
if [ -d "frontend/node_modules" ]; then
    echo "Removing frontend node_modules..."
    rm -rf frontend/node_modules
fi

if [ -f "frontend/package-lock.json" ]; then
    echo "Removing frontend package-lock.json..."
    rm -f frontend/package-lock.json
fi

# Clean Go cache
if [ -d "backend" ]; then
    echo "Cleaning Go module cache..."
    cd backend
    go clean -modcache 2>/dev/null || true
    cd ..
fi

echo ""
echo "✨ Installing fresh dependencies..."
echo ""

# Install Go dependencies
if command -v go &> /dev/null; then
    echo "📦 Installing Go dependencies..."
    cd backend
    go mod download
    go mod tidy
    cd ..
    echo "✅ Go dependencies installed"
else
    echo "⚠️  Go not found, skipping backend dependencies"
fi

echo ""

# Install frontend dependencies
echo "📦 Installing frontend dependencies..."
cd frontend
npm install
cd ..
echo "✅ Frontend dependencies installed"

echo ""

# Install root dependencies (Electron)
echo "📦 Installing Electron dependencies..."
npm install
echo "✅ Electron dependencies installed"

echo ""
echo "🎉 Clean install completed!"
echo ""
echo "You can now run:"
echo "  npm run dev          # Start development environment"
echo "  ./scripts/build.sh   # Build for production"
