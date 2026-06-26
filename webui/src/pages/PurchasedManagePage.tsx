import { type ColumnDef, flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table"
import {
  BookMarked,
  Compass,
  Download,
  Eye,
  FileText,
  GraduationCap,
  Headphones,
  Loader2,
  Play,
  Rows3,
} from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate } from "react-router-dom"
import { toast } from "sonner"
import {
  api,
  type AudioDetailResponse,
  type CourseListItem,
  type PurchasedNavbarChild,
  type PurchasedNavbarItem,
} from "@/api"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/shadcn/dialog"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/shadcn/pagination"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/shadcn/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/shadcn/table"
import { Tabs, TabsList, TabsTrigger } from "@/components/shadcn/tabs"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/shadcn/tooltip"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"
import { useDownloadProgress } from "@/providers/DownloadProgressProvider"

type ManageTabKey = "course" | "audio" | "ebook" | "compass"

type ManageTabConfig = {
  key: ManageTabKey
  label: string
  category: string
  icon: typeof GraduationCap
  description: string
}

type GroupMode = {
  active: boolean
  groupId: number
  title: string
}

type DownloadOption = {
  value: number
  label: string
}

type DownloadTarget = {
  item: CourseListItem
  tab: ManageTabConfig
}

const MANAGE_TABS: ManageTabConfig[] = [
  {
    key: "course",
    label: "课程",
    category: "bauhinia",
    icon: GraduationCap,
    description: "统一查看课程进度、更新状态和下载导出。",
  },
  {
    key: "audio",
    label: "听书",
    category: "odob",
    icon: Headphones,
    description: "集中管理听书播放、文稿查看和下载导出。",
  },
  {
    key: "ebook",
    label: "电子书",
    category: "ebook",
    icon: BookMarked,
    description: "统一查看电子书并执行下载操作。",
  },
  {
    key: "compass",
    label: "锦囊",
    category: "compass",
    icon: Compass,
    description: "锦囊当前仅支持统一浏览，不提供任何操作。",
  },
]

const DOWNLOAD_OPTIONS: Record<Exclude<ManageTabKey, "compass">, DownloadOption[]> = {
  course: [
    { value: 1, label: "下载音频 MP3" },
    { value: 2, label: "下载课程 PDF" },
    { value: 3, label: "下载课程 Markdown" },
  ],
  audio: [
    { value: 1, label: "下载音频 MP3" },
    { value: 2, label: "下载文稿 PDF" },
    { value: 3, label: "下载文稿 Markdown" },
  ],
  ebook: [
    { value: 1, label: "下载 HTML" },
    { value: 2, label: "下载 PDF" },
    { value: 3, label: "下载 EPUB" },
  ],
}

const PAGE_SIZE_OPTIONS = [10, 20, 50]

function ensureAllFilter(options: PurchasedNavbarChild[]) {
  if (options.some((item) => item.filter === "all")) {
    return options
  }

  return [{ name: "全部", count: 0, filter: "all", show_count: false }, ...options]
}

function resolveItemTitle(item: CourseListItem) {
  return item.title || item.name || "未命名内容"
}

function resolveItemCover(item: CourseListItem) {
  return item.cover || item.icon || item.index_img || "https://placehold.co/320x320/e2e8f0/334155?text=Dedao"
}

function resolveItemSummary(item: CourseListItem) {
  return item.subtitle || item.intro || ""
}

function truncateText(value: string, maxLength = 68) {
  if (!value || value.length <= maxLength) {
    return value
  }

  return `${value.slice(0, maxLength).trim()}...`
}

function formatMinutes(duration?: number) {
  if (!duration) {
    return "0 分钟"
  }

  return `${Math.round(duration / 60)} 分钟`
}

function getPrimaryMeta(item: CourseListItem, key: ManageTabKey) {
  if (key === "course") {
    return item.lecturer_name || item.author || "课程内容"
  }

  if (key === "audio") {
    return formatMinutes(item.duration)
  }

  return item.author || item.lecturer_name || "内容"
}

