import { BookOpen, Clock3, Loader2, NotebookPen, Star } from "lucide-react"
import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { toast } from "sonner"
import { api, type EbookDetailResponse } from "@/api"
import { DownloadActionsPanel, type DownloadOption } from "@/components/download/DownloadActionsPanel"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"
import {
  semanticSecondaryTextClass,
  semanticSubtlePanelClass,
} from "@/lib/semanticStyles"

const ebookDownloadOptions: DownloadOption[] = [
  { value: 1, label: "下载 HTML" },
  { value: 2, label: "下载 PDF" },
  { value: 3, label: "下载 EPUB" },
]

function formatTime(value?: number) {
  if (!value) {
    return "未知"
  }

  return new Date(value * 1000).toLocaleString("zh-CN")
}

function normalizeCatalogLevel(level?: number) {
  return Math.max(level ?? 0, 0)
}

function getCatalogIndent(level: number) {
  return Math.min(level, 6) * 18
}

function getCatalogTextClass(level: number) {
  if (level === 0) {
    return "text-sm font-semibold text-text-primary"
  }

  if (level === 1) {
    return "text-sm font-medium text-text-primary"
  }

  if (level === 2) {
    return "text-sm text-text-secondary"
  }

  return "text-sm text-text-muted"
}

function getCatalogRowClass(level: number) {
  if (level === 0) {
    return "bg-surface-soft/90"
  }

  if (level === 1) {
    return "bg-surface-panel"
  }

  return "bg-surface-panel/70"
}

function getCatalogDotClass(level: number) {
  if (level === 0) {
    return "bg-text-muted"
  }

  if (level === 1) {
    return "bg-text-muted/80"
  }

  return "bg-border"
}

