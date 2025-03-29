# React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) uses [Babel](https://babeljs.io/) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## Expanding the ESLint configuration

If you are developing a production application, we recommend using TypeScript and enable type-aware lint rules. Check out the [TS template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) to integrate TypeScript and [`typescript-eslint`](https://typescript-eslint.io) in your project.

<!--  -->

## init

hướng dẫn khởi tạo đầy đủ project Go (Gin) + ReactJS, hỗ trợ:

✅ Chạy dev từng phần riêng

✅ Build chung một lần, server Go sẽ phục vụ cả React + API

<pre lang="md"> <code> <details> <summary>📁 Project structure</summary> ``` go-react-app/ ├── backend/ │ └── main.go ├── frontend/ │ ├── package.json │ └── src/ ├── .gitignore ├── README.md ``` </details> </code> </pre>
