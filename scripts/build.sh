#!/bin/bash

# 脚本出错时立即退出
set -e

# 项目根目录
ROOT_DIR=$(dirname "$0")/..
cd "$ROOT_DIR"

# --- 清理旧的构建产物 ---
echo "🧹 Cleaning up old build directory..."
rm -rf build
mkdir -p build/static

# --- 构建前端 ---
echo "📦 Building frontend..."
cd frontend

# 如果 node_modules 不存在，则安装依赖
if [ ! -d "node_modules" ]; then
  echo "Installing frontend dependencies..."
  npm install
fi

npm run build
echo "✅ Frontend built successfully."

# --- 移动前端产物 ---
echo "🚚 Moving frontend assets..."
cd ..
cp -r frontend/build/* build/static/
rm -rf frontend/build

# --- 构建后端 ---
echo "🏗️ Building backend..."
# 编译 Go 应用，输出到 build/server（为 Linux 平台编译）
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/server main.go
echo "✅ Backend built successfully."

echo "🚀 Build complete! All artifacts are in the 'build' directory."
echo "   - Backend executable: build/server"
echo "   - Frontend assets: build/static"
echo "To start the application, run: ./scripts/start.sh"