function getSecondaryMeta(item: CourseListItem, key: ManageTabKey) {
  if (key === "course") {
    return `已更 ${item.publish_num || 0}/${item.course_num || 0}`
  }

  if (key === "audio") {
    if (item.type === 1013) {
      return "听书合集"
    }

    return item.in_bookrack ? "已加入书架" : "单本听书"
  }

  if (key === "ebook") {
    return item.price ? `¥${item.price}` : "电子书"
  }

  return item.is_group ? "分组内容" : "仅展示"
}

function getStatusText(item: CourseListItem, key: ManageTabKey) {
  if (key === "course") {
    return `${item.progress || 0}%`
  }

  if (key === "audio") {
    if (item.is_group || item.type === 1013) {
      return "合集"
    }

    return item.has_play_auth ? "可播放" : "待详情授权"
  }

  if (key === "ebook") {
    return item.in_bookrack || item.is_on_bookshelf ? "书架中" : "可查看"
  }

  return ""
}

function getBadgeLabels(item: CourseListItem, key: ManageTabKey) {
  const badges: string[] = []

  if (item.is_group) {
    badges.push("分组")
  }

  if (key === "audio" && item.type === 1013) {
    badges.push("合集")
  }

  if (key === "audio" && item.in_bookrack) {
    badges.push("书架")
  }

  if (key === "course" && (item.progress || 0) > 0) {
    badges.push("学习中")
  }

  if (key === "compass") {
    badges.push("")
  }

  return badges
}

function getDownloadDisabledReason(item: CourseListItem, key: ManageTabKey) {
  if (key === "compass") {
    return "锦囊当前仅支持展示，不提供任何操作。"
  }

  if (!item.enid) {
    return "当前条目缺少 enid，暂不可下载。"
  }

  return null
}

function buildPageItems(currentPage: number, totalPages: number) {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, index) => index + 1)
  }

  const pages: Array<number | "ellipsis-start" | "ellipsis-end"> = [1]
  const start = Math.max(2, currentPage - 1)
  const end = Math.min(totalPages - 1, currentPage + 1)

  if (start > 2) {
    pages.push("ellipsis-start")
  }

  for (let page = start; page <= end; page += 1) {
    pages.push(page)
  }

  if (end < totalPages - 1) {
    pages.push("ellipsis-end")
  }

  pages.push(totalPages)
  return pages
}

function parsePositiveId(value?: number | string) {
  const parsed = typeof value === "number" ? value : Number(value)
  return Number.isFinite(parsed) && parsed > 0 ? parsed : undefined
}

function buildAudioTaskTitle(detail: AudioDetailResponse, fallbackTitle: string) {
  const packageTitle = detail.package_title?.trim()
  const itemTitle = detail.title?.trim()

  if (packageTitle && itemTitle) {
    return `${packageTitle} - ${itemTitle}`
  }

  return itemTitle || packageTitle || fallbackTitle
}

function createDownloadRequest(item: CourseListItem, key: Exclude<ManageTabKey, "compass">, downloadType: number) {
  if (key === "course") {
    return api.download.course({
      enid: item.enid,
      title: resolveItemTitle(item),
      downloadType,
      isOrder: true,
    })
  }

  if (key === "audio") {
    return api.download.audio({
      enid: item.enid,
      title: resolveItemTitle(item),
      downloadType,
    })
  }

  return api.download.ebook({
    enid: item.enid,
    title: resolveItemTitle(item),
    downloadType,
  })
}

