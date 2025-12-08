#!/bin/bash

# Start React Frontend

set -e

echo "🚀 Starting React frontend..."
echo "Node version: $(node -v)"
echo ""

cd frontend

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Start Vite dev server
echo "✅ Starting frontend on http://localhost:5173"
npm run dev
