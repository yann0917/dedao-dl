import { Clock3, Headphones, Loader2, Radio, Waves } from "lucide-react"
import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { toast } from "sonner"
import { api, type AudioDetailResponse } from "@/api"
import { DownloadActionsPanel, type DownloadOption } from "@/components/download/DownloadActionsPanel"
import { InfoBlock, StatCard, StatusBadge } from "@/components/ui/Semantic"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  semanticMetaTextClass,
  semanticPageSectionClass,
  semanticSecondaryTextClass,
} from "@/lib/semanticStyles"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

const audioDownloadOptions: DownloadOption[] = [
  { value: 1, label: "下载音频 MP3" },
  { value: 2, label: "下载文稿 PDF" },
  { value: 3, label: "下载文稿 Markdown" },
]

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
    try {
      await api.audio.addToShelf([enid])
      setData((current) => (current ? { ...current, book_shelf_status: 1 } : current))
      toast.success("已加入书架", {
        description: data.title || data.package_title || "当前听书已加入书架",
      })
    } catch (err) {
      toast.error("加入书架失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setShelfLoading(false)
    }
  }

  const handleAudioDownload = (downloadType: number) =>
    api.download.audio({
      enid,
      title: data.title || data.package_title || "听书下载",
      downloadType,
    })

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
            <StatusBadge variant="accent">每天听本书</StatusBadge>
            {data.is_vip ? (
              <StatusBadge variant="warning">会员内容</StatusBadge>
            ) : null}
            {data.has_play_auth ? (
              <StatusBadge variant="success">可播放</StatusBadge>
            ) : null}
            {isInBookshelf ? (
              <StatusBadge>已加入书架</StatusBadge>
            ) : null}
          </div>

          <div className="mt-6 space-y-3">
            <InfoBlock>演播：{data.reader_name || "未知"}</InfoBlock>
            <InfoBlock>时长：{formatDuration(data.duration)}</InfoBlock>
            <InfoBlock>来源：{data.source_name || "得到"}</InfoBlock>
            <InfoBlock>播放次数：{data.play_count || 0}</InfoBlock>
          </div>
        </Card>

        <div className="space-y-6">
          <DownloadActionsPanel
            description="将当前听书内容直接下载到服务端本地 output 目录，支持音频与文稿两类导出。"
            disabled={!data.has_play_auth}
            disabledReason="当前账号暂无播放权限，暂不支持直接下载。"
            onDownload={handleAudioDownload}
            options={audioDownloadOptions}
            title="听书下载"
          />

          <Card className="p-6">
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <StatCard icon={<Clock3 className="h-4 w-4" />} label="时长" value={formatDuration(data.duration)} />
              <StatCard icon={<Headphones className="h-4 w-4" />} label="播放次数" value={data.play_count || 0} />
              <StatCard icon={<Radio className="h-4 w-4" />} label="权限状态" value={data.has_play_auth ? "可播放" : "需权限"} />
              <StatCard
                icon={<Waves className="h-4 w-4" />}
                label="音频地址"
                value={data.mp3_play_url ? "已获取" : "未获取"}
                valueClassName="text-lg"
              />
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">内容简介</h3>
            <p className={`mt-4 whitespace-pre-wrap text-sm leading-7 ${semanticSecondaryTextClass}`}>
              {data.summary || data.slogan || data.update_tips || "暂无简介"}
            </p>
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
