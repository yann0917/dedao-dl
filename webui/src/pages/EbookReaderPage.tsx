import {
  ArrowLeft,
  ChevronLeft,
  ChevronRight,
  Loader2,
  Settings2,
  SlidersHorizontal,
} from "lucide-react"
import { useCallback, useEffect, useMemo, useRef, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { toast } from "sonner"
import {
  api,
  type EbookReaderChapter,
  type EbookReaderPageData,
} from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"

// ===== 阅读器设置（对应 get_pages config 的排版参数 + 前端呈现） =====

type ReaderTheme = "light" | "green" | "paper" | "dark"
type ReaderPageMode = "scroll" | "page"

type ReaderSettings = {
  fontSize: number // px，get_pages config.font_size
  lineHeight: number // em，config.line_height
  fontFamily: string // config.font_name（服务端可选字体白名单）
  theme: ReaderTheme
  pageMode: ReaderPageMode
}

const DEFAULT_SETTINGS: ReaderSettings = {
  fontSize: 16,
  lineHeight: 2,
  fontFamily: "pingfang",
  theme: "light",
  pageMode: "scroll",
}

const SETTINGS_STORAGE_KEY = "dedao-reader-settings"

const FONT_OPTIONS: { value: string; label: string }[] = [
  { value: "pingfang", label: "黑体（默认）" },
  { value: "FZKai-Z03", label: "楷体" },
  { value: "FZFangSong-Z02", label: "仿宋" },
]

const THEME_OPTIONS: { value: ReaderTheme; label: string; bg: string }[] = [
  { value: "light", label: "白色", bg: "#ffffff" },
  { value: "green", label: "护眼", bg: "#cce8cf" },
  { value: "paper", label: "羊皮纸", bg: "#f3e9d2" },
  { value: "dark", label: "夜间", bg: "#16181d" },
]

function loadSettings(): ReaderSettings {
  try {
    const raw = window.localStorage.getItem(SETTINGS_STORAGE_KEY)
    if (!raw) {
      return DEFAULT_SETTINGS
    }
    const parsed = JSON.parse(raw) as Partial<ReaderSettings>
    return {
      ...DEFAULT_SETTINGS,
      ...parsed,
      fontSize: clampNumber(Number(parsed.fontSize) || DEFAULT_SETTINGS.fontSize, 12, 28),
      lineHeight: clampNumber(Number(parsed.lineHeight) || DEFAULT_SETTINGS.lineHeight, 1.2, 3),
    }
  } catch {
    return DEFAULT_SETTINGS
  }
}

function clampNumber(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max)
}

// ===== 加载批与页面游标 =====

type ReaderBatch = { start: number; pages: EbookReaderPageData[]; isEnd: boolean }

// 一批最多连续请求 5 次 get_pages（每次 20 页）；
// scroll 模式距页底不足该距离时取下一批
const BATCH_REQUESTS = 5
const PAGE_SIZE = 20
const LOAD_MORE_THRESHOLD_PX = 1200
// page 模式当前页距批尾多少页内时预取下一批
const PREFETCH_REMAINING = 3
// 距视口多少像素内才真正挂载 SVG DOM（单页可能含上千个 text 节点，必须懒挂载）
const MOUNT_ROOT_MARGIN_PX = 2000

function batchStartOf(pageCursor: number) {
  return Math.floor(pageCursor / PAGE_SIZE) * PAGE_SIZE
}

function clampDim(value: number, min: number, max: number) {
  return Math.round(Math.min(Math.max(value, min), max))
}

// 传给 get_pages 的布局尺寸：宽取阅读栏实际宽度（页面 1:1 渲染不缩放），
// 高取浏览器可视区高度减去页头/状态条的大致占用。
// 阅读中途改窗口大小不重新分页，重进页面即生效
function measureReaderDims(column: HTMLElement | null): { width: number; height: number } {
  const width = clampDim(column?.clientWidth || window.innerWidth - 340, 320, 2400)
  const height = clampDim(window.innerHeight - 170, 480, 2400)
  return { width, height }
}

// 根标签的 width/height 只是请求布局的回显，换成内容包围盒 viewBox
// 让每页按宽度自适应渲染，既不裁内容也不留画布空白
function svgWithViewBox(page: EbookReaderPageData): string {
  const gt = page.svg.indexOf(">")
  if (gt < 0) {
    return page.svg
  }
  const root = page.svg
    .slice(0, gt)
    .replace(/\s+width="[^"]*"/, "")
    .replace(/\s+height="[^"]*"/, "")
  return `${root} viewBox="0 0 ${page.width} ${page.height}" preserveAspectRatio="xMidYMin meet"${page.svg.slice(gt)}`
}

