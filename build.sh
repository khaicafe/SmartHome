#!/bin/bash
echo "📦 Building frontend..."
cd frontend && npm run build && cd ..

echo "🔨 Building backend..."
cd backend && go build -o app && cd ..

echo "🚀 Done! Run ./backend/app to start full web + API server."
