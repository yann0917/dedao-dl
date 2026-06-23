import { RouterProvider } from "react-router-dom"
import { GlobalAudioPlayer } from "@/components/player/GlobalAudioPlayer"
import { AuthProvider } from "@/providers/AuthProvider"
import { AudioPlayerProvider } from "@/providers/AudioPlayerProvider"
import { ThemeProvider } from "@/providers/ThemeProvider"
import { router } from "@/router"

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <AudioPlayerProvider>
          <RouterProvider router={router} />
          <GlobalAudioPlayer />
        </AudioPlayerProvider>
      </AuthProvider>
    </ThemeProvider>
  )
}
