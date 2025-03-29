# React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## Expanding the ESLint configuration

If you are developing a production application, we recommend using TypeScript and enable type-aware lint rules. Check out the [TS template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) to integrate TypeScript and [`typescript-eslint`](https://typescript-eslint.io) in your project.

<!--  -->

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
