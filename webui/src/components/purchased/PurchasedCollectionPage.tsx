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
import { api, type CourseListItem, type PurchasedNavbarChild } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  getSemanticChipClass,
  getSemanticStatusBadgeClass,
  semanticMetaTextClass,
  semanticPageSectionClass,
  semanticSecondaryTextClass,
} from "@/lib/semanticStyles"

type PurchasedCollectionConfig = {
  category: string
  title: string
  description?: string
  loadingText: string
  emptyTitle: string
  emptyDescription?: string
  externalError?: string | null
  itemLabel: string
  icon: "course" | "ebook" | "audio" | "compass"
  onOpenItem: (item: CourseListItem, navigate: ReturnType<typeof useNavigate>) => void
  getSummary?: (item: CourseListItem) => string
  getPrimaryMeta?: (item: CourseListItem) => string
  getSecondaryMeta?: (item: CourseListItem) => string | null | undefined
  getProgress?: (item: CourseListItem) => number | null
  coverContainerClassName?: string
  coverImageClassName?: string
  renderActions?: (
    item: CourseListItem,
    helpers: {
      navigate: ReturnType<typeof useNavigate>
      openItem: (item: CourseListItem) => void
      updateItem: (item: CourseListItem, patch: Partial<CourseListItem>) => void
    },
  ) => ReactNode
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

function ensureAllFilter(options: PurchasedNavbarChild[]) {
  if (options.some((item) => item.filter === "all")) {
    return options
  }

  return [{ name: "全部", count: 0, filter: "all", show_count: false }, ...options]
}

export function PurchasedCollectionPage(config: PurchasedCollectionConfig) {
  const navigate = useNavigate()
  const [items, setItems] = useState<CourseListItem[]>([])
  const [filterOptions, setFilterOptions] = useState<PurchasedNavbarChild[]>([])
  const [currentFilter, setCurrentFilter] = useState("all")
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
  }, [config.category, currentFilter, groupMode.active, groupMode.groupId])

  useEffect(() => {
    let cancelled = false

    const loadNavbar = async () => {
      try {
        const data = await api.course.navbar()
        if (cancelled) {
          return
        }

        const matched = data.list.find((item) => item.category === config.category)
        setFilterOptions(ensureAllFilter(matched?.children ?? []))
      } catch {
        if (!cancelled) {
          setFilterOptions([])
        }
      }
    }

    setCurrentFilter("all")
    setGroupMode({
      active: false,
      groupId: 0,
      title: "",
    })
    void loadNavbar()

    return () => {
      cancelled = true
    }
  }, [config.category])

  const hasMore = useMemo(() => items.length < total, [items.length, total])
  const PageIcon = resolvePageIcon(config.icon)
  const FallbackIcon = resolveFallbackIcon(config.icon)
  const coverContainerClassName = config.coverContainerClassName || "aspect-square"
  const coverImageClassName = config.coverImageClassName || "object-cover"

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
        filter: currentFilter,
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

  const handleFilterChange = (filter: string) => {
    if (filter === currentFilter) {
      return
    }

    setItems([])
    setPage(1)
    setCurrentFilter(filter)
  }

  const openItem = (item: CourseListItem) => {
    if (item.is_group) {
      enterGroup(item)
      return
    }

    config.onOpenItem(item, navigate)
  }

  const updateItem = (target: CourseListItem, patch: Partial<CourseListItem>) => {
    setItems((current) =>
      current.map((item) => {
        const sameById = item.id > 0 && target.id > 0 && item.id === target.id
        const sameByEnid = item.enid && target.enid && item.enid === target.enid
        if (!sameById && !sameByEnid) {
          return item
        }
        return {
          ...item,
          ...patch,
        }
      }),
    )
  }

  const actionHelpers = {
    navigate,
    openItem,
    updateItem,
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
      <section className={`${semanticPageSectionClass} p-6`}>
        <p className={semanticMetaTextClass}>{config.title}</p>
        <h2 className="mt-2 text-3xl font-semibold text-text-primary">{groupMode.active ? groupMode.title : config.title}</h2>
        {config.description ? <p className={`mt-3 max-w-3xl text-sm leading-7 ${semanticSecondaryTextClass}`}>{config.description}</p> : null}
        <div className="mt-5 flex flex-wrap items-center gap-3 text-sm text-text-muted">
          <span className={getSemanticStatusBadgeClass("neutral", "inline-flex items-center gap-2 px-3 py-1.5 text-sm")}>
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

      {!groupMode.active && filterOptions.length > 0 ? (
        <Card className="p-4">
          <div className="flex flex-wrap gap-2">
            {filterOptions.map((item) => {
              const active = currentFilter === item.filter

              return (
                <button
                  className={getSemanticChipClass(active)}
                  key={item.filter}
                  onClick={() => handleFilterChange(item.filter)}
                  type="button"
                >
                  {item.name}
                  {item.show_count ? ` (${item.count})` : ""}
                </button>
              )
            })}
          </div>
        </Card>
      ) : null}

      {error ? (
        <Card className="border-danger bg-danger-soft">
          <div className="p-4 text-sm text-danger">{error}</div>
        </Card>
      ) : null}

      {config.externalError ? (
        <Card className="border-danger bg-danger-soft">
          <div className="p-4 text-sm text-danger">{config.externalError}</div>
        </Card>
      ) : null}

      {loading ? (
        <Card className="flex items-center justify-center p-8 text-text-muted">
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
            const secondaryMeta = item.is_group ? "进入分组" : (config.getSecondaryMeta?.(item) ?? "打开详情")
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
                  <div className={`relative overflow-hidden bg-surface-soft ${coverContainerClassName}`}>
                    {item.icon || item.cover || item.index_img ? (
                      <img
                        alt={title}
                        className={`h-full w-full transition duration-500 hover:scale-105 ${coverImageClassName}`}
                        src={item.icon || item.cover || item.index_img}
                      />
                    ) : (
                      <div className="flex h-full items-center justify-center text-text-muted">
                        {item.is_group ? <FolderOpen className="h-10 w-10" /> : <FallbackIcon className="h-10 w-10" />}
                      </div>
                    )}
                    <div className="absolute left-3 top-3">
                      <span
                        className={[
                          "rounded-full px-3 py-1 text-xs font-medium",
                          item.is_group ? "bg-secondary/85 text-secondary-foreground" : "bg-surface-panel/90 text-text-secondary",
                        ].join(" ")}
                      >
                        {item.is_group ? "内容分组" : config.itemLabel}
                      </span>
                    </div>
                  </div>

                  <div className="space-y-3 p-4">
                    <h3 className="line-clamp-2 text-base font-semibold text-text-primary">{title}</h3>
                    <p className="line-clamp-2 text-sm leading-6 text-text-secondary">{summary}</p>
                    <div className="flex items-center justify-between gap-3 text-xs text-text-muted">
                      <span>{primaryMeta}</span>
                      {secondaryMeta ? <span>{secondaryMeta}</span> : null}
                    </div>
                    {typeof progress === "number" ? (
                      <div className="h-1.5 overflow-hidden rounded-full bg-surface-soft">
                        <div
                          className="h-full rounded-full bg-primary transition-[width]"
                          style={{ width: `${Math.max(0, Math.min(progress, 100))}%` }}
                        />
                      </div>
                    ) : null}
                    {config.renderActions ? (
                      <div
                        className="flex flex-wrap items-center gap-2 border-t border-border pt-1"
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
        <Card className="p-10 text-center text-text-muted">
          <p className="text-lg font-medium text-text-primary">{config.emptyTitle}</p>
          {config.emptyDescription ? <p className="mt-2 text-sm">{config.emptyDescription}</p> : null}
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
