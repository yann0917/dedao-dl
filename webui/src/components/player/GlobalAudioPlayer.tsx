import { Pause, Play, SkipBack, SkipForward, Volume2 } from "lucide-react"
import { Card } from "@/components/ui/Card"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

function formatTime(value: number) {
  if (!Number.isFinite(value) || value <= 0) {
    return "00:00"
  }

  const minutes = Math.floor(value / 60)
  const seconds = Math.floor(value % 60)
  return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`
}

export function GlobalAudioPlayer() {
  const { currentTrack, currentIndex, queue, playing, currentTime, duration, togglePlay, playNext, playPrev, seek } =
    useAudioPlayer()

  if (!currentTrack) {
    return null
  }

  return (
    <div className="fixed bottom-4 left-4 right-4 z-50 lg:left-[300px]">
      <Card className="border border-slate-200/80 bg-white/95 p-4 shadow-soft backdrop-blur">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex min-w-0 items-center gap-4">
            <img
              alt={currentTrack.title}
              className="h-14 w-14 rounded-2xl object-cover"
              src={currentTrack.poster || "https://placehold.co/120x120/e2e8f0/334155?text=Audio"}
            />
            <div className="min-w-0">
              <p className="truncate text-sm text-slate-500">正在播放</p>
              <p className="truncate text-base font-semibold text-slate-950">{currentTrack.title}</p>
              <p className="truncate text-xs text-slate-500">{currentTrack.subtitle || `队列 ${currentIndex + 1}/${queue.length}`}</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-slate-100 text-slate-700 transition hover:bg-slate-200"
              onClick={playPrev}
              type="button"
            >
              <SkipBack className="h-4 w-4" />
            </button>
            <button
              className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary text-white transition hover:bg-primary/90"
              onClick={togglePlay}
              type="button"
            >
              {playing ? <Pause className="h-5 w-5" /> : <Play className="h-5 w-5" />}
            </button>
            <button
              className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-slate-100 text-slate-700 transition hover:bg-slate-200"
              onClick={playNext}
              type="button"
            >
              <SkipForward className="h-4 w-4" />
            </button>
          </div>

          <div className="flex min-w-0 flex-1 items-center gap-3 lg:max-w-xl">
            <Volume2 className="h-4 w-4 text-slate-400" />
            <span className="w-11 text-xs text-slate-500">{formatTime(currentTime)}</span>
            <input
              className="h-2 flex-1 cursor-pointer accent-primary"
              max={duration || 0}
              min={0}
              onChange={(event) => seek(Number(event.target.value))}
              type="range"
              value={Math.min(currentTime, duration || 0)}
            />
            <span className="w-11 text-right text-xs text-slate-500">{formatTime(duration)}</span>
          </div>
        </div>
      </Card>
    </div>
  )
}
