## init

# 🧩 Go + React Fullstack Starter

Dự án mẫu kết hợp giữa **Golang (Gin)** và **ReactJS (Vite)**.

## 🚀 Tính năng

- 🧠 Backend API bằng Gin framework
- ⚛️ Frontend ReactJS với Vite
- 🔄 Proxy API khi dev
- 🏗 Tự động build React vào `dist/` và Go server sẽ phục vụ giao diện đó
- 📦 Tích hợp build script

---

## 📁 Cấu trúc thư mục

```bash
go-react-app/
├── backend/              # Golang API server
│   └── main.go
├── frontend/             # ReactJS frontend (Vite)
│   ├── package.json
│   ├── vite.config.js
│   └── src/
├── build.sh              # Script build fullstack
├── .gitignore
├── README.md

```

---

### Cách chạy Backend, Frontend và Build

### 🏗 Build toàn bộ fullstack

./build.sh

#### Chạy backend API

bash
cd backend
go run main.go

#### Chạy frontend dev

cd frontend
npm install # Chạy 1 lần để cài dependencies
npm run dev # Start Vite dev server
