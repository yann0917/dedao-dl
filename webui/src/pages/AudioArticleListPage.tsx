import { FileText, Play, Rows4, Loader2 } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { api, type AudioGroupResponse } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { type PlayerTrack, useAudioPlayer } from "@/providers/AudioPlayerProvider"

function formatMinutes(value?: number) {
  if (!value) {
    return "0 分钟"
  }

  return `${Math.round(value / 60)} 分钟`
}

export function AudioArticleListPage() {
  const navigate = useNavigate()
  const { enid = "" } = useParams()
  const { setQueue } = useAudioPlayer()
  const [data, setData] = useState<AudioGroupResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      setError(null)
      try {
        const result = await api.audio.group(enid)
        if (!cancelled) {
          setData(result)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "听书文稿列表加载失败")
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

  const items = data?.outside?.items ?? []

  const tracks = useMemo<PlayerTrack[]>(
    () =>
      items
        .map((item) => ({
          id: item.extra.audio_alias_id || item.extra.enid,
          title: item.title,
          src: item.extra.odob_audio_detail?.mp3_play_url || "",
          poster: item.icon || item.extra.odob_audio_detail?.icon,
          subtitle: data?.outside?.spu.title || "听书文稿",
        }))
        .filter((item) => !!item.id && !!item.src),
    [data?.outside?.spu.title, items],
  )

  const handlePlay = (articleEnid: string) => {
    const index = tracks.findIndex((item) => item.id === articleEnid)
    if (index < 0) {
      return
    }
    setQueue(tracks, index)
  }

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-slate-500">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载文稿列表...
        </div>
      </main>
    )
  }

  if (error || !data?.outside?.spu) {
    return (
      <main className="space-y-6">
        <Card className="border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700">
          {error || "未找到听书文稿列表"}
        </Card>
      </main>
    )
  }

  return (
    <main className="space-y-6">
      <section className="rounded-3xl border border-white/70 bg-white/90 p-6 shadow-soft backdrop-blur">
        <p className="text-sm text-slate-500">听书文稿</p>
        <h2 className="mt-2 text-3xl font-semibold text-slate-950">{data.outside.spu.title}</h2>
        <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">
          这里承接听书合集里的文稿列表。每一项都可以直接播放对应音频，也可以进入单篇文稿详情。
        </p>
      </section>

      <div className="flex flex-wrap gap-3">
        <Button onClick={() => setQueue(tracks, 0)}>
          <Play className="mr-2 h-4 w-4" />
          播放整组
        </Button>
        <Button onClick={() => navigate(`/audio-groups/${encodeURIComponent(enid)}`)} variant="outline">
          返回合集详情
        </Button>
      </div>

      <section className="space-y-4">
        {items.map((item) => {
          const articleEnid = item.extra.audio_alias_id || item.extra.enid
          return (
            <Card className="p-5" key={articleEnid}>
              <div className="flex flex-col gap-4 md:flex-row md:items-start">
                <img
                  alt={item.title}
                  className="h-24 w-24 rounded-3xl object-cover"
                  src={item.icon || item.extra.odob_audio_detail?.icon || "https://placehold.co/200x200/e2e8f0/334155?text=Doc"}
                />
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="rounded-full bg-orange-100 px-3 py-1 text-xs font-medium text-orange-700">文稿条目</span>
                    <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
                      {formatMinutes(item.extra.duration)}
                    </span>
                  </div>
                  <h3 className="mt-3 text-xl font-semibold text-slate-950">{item.title}</h3>
                  <p className="mt-2 text-sm leading-7 text-slate-600">{item.summary || item.intro || "暂无摘要"}</p>
                </div>
                <div className="flex gap-3">
                  <Button onClick={() => handlePlay(articleEnid)} variant="outline">
                    <Play className="mr-2 h-4 w-4" />
                    播放
                  </Button>
                  <Button
                    onClick={() =>
                      navigate(
                        `/articles/2/${encodeURIComponent(articleEnid)}?from=audio-group&groupEnid=${encodeURIComponent(enid)}&parentTitle=${encodeURIComponent(
                          data.outside?.spu.title || "听书文稿",
                        )}`,
                      )
                    }
                  >
                    <FileText className="mr-2 h-4 w-4" />
                    查看文稿
                  </Button>
                </div>
              </div>
            </Card>
          )
        })}
      </section>

      {items.length === 0 ? (
        <Card className="p-10 text-center text-slate-500">
          <Rows4 className="mx-auto h-8 w-8" />
          <p className="mt-3">当前合集还没有可展示的文稿条目。</p>
        </Card>
      ) : null}
    </main>
  )
}
