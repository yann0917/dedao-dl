import React from "react"
import ReactDOM from "react-dom/client"
import App from "@/App"
import { initializeTheme } from "@/providers/ThemeProvider"
import "./index.css"

initializeTheme()

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
