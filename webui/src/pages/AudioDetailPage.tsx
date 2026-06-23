import { Clock3, Headphones, Loader2, Radio, Waves } from "lucide-react"
import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { api, type AudioDetailResponse } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  getSemanticStatusBadgeClass,
  semanticInfoBlockClass,
  semanticMetaTextClass,
  semanticPageSectionClass,
  semanticSecondaryTextClass,
  semanticStatCardClass,
} from "@/lib/semanticStyles"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

function formatDuration(value?: number) {
  if (!value) {
    return "0 分钟"
  }

  const minutes = Math.round(value / 60)
  return `${minutes} 分钟`
}

export function AudioDetailPage() {
  const navigate = useNavigate()
  const { enid = "" } = useParams()
  const { setQueue } = useAudioPlayer()
  const [data, setData] = useState<AudioDetailResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [shelfLoading, setShelfLoading] = useState(false)
  const [shelfError, setShelfError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      setError(null)
      try {
        const result = await api.audio.detail(enid)
        if (!cancelled) {
          setData(result)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "听书详情加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    if (enid) {
      void load()
    }

    return () => {
      cancelled = true
    }
  }, [enid])

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载听书详情...
        </div>
      </main>
    )
  }

  if (error || !data) {
    return (
      <main className="space-y-6">
        <Card className="border-danger bg-danger-soft p-6 text-sm text-danger">
          {error || "未找到听书详情"}
        </Card>
      </main>
    )
  }

  const handlePlay = () => {
    if (!data.mp3_play_url) {
      return
    }

    setQueue(
      [
        {
          id: data.alias_id || enid,
          title: data.title || data.package_title || "每天听本书",
          src: data.mp3_play_url,
          poster: data.index_img || data.icon,
          subtitle: data.source_name || "每天听本书",
        },
      ],
      0,
    )
  }

  const handleAddToShelf = async () => {
    setShelfLoading(true)
    setShelfError(null)
    try {
      await api.audio.addToShelf([enid])
      setData((current) => (current ? { ...current, book_shelf_status: 1 } : current))
    } catch (err) {
      setShelfError(err instanceof Error ? err.message : "加入书架失败")
    } finally {
      setShelfLoading(false)
    }
  }

  const isInBookshelf = (data.book_shelf_status ?? 0) > 0

  return (
    <main className="space-y-6">
      <section className={`${semanticPageSectionClass} p-6`}>
        <p className={semanticMetaTextClass}>听书详情</p>
        <h2 className="mt-2 text-3xl font-semibold text-text-primary">{data.title || data.package_title || "每天听本书"}</h2>
      </section>

      <section className="grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)]">
        <Card className="p-6">
          <img
            alt={data.title}
            className="mx-auto aspect-square w-full max-w-[260px] rounded-3xl object-cover shadow-lg"
            src={data.index_img || data.icon || "https://placehold.co/600x600/e2e8f0/334155?text=Audio"}
          />
          <div className="mt-6 flex flex-wrap gap-2">
            <span className={getSemanticStatusBadgeClass("accent")}>每天听本书</span>
            {data.is_vip ? (
              <span className={getSemanticStatusBadgeClass("warning")}>会员内容</span>
            ) : null}
            {data.has_play_auth ? (
              <span className={getSemanticStatusBadgeClass("success")}>可播放</span>
            ) : null}
            {isInBookshelf ? (
              <span className={getSemanticStatusBadgeClass("neutral")}>已加入书架</span>
            ) : null}
          </div>

          <div className="mt-6 space-y-3">
            <div className={semanticInfoBlockClass}>演播：{data.reader_name || "未知"}</div>
            <div className={semanticInfoBlockClass}>时长：{formatDuration(data.duration)}</div>
            <div className={semanticInfoBlockClass}>来源：{data.source_name || "得到"}</div>
            <div className={semanticInfoBlockClass}>播放次数：{data.play_count || 0}</div>
          </div>
        </Card>

        <div className="space-y-6">
          <Card className="p-6">
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <div className={semanticStatCardClass}>
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <Clock3 className="h-4 w-4" />
                  时长
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">{formatDuration(data.duration)}</p>
              </div>
              <div className={semanticStatCardClass}>
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <Headphones className="h-4 w-4" />
                  播放次数
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">{data.play_count || 0}</p>
              </div>
              <div className={semanticStatCardClass}>
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <Radio className="h-4 w-4" />
                  权限状态
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">{data.has_play_auth ? "可播放" : "需权限"}</p>
              </div>
              <div className={semanticStatCardClass}>
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <Waves className="h-4 w-4" />
                  音频地址
                </div>
                <p className="mt-3 text-lg font-semibold text-text-primary">{data.mp3_play_url ? "已获取" : "未获取"}</p>
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">内容简介</h3>
            <p className={`mt-4 whitespace-pre-wrap text-sm leading-7 ${semanticSecondaryTextClass}`}>
              {data.summary || data.slogan || data.update_tips || "暂无简介"}
            </p>
            {shelfError ? (
              <div className="mt-5 rounded-2xl border border-danger bg-danger-soft p-4 text-sm text-danger">{shelfError}</div>
            ) : null}
            <div className="mt-5 flex flex-wrap gap-3">
              {!isInBookshelf ? (
                <Button disabled={shelfLoading} onClick={handleAddToShelf}>
                  {shelfLoading ? "处理中..." : "加入书架"}
                </Button>
              ) : null}
              {data.mp3_play_url ? <Button onClick={handlePlay}>开始播放</Button> : null}
              {data.alias_id ? (
                <Button
                  onClick={() =>
                    navigate(
                      `/articles/2/${encodeURIComponent(data.alias_id)}?from=audio&parentEnid=${encodeURIComponent(enid)}&parentTitle=${encodeURIComponent(
                        data.title || data.package_title || "听书详情",
                      )}`,
                    )
                  }
                  variant="outline"
                >
                  查看文稿
                </Button>
              ) : null}
              {data.share_url ? (
                <Button onClick={() => window.open(data.share_url, "_blank", "noopener,noreferrer")} variant="outline">
                  分享链接
                </Button>
              ) : null}
            </div>
          </Card>

          {(data.trial_listen_tips || data.update_tips) ? (
            <Card className="p-6">
              <h3 className="text-xl font-semibold text-text-primary">补充说明</h3>
              <div className={`mt-4 space-y-3 text-sm leading-7 ${semanticSecondaryTextClass}`}>
                {data.trial_listen_tips ? <p>{data.trial_listen_tips}</p> : null}
                {data.update_tips ? <p>{data.update_tips}</p> : null}
              </div>
            </Card>
          ) : null}
        </div>
      </section>
    </main>
  )
}