export function EbookDetailPage() {
  const navigate = useNavigate()
  const { enid = "" } = useParams()
  const [data, setData] = useState<EbookDetailResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [shelfLoading, setShelfLoading] = useState(false)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      setError(null)
      try {
        const result = await api.ebook.detail(enid)
        if (!cancelled) {
          setData(result)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "电子书详情加载失败")
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
          正在加载电子书详情...
        </div>
      </main>
    )
  }

  if (error || !data?.detail) {
    return (
      <main className="space-y-6">
        <Card className="border border-danger bg-danger-soft p-6 text-sm text-danger">
          {error || "未找到电子书详情"}
        </Card>
      </main>
    )
  }

  const detail = data.detail
  const notes = data.notes?.list.slice(0, 6) ?? []
  const catalogList = detail.catalog_list ?? []
  const chapterCount = catalogList.length || detail.count || 0
  const pressName = detail.press?.name || "未知"
  const pressBrief = detail.press?.brief || ""
  const isOnBookshelf = Boolean(detail.is_on_bookshelf)

  const handleAddToShelf = async () => {
    setShelfLoading(true)
    try {
      await api.ebook.addToShelf([enid])
      setData((current) => {
        if (!current?.detail) {
          return current
        }

        return {
          ...current,
          detail: {
            ...current.detail,
            is_on_bookshelf: true,
          },
        }
      })
      toast.success("已加入书架", {
        description: detail.title || "当前电子书已加入书架",
      })
    } catch (err) {
      toast.error("加入书架失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setShelfLoading(false)
    }
  }

  const handleRemoveFromShelf = async () => {
    if (!window.confirm("确定要将这本电子书移出书架吗？")) {
      return
    }

    setShelfLoading(true)
    try {
      await api.ebook.removeFromShelf([enid])
      setData((current) => {
        if (!current?.detail) {
          return current
        }

        return {
          ...current,
          detail: {
            ...current.detail,
            is_on_bookshelf: false,
          },
        }
      })
      toast.success("已移出书架", {
        description: detail.title || "当前电子书已从书架移除",
      })
    } catch (err) {
      toast.error("移出书架失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setShelfLoading(false)
    }
  }

  const handleEbookDownload = (downloadType: number) =>
    api.download.ebook({
      enid,
      title: detail.title || "电子书下载",
      downloadType,
    })

  return (
    <main className="space-y-6">
      <section className="rounded-3xl border border-border bg-surface-panel p-6 shadow-soft backdrop-blur">
        <p className="text-sm text-text-muted">电子书详情</p>
        <h2 className="mt-2 text-3xl font-semibold text-text-primary">{detail.title}</h2>
      </section>

      <section className="grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)]">
        <Card className="p-6">
          <img
            alt={detail.title}
            className="mx-auto aspect-[3/4] w-full max-w-[260px] rounded-3xl object-cover shadow-lg"
            src={detail.cover || "https://placehold.co/600x800/e2e8f0/334155?text=Book"}
          />
          <div className="mt-6 flex flex-wrap gap-2">
            <span className="rounded-full bg-accent-soft px-3 py-1 text-xs font-medium text-accent">电子书</span>
            {detail.is_vip_book ? (
              <span className="rounded-full bg-warning-soft px-3 py-1 text-xs font-medium text-warning">会员可读</span>
            ) : null}
            {detail.is_on_bookshelf ? (
              <span className="rounded-full bg-success-soft px-3 py-1 text-xs font-medium text-success">已在书架</span>
            ) : null}
          </div>
          <div className="mt-6 space-y-3 text-sm text-text-secondary">
            <div className="rounded-2xl bg-surface-soft p-4">作者：{detail.book_author || detail.author_list.join(" / ") || "未知"}</div>
            <div className="rounded-2xl bg-surface-soft p-4">出版社：{pressName}</div>
            <div className="rounded-2xl bg-surface-soft p-4">分类：{detail.classify_name || "未分类"}</div>
            <div className="rounded-2xl bg-surface-soft p-4">出版时间：{detail.publish_time || "未知"}</div>
            <div className="rounded-2xl bg-surface-soft p-4">阅读时长：{detail.read_time || 0} 分钟</div>
          </div>
        </Card>

        <div className="space-y-6">
          <DownloadActionsPanel
            description="将当前电子书直接下载到服务端本地 output 目录，支持 HTML、PDF 与 EPUB 三种格式。"
            onDownload={handleEbookDownload}
            options={ebookDownloadOptions}
            title="电子书下载"
          />

          <Card className="p-6">
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <div className="rounded-2xl bg-surface-soft p-4">
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <BookOpen className="h-4 w-4" />
                  章节数
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">{chapterCount}</p>
              </div>
              <div className="rounded-2xl bg-surface-soft p-4">
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <Star className="h-4 w-4" />
                  评分
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">{detail.product_score || detail.douban_score || "暂无"}</p>
              </div>
              <div className="rounded-2xl bg-surface-soft p-4">
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <Clock3 className="h-4 w-4" />
                  试读
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">
                  {detail.can_trial_read ? `${detail.trial_read_proportion || "可试读"}` : "不可试读"}
                </p>
              </div>
              <div className="rounded-2xl bg-surface-soft p-4">
                <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                  <NotebookPen className="h-4 w-4" />
                  笔记数
                </div>
                <p className="mt-3 text-2xl font-semibold text-text-primary">{data.notes?.list.length ?? 0}</p>
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">内容简介</h3>
            <p className="mt-4 whitespace-pre-wrap text-sm leading-7 text-text-secondary">
              {detail.book_intro || detail.author_info || "暂无简介"}
            </p>
            {pressBrief ? (
              <div className="mt-5 rounded-2xl bg-surface-soft p-4">
                <h4 className="text-sm font-medium text-text-primary">出版社简介</h4>
                <p className="mt-2 whitespace-pre-wrap text-sm leading-7 text-text-secondary">{pressBrief}</p>
              </div>
            ) : null}
            <div className="mt-5 flex flex-wrap gap-3">
              {isOnBookshelf ? (
                <Button
                  disabled={shelfLoading}
                  onClick={handleRemoveFromShelf}
                  variant="danger"
                >
                  {shelfLoading ? "处理中..." : "移出书架"}
                </Button>
              ) : (
                <Button disabled={shelfLoading} onClick={handleAddToShelf}>
                  {shelfLoading ? "处理中..." : "加入书架"}
                </Button>
              )}
              <Button
                onClick={() => navigate(`/ebooks/${encodeURIComponent(enid)}/comments?title=${encodeURIComponent(detail.title || "电子书书评")}`)}
                variant="outline"
              >
                查看书评
              </Button>
              {detail.add_studylist_dd_url ? (
                <Button onClick={() => window.open(detail.add_studylist_dd_url, "_blank", "noopener,noreferrer")}>
                  去得到查看完整电子书
                </Button>
              ) : null}
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">目录</h3>
            <div className="mt-4">
              {catalogList.length > 0 ? (
                <div className="max-h-[36rem] overflow-y-auto rounded-2xl border border-border bg-surface-soft/60">
                  {catalogList.map((item, index) => {
                    const level = normalizeCatalogLevel(item.level)

                    return (
                      <div
                        className={cn(
                          "border-b border-border/80 px-3 py-2.5 last:border-b-0",
                          getCatalogRowClass(level),
                        )}
                        key={`${item.playOrder}-${item.href || index}`}
                      >
                        <div
                          className="flex items-start gap-2"
                          style={{ paddingLeft: `${getCatalogIndent(level)}px` }}
                        >
                          <span
                            className={cn(
                              "mt-[3px] h-1.5 w-1.5 flex-none rounded-full",
                              getCatalogDotClass(level),
                            )}
                          />
                          <p className={cn("leading-6", getCatalogTextClass(level))}>
                            {item.text || `第 ${index + 1} 章`}
                          </p>
                        </div>
                      </div>
                    )
                  })}
                </div>
              ) : (
                <div className="rounded-2xl border border-dashed border-border p-6 text-sm text-text-muted">
                  当前电子书暂无目录信息。
                </div>
              )}
            </div>
          </Card>

          <Card className="p-6">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-semibold text-text-primary">笔记预览</h3>
              <Button onClick={() => navigate("/category")} variant="ghost">
                返回分类页
              </Button>
            </div>

            {data.notesError ? (
              <div className="mt-4 rounded-2xl border border-warning bg-warning-soft p-4 text-sm text-warning">
                {data.notesError}
              </div>
            ) : null}

            <div className="mt-4 space-y-3">
              {notes.length > 0 ? (
                notes.map((note) => (
                  <div className={`${semanticSubtlePanelClass} p-4`} key={note.note_id}>
                    <div className="flex items-center gap-3">
                      <img
                        alt={note.notes_owner?.nickname || "avatar"}
                        className="h-9 w-9 rounded-full object-cover"
                        src={note.notes_owner?.avatar || "https://placehold.co/80x80/e2e8f0/334155?text=DD"}
                      />
                      <div>
                        <p className="text-sm font-medium text-text-primary">{note.notes_owner?.nickname || "得到用户"}</p>
                        <p className="text-xs text-text-muted">{formatTime(note.create_time)}</p>
                      </div>
                    </div>
                    <p className={`mt-3 text-sm leading-7 ${semanticSecondaryTextClass}`}>
                      {note.note || note.content || note.note_line || "暂无笔记内容"}
                    </p>
                  </div>
                ))
              ) : (
                <div className="rounded-2xl border border-dashed border-border p-6 text-sm text-text-muted">
                  当前电子书还没有可展示的笔记。
                </div>
              )}
            </div>
          </Card>
        </div>
      </section>
    </main>
  )
}