export function PurchasedManagePage() {
  const navigate = useNavigate()
  const { setQueue } = useAudioPlayer()
  const { beginSession } = useDownloadProgress()
  const [activeTabKey, setActiveTabKey] = useState<ManageTabKey>("course")
  const [navbarItems, setNavbarItems] = useState<PurchasedNavbarItem[]>([])
  const [items, setItems] = useState<CourseListItem[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [currentFilter, setCurrentFilter] = useState("all")
  const [groupMode, setGroupMode] = useState<GroupMode>({
    active: false,
    groupId: 0,
    title: "",
  })
  const [loadingNavbar, setLoadingNavbar] = useState(true)
  const [loadingItems, setLoadingItems] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [playingEnid, setPlayingEnid] = useState<string | null>(null)
  const [downloadTarget, setDownloadTarget] = useState<DownloadTarget | null>(null)
  const [pendingDownloadType, setPendingDownloadType] = useState<number | null>(null)

  const activeTab = useMemo(
    () => MANAGE_TABS.find((tab) => tab.key === activeTabKey) || MANAGE_TABS[0],
    [activeTabKey],
  )

  useEffect(() => {
    let cancelled = false

    const loadNavbar = async () => {
      setLoadingNavbar(true)

      try {
        const data = await api.course.navbar()
        if (cancelled) {
          return
        }

        setNavbarItems(data.list)
      } catch (err) {
        if (!cancelled) {
          setNavbarItems([])
          setError(err instanceof Error ? err.message : "已购筛选加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoadingNavbar(false)
        }
      }
    }

    void loadNavbar()

    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    setItems([])
    setTotal(0)
    setPage(1)
    setCurrentFilter("all")
    setGroupMode({
      active: false,
      groupId: 0,
      title: "",
    })
  }, [activeTabKey])

  useEffect(() => {
    let cancelled = false

    const loadItems = async () => {
      setLoadingItems(true)
      setError(null)

      try {
        const params = new URLSearchParams({
          category: activeTab.category,
          filter: currentFilter,
          order: "study",
          page: String(page),
          limit: String(pageSize),
        })

        if (groupMode.active && groupMode.groupId > 0) {
          params.set("groupId", String(groupMode.groupId))
        }

        const result = await api.course.list(params)
        if (cancelled) {
          return
        }

        setItems(result.list)
        setTotal(result.total || 0)
      } catch (err) {
        if (!cancelled) {
          setItems([])
          setTotal(0)
          setError(err instanceof Error ? err.message : "已购内容加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoadingItems(false)
        }
      }
    }

    void loadItems()

    return () => {
      cancelled = true
    }
  }, [activeTab.category, currentFilter, groupMode.active, groupMode.groupId, page, pageSize])

  const filterOptions = useMemo(() => {
    const matched = navbarItems.find((item) => item.category === activeTab.category)
    return ensureAllFilter(matched?.children ?? [])
  }, [activeTab.category, navbarItems])

  const totalPages = Math.max(1, Math.ceil(total / pageSize))
  const pageItems = useMemo(() => buildPageItems(page, totalPages), [page, totalPages])
  const currentDownloadOptions = downloadTarget && downloadTarget.tab.key !== "compass"
    ? DOWNLOAD_OPTIONS[downloadTarget.tab.key]
    : []

  const openItem = (item: CourseListItem) => {
    if (item.is_group) {
      const groupId = Number(item.group_id || item.id || 0)
      if (!groupId) {
        return
      }

      setGroupMode({
        active: true,
        groupId,
        title: resolveItemTitle(item),
      })
      setPage(1)
      return
    }

    if (!item.enid) {
      return
    }

    if (activeTab.key === "course") {
      navigate(
        `/courses/${encodeURIComponent(item.enid)}?from=purchased-manage&parentTitle=${encodeURIComponent(resolveItemTitle(item))}`,
      )
      return
    }

    if (activeTab.key === "audio") {
      if (item.type === 1013) {
        navigate(
          `/audio-groups/${encodeURIComponent(item.enid)}?from=purchased-manage&parentTitle=${encodeURIComponent(resolveItemTitle(item))}`,
        )
        return
      }

      navigate(
        `/audios/${encodeURIComponent(item.enid)}?from=purchased-manage&parentTitle=${encodeURIComponent(resolveItemTitle(item))}`,
      )
      return
    }

    if (activeTab.key === "ebook") {
      navigate(
        `/ebooks/${encodeURIComponent(item.enid)}?from=purchased-manage&parentTitle=${encodeURIComponent(resolveItemTitle(item))}`,
      )
    }
  }

  const handlePlayAudio = async (item: CourseListItem) => {
    if (!item.enid) {
      return
    }

    setPlayingEnid(item.enid)

    try {
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
            title: detail.title || detail.package_title || resolveItemTitle(item),
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

  const handleDownload = async (downloadType: number) => {
    if (!downloadTarget || downloadTarget.tab.key === "compass") {
      return
    }

    setPendingDownloadType(downloadType)

    try {
      if (downloadTarget.tab.key === "audio" && downloadTarget.item.type === 1013) {
        // 合集本身不是可下载音频，先解析合集详情，再按子条目逐个创建下载任务。
        const groupResult = await api.audio.group(downloadTarget.item.enid)
        const detailList = groupResult.group?.odob_audio_detail_list ?? []
        const validItems = detailList.filter((detail) => detail.topic_encode_id && detail.has_play_auth)

        if (validItems.length === 0) {
          throw new Error("当前合集没有可下载的听书条目，请先确认播放权限是否齐全")
        }

        let successCount = 0
        let failureCount = 0
        const skippedCount = detailList.length - validItems.length
        let lastErrorMessage = ""

        for (const detail of validItems) {
          try {
            const result = await api.download.audio({
              id: parsePositiveId(detail.audio_id),
              enid: detail.topic_encode_id,
              title: buildAudioTaskTitle(detail, resolveItemTitle(downloadTarget.item)),
              downloadType,
            })
            beginSession(result)
            successCount += 1
          } catch (err) {
            failureCount += 1
            lastErrorMessage = err instanceof Error ? err.message : "请稍后重试"
          }
        }

        if (successCount === 0) {
          throw new Error(lastErrorMessage || "合集下载任务创建失败")
        }

        const summaryParts = [`已加入 ${successCount} 个任务`]
        if (skippedCount > 0) {
          summaryParts.push(`跳过 ${skippedCount} 个无权限或缺少 enid 的条目`)
        }
        if (failureCount > 0) {
          summaryParts.push(`失败 ${failureCount} 个`)
        }

        toast.success("合集下载任务已加入队列", {
          description: `${summaryParts.join("，")}。`,
        })
        setDownloadTarget(null)
        return
      }

      const result = await createDownloadRequest(downloadTarget.item, downloadTarget.tab.key, downloadType)
      beginSession(result)
      toast.success("下载已开始", {
        description: `正在连接进度流。输出目录：${result.outputDir}`,
      })
      setDownloadTarget(null)
    } catch (err) {
      toast.error("下载失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setPendingDownloadType(null)
    }
  }

  const columns = useMemo<ColumnDef<CourseListItem>[]>(
    () => [
      {
        id: "content",
        header: "内容",
        cell: ({ row }) => {
          const item = row.original
          const summary = truncateText(resolveItemSummary(item))
          const badges = getBadgeLabels(item, activeTab.key)
          const coverClassName = activeTab.key === "ebook" ? "aspect-[3/4]" : "aspect-square"

          return (
            <div className="flex min-w-[320px] items-start gap-4">
              <div className={cn("w-20 shrink-0 rounded-2xl bg-surface-soft p-2", coverClassName)}>
                <img
                  alt={resolveItemTitle(item)}
                  className="h-full w-full rounded-xl object-contain"
                  src={resolveItemCover(item)}
                />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-semibold text-text-primary">{resolveItemTitle(item)}</p>
                <p className="mt-1 text-sm text-text-secondary">{getPrimaryMeta(item, activeTab.key)}</p>
                <p className="mt-1 text-xs text-text-muted">{getSecondaryMeta(item, activeTab.key)}</p>
                {summary ? (
                  <p className="mt-2 text-xs leading-6 text-text-muted">{summary}</p>
                ) : null}
                {badges.length > 0 ? (
                  <div className="mt-3 flex flex-wrap gap-2">
                    {badges.map((badge) => (
                      <span
                        className="rounded-full bg-surface-soft px-2.5 py-1 text-[11px] font-medium text-text-secondary"
                        key={badge}
                      >
                        {badge}
                      </span>
                    ))}
                  </div>
                ) : null}
              </div>
            </div>
          )
        },
      },
      {
        id: "status",
        header: "状态",
        cell: ({ row }) => {
          const item = row.original

          return (
            <div className="min-w-[160px] space-y-2">
              <p className="text-sm font-medium text-text-primary">{getStatusText(item, activeTab.key)}</p>
              {activeTab.key === "course" ? (
                <div className="space-y-2">
                  <div className="h-2 rounded-full bg-surface-soft">
                    <div
                      className="h-2 rounded-full bg-accent transition-[width]"
                      style={{ width: `${Math.min(item.progress || 0, 100)}%` }}
                    />
                  </div>
                  <p className="text-xs text-text-muted">学习进度 {item.progress || 0}%</p>
                </div>
              ) : (
                <p className="text-xs leading-6 text-text-muted">
                  {activeTab.key === "audio"
                    ? item.audio_detail?.alias_id
                      ? "含文稿入口"
                      : "文稿入口以详情数据为准"
                    : activeTab.key === "ebook"
                      ? item.is_on_bookshelf || item.in_bookrack
                        ? "已在书架"
                        : "可查看详情"
                      : "锦囊暂不支持下载或跳转"}
                </p>
              )}
            </div>
          )
        },
      },
      {
        id: "actions",
        header: "操作",
        cell: ({ row }) => {
          const item = row.original
          const downloadDisabledReason = getDownloadDisabledReason(item, activeTab.key)
          const canDownload = !downloadDisabledReason

          if (activeTab.key === "compass") {
            return (
              <Tooltip>
                <TooltipTrigger asChild>
                  <span className="inline-flex cursor-help rounded-xl bg-surface-soft px-3 py-2 text-sm text-text-muted">
                    仅展示
                  </span>
                </TooltipTrigger>
                <TooltipContent>锦囊当前仅支持统一浏览。</TooltipContent>
              </Tooltip>
            )
          }

          return (
            <div className="flex min-w-[280px] flex-wrap items-center gap-2">
              <Button className="h-9 px-3" onClick={() => openItem(item)} variant="outline">
                {item.is_group ? (
                  <Rows3 className="mr-2 h-4 w-4" />
                ) : (
                  <Eye className="mr-2 h-4 w-4" />
                )}
                {item.is_group ? "进入分组" : activeTab.key === "audio" && item.type === 1013 ? "查看合集" : "查看详情"}
              </Button>

              {activeTab.key === "audio" && !item.is_group && item.type !== 1013 ? (
                <Button
                  className="h-9 px-3"
                  disabled={!item.enid || playingEnid === item.enid}
                  onClick={() => void handlePlayAudio(item)}
                >
                  {playingEnid === item.enid ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <Play className="mr-2 h-4 w-4" />
                  )}
                  {playingEnid === item.enid ? "获取中..." : "播放"}
                </Button>
              ) : null}

              {activeTab.key === "audio" && !item.is_group && item.type !== 1013 && item.audio_detail?.alias_id ? (
                <Button
                  className="h-9 px-3"
                  onClick={() =>
                    navigate(
                      `/articles/2/${encodeURIComponent(item.audio_detail?.alias_id || "")}?from=audio&parentEnid=${encodeURIComponent(item.enid)}&parentTitle=${encodeURIComponent(resolveItemTitle(item))}`,
                    )
                  }
                  variant="ghost"
                >
                  <FileText className="mr-2 h-4 w-4" />
                  查看文稿
                </Button>
              ) : null}

              {canDownload ? (
                <Button
                  className="h-9 px-3"
                  onClick={() => setDownloadTarget({ item, tab: activeTab })}
                  variant="ghost"
                >
                  <Download className="mr-2 h-4 w-4" />
                  下载
                </Button>
              ) : (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <span>
                      <Button className="h-9 px-3" disabled variant="ghost">
                        <Download className="mr-2 h-4 w-4" />
                        下载
                      </Button>
                    </span>
                  </TooltipTrigger>
                  <TooltipContent>{downloadDisabledReason}</TooltipContent>
                </Tooltip>
              )}
            </div>
          )
        },
      },
    ],
    [activeTab, navigate, playingEnid],
  )

  const table = useReactTable({
    data: items,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <TooltipProvider delayDuration={150}>
      <main className="space-y-6">
        {groupMode.active ? (
          <section className="rounded-3xl border border-border bg-surface-panel p-4 shadow-soft backdrop-blur">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <p className="text-sm text-text-muted">当前分组</p>
                <p className="text-lg font-semibold text-text-primary">{groupMode.title}</p>
              </div>
              <Button
                className="h-10 px-4"
                onClick={() => {
                  setGroupMode({
                    active: false,
                    groupId: 0,
                    title: "",
                  })
                  setPage(1)
                }}
                variant="outline"
              >
                返回全部
              </Button>
            </div>
          </section>
        ) : null}

        <Card className="p-3">
          <Tabs onValueChange={(value) => setActiveTabKey(value as ManageTabKey)} value={activeTab.key}>
            <TabsList variant="line" className="w-full flex-wrap">
              {MANAGE_TABS.map((tab) => (
                <TabsTrigger
                  className="h-10 rounded-t-xl rounded-b-none border border-transparent border-b-2 border-b-transparent px-4 text-text-secondary data-[state=active]:border-border data-[state=active]:border-b-accent data-[state=active]:bg-surface-panel data-[state=active]:text-text-primary data-[state=active]:shadow-none"
                  key={tab.key}
                  value={tab.key}
                >
                  <tab.icon className="mr-2 h-4 w-4" />
                  {tab.label}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
        </Card>

        <Card className="p-4">
          <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
            <div>
              <p className="text-sm font-medium text-text-primary">{activeTab.label}管理台</p>
              <p className="mt-1 text-sm text-text-muted">{activeTab.description}</p>
            </div>

            <div className="flex flex-col gap-3 sm:flex-row">
              <div className="space-y-1">
                <p className="text-xs font-medium text-text-muted">二级筛选</p>
                <Select
                  onValueChange={(value) => {
                    setCurrentFilter(value)
                    setPage(1)
                  }}
                  value={currentFilter}
                >
                  <SelectTrigger className="h-10 w-[220px] rounded-xl border-border bg-surface-panel text-text-primary">
                    <SelectValue placeholder={loadingNavbar ? "加载筛选中..." : "选择筛选"} />
                  </SelectTrigger>
                  <SelectContent>
                    {filterOptions.map((option) => (
                      <SelectItem key={option.filter} value={option.filter}>
                        {option.show_count ? `${option.name} (${option.count})` : option.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-1">
                <p className="text-xs font-medium text-text-muted">每页数量</p>
                <Select
                  onValueChange={(value) => {
                    setPageSize(Number(value))
                    setPage(1)
                  }}
                  value={String(pageSize)}
                >
                  <SelectTrigger className="h-10 w-[140px] rounded-xl border-border bg-surface-panel text-text-primary">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PAGE_SIZE_OPTIONS.map((size) => (
                      <SelectItem key={size} value={String(size)}>
                        {size} 条
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
        </Card>

        <Card className="overflow-hidden">
          {loadingItems ? (
            <div className="flex min-h-[320px] items-center justify-center gap-3 text-text-muted">
              <Loader2 className="h-5 w-5 animate-spin" />
              正在加载{activeTab.label}列表...
            </div>
          ) : error ? (
            <div className="flex min-h-[240px] items-center justify-center px-6 text-center text-sm text-danger">
              {error}
            </div>
          ) : items.length === 0 ? (
            <div className="flex min-h-[240px] items-center justify-center px-6 text-center text-sm text-text-muted">
              当前筛选下没有可展示的{activeTab.label}内容。
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  {table.getHeaderGroups().map((headerGroup) => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map((header) => (
                        <TableHead key={header.id}>
                          {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                        </TableHead>
                      ))}
                    </TableRow>
                  ))}
                </TableHeader>
                <TableBody>
                  {table.getRowModel().rows.map((row) => (
                    <TableRow key={row.id}>
                      {row.getVisibleCells().map((cell) => (
                        <TableCell className="align-top" key={cell.id}>
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>

              <div className="flex flex-col gap-4 border-t border-border px-4 py-4 lg:flex-row lg:items-center lg:justify-between">
                <p className="text-sm text-text-muted">
                  共 {total} 条，当前第 {page} / {totalPages} 页
                </p>

                <div className="overflow-x-auto pb-1">
                  <Pagination className="mx-0 w-max min-w-full justify-start lg:justify-end">
                    <PaginationContent className="w-max min-w-max">
                      <PaginationItem>
                        <PaginationPrevious
                          aria-disabled={page <= 1}
                          className={page <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
                          href="#"
                          onClick={(event) => {
                            event.preventDefault()
                            if (page > 1) {
                              setPage(page - 1)
                            }
                          }}
                        />
                      </PaginationItem>

                      {pageItems.map((pageItem) => {
                        if (typeof pageItem !== "number") {
                          return (
                            <PaginationItem key={pageItem}>
                              <PaginationEllipsis />
                            </PaginationItem>
                          )
                        }

                        return (
                          <PaginationItem key={pageItem}>
                            <PaginationLink
                              className="cursor-pointer"
                              href="#"
                              isActive={pageItem === page}
                              onClick={(event) => {
                                event.preventDefault()
                                setPage(pageItem)
                              }}
                            >
                              {pageItem}
                            </PaginationLink>
                          </PaginationItem>
                        )
                      })}

                      <PaginationItem>
                        <PaginationNext
                          aria-disabled={page >= totalPages}
                          className={page >= totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
                          href="#"
                          onClick={(event) => {
                            event.preventDefault()
                            if (page < totalPages) {
                              setPage(page + 1)
                            }
                          }}
                        />
                      </PaginationItem>
                    </PaginationContent>
                  </Pagination>
                </div>
              </div>
            </>
          )}
        </Card>

        <Dialog
          onOpenChange={(open) => {
            if (!open) {
              setDownloadTarget(null)
              setPendingDownloadType(null)
            }
          }}
          open={downloadTarget !== null}
        >
          <DialogContent className="max-w-xl rounded-3xl border-border bg-surface-panel text-text-primary">
            <DialogHeader>
              <DialogTitle>下载 {downloadTarget ? resolveItemTitle(downloadTarget.item) : ""}</DialogTitle>
              <DialogDescription className="text-text-muted">
                {downloadTarget?.tab.key === "audio" && downloadTarget.item.type === 1013
                  ? "将先解析合集中的每本听书，再按所选格式批量创建下载任务。"
                  : "根据当前内容类型选择导出格式。下载任务会在服务端执行，并同步显示进度。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-3">
              {currentDownloadOptions.map((option) => (
                <button
                  className={cn(
                    "flex w-full items-center justify-between rounded-2xl border border-border bg-surface-soft px-4 py-4 text-left transition hover:border-border hover:bg-surface-panel",
                    pendingDownloadType === option.value ? "cursor-wait opacity-70" : "",
                  )}
                  disabled={pendingDownloadType !== null}
                  key={option.value}
                  onClick={() => void handleDownload(option.value)}
                  type="button"
                >
                  <span className="flex items-center gap-3 text-sm font-medium text-text-primary">
                    {pendingDownloadType === option.value ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Download className="h-4 w-4" />
                    )}
                    {option.label}
                  </span>
                  <span className="text-xs text-text-muted">
                    {pendingDownloadType === option.value ? "正在提交..." : "开始下载"}
                  </span>
                </button>
              ))}
            </div>

            <DialogFooter>
              <Button onClick={() => setDownloadTarget(null)} variant="outline">
                关闭
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </main>
    </TooltipProvider>
  )
}
