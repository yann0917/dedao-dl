import { RouterProvider } from "react-router-dom"
import { GlobalAudioPlayer } from "@/components/player/GlobalAudioPlayer"
import { AuthProvider } from "@/providers/AuthProvider"
import { AudioPlayerProvider } from "@/providers/AudioPlayerProvider"
import { router } from "@/router"

export default function App() {
  return (
    <AuthProvider>
      <AudioPlayerProvider>
        <RouterProvider router={router} />
        <GlobalAudioPlayer />
      </AudioPlayerProvider>
    </AuthProvider>
  )
}
