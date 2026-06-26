import { useEffect, useMemo, useState } from "react"
import { Loader2, RefreshCcw, ScanLine } from "lucide-react"
import { toast } from "sonner"
import { api, type QRCodeSession, type UserInfo } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { getSemanticStatusBadgeClass, semanticMetaTextClass } from "@/lib/semanticStyles"

type QrLoginCardProps = {
  onLoginSuccess: (user?: UserInfo | null) => void | Promise<void>
}

export function QrLoginCard({ onLoginSuccess }: QrLoginCardProps) {
  const [session, setSession] = useState<QRCodeSession | null>(null)
  const [loading, setLoading] = useState(false)
  const [polling, setPolling] = useState(false)

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

  return (
    <Card>
      <div className="space-y-6 p-6">
        <div className="flex items-center justify-between gap-3">
          <div>
            <h2 className="text-xl font-semibold text-text-primary">扫码登录</h2>
            <p className={`mt-1 ${semanticMetaTextClass}`}>使用得到 App 或微信扫码登录。</p>
          </div>
          <span className={getSemanticStatusBadgeClass("neutral")}>{remaining ? `剩余 ${remaining}` : "待生成"}</span>
        </div>

        <div className="rounded-3xl border border-dashed border-border bg-surface-soft p-6">
          <div className="mx-auto flex h-64 w-64 items-center justify-center overflow-hidden rounded-3xl bg-surface-panel shadow-sm">
            {session?.qrCode ? (
              <img alt="二维码" className="h-full w-full object-cover" src={session.qrCode} />
            ) : (
              <div className="flex flex-col items-center gap-3 text-sm text-text-muted">
                <Loader2 className="h-6 w-6 animate-spin" />
                正在准备二维码
              </div>
            )}
          </div>
        </div>

        <div className="rounded-2xl border border-border bg-surface-soft p-4 text-text-primary">
          <div className="flex items-center gap-2 text-sm font-medium text-text-primary">
            <ScanLine className="h-4 w-4 text-accent" />
            得到 App / 微信扫码
          </div>
          <p className={`mt-2 text-sm ${semanticMetaTextClass}`}>使用得到 App 或微信扫码，手机确认后会自动进入工作台。</p>
        </div>

        <div className="flex gap-3">
          <Button className="flex-1" disabled={loading} onClick={() => void loadQRCode()}>
            {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <RefreshCcw className="mr-2 h-4 w-4" />}
            刷新二维码
          </Button>
          <Button className="flex-1" onClick={() => window.open("https://www.dedao.cn", "_blank")?.focus()} variant="outline">
            打开官网
          </Button>
        </div>
      </div>
    </Card>
  )
}