function LazySvgPage({ page }: { page: EbookReaderPageData }) {
  const hostRef = useRef<HTMLDivElement | null>(null)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    const el = hostRef.current
    if (!el) {
      return
    }
    // 只保留视口附近的页面 DOM，滚远即卸载（重新滚回时会重挂载）
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          setMounted(entry.isIntersecting)
        }
      },
      { rootMargin: `${MOUNT_ROOT_MARGIN_PX}px 0px` },
    )
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  return (
    <div
      className="reader-page"
      ref={hostRef}
      style={{ aspectRatio: `${page.width} / ${page.height}` }}
    >
      {mounted ? (
        <div className="reader-page-svg" dangerouslySetInnerHTML={{ __html: svgWithViewBox(page) }} />
      ) : null}
    </div>
  )
}

// ===== 设置面板 =====

function SettingsPanel({
  settings,
  onChange,
}: {
  settings: ReaderSettings
  onChange: (next: ReaderSettings) => void
}) {
  return (
    <div className="absolute right-0 top-11 z-30 w-64 space-y-4 rounded-lg border border-border bg-surface-panel p-4 shadow-soft">
      <div>
        <p className="mb-2 text-xs font-medium text-text-muted">字号</p>
        <div className="flex items-center gap-2">
          <Button className="h-8 w-8 px-0" onClick={() => onChange({ ...settings, fontSize: clampNumber(settings.fontSize - 1, 12, 28) })} variant="outline">-</Button>
          <span className="min-w-10 text-center font-mono text-sm text-text-primary">{settings.fontSize}</span>
          <Button className="h-8 w-8 px-0" onClick={() => onChange({ ...settings, fontSize: clampNumber(settings.fontSize + 1, 12, 28) })} variant="outline">+</Button>
        </div>
      </div>

      <div>
        <p className="mb-2 text-xs font-medium text-text-muted">行间距</p>
        <div className="flex items-center gap-2">
          <Button className="h-8 w-8 px-0" onClick={() => onChange({ ...settings, lineHeight: Math.max(1.2, Math.round((settings.lineHeight - 0.25) * 100) / 100) })} variant="outline">-</Button>
          <span className="min-w-10 text-center font-mono text-sm text-text-primary">{settings.lineHeight.toFixed(2)}</span>
          <Button className="h-8 w-8 px-0" onClick={() => onChange({ ...settings, lineHeight: Math.min(3, Math.round((settings.lineHeight + 0.25) * 100) / 100) })} variant="outline">+</Button>
        </div>
      </div>

      <div>
        <p className="mb-2 text-xs font-medium text-text-muted">字体</p>
        <div className="grid grid-cols-3 gap-1">
          {FONT_OPTIONS.map((font) => (
            <button
              className={cn(
                "rounded-md border px-2 py-1.5 text-xs transition",
                settings.fontFamily === font.value
                  ? "border-accent bg-accent-soft font-medium text-accent"
                  : "border-border text-text-secondary hover:bg-surface-soft",
              )}
              key={font.value}
              onClick={() => onChange({ ...settings, fontFamily: font.value })}
              type="button"
            >
              {font.label.replace("（默认）", "")}
            </button>
          ))}
        </div>
        <p className="mt-1.5 text-[10px] text-text-muted">字体由服务端排版，若不支持会回退默认</p>
      </div>

      <div>
        <p className="mb-2 text-xs font-medium text-text-muted">主题 / 背景色</p>
        <div className="grid grid-cols-4 gap-1.5">
          {THEME_OPTIONS.map((theme) => (
            <button
              aria-label={theme.label}
              className={cn(
                "h-8 rounded-md border transition",
                settings.theme === theme.value ? "border-accent ring-2 ring-ring/40" : "border-border",
              )}
              key={theme.value}
              onClick={() => onChange({ ...settings, theme: theme.value })}
              style={{ backgroundColor: theme.bg }}
              title={theme.label}
              type="button"
            />
          ))}
        </div>
      </div>

      <div>
        <p className="mb-2 text-xs font-medium text-text-muted">翻页方式</p>
        <div className="grid grid-cols-2 gap-1">
          {(
            [
              { value: "scroll", label: "上下滚动" },
              { value: "page", label: "左右翻页" },
            ] as { value: ReaderPageMode; label: string }[]
          ).map((mode) => (
            <button
              className={cn(
                "rounded-md border px-2 py-1.5 text-xs transition",
                settings.pageMode === mode.value
                  ? "border-accent bg-accent-soft font-medium text-accent"
                  : "border-border text-text-secondary hover:bg-surface-soft",
              )}
              key={mode.value}
              onClick={() => onChange({ ...settings, pageMode: mode.value })}
              type="button"
            >
              {mode.label}
            </button>
          ))}
        </div>
        <p className="mt-1.5 text-[10px] text-text-muted">改字号/行距/字体会重新分页，已读进度回到当前章开头</p>
      </div>
    </div>
  )
}

