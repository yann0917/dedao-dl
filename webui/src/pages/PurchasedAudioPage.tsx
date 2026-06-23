import { FileText, PanelRightOpen, Play, Rows3 } from "lucide-react"
import { useState } from "react"
import { api, type CourseListItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

function formatMinutes(duration?: number) {
  if (!duration) {
    return "0 分钟"
  }

  return `${Math.round(duration / 60)} 分钟`
}

export function PurchasedAudioPage() {
  const { setQueue } = useAudioPlayer()
  const [actionError, setActionError] = useState<string | null>(null)

  const playAudio = (item: CourseListItem) => {
    const src = item.audio_detail?.mp3_play_url || item.odob_group_ext_info?.audio_detail?.mp3_play_url || ""
    if (!src) {
      return
    }

    setQueue(
      [
        {
          id: item.audio_detail?.alias_id || item.odob_group_ext_info?.audio_detail?.alias_id || item.enid || String(item.id),
          title: item.title || item.name || "每天听本书",
          src,
          poster: item.audio_detail?.icon || item.odob_group_ext_info?.audio_detail?.icon || item.icon || item.cover,
          subtitle: "已购听书",
        },
      ],
      0,
    )
  }

  return (
    <PurchasedCollectionPage
      category="odob"
      emptyTitle="当前没有可展示的已购听书"
      externalError={actionError}
      getPrimaryMeta={(item: CourseListItem) => formatMinutes(item.duration)}
      getSecondaryMeta={(item: CourseListItem) => (item.type === 1013 ? "打开合集" : "打开详情")}
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
            {item.in_bookrack ? (
              <span className="inline-flex h-9 items-center rounded-xl bg-slate-100 px-3 text-sm font-medium text-slate-600">
                已加入书架
              </span>
            ) : (
              <Button
                className="h-9 px-3"
                onClick={async () => {
                  if (!item.enid) {
                    return
                  }

                  setActionError(null)
                  try {
                    await api.audio.addToShelf([item.enid])
                    helpers.updateItem(item, { in_bookrack: true })
                  } catch (err) {
                    setActionError(err instanceof Error ? err.message : "听书加入书架失败")
                  }
                }}
              >
                加入书架
              </Button>
            )}
            <Button
              className="h-9 px-3"
              disabled={!item.audio_detail?.mp3_play_url && !item.odob_group_ext_info?.audio_detail?.mp3_play_url}
              onClick={() => playAudio(item)}
            >
              <Play className="mr-2 h-4 w-4" />
              播放
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
            <Button className="h-9 px-3" onClick={() => helpers.openItem(item)} variant="ghost">
              <PanelRightOpen className="mr-2 h-4 w-4" />
              详情
            </Button>
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
