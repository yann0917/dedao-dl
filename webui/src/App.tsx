import { DownloadProgressPanel } from "@/components/download/DownloadProgressPanel"
import { Sonner } from "@/components/ui/Sonner"
import { RouterProvider } from "react-router-dom"
import { GlobalAudioPlayer } from "@/components/player/GlobalAudioPlayer"
import { DownloadProgressProvider } from "@/providers/DownloadProgressProvider"
import { AuthProvider } from "@/providers/AuthProvider"
import { AudioPlayerProvider } from "@/providers/AudioPlayerProvider"
import { ThemeProvider } from "@/providers/ThemeProvider"
import { router } from "@/router"

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <DownloadProgressProvider>
          <AudioPlayerProvider>
            <RouterProvider router={router} />
            <DownloadProgressPanel />
            <GlobalAudioPlayer />
            <Sonner />
          </AudioPlayerProvider>
        </DownloadProgressProvider>
      </AuthProvider>
    </ThemeProvider>
  )
}
