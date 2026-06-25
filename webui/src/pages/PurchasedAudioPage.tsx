import { FileText, Loader2, Play, Rows3 } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"
import { api, type CourseListItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"
import { getSemanticStatusBadgeClass } from "@/lib/semanticStyles"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

function formatMinutes(duration?: number) {
  if (!duration) {
    return "0 分钟"
  }

  return `${Math.round(duration / 60)} 分钟`
}

export function PurchasedAudioPage() {
  const { setQueue } = useAudioPlayer()
  const [playingEnid, setPlayingEnid] = useState<string | null>(null)

  const playAudio = async (item: CourseListItem) => {
    if (!item.enid) {
      return
    }

    setPlayingEnid(item.enid)

    try {
      // 播放权限和授权后的音频地址都依赖详情接口，列表数据本身不可靠。
      const detail = await api.audio.detail(item.enid)
      if (!detail.has_play_auth || !detail.mp3_play_url) {
        toast.error("当前内容暂不可播放", {
          description: detail.trial_listen_tips || detail.update_tips || "当前账号尚未获得播放授权",
        })
        return
      }

      setQueue(
        [
          {
            id: detail.alias_id || item.enid || String(item.id),
            title: detail.title || detail.package_title || item.title || item.name || "每天听本书",
            src: detail.mp3_play_url,
            poster: detail.index_img || detail.icon || item.icon || item.cover,
            subtitle: detail.source_name || "每天听本书",
          },
        ],
        0,
      )
    } catch (err) {
      toast.error("播放失败", {
        description: err instanceof Error ? err.message : "听书详情获取失败，请稍后重试",
      })
    } finally {
      setPlayingEnid((current) => (current === item.enid ? null : current))
    }
  }

  return (
    <PurchasedCollectionPage
      category="odob"
      coverContainerClassName="aspect-square p-3"
      coverImageClassName="object-contain"
      emptyTitle="当前没有可展示的已购听书"
      getPrimaryMeta={(item: CourseListItem) => formatMinutes(item.duration)}
      getSecondaryMeta={() => undefined}
      icon="audio"
      itemLabel="听书"
      loadingText="正在加载已购听书..."
      renderActions={(item, helpers) =>
        item.is_group || item.type === 1013 ? (
          <>
            <Button className="h-9 px-3" onClick={() => helpers.openItem(item)} variant="outline">
              <Rows3 className="mr-2 h-4 w-4" />
              查看合集
            </Button>
            {item.enid ? (
              <Button
                className="h-9 px-3"
                onClick={() => helpers.navigate(`/audio-groups/${encodeURIComponent(item.enid)}/articles`)}
                variant="ghost"
              >
                <FileText className="mr-2 h-4 w-4" />
                文稿列表
              </Button>
            ) : null}
          </>
        ) : (
          <>
            {item.in_bookrack ? <span className={getSemanticStatusBadgeClass("neutral", "inline-flex h-9 items-center rounded-xl px-3 text-sm")}>已加入书架</span> : null}
            <Button
              className="h-9 px-3"
              disabled={!item.enid || playingEnid === item.enid}
              onClick={() => void playAudio(item)}
            >
              {playingEnid === item.enid ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Play className="mr-2 h-4 w-4" />}
              {playingEnid === item.enid ? "获取中..." : "播放"}
            </Button>
            {item.audio_detail?.alias_id ? (
              <Button
                className="h-9 px-3"
                onClick={() =>
                  helpers.navigate(
                    `/articles/2/${encodeURIComponent(item.audio_detail?.alias_id || "")}?from=audio&parentEnid=${encodeURIComponent(item.enid)}&parentTitle=${encodeURIComponent(item.title || item.name || "听书详情")}`,
                  )
                }
                variant="outline"
              >
                <FileText className="mr-2 h-4 w-4" />
                查看文稿
              </Button>
            ) : null}
          </>
        )
      }
      onOpenItem={(item, navigate) => {
        if (!item.enid) {
          return
        }

        if (item.type === 1013) {
          navigate(`/audio-groups/${encodeURIComponent(item.enid)}?from=purchased-audio&parentTitle=${encodeURIComponent(item.title || item.name || "听书合集")}`)
          return
        }

        navigate(`/audios/${encodeURIComponent(item.enid)}?from=purchased-audio&parentTitle=${encodeURIComponent(item.title || item.name || "听书详情")}`)
      }}
      title="已购听书"
    />
  )
}
