import { Library, Loader2, Radio, Users } from "lucide-react"
import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { api, type AudioGroupResponse } from "@/api"
import { InfoBlock, StatCard, StatusBadge } from "@/components/ui/Semantic"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  semanticMetaTextClass,
  semanticPageSectionClass,
  semanticSecondaryTextClass,
} from "@/lib/semanticStyles"
import { type PlayerTrack, useAudioPlayer } from "@/providers/AudioPlayerProvider"

function formatMinutes(value?: number) {
  if (!value) {
    return "0 分钟"
  }

  return `${Math.round(value / 60)} 分钟`
}

export function AudioGroupPage() {
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
          setError(err instanceof Error ? err.message : "听书合集加载失败")
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
          正在加载听书合集...
        </div>
      </main>
    )
  }

  if (error || !data?.outside?.spu) {
    return (
      <main className="space-y-6">
        <Card className="border-danger bg-danger-soft p-6 text-sm text-danger">
          {error || "未找到听书合集详情"}
        </Card>
      </main>
    )
  }

  const spu = data.outside.spu
  const audioList = data.group?.odob_audio_detail_list ?? []
  const tracks: PlayerTrack[] = audioList
    .map((item) => ({
      id: item.alias_id || item.audio_id,
      title: item.title,
      src: item.mp3_play_url || "",
      poster: item.index_img || item.icon,
      subtitle: spu.title,
    }))
    .filter((item) => !!item.id && !!item.src)

  return (
    <main className="space-y-6">
      <section className={`${semanticPageSectionClass} p-6`}>
        <p className={semanticMetaTextClass}>听书合集</p>
        <h2 className="mt-2 text-3xl font-semibold text-text-primary">{spu.title}</h2>
      </section>

      <section className="grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)]">
        <Card className="p-6">
          <img
            alt={spu.title}
            className="mx-auto aspect-square w-full max-w-[260px] rounded-3xl object-cover shadow-lg"
            src={spu.icon || "https://placehold.co/600x600/e2e8f0/334155?text=Group"}
          />
          <div className="mt-6 flex flex-wrap gap-2">
            <StatusBadge variant="accent">名家讲书</StatusBadge>
            {spu.extra?.teacher_name ? (
              <StatusBadge>{spu.extra.teacher_name}</StatusBadge>
            ) : null}
          </div>

          <div className="mt-6 space-y-3">
            <InfoBlock>老师：{spu.extra?.teacher_name || "未知"}</InfoBlock>
            <InfoBlock>学习人数：{spu.extra?.odob_consumer_num || 0}</InfoBlock>
            <InfoBlock>学习描述：{spu.extra?.rn_learn_count_desc || "暂无"}</InfoBlock>
          </div>
        </Card>

        <div className="space-y-6">
          <Card className="p-6">
            <div className="grid gap-4 md:grid-cols-3">
              <StatCard icon={<Library className="h-4 w-4" />} label="合集条目" value={data.outside.count || audioList.length} />
              <StatCard icon={<Users className="h-4 w-4" />} label="学习人数" value={spu.extra?.odob_consumer_num || 0} />
              <StatCard icon={<Radio className="h-4 w-4" />} label="音频解析" value={audioList.length > 0 ? "已就绪" : "待补齐"} />
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">合集简介</h3>
            <p className={`mt-4 whitespace-pre-wrap text-sm leading-7 ${semanticSecondaryTextClass}`}>
              {spu.intro || spu.extra?.intro_text || spu.summary || "暂无简介"}
            </p>
            <div className="mt-5 flex flex-wrap gap-3">
              {tracks.length > 0 ? <Button onClick={() => setQueue(tracks, 0)}>播放整组</Button> : null}
              <Button onClick={() => navigate(`/audio-groups/${encodeURIComponent(enid)}/articles`)} variant="outline">
                查看文稿列表
              </Button>
            </div>
          </Card>

          {data.groupError ? (
            <Card className="border-warning bg-warning-soft p-4 text-sm text-warning">{data.groupError}</Card>
          ) : null}

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">音频列表</h3>
            <div className="mt-4 space-y-3">
              {audioList.length > 0 ? (
                audioList.map((item) => (
                  <InfoBlock key={item.alias_id || item.audio_id}>
                    <div className="flex items-start gap-4">
                      <img
                        alt={item.title}
                        className="h-16 w-16 rounded-2xl object-cover"
                        src={item.index_img || item.icon || "https://placehold.co/200x200/e2e8f0/334155?text=Audio"}
                      />
                      <div className="min-w-0 flex-1">
                        <p className="font-medium text-text-primary">{item.title}</p>
                        <p className="mt-1 line-clamp-2 text-sm leading-6 text-text-secondary">{item.summary || item.slogan || "暂无简介"}</p>
                        <p className="mt-2 text-xs text-text-muted">时长：{formatMinutes(item.duration)}</p>
                      </div>
                      {tracks.some((track) => track.id === (item.alias_id || item.audio_id)) ? (
                        <Button
                          onClick={() => {
                            const index = tracks.findIndex((track) => track.id === (item.alias_id || item.audio_id))
                            if (index >= 0) {
                              setQueue(tracks, index)
                            }
                          }}
                          variant="outline"
                        >
                          播放
                        </Button>
                      ) : null}
                    </div>
                  </InfoBlock>
                ))
              ) : (
                <div className="rounded-2xl border border-dashed border-border p-6 text-sm text-text-muted">
                  当前还没有可展示的音频列表。
                </div>
              )}
            </div>
          </Card>
        </div>
      </section>
    </main>
  )
}
