import { ListMusic, Pause, Play, SkipBack, SkipForward, Volume2, X } from "lucide-react"
import { useState } from "react"
import { cn } from "@/lib/cn"
import { Card } from "@/components/ui/Card"
import { semanticMetaTextClass } from "@/lib/semanticStyles"
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
  const [showQueue, setShowQueue] = useState(false)
  const { currentTrack, currentIndex, queue, playing, currentTime, duration, togglePlay, playNext, playPrev, playQueueIndex, seek, clearQueue } =
    useAudioPlayer()

  if (!currentTrack) {
    return null
  }

  return (
    <div className="fixed bottom-4 left-4 right-4 z-50 lg:left-[300px]">
      <Card className="relative border-border/80 bg-surface-panel/95 p-4 shadow-soft backdrop-blur">
        <button
          className="absolute right-3 top-3 inline-flex h-8 w-8 items-center justify-center rounded-full text-text-muted transition hover:bg-danger-soft hover:text-danger"
          onClick={clearQueue}
          type="button"
        >
          <X className="h-4 w-4" />
        </button>

        <div className="flex flex-col gap-4 pr-10">
          <div className="flex min-w-0 items-center gap-4">
            <img
              alt={currentTrack.title}
              className="h-14 w-14 rounded-2xl object-cover"
              src={currentTrack.poster || "https://placehold.co/120x120/e2e8f0/334155?text=Audio"}
            />
            <div className="min-w-0">
              <p className={`truncate text-sm ${semanticMetaTextClass}`}>正在播放</p>
              <p className="truncate text-base font-semibold text-text-primary">{currentTrack.title}</p>
              <p className={`truncate text-xs ${semanticMetaTextClass}`}>{currentTrack.subtitle || `队列 ${currentIndex + 1}/${queue.length}`}</p>
            </div>
          </div>

          <div className="flex flex-wrap items-center justify-center gap-3">
            <button
              className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-surface-soft text-text-secondary transition hover:bg-accent-soft hover:text-accent"
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
              className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-surface-soft text-text-secondary transition hover:bg-accent-soft hover:text-accent"
              onClick={playNext}
              type="button"
            >
              <SkipForward className="h-4 w-4" />
            </button>
            <button
              aria-label={showQueue ? "折叠播放列表" : "展开播放列表"}
              className={cn(
                "inline-flex h-10 w-10 items-center justify-center rounded-full transition",
                showQueue
                  ? "bg-primary/10 text-primary"
                  : "bg-surface-soft text-text-secondary hover:bg-surface-panel hover:text-text-primary",
              )}
              onClick={() => setShowQueue((value) => !value)}
              type="button"
            >
              <ListMusic className="h-4 w-4" />
            </button>
          </div>

          <div className="flex min-w-0 items-center gap-3">
            <Volume2 className="h-4 w-4 text-text-muted" />
            <span className="w-11 text-xs text-text-muted">{formatTime(currentTime)}</span>
            <input
              className="h-2 flex-1 cursor-pointer accent-primary"
              max={duration || 0}
              min={0}
              onChange={(event) => seek(Number(event.target.value))}
              type="range"
              value={Math.min(currentTime, duration || 0)}
            />
            <span className="w-11 text-right text-xs text-text-muted">{formatTime(duration)}</span>
          </div>
        </div>

        {showQueue ? (
          <div className="mt-4 max-h-72 space-y-2 overflow-y-auto border-t border-border/70 pt-4">
            {queue.map((track, index) => {
              const active = index === currentIndex

              return (
                <button
                  className={cn(
                    "flex w-full items-center justify-between rounded-2xl px-3 py-3 text-left transition",
                    active ? "bg-primary/10 text-primary" : "bg-surface-soft/80 text-text-secondary hover:bg-surface-soft",
                  )}
                  key={`${track.id}-${index}`}
                  onClick={() => playQueueIndex(index)}
                  type="button"
                >
                  <div className="min-w-0">
                    <p className={cn("truncate text-sm font-medium", active ? "text-primary" : "text-text-primary")}>{track.title}</p>
                    <p className={cn("truncate text-xs", active ? "text-primary/80" : semanticMetaTextClass)}>
                      {track.subtitle || `第 ${index + 1} 首`}
                    </p>
                  </div>
                  <span className={cn("shrink-0 text-xs", active ? "text-primary" : "text-text-muted")}>{active ? "播放中" : `${index + 1}/${queue.length}`}</span>
                </button>
              )
            })}
          </div>
        ) : null}
      </Card>
    </div>
  )
}
