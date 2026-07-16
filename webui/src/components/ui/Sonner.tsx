import { Toaster } from "sonner"
import { useTheme } from "@/providers/ThemeProvider"

export function Sonner() {
  const { theme } = useTheme()

  return (
    <Toaster
      closeButton
      position="top-right"
      richColors
      theme={theme}
      toastOptions={{
        duration: 3000,
      }}
    />
  )
}
