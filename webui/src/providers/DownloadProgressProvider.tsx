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
import type { DownloadSessionResponse, DownloadStreamEvent } from "@/api"

type DownloadLogEntry = {
  id: string
  level: string
  message: string
  timestamp: number
}

type ActiveDownloadSession = {
  sessionId: string
  target: string
  title: string
  status: string
  progress: number
  current: number
  total: number
  currentName: string
  outputDir: string
  streamUrl: string
  logs: DownloadLogEntry[]
  createdAt: number
  updatedAt: number
}

type DownloadProgressContextValue = {
  sessions: ActiveDownloadSession[]
  activeSession: ActiveDownloadSession | null
  isPanelOpen: boolean
  beginSession: (session: DownloadSessionResponse) => void
  closePanel: () => void
  selectSession: (sessionId: string) => void
}

const DownloadProgressContext = createContext<DownloadProgressContextValue | null>(null)

function createBaseSession(session: DownloadSessionResponse): ActiveDownloadSession {
  const now = Date.now()
  return {
    sessionId: session.sessionId,
    target: session.target,
    title: session.title,
    status: session.status,
    progress: 0,
    current: 0,
    total: 0,
    currentName: "",
    outputDir: session.outputDir,
    streamUrl: session.streamUrl,
    logs: [],
    createdAt: now,
    updatedAt: now,
  }
}

function toLogEntry(event: DownloadStreamEvent): DownloadLogEntry {
  return {
    id: `${event.timestamp}-${event.type}-${event.message ?? ""}`,
    level: event.level || (event.type === "error" ? "error" : "info"),
    message: event.message || "",
    timestamp: event.timestamp,
  }
}

export function DownloadProgressProvider({ children }: PropsWithChildren) {
  const [sessions, setSessions] = useState<ActiveDownloadSession[]>([])
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null)
  const [isPanelOpen, setIsPanelOpen] = useState(false)
  const eventSourcesRef = useRef<Record<string, EventSource>>({})

  const cleanupSession = useCallback((sessionId: string) => {
    const source = eventSourcesRef.current[sessionId]
    if (source) {
      source.close()
      delete eventSourcesRef.current[sessionId]
    }
  }, [])

  const cleanupAll = useCallback(() => {
    Object.keys(eventSourcesRef.current).forEach((sessionId) => {
      cleanupSession(sessionId)
    })
  }, [cleanupSession])

  useEffect(() => cleanupAll, [cleanupAll])

  const closePanel = useCallback(() => {
    setIsPanelOpen(false)
  }, [])

  const selectSession = useCallback((sessionId: string) => {
    setActiveSessionId(sessionId)
    setIsPanelOpen(true)
  }, [])

  const appendLog = useCallback((event: DownloadStreamEvent) => {
    if (!event.message) {
      return
    }

    setSessions((currentSessions) =>
      currentSessions.map((session) => {
        if (session.sessionId !== event.sessionId) {
          return session
        }

        const nextLogs = [...session.logs, toLogEntry(event)].slice(-200)
        return {
          ...session,
          logs: nextLogs,
          updatedAt: event.timestamp || Date.now(),
        }
      }),
    )
  }, [])

  const applyEvent = useCallback(
    (event: DownloadStreamEvent) => {
      if (event.type === "heartbeat") {
        return
      }

      setSessions((currentSessions) =>
        currentSessions.map((session) => {
          if (session.sessionId !== event.sessionId) {
            return session
          }

          return {
            ...session,
            status: event.status || session.status,
            progress: typeof event.progress === "number" ? event.progress : session.progress,
            current: typeof event.current === "number" ? event.current : session.current,
            total: typeof event.total === "number" ? event.total : session.total,
            currentName: event.currentName ?? session.currentName,
            outputDir: event.outputDir || session.outputDir,
            title: event.title || session.title,
            updatedAt: event.timestamp || Date.now(),
          }
        }),
      )

      if (event.type === "log" || event.type === "error" || event.type === "done" || event.type === "start") {
        appendLog(event)
      }

      if (event.type === "error") {
        cleanupSession(event.sessionId)
      }
      if (event.type === "done") {
        cleanupSession(event.sessionId)
      }
    },
    [appendLog, cleanupSession],
  )

  const beginSession = useCallback(
    (nextSession: DownloadSessionResponse) => {
      cleanupSession(nextSession.sessionId)
      setIsPanelOpen(true)
      setActiveSessionId(nextSession.sessionId)
      setSessions((currentSessions) => {
        const nextBaseSession = createBaseSession(nextSession)
        const nextSessions = [nextBaseSession, ...currentSessions.filter((session) => session.sessionId !== nextSession.sessionId)]
        return nextSessions.slice(0, 8)
      })

      const source = new EventSource(nextSession.streamUrl)
      eventSourcesRef.current[nextSession.sessionId] = source

      source.onmessage = (messageEvent) => {
        try {
          const payload = JSON.parse(messageEvent.data) as DownloadStreamEvent
          applyEvent(payload)
        } catch {
          // Ignore malformed SSE payloads from interrupted responses.
        }
      }

      source.onerror = () => {
        if (eventSourcesRef.current[nextSession.sessionId] !== source) {
          return
        }

        setSessions((currentSessions) =>
          currentSessions.map((session) => {
            if (session.sessionId !== nextSession.sessionId) {
              return session
            }
            return {
              ...session,
              logs: [
                ...session.logs,
                {
                  id: `${Date.now()}-stream-error-${nextSession.sessionId}`,
                  level: "warning",
                  message: "进度连接已断开，若服务端仍在执行，可重新发起下载查看最新状态。",
                  timestamp: Date.now(),
                },
              ].slice(-200),
              updatedAt: Date.now(),
            }
          }),
        )
        cleanupSession(nextSession.sessionId)
      }
    },
    [applyEvent, cleanupSession],
  )

  const activeSession = useMemo(() => {
    if (sessions.length === 0) {
      return null
    }
    return sessions.find((session) => session.sessionId === activeSessionId) ?? sessions[0]
  }, [activeSessionId, sessions])

  const value = useMemo<DownloadProgressContextValue>(
    () => ({
      sessions,
      activeSession,
      isPanelOpen,
      beginSession,
      closePanel,
      selectSession,
    }),
    [activeSession, beginSession, closePanel, isPanelOpen, selectSession, sessions],
  )

  return <DownloadProgressContext.Provider value={value}>{children}</DownloadProgressContext.Provider>
}

export function useDownloadProgress() {
  const context = useContext(DownloadProgressContext)
  if (!context) {
    throw new Error("useDownloadProgress 必须在 DownloadProgressProvider 内部使用")
  }

  return context
}
