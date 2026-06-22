import {
  ArrowLeft,
  BookMarked,
  Compass,
  FolderOpen,
  GraduationCap,
  Headphones,
  LibraryBig,
  Loader2,
} from "lucide-react"
import { useEffect, useMemo, useState, type KeyboardEvent, type ReactNode } from "react"
import { useNavigate } from "react-router-dom"
import { api, type CourseListItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"

type PurchasedCollectionConfig = {
  category: string
  title: string
  description: string
  loadingText: string
  emptyTitle: string
  emptyDescription: string
  itemLabel: string
  icon: "course" | "ebook" | "audio" | "compass"
  onOpenItem: (item: CourseListItem, navigate: ReturnType<typeof useNavigate>) => void
  getSummary?: (item: CourseListItem) => string
  getPrimaryMeta?: (item: CourseListItem) => string
  getSecondaryMeta?: (item: CourseListItem) => string
  getProgress?: (item: CourseListItem) => number | null
  renderActions?: (item: CourseListItem, helpers: { navigate: ReturnType<typeof useNavigate>; openItem: (item: CourseListItem) => void }) => ReactNode
}

function resolvePageIcon(icon: PurchasedCollectionConfig["icon"]) {
  if (icon === "course") {
    return GraduationCap
  }

  if (icon === "ebook") {
    return BookMarked
  }

  if (icon === "audio") {
    return Headphones
  }

  return Compass
}

function resolveFallbackIcon(icon: PurchasedCollectionConfig["icon"]) {
  if (icon === "course") {
    return GraduationCap
  }

  if (icon === "ebook") {
    return LibraryBig
  }

  if (icon === "audio") {
    return Headphones
  }

  return Compass
}

export function PurchasedCollectionPage(config: PurchasedCollectionConfig) {
  const navigate = useNavigate()
  const [items, setItems] = useState<CourseListItem[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [groupMode, setGroupMode] = useState({
    active: false,
    groupId: 0,
    title: "",
  })

  useEffect(() => {
    void loadItems(false, 1)
  }, [config.category, groupMode.active, groupMode.groupId])

  const hasMore = useMemo(() => items.length < total, [items.length, total])
  const PageIcon = resolvePageIcon(config.icon)
  const FallbackIcon = resolveFallbackIcon(config.icon)

  async function loadItems(append: boolean, targetPage: number) {
    if (append) {
      setLoadingMore(true)
    } else {
      setLoading(true)
    }
    setError(null)

    try {
      const params = new URLSearchParams({
        category: config.category,
        order: "study",
        page: String(targetPage),
        limit: "20",
      })

      if (groupMode.active && groupMode.groupId > 0) {
        params.set("groupId", String(groupMode.groupId))
      }

      const result = await api.course.list(params)
      setPage(targetPage)
      setTotal(result.total || 0)
      setItems((current) => (append ? [...current, ...result.list] : result.list))
    } catch (err) {
      setError(err instanceof Error ? err.message : `${config.title}加载失败`)
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }

  const enterGroup = (item: CourseListItem) => {
    const groupId = Number(item.group_id || item.id || 0)
    if (!groupId) {
      return
    }

    setItems([])
    setGroupMode({
      active: true,
      groupId,
      title: item.title || item.name || `${config.title}分组`,
    })
  }

  const exitGroup = () => {
    setItems([])
    setGroupMode({
      active: false,
      groupId: 0,
      title: "",
    })
  }

  const openItem = (item: CourseListItem) => {
    if (item.is_group) {
      enterGroup(item)
      return
    }

    config.onOpenItem(item, navigate)
  }

  const actionHelpers = {
    navigate,
    openItem,
  }

  const handleCardKeyDown = (event: KeyboardEvent<HTMLDivElement>, item: CourseListItem) => {
    if (event.key !== "Enter" && event.key !== " ") {
      return
    }

    event.preventDefault()
    openItem(item)
  }

  return (
    <main className="space-y-6">
      <section className="rounded-3xl border border-white/70 bg-white/90 p-6 shadow-soft backdrop-blur">
        <p className="text-sm text-slate-500">{config.title}</p>
        <h2 className="mt-2 text-3xl font-semibold text-slate-950">{groupMode.active ? groupMode.title : config.title}</h2>
        <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">{config.description}</p>
        <div className="mt-5 flex flex-wrap items-center gap-3 text-sm text-slate-500">
          <span className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1.5">
            <PageIcon className="h-4 w-4" />
            当前共 {total} 项
          </span>
          {groupMode.active ? (
            <Button onClick={exitGroup} variant="outline">
              <ArrowLeft className="mr-2 h-4 w-4" />
              返回{config.title}
            </Button>
          ) : null}
        </div>
      </section>

      {error ? (
        <Card className="border border-rose-200 bg-rose-50">
          <div className="p-4 text-sm text-rose-700">{error}</div>
        </Card>
      ) : null}

      {loading ? (
        <Card className="flex items-center justify-center p-8 text-slate-500">
          <Loader2 className="mr-2 h-5 w-5 animate-spin" />
          {config.loadingText}
        </Card>
      ) : null}

      {!loading ? (
        <section className="grid gap-5 md:grid-cols-2 xl:grid-cols-5">
          {items.map((item) => {
            const title = item.title || item.name || `${config.itemLabel}内容`
            const summary = config.getSummary?.(item) || item.intro || item.subtitle || "暂无简介"
            const primaryMeta = item.is_group ? `${item.course_num || 0} 项内容` : (config.getPrimaryMeta?.(item) || "查看内容")
            const secondaryMeta = item.is_group ? "进入分组" : (config.getSecondaryMeta?.(item) || "打开详情")
            const progress = item.is_group ? null : config.getProgress?.(item)

            return (
              <div
                className="text-left"
                key={`${item.group_id || item.id}-${item.enid || title}`}
                onClick={() => openItem(item)}
                onKeyDown={(event) => handleCardKeyDown(event, item)}
                role="button"
                tabIndex={0}
              >
                <Card className="h-full overflow-hidden transition hover:-translate-y-1 hover:shadow-lg">
                  <div className="relative aspect-square overflow-hidden bg-slate-100">
                    {item.icon || item.cover || item.index_img ? (
                      <img
                        alt={title}
                        className="h-full w-full object-cover transition duration-500 hover:scale-105"
                        src={item.icon || item.cover || item.index_img}
                      />
                    ) : (
                      <div className="flex h-full items-center justify-center text-slate-400">
                        {item.is_group ? <FolderOpen className="h-10 w-10" /> : <FallbackIcon className="h-10 w-10" />}
                      </div>
                    )}
                    <div className="absolute left-3 top-3">
                      <span
                        className={[
                          "rounded-full px-3 py-1 text-xs font-medium",
                          item.is_group ? "bg-slate-900/80 text-white" : "bg-white/90 text-slate-700",
                        ].join(" ")}
                      >
                        {item.is_group ? "内容分组" : config.itemLabel}
                      </span>
                    </div>
                  </div>

                  <div className="space-y-3 p-4">
                    <h3 className="line-clamp-2 text-base font-semibold text-slate-950">{title}</h3>
                    <p className="line-clamp-2 text-sm leading-6 text-slate-500">{summary}</p>
                    <div className="flex items-center justify-between text-xs text-slate-400">
                      <span>{primaryMeta}</span>
                      <span>{secondaryMeta}</span>
                    </div>
                    {typeof progress === "number" ? (
                      <div className="h-1.5 overflow-hidden rounded-full bg-slate-100">
                        <div
                          className="h-full rounded-full bg-primary transition-[width]"
                          style={{ width: `${Math.max(0, Math.min(progress, 100))}%` }}
                        />
                      </div>
                    ) : null}
                    {config.renderActions ? (
                      <div
                        className="flex flex-wrap gap-2 border-t border-slate-100 pt-1"
                        onClick={(event) => event.stopPropagation()}
                        onKeyDown={(event) => event.stopPropagation()}
                      >
                        {config.renderActions(item, actionHelpers)}
                      </div>
                    ) : null}
                  </div>
                </Card>
              </div>
            )
          })}
        </section>
      ) : null}

      {!loading && items.length === 0 ? (
        <Card className="p-10 text-center text-slate-500">
          <p className="text-lg font-medium text-slate-900">{config.emptyTitle}</p>
          <p className="mt-2 text-sm">{config.emptyDescription}</p>
        </Card>
      ) : null}

      {!loading && hasMore ? (
        <div className="flex justify-center">
          <Button disabled={loadingMore} onClick={() => void loadItems(true, page + 1)} variant="outline">
            {loadingMore ? "加载中..." : "加载更多"}
          </Button>
        </div>
      ) : null}
    </main>
  )
}
