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

const AUDIO_PLAYER_STORAGE_KEY = "dedao:audio-player"

export type PlayerTrack = {
  id: string
  title: string
  src: string
  poster?: string
  subtitle?: string
}

type PersistedAudioPlayerState = {
  queue: PlayerTrack[]
  currentIndex: number
  currentTime: number
}

function readPersistedAudioPlayerState(): PersistedAudioPlayerState {
  if (typeof window === "undefined") {
    return {
      queue: [],
      currentIndex: -1,
      currentTime: 0,
    }
  }

  try {
    const raw = window.localStorage.getItem(AUDIO_PLAYER_STORAGE_KEY)
    if (!raw) {
      return {
        queue: [],
        currentIndex: -1,
        currentTime: 0,
      }
    }

    const parsed = JSON.parse(raw) as Partial<PersistedAudioPlayerState>
    const queue = Array.isArray(parsed.queue)
      ? parsed.queue.filter(
          (track): track is PlayerTrack =>
            !!track &&
            typeof track.id === "string" &&
            typeof track.title === "string" &&
            typeof track.src === "string" &&
            track.src.length > 0,
        )
      : []
    const currentIndex =
      typeof parsed.currentIndex === "number" && Number.isInteger(parsed.currentIndex) ? parsed.currentIndex : -1
    const currentTime = typeof parsed.currentTime === "number" && Number.isFinite(parsed.currentTime) ? parsed.currentTime : 0

    return {
      queue,
      currentIndex: queue.length > 0 ? Math.min(Math.max(currentIndex, 0), queue.length - 1) : -1,
      currentTime: Math.max(currentTime, 0),
    }
  } catch {
    return {
      queue: [],
      currentIndex: -1,
      currentTime: 0,
    }
  }
}

function clearPersistedAudioPlayerState() {
  if (typeof window === "undefined") {
    return
  }

  window.localStorage.removeItem(AUDIO_PLAYER_STORAGE_KEY)
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
  playQueueIndex: (index: number) => void
  togglePlay: () => void
  playNext: () => void
  playPrev: () => void
  seek: (time: number) => void
  clearQueue: () => void
}

const AudioPlayerContext = createContext<AudioPlayerContextValue | null>(null)

export function AudioPlayerProvider({ children }: PropsWithChildren) {
  const persistedState = useMemo(readPersistedAudioPlayerState, [])
  const audioRef = useRef<HTMLAudioElement | null>(null)
  const [queue, setQueueState] = useState<PlayerTrack[]>(persistedState.queue)
  const [currentIndex, setCurrentIndex] = useState(persistedState.currentIndex)
  const [playing, setPlaying] = useState(false)
  const [currentTime, setCurrentTime] = useState(persistedState.currentTime)
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

    if (currentTime > 0) {
      audio.currentTime = currentTime
    }

    void audio.play().catch(() => {
      setPlaying(false)
    })
  }, [currentTime, currentTrack])

  useEffect(() => {
    if (queue.length === 0 || currentIndex < 0) {
      clearPersistedAudioPlayerState()
      return
    }

    if (typeof window === "undefined") {
      return
    }

    const payload: PersistedAudioPlayerState = {
      queue,
      currentIndex: Math.min(Math.max(currentIndex, 0), queue.length - 1),
      currentTime,
    }
    window.localStorage.setItem(AUDIO_PLAYER_STORAGE_KEY, JSON.stringify(payload))
  }, [currentIndex, currentTime, queue])

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

  const playQueueIndex = useCallback(
    (index: number) => {
      if (index < 0 || index >= queue.length) {
        return
      }
      setCurrentIndex(index)
      setCurrentTime(0)
    },
    [queue.length],
  )

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

  const clearQueue = useCallback(() => {
    const audio = audioRef.current
    if (audio) {
      audio.pause()
      audio.removeAttribute("src")
      audio.load()
    }

    setQueueState([])
    setCurrentIndex(-1)
    setCurrentTime(0)
    setDuration(0)
    setPlaying(false)
    clearPersistedAudioPlayerState()
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
      playQueueIndex,
      togglePlay,
      playNext,
      playPrev,
      seek,
      clearQueue,
    }),
    [clearQueue, currentIndex, currentTime, currentTrack, duration, playNext, playPrev, playQueueIndex, playTrack, playing, queue, seek, setQueue, togglePlay],
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
