import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type PropsWithChildren,
} from "react"

export type PlayerTrack = {
  id: string
  title: string
  src: string
  poster?: string
  subtitle?: string
}

type AudioPlayerContextValue = {
  queue: PlayerTrack[]
  currentTrack: PlayerTrack | null
  currentIndex: number
  playing: boolean
  currentTime: number
  duration: number
  setQueue: (tracks: PlayerTrack[], startIndex?: number) => void
  playTrack: (track: PlayerTrack) => void
  togglePlay: () => void
  playNext: () => void
  playPrev: () => void
  seek: (time: number) => void
}

const AudioPlayerContext = createContext<AudioPlayerContextValue | null>(null)

export function AudioPlayerProvider({ children }: PropsWithChildren) {
  const audioRef = useRef<HTMLAudioElement | null>(null)
  const [queue, setQueueState] = useState<PlayerTrack[]>([])
  const [currentIndex, setCurrentIndex] = useState(-1)
  const [playing, setPlaying] = useState(false)
  const [currentTime, setCurrentTime] = useState(0)
  const [duration, setDuration] = useState(0)

  const currentTrack = currentIndex >= 0 ? queue[currentIndex] ?? null : null

  useEffect(() => {
    const audio = audioRef.current
    if (!audio) {
      return
    }

    const handleTimeUpdate = () => setCurrentTime(audio.currentTime)
    const handleLoadedMetadata = () => setDuration(Number.isFinite(audio.duration) ? audio.duration : 0)
    const handlePlay = () => setPlaying(true)
    const handlePause = () => setPlaying(false)
    const handleEnded = () => {
      setCurrentTime(0)
      setCurrentIndex((index) => {
        if (index + 1 < queue.length) {
          return index + 1
        }
        return index
      })
    }

    audio.addEventListener("timeupdate", handleTimeUpdate)
    audio.addEventListener("loadedmetadata", handleLoadedMetadata)
    audio.addEventListener("play", handlePlay)
    audio.addEventListener("pause", handlePause)
    audio.addEventListener("ended", handleEnded)

    return () => {
      audio.removeEventListener("timeupdate", handleTimeUpdate)
      audio.removeEventListener("loadedmetadata", handleLoadedMetadata)
      audio.removeEventListener("play", handlePlay)
      audio.removeEventListener("pause", handlePause)
      audio.removeEventListener("ended", handleEnded)
    }
  }, [queue.length])

  useEffect(() => {
    const audio = audioRef.current
    if (!audio || !currentTrack?.src) {
      return
    }

    if (audio.src !== currentTrack.src) {
      audio.src = currentTrack.src
    }

    void audio.play().catch(() => {
      setPlaying(false)
    })
  }, [currentTrack])

  const setQueue = useCallback((tracks: PlayerTrack[], startIndex = 0) => {
    const validTracks = tracks.filter((track) => !!track.src)
    if (validTracks.length === 0) {
      return
    }

    setQueueState(validTracks)
    setCurrentIndex(Math.min(Math.max(startIndex, 0), validTracks.length - 1))
    setCurrentTime(0)
  }, [])

  const playTrack = useCallback((track: PlayerTrack) => {
    setQueueState([track])
    setCurrentIndex(0)
    setCurrentTime(0)
  }, [])

  const togglePlay = useCallback(() => {
    const audio = audioRef.current
    if (!audio || !currentTrack) {
      return
    }

    if (audio.paused) {
      void audio.play()
    } else {
      audio.pause()
    }
  }, [currentTrack])

  const playNext = useCallback(() => {
    setCurrentIndex((index) => {
      if (index + 1 < queue.length) {
        return index + 1
      }
      return index
    })
  }, [queue.length])

  const playPrev = useCallback(() => {
    setCurrentIndex((index) => Math.max(index - 1, 0))
  }, [])

  const seek = useCallback((time: number) => {
    const audio = audioRef.current
    if (!audio) {
      return
    }

    audio.currentTime = time
    setCurrentTime(time)
  }, [])

  const value = useMemo<AudioPlayerContextValue>(
    () => ({
      queue,
      currentTrack,
      currentIndex,
      playing,
      currentTime,
      duration,
      setQueue,
      playTrack,
      togglePlay,
      playNext,
      playPrev,
      seek,
    }),
    [currentIndex, currentTime, currentTrack, duration, playNext, playPrev, playTrack, playing, queue, seek, setQueue, togglePlay],
  )

  return (
    <AudioPlayerContext.Provider value={value}>
      {children}
      <audio preload="metadata" ref={audioRef} />
    </AudioPlayerContext.Provider>
  )
}

export function useAudioPlayer() {
  const context = useContext(AudioPlayerContext)
  if (!context) {
    throw new Error("useAudioPlayer 必须在 AudioPlayerProvider 内部使用")
  }

  return context
}