// ===== 阅读页 =====

export function EbookReaderPage() {
  const navigate = useNavigate()
  const { enid = "" } = useParams()
  const [searchParams] = useSearchParams()
  const bookTitle = searchParams.get("title") || "电子书阅读"

  const [settings, setSettings] = useState<ReaderSettings>(() => loadSettings())
  const [showSettings, setShowSettings] = useState(false)
  const [readerDims, setReaderDims] = useState<{ width: number; height: number } | null>(null)

  const [chapters, setChapters] = useState<EbookReaderChapter[]>([])
  const [chapterIdx, setChapterIdx] = useState(0)
  const [batches, setBatches] = useState<ReaderBatch[]>([])
  const [pageCursor, setPageCursor] = useState(0) // 章内页号
  const [loadingInfo, setLoadingInfo] = useState(true)
  const [loadingChunk, setLoadingChunk] = useState(false)
  const [chapterEnd, setChapterEnd] = useState(false)
  const [bookEnd, setBookEnd] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const columnRef = useRef<HTMLDivElement | null>(null)
  // StrictMode 下 effect 会双跑，加载状态统一走 ref 防并发重复请求
  const loadingRef = useRef(false)
  const chapterEndRef = useRef(false)
  const bookEndRef = useRef(false)
  const chaptersRef = useRef<EbookReaderChapter[]>([])
  const chapterIdxRef = useRef(0)
  const batchesRef = useRef<ReaderBatch[]>([])
  const pageCursorRef = useRef(0)
  const settingsRef = useRef(settings)
  const dimsRef = useRef<{ width: number; height: number } | null>(null)

  const activeChapter = chapters[chapterIdx] ?? null

  // 排版设置变化时更新 ref 并重新分页（主题/翻页方式不触发重载）
  const layoutKey = `${readerDims?.width}x${readerDims?.height}:${settings.fontSize}:${settings.lineHeight}:${settings.fontFamily}`

  useEffect(() => {
    settingsRef.current = settings
    window.localStorage.setItem(SETTINGS_STORAGE_KEY, JSON.stringify(settings))
  }, [settings])

  const syncBatches = (next: ReaderBatch[]) => {
    batchesRef.current = next
    setBatches(next)
  }

  const syncCursor = (next: number) => {
    pageCursorRef.current = next
    setPageCursor(next)
  }

  const applySettings = (next: ReaderSettings) => {
    const layoutChanged =
      next.fontSize !== settingsRef.current.fontSize ||
      next.lineHeight !== settingsRef.current.lineHeight ||
      next.fontFamily !== settingsRef.current.fontFamily
    setSettings(next)
    if (layoutChanged) {
      // 分页随版式失效：清空当前章已加载内容，从头重载
      setError(null)
      syncBatches([])
      syncCursor(0)
      chapterEndRef.current = false
      bookEndRef.current = false
      setChapterEnd(false)
      setBookEnd(false)
      window.scrollTo({ top: 0 })
      setLayoutVersion((v) => v + 1)
    }
  }

  const [layoutVersion, setLayoutVersion] = useState(0)

  // 首次进入：拉目录，测尺寸
  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoadingInfo(true)
      setError(null)
      try {
        const info = await api.ebook.readerInfo(enid)
        if (cancelled) {
          return
        }
        const list = info.chapters ?? []
        chaptersRef.current = list
        setChapters(list)
        dimsRef.current = measureReaderDims(columnRef.current)
        setReaderDims(dimsRef.current)
        if (list.length === 0) {
          setError("没有拿到这本书的章节目录，可能没有阅读权限")
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "章节目录加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoadingInfo(false)
        }
      }
    }

    if (enid) {
      void load()
    }

    return () => {
      cancelled = true
      chaptersRef.current = []
      chapterIdxRef.current = 0
      batchesRef.current = []
      loadingRef.current = false
      chapterEndRef.current = false
      bookEndRef.current = false
    }
  }, [enid])

  // 加载指定批（一批 1 次请求 20 页；isEnd 标记当前章结束）
  const loadBatchAt = useCallback(
    async (start: number): Promise<boolean> => {
      if (loadingRef.current || !enid) {
        return false
      }
      const dims = dimsRef.current
      const chapter = chaptersRef.current[chapterIdxRef.current]
      if (!dims || !chapter) {
        return false
      }
      if (batchesRef.current.some((batch) => batch.start === start)) {
        return true
      }

      loadingRef.current = true
      setLoadingChunk(true)
      setError(null)
      try {
        const s = settingsRef.current
        const chunk = await api.ebook.readerChunk(enid, chapter.chapterId, start, {
          width: dims.width,
          height: dims.height,
          fontSize: s.fontSize,
          lineHeight: s.lineHeight,
          fontFamily: s.fontFamily,
        })

        const next = [
          ...batchesRef.current.filter((batch) => batch.start !== start),
          { start, pages: chunk.pages, isEnd: chunk.isEnd },
        ].sort((a, b) => a.start - b.start)
        syncBatches(next)

        if (chunk.isEnd) {
          chapterEndRef.current = true
          setChapterEnd(true)
          if (chapterIdxRef.current >= chaptersRef.current.length - 1) {
            bookEndRef.current = true
            setBookEnd(true)
          }
        }
        return true
      } catch (err) {
        const message = err instanceof Error ? err.message : "加载失败，请稍后重试"
        setError(message)
        toast.error("内容加载失败", { description: message })
        return false
      } finally {
        loadingRef.current = false
        setLoadingChunk(false)
      }
    },
    [enid],
  )

  // 版式/章节变化后（batches 已清空）重载当前章开头
  useEffect(() => {
    if (loadingInfo || !readerDims || chapters.length === 0 || error) {
      return
    }
    if (batches.length === 0 && !chapterEnd) {
      void loadBatchAt(0)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [loadingInfo, readerDims, chapters, layoutVersion, chapterIdx, error])

  // scroll 模式：内容不足一屏时继续补批，保证可滚动
  useEffect(() => {
    if (settings.pageMode !== "scroll" || loadingChunk || error || chapterEnd || batches.length === 0) {
      return
    }
    const remaining = document.documentElement.scrollHeight - window.innerHeight - window.scrollY
    if (remaining < LOAD_MORE_THRESHOLD_PX) {
      const last = batches[batches.length - 1]
      const timer = window.setTimeout(() => void loadBatchAt(last.start + PAGE_SIZE), 200)
      return () => window.clearTimeout(timer)
    }
  }, [batches, settings.pageMode, loadingChunk, error, chapterEnd, loadBatchAt])

  // scroll 模式：接近页底时预取下一批
  useEffect(() => {
    if (settings.pageMode !== "scroll") {
      return
    }
    const onScroll = () => {
      if (loadingRef.current || error || chapterEndRef.current) {
        return
      }
      const remaining = document.documentElement.scrollHeight - window.innerHeight - window.scrollY
      if (remaining < LOAD_MORE_THRESHOLD_PX) {
        const last = batchesRef.current[batchesRef.current.length - 1]
        if (last) {
          void loadBatchAt(last.start + PAGE_SIZE)
        }
      }
    }

    window.addEventListener("scroll", onScroll, { passive: true })
    return () => window.removeEventListener("scroll", onScroll)
  }, [settings.pageMode, error, loadBatchAt])

  const jumpToChapter = useCallback(
    (idx: number) => {
      if (loadingRef.current || loadingInfo || idx === chapterIdxRef.current) {
        if (idx === chapterIdxRef.current && batchesRef.current.length > 0) {
          window.scrollTo({ top: 0 })
        }
        return
      }
      chapterIdxRef.current = idx
      setChapterIdx(idx)
      syncBatches([])
      syncCursor(0)
      chapterEndRef.current = false
      bookEndRef.current = false
      setChapterEnd(false)
      setBookEnd(false)
      setError(null)
      window.scrollTo({ top: 0 })
      setLayoutVersion((v) => v + 1)
    },
    [loadingInfo],
  )

  // page 模式：章内已加载总页数（有 isEnd 批时才可信）
  const knownPageCount = useMemo(() => {
    const lastEnd = [...batches].reverse().find((batch) => batch.isEnd)
    if (!lastEnd) {
      return null
    }
    return lastEnd.start + lastEnd.pages.length
  }, [batches])

  const goPage = useCallback(
    async (delta: 1 | -1) => {
      if (loadingRef.current) {
        return
      }
      const target = pageCursorRef.current + delta
      if (target < 0) {
        return
      }

      // 章末翻下一页 → 下一章
      if (knownPageCount !== null && target >= knownPageCount) {
        if (chapterIdxRef.current < chaptersRef.current.length - 1) {
          jumpToChapter(chapterIdxRef.current + 1)
        }
        return
      }

      const start = batchStartOf(target)
      if (!batchesRef.current.some((batch) => batch.start === start)) {
        const loaded = await loadBatchAt(start)
        if (!loaded || !batchesRef.current.some((batch) => batch.start === start)) {
          return
        }
      }
      syncCursor(target)
      window.scrollTo({ top: 0 })

      // 接近批尾预取下一批
      if (target % PAGE_SIZE >= PAGE_SIZE - PREFETCH_REMAINING && chapterEndRef.current !== true) {
        const nextStart = start + PAGE_SIZE
        if (!batchesRef.current.some((batch) => batch.start === nextStart)) {
          void loadBatchAt(nextStart)
        }
      }
    },
    [knownPageCount, jumpToChapter, loadBatchAt],
  )

  // page 模式键盘翻页
  useEffect(() => {
    if (settings.pageMode !== "page") {
      return
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "ArrowRight") {
        void goPage(1)
      } else if (event.key === "ArrowLeft") {
        void goPage(-1)
      }
    }
    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [settings.pageMode, goPage])

  const retry = () => {
    setError(null)
    void loadBatchAt(batchStartOf(pageCursorRef.current) === 0 ? 0 : batchStartOf(pageCursorRef.current))
  }

  const goBack = () => {
    if (window.history.length > 1) {
      navigate(-1)
    } else {
      navigate("/purchased/ebooks")
    }
  }

  const allPages = useMemo(
    () => batches.flatMap((batch) => batch.pages),
    [batches],
  )

  // page 模式当前页
  const currentBatch = batches.find((batch) => batch.start === batchStartOf(pageCursor))
  const currentPage =
    currentBatch && pageCursor >= currentBatch.start
      ? currentBatch.pages[pageCursor - currentBatch.start]
      : undefined

  const themeClass = `reader-theme-${settings.theme}`

  return (
    <div className="grid gap-6 xl:grid-cols-[264px_minmax(0,1fr)]">
      <aside className="hidden xl:block">
        <div className="sticky top-24 flex max-h-[calc(100vh-7rem)] flex-col rounded-lg border border-border bg-surface-panel">
          <p className="border-b border-border px-4 py-3 text-sm font-semibold text-text-primary">目录</p>
          <div className="min-h-0 flex-1 overflow-y-auto p-2">
            {loadingInfo ? (
              <div className="flex items-center gap-2 px-3 py-2 text-sm text-text-muted">
                <Loader2 className="size-icon-md animate-spin" />
                正在加载目录...
              </div>
            ) : (
              chapters.map((chapter, idx) => (
                <button
                  className={cn(
                    "block w-full truncate rounded-md px-3 py-2 text-left text-sm transition",
                    idx === chapterIdx
                      ? "bg-accent-soft font-medium text-accent"
                      : "text-text-secondary hover:bg-surface-soft hover:text-text-primary",
                  )}
                  key={chapter.chapterId}
                  onClick={() => jumpToChapter(idx)}
                  title={chapter.title || `第 ${idx + 1} 章`}
                  type="button"
                >
                  {chapter.title || `第 ${idx + 1} 章`}
                </button>
              ))
            )}
          </div>
        </div>
      </aside>

      <div className="min-w-0 space-y-6" ref={columnRef}>
        <div className="flex flex-wrap items-center gap-3">
          <Button onClick={goBack} variant="outline">
            <ArrowLeft className="mr-2 size-icon-md" />
            返回
          </Button>
          <div className="min-w-0">
            <h2 className="truncate text-xl font-semibold text-text-primary">{bookTitle}</h2>
            <p className="truncate text-sm text-text-muted">
              {activeChapter?.title || "正在打开..."}
            </p>
          </div>
          <div className="relative ml-auto">
            <Button
              onClick={() => setShowSettings((v) => !v)}
              variant="outline"
            >
              {showSettings ? <SlidersHorizontal className="mr-2 size-icon-md" /> : <Settings2 className="mr-2 size-icon-md" />}
              阅读设置
            </Button>
            {showSettings ? (
              <SettingsPanel settings={settings} onChange={applySettings} />
            ) : null}
          </div>
        </div>

        {error && allPages.length === 0 ? (
          <Card className="border border-danger bg-danger-soft p-6 text-sm text-danger">
            {error}
            <div className="mt-4">
              <Button onClick={retry} variant="outline">
                重试
              </Button>
            </div>
          </Card>
        ) : (
          <>
            {settings.pageMode === "scroll" ? (
              <Card className={cn("overflow-hidden", themeClass)}>
                <div className="flex items-baseline justify-between gap-4 border-b border-border px-6 py-3">
                  <p className="truncate text-sm font-medium text-text-primary">
                    {activeChapter?.title || `第 ${chapterIdx + 1} 章`}
                  </p>
                  <p className="flex-none font-mono text-[10px] uppercase tracking-[0.2em] text-text-muted">
                    已加载 {allPages.length} 页
                  </p>
                </div>
                <div className="space-y-3 p-3 lg:p-5">
                  {allPages.map((page, idx) => (
                    <LazySvgPage key={`${chapterIdx}-${idx}-${page.width}x${page.height}`} page={page} />
                  ))}
                </div>
              </Card>
            ) : (
              <div>
                <div
                  className={cn(
                    "reader-page-host relative flex items-start justify-center overflow-hidden rounded-lg border border-border p-3",
                    themeClass,
                  )}
                  style={{ height: readerDims ? `${readerDims.height}px` : "70vh" }}
                >
                  {currentPage ? (
                    <LazySvgPage page={currentPage} />
                  ) : loadingChunk ? (
                    <div className="flex h-full items-center gap-3 text-sm text-text-muted">
                      <Loader2 className="size-icon-lg animate-spin" />
                      正在加载...
                    </div>
                  ) : (
                    <div className="flex h-full items-center text-sm text-text-muted">
                      {chapterEnd ? "本章到此结束，按 → 进入下一章" : "按 → 翻到下一页"}
                    </div>
                  )}
                </div>
                <div className="mt-3 flex items-center justify-center gap-4">
                  <Button disabled={pageCursor <= 0 || loadingChunk} onClick={() => void goPage(-1)} variant="outline">
                    <ChevronLeft className="mr-1 size-icon-md" />
                    上一页
                  </Button>
                  <span className="min-w-24 text-center font-mono text-sm text-text-muted">
                    {pageCursor + 1}
                    {knownPageCount !== null ? ` / ${knownPageCount}` : ""}
                  </span>
                  <Button
                    disabled={
                      loadingChunk ||
                      (bookEnd && knownPageCount !== null && pageCursor + 1 >= knownPageCount)
                    }
                    onClick={() => void goPage(1)}
                    variant="outline"
                  >
                    下一页
                    <ChevronRight className="ml-1 size-icon-md" />
                  </Button>
                </div>
              </div>
            )}

            {error ? (
              <Card className="flex flex-wrap items-center justify-between gap-3 border border-warning bg-warning-soft p-4 text-sm text-warning">
                <span>{error}</span>
                <Button onClick={retry} variant="outline">
                  重试加载
                </Button>
              </Card>
            ) : null}

            {loadingChunk ? (
              <div className="flex items-center justify-center gap-3 rounded-lg border border-dashed border-border p-8 text-sm text-text-muted">
                <Loader2 className="size-icon-lg animate-spin" />
                正在加载（每次 1 次请求，约每 5 秒 1 次）...
              </div>
            ) : null}

            {bookEnd && (settings.pageMode === "scroll" || (knownPageCount !== null && pageCursor + 1 >= knownPageCount)) && !loadingChunk ? (
              <div className="rounded-lg border border-dashed border-border p-8 text-center text-sm text-text-muted">
                已经读到最后一章了
              </div>
            ) : null}
          </>
        )}
      </div>
    </div>
  )
}
