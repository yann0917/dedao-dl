import { BookOpen, Loader2 } from "lucide-react"
import { useEffect, useState } from "react"
import { api, type RecentItem, type RecentStats } from "@/api"
import { StudyStatsPanel } from "@/components/study/StudyStatsPanel"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { getSemanticStatusBadgeClass, semanticMetaTextClass, semanticPageSectionClass } from "@/lib/semanticStyles"

const PAGE_SIZE = 20

// 上游多个文本字段（title/last_info/intro）内嵌 <hl>、<h1> 等高亮标签，展示前统一剥掉
function stripHtml(raw: string) {
  return raw
    .replace(/<[^>]*>/g, " ")
    .replace(/\s+/g, " ")
    .trim()
}

// 记录时间为毫秒时间戳；last_info 为空时（听书条目常见）直接格式化时间戳
function formatStudyTime(item: RecentItem) {
  if (item.last_info) {
    return stripHtml(item.last_info)
  }
  if (!item.timestamp) {
    return ""
  }
  return new Date(item.timestamp).toLocaleString("zh-CN")
}

// 进度展示优先级与 CLI recent 一致：max_progress 百分比 > intro 文本 > 原始 progress 数值
// 特例：上游 recent 接口对课程类条目的 progress 恒为 0，真实学习进度百分比放在 max_progress
// （已与课程列表接口的 progress 字段交叉验证：76/50/100 三个值完全吻合）
function resolveProgress(item: RecentItem): number | string {
  const { progress, max_progress: maxProgress, intro } = item.progress_intro || {}
  if (item.type_name === "课程" && progress === 0 && maxProgress > 0) {
    return maxProgress
  }
  if (maxProgress > 0) {
    return Math.round((progress / maxProgress) * 100)
  }
  if (intro) {
    return stripHtml(intro)
  }
  return progress
}

function RecordRow({ item }: { item: RecentItem }) {
  const cover = item.square_img || item.index_img
  const title = stripHtml(item.title)
  const progress = resolveProgress(item)
  const studyTime = formatStudyTime(item)

  return (
    <Card className="p-0 shadow-soft">
      <div className="flex items-center gap-4 p-4">
        <div className="flex h-14 w-14 shrink-0 items-center justify-center overflow-hidden rounded-md bg-surface-soft">
          {cover ? (
            <img alt={title} className="h-full w-full object-cover" decoding="async" loading="lazy" src={cover} />
          ) : (
            <BookOpen className="size-icon-lg text-text-muted" />
          )}
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <span className={getSemanticStatusBadgeClass("neutral", "px-2 py-0.5 text-[11px]")}>{item.type_name || "内容"}</span>
            <h3 className="truncate text-base font-semibold text-text-primary">{title}</h3>
          </div>
          <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-text-muted">
            {item.author ? <span className="truncate">{item.author}</span> : null}
            {typeof progress === "number" ? (
              <span className="flex items-center gap-2">
                <span className="h-1.5 w-24 overflow-hidden rounded-full bg-surface-soft">
                  <span
                    className="block h-full rounded-full bg-primary transition-[width]"
                    style={{ width: `${Math.max(0, Math.min(progress, 100))}%` }}
                  />
                </span>
                {progress}%
              </span>
            ) : (
              progress
            )}
            {studyTime ? <span>{studyTime}</span> : null}
          </div>
        </div>
      </div>
    </Card>
  )
}

export function StudyRecordPage() {
  const [items, setItems] = useState<RecentItem[]>([])
  const [hasMore, setHasMore] = useState(false)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [stats, setStats] = useState<RecentStats | null>(null)
  const [statsError, setStatsError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    api.recent
      .stats()
      .then((result) => {
        if (!cancelled) {
          setStats(result)
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setStatsError(err instanceof Error ? err.message : "学习统计加载失败")
        }
      })

    return () => {
      cancelled = true
    }
  }, [])

  const load = async (append: boolean) => {
    if (append) {
      setLoadingMore(true)
    } else {
      setLoading(true)
    }
    setError(null)

    try {
      // 游标翻页：下一页 maxId 传已加载最后一条的 timestamp（毫秒）
      const cursor = append && items.length > 0 ? items[items.length - 1].timestamp : 0
      const result = await api.recent.list(cursor, PAGE_SIZE)
      setItems((current) => (append ? [...current, ...result.list] : result.list))
      setHasMore(result.has_more)
    } catch (err) {
      setError(err instanceof Error ? err.message : "学习记录加载失败")
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }

  useEffect(() => {
    void load(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="size-icon-lg animate-spin" />
          正在加载学习记录...
        </div>
      </main>
    )
  }

  return (
    <main className="space-y-6">
      <section className={`${semanticPageSectionClass} p-6`}>
        <p className={semanticMetaTextClass}>最近学习</p>
        <h2 className="mt-2 text-3xl font-semibold text-text-primary">学习记录</h2>
        <p className="mt-2 text-sm text-text-muted">当前登录账号的最近学习内容，按学习时间倒序。</p>
      </section>

      {statsError ? (
        <p className={`text-sm text-text-muted`}>统计加载失败：{statsError}</p>
      ) : stats && stats.total > 0 ? (
        <StudyStatsPanel stats={stats} />
      ) : null}

      <p className={semanticMetaTextClass}>全部记录</p>

      {error ? <Card className="border-danger bg-danger-soft p-6 text-sm text-danger">{error}</Card> : null}

      {items.length ? (
        <section className="space-y-3">
          {items.map((item) => (
            <RecordRow item={item} key={`${item.product_type}-${item.product_id}-${item.timestamp}`} />
          ))}
        </section>
      ) : (
        !error && (
          <Card className="p-10 text-center text-text-muted">
            <p className="text-lg font-medium text-text-primary">还没有学习记录</p>
            <p className="mt-2 text-sm">在得到 App 或网页版学习后，这里会显示最近的学习内容。</p>
          </Card>
        )
      )}

      {!error && hasMore ? (
        <div className="flex justify-center">
          <Button disabled={loadingMore} onClick={() => void load(true)} variant="outline">
            {loadingMore ? "加载中..." : "加载更多"}
          </Button>
        </div>
      ) : null}
    </main>
  )
}
