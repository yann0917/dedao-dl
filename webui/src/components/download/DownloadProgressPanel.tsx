import { X } from "lucide-react"
import { useEffect, useRef } from "react"
import { useDownloadProgress } from "@/providers/DownloadProgressProvider"

function formatTime(timestamp: number) {
  if (!timestamp) {
    return ""
  }
  return new Date(timestamp).toLocaleTimeString("zh-CN", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  })
}

function statusLabel(status: string) {
  if (status === "running") return "下载中"
  if (status === "success") return "已完成"
  if (status === "error") return "失败"
  return "准备中"
}

function statusClassName(status: string) {
  if (status === "success") return "bg-success-soft text-success"
  if (status === "error") return "bg-danger-soft text-danger"
  if (status === "running") return "bg-warning-soft text-warning"
  return "bg-surface-soft text-text-secondary"
}

function logClassName(level: string) {
  if (level === "error") return "text-danger"
  if (level === "warning") return "text-warning"
  return "text-text-secondary"
}

export function DownloadProgressPanel() {
  const { activeSession, closePanel, isPanelOpen, selectSession, sessions } = useDownloadProgress()
  const logContainerRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    if (!activeSession) {
      return
    }

    const container = logContainerRef.current
    if (!container) {
      return
    }

    container.scrollTop = container.scrollHeight
  }, [activeSession?.sessionId, activeSession?.logs.length])

  if (!isPanelOpen || !activeSession) {
    return null
  }

  return (
    <section className="fixed bottom-4 right-4 z-[80] w-[min(95vw,840px)] rounded-3xl border border-border bg-surface-panel p-5 shadow-soft backdrop-blur">
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-sm text-text-muted">下载进度</p>
          <h3 className="mt-1 text-lg font-semibold text-text-primary">{activeSession.title || "下载任务"}</h3>
          <div className="mt-3 flex flex-wrap items-center gap-2 text-xs">
            <span className={`rounded-full px-3 py-1 font-medium ${statusClassName(activeSession.status)}`}>
              {statusLabel(activeSession.status)}
            </span>
            <span className="rounded-full bg-surface-soft px-3 py-1 text-text-secondary">{activeSession.target}</span>
            <span className="rounded-full bg-surface-soft px-3 py-1 text-text-secondary">最近 {sessions.length} 次</span>
            {activeSession.currentName ? (
              <span className="rounded-full bg-surface-soft px-3 py-1 text-text-secondary">
                当前：{activeSession.currentName}
              </span>
            ) : null}
          </div>
        </div>

        <button
          className="flex h-10 w-10 items-center justify-center rounded-2xl text-text-muted transition hover:bg-surface-soft hover:text-text-primary"
          onClick={closePanel}
          type="button"
        >
          <X className="h-5 w-5" />
        </button>
      </div>

      <div className="mt-5 grid gap-4 lg:grid-cols-[220px_minmax(0,1fr)]">
        <div className="rounded-2xl border border-border bg-surface-page/70">
          <div className="border-b border-border px-4 py-3 text-sm font-medium text-text-primary">最近会话</div>
          <div className="max-h-[28rem] space-y-2 overflow-y-auto p-3">
            {sessions.map((session) => {
              const isActive = session.sessionId === activeSession.sessionId
              return (
                <button
                  className={`w-full rounded-2xl border px-3 py-3 text-left transition ${
                    isActive
                      ? "border-accent bg-accent/10"
                      : "border-border bg-surface-panel hover:bg-surface-soft"
                  }`}
                  key={session.sessionId}
                  onClick={() => selectSession(session.sessionId)}
                  type="button"
                >
                  <div className="flex items-center justify-between gap-2">
                    <span className="line-clamp-2 text-sm font-medium text-text-primary">{session.title || "下载任务"}</span>
                    <span className={`shrink-0 rounded-full px-2 py-1 text-[11px] font-medium ${statusClassName(session.status)}`}>
                      {statusLabel(session.status)}
                    </span>
                  </div>
                  <div className="mt-2 flex items-center justify-between text-xs text-text-muted">
                    <span>{session.progress}%</span>
                    <span>{formatTime(session.updatedAt)}</span>
                  </div>
                </button>
              )
            })}
          </div>
        </div>

        <div>
          <div>
            <div className="h-2 overflow-hidden rounded-full bg-surface-soft">
              <div
                className="h-full rounded-full bg-accent transition-[width] duration-300"
                style={{ width: `${Math.max(0, Math.min(100, activeSession.progress))}%` }}
              />
            </div>
            <div className="mt-2 flex items-center justify-between text-xs text-text-muted">
              <span>{activeSession.progress}%</span>
              <span>
                {activeSession.total > 0 ? `${activeSession.current}/${activeSession.total}` : "等待进度数据"}
              </span>
            </div>
          </div>

          <div className="mt-5 rounded-2xl bg-surface-soft px-4 py-3 text-sm text-text-secondary">
            输出目录：{activeSession.outputDir}
          </div>

          <div className="mt-5 rounded-2xl border border-border bg-surface-page/80">
            <div className="border-b border-border px-4 py-3 text-sm font-medium text-text-primary">实时日志</div>
            <div className="max-h-72 space-y-2 overflow-y-auto px-4 py-3 text-sm" ref={logContainerRef}>
              {activeSession.logs.length > 0 ? (
                activeSession.logs.map((log) => (
                  <div className={`leading-6 ${logClassName(log.level)}`} key={log.id}>
                    <span className="mr-2 text-xs text-text-muted">{formatTime(log.timestamp)}</span>
                    <span>{log.message}</span>
                  </div>
                ))
              ) : (
                <div className="text-text-muted">等待下载日志...</div>
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
