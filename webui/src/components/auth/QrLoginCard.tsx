import { useEffect, useMemo, useState } from "react"
import { Loader2, RefreshCcw } from "lucide-react"
import { toast } from "sonner"
import { api, type QRCodeSession, type UserInfo } from "@/api"

type QrLoginCardProps = {
  onLoginSuccess: (user?: UserInfo | null) => void | Promise<void>
}

const cropMarkPositions = [
  "-left-2 -top-2 border-l border-t",
  "-right-2 -top-2 border-r border-t",
  "-bottom-2 -left-2 border-b border-l",
  "-bottom-2 -right-2 border-b border-r",
]

export function QrLoginCard({ onLoginSuccess }: QrLoginCardProps) {
  const [session, setSession] = useState<QRCodeSession | null>(null)
  const [loading, setLoading] = useState(false)
  const [polling, setPolling] = useState(false)
  const [expired, setExpired] = useState(false)

  const remaining = useMemo(() => {
    if (!session?.expiresAt) {
      return ""
    }

    const seconds = Math.max(session.expiresAt - Math.floor(Date.now() / 1000), 0)
    const minutes = Math.floor(seconds / 60)
    return `${minutes}:${String(seconds % 60).padStart(2, "0")}`
  }, [session?.expiresAt])

  const loadQRCode = async () => {
    setLoading(true)

    try {
      const next = await api.auth.createQRCode()
      setSession(next)
      setExpired(false)
      setPolling(true)
    } catch (err) {
      setPolling(false)
      toast.error("二维码生成失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadQRCode()
  }, [])

  useEffect(() => {
    if (!polling || !session) {
      return
    }

    const timer = window.setInterval(async () => {
      try {
        const result = await api.auth.getQRCodeStatus(session.sessionId)
        if (result.status === 1) {
          window.clearInterval(timer)
          setPolling(false)
          toast.success("扫码成功", {
            description: "正在进入工作台...",
          })
          await onLoginSuccess(result.user ?? null)
          return
        }

        if (result.status === 2) {
          window.clearInterval(timer)
          setPolling(false)
          setExpired(true)
          toast.error("二维码已过期", {
            description: "请刷新后重新扫码",
          })
          return
        }
      } catch (err) {
        setPolling(false)
        window.clearInterval(timer)
        toast.error("轮询登录状态失败", {
          description: err instanceof Error ? err.message : "请稍后重试",
        })
      }
    }, 2000)

    return () => window.clearInterval(timer)
  }, [onLoginSuccess, polling, session])

  const statusText = expired
    ? "QR EXPIRED"
    : polling
      ? "WAITING FOR SCAN"
      : loading
        ? "PREPARING"
        : ""

  return (
    <div className="ed-panel relative">
      <div className="flex items-baseline justify-between border-b border-border px-7 pb-4 pt-6">
        <h2 className="font-display text-xl font-semibold tracking-[0.06em] text-text-primary">
          扫码登录
        </h2>
        <span className="font-mono text-xs tabular-nums text-text-muted">
          {expired ? "已过期" : remaining ? `有效期 ${remaining}` : "生成中"}
        </span>
      </div>

      <div className="px-7 pb-6 pt-7">
        <div className="relative mx-auto w-fit">
          <div className="relative bg-[#f7f4ed] p-4">
            {session?.qrCode ? (
              <img
                alt="登录二维码"
                className="block h-56 w-56"
                src={session.qrCode}
              />
            ) : (
              <div className="flex h-56 w-56 flex-col items-center justify-center gap-3 bg-[#f6f4ef] text-xs text-neutral-500">
                <Loader2 className="size-icon-lg animate-spin" />
                正在生成二维码
              </div>
            )}
          </div>
          {cropMarkPositions.map((pos) => (
            <span
              key={pos}
              aria-hidden="true"
              className={`pointer-events-none absolute size-icon-sm border-text-primary/40 ${pos}`}
            />
          ))}
          {expired && (
            <div className="absolute inset-0 flex flex-col items-center justify-center gap-4 bg-surface-page/95">
              <p className="text-sm text-text-secondary">二维码已过期</p>
              <button
                className="inline-flex h-9 items-center gap-2 border border-border-strong px-4 text-sm text-text-primary transition hover:bg-surface-soft"
                onClick={() => void loadQRCode()}
                type="button"
              >
                <RefreshCcw className="size-icon-sm" />
                重新生成
              </button>
            </div>
          )}
        </div>

        <p className="mt-6 text-center text-sm leading-6 text-text-secondary">
          打开得到 App 或微信「扫一扫」，确认后自动进入工作台
        </p>
        <p className="mt-2.5 flex h-4 items-center justify-center gap-2 font-mono text-[10px] tracking-[0.28em] text-text-muted">
          {polling && !expired && (
            <span
              aria-hidden="true"
              className="ed-pulse inline-block h-1.5 w-1.5 rounded-full bg-accent"
            />
          )}
          {statusText}
        </p>
      </div>

      <div className="flex items-center justify-between border-t border-border px-7 py-4">
        <button
          className="inline-flex h-9 items-center gap-2 border border-border-strong px-4 text-sm text-text-primary transition hover:bg-surface-soft disabled:cursor-not-allowed disabled:opacity-50"
          disabled={loading}
          onClick={() => void loadQRCode()}
          type="button"
        >
          <RefreshCcw className="size-icon-sm" />
          刷新二维码
        </button>
        <a
          className="text-sm text-text-muted underline decoration-border-strong underline-offset-4 transition hover:text-accent hover:decoration-accent"
          href="https://www.dedao.cn"
          rel="noreferrer"
          target="_blank"
        >
          打开得到官网
        </a>
      </div>
    </div>
  )
}
