import { ArrowLeft, Loader2, MessageSquare, Star } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { api, type EbookCommentItem, type EbookCommentResponse } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"

function formatTime(value?: number) {
  if (!value) {
    return "未知时间"
  }

  return new Date(value * 1000).toLocaleString("zh-CN")
}

function renderStyledNote(raw?: string) {
  if (!raw) {
    return ""
  }

  try {
    const blocks = JSON.parse(raw) as Array<{
      type: string
      ordered?: boolean
      contents?: Array<any>
    }>

    return blocks
      .map((block) => {
        if (block.type === "paragraph") {
          const content = (block.contents || [])
            .map((item) => {
              if (item.type !== "text") {
                return ""
              }

              const text = String(item.text?.content || "").trim()
              if (!text) {
                return ""
              }

              return item.text?.bold ? `<strong>${text}</strong>` : text
            })
            .join("")

          return content ? `<p>${content}</p>` : ""
        }

        if (block.type === "list") {
          const tag = block.ordered ? "ol" : "ul"
          const items = (block.contents || [])
            .map((row) => {
              const rowText = (row || [])
                .map((item: any) => {
                  if (item.type !== "text") {
                    return ""
                  }

                  const text = String(item.text?.content || "").trim()
                  if (!text) {
                    return ""
                  }

                  return item.text?.bold ? `<strong>${text}</strong>` : text
                })
                .join("")

              return rowText ? `<li>${rowText}</li>` : ""
            })
            .join("")

          return items ? `<${tag}>${items}</${tag}>` : ""
        }

        return ""
      })
      .join("")
  } catch {
    return raw
  }
}

function CommentCard({ item }: { item: EbookCommentItem }) {
  const content = renderStyledNote(item.note_line_style || item.note_line || "")

  return (
    <div className="break-inside-avoid pb-5">
      <Card className="overflow-hidden border-slate-200 bg-white/95 p-0 shadow-soft">
        <div className="border-b border-slate-100 p-4">
          <div className="flex items-start gap-3">
            <img
              alt={item.notes_owner?.name || "得到用户"}
              className="h-10 w-10 rounded-full object-cover"
              src={item.notes_owner?.avatar || "https://placehold.co/80x80/e2e8f0/334155?text=DD"}
            />
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <p className="truncate text-sm font-medium text-slate-900">{item.notes_owner?.name || "得到用户"}</p>
                {item.notes_owner?.slogan ? (
                  <span className="rounded-full bg-sky-50 px-2 py-0.5 text-[11px] text-sky-700">{item.notes_owner.slogan}</span>
                ) : null}
              </div>
              <p className="mt-1 text-xs text-slate-500">{formatTime(item.create_time)}</p>
            </div>
          </div>
          {item.note_title ? <h3 className="mt-3 line-clamp-2 text-base font-semibold text-slate-950">{item.note_title}</h3> : null}
        </div>

        <div className="space-y-3 p-4">
          {content ? (
            <div
              className="prose prose-sm max-w-none text-slate-700 prose-p:my-2 prose-li:my-1 prose-ol:pl-5 prose-ul:pl-5"
              dangerouslySetInnerHTML={{ __html: content }}
            />
          ) : (
            <p className="text-sm leading-7 text-slate-600">{item.note_line || "暂无内容"}</p>
          )}

          <div className="flex items-center justify-between text-xs text-slate-400">
            <span>点赞 {item.notes_count?.like_count || 0}</span>
            <span>评论 {item.notes_count?.comment_count || 0}</span>
          </div>
        </div>
      </Card>
    </div>
  )
}

export function EbookCommentPage() {
  const navigate = useNavigate()
  const { enid = "" } = useParams()
  const [searchParams] = useSearchParams()
  const [data, setData] = useState<EbookCommentResponse | null>(null)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const title = searchParams.get("title") || "电子书书评"

  useEffect(() => {
    setData(null)
    setPage(1)
  }, [enid])

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      if (!enid) {
        return
      }

      if (page === 1) {
        setLoading(true)
      } else {
        setLoadingMore(true)
      }
      setError(null)

      try {
        const result = await api.ebook.comments(enid, page, 15)
        if (cancelled) {
          return
        }

        setData((current) => {
          if (!current || page === 1) {
            return result
          }

          return {
            ...result,
            list: [...current.list, ...result.list],
          }
        })
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "电子书书评加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
          setLoadingMore(false)
        }
      }
    }

    void load()

    return () => {
      cancelled = true
    }
  }, [enid, page])

  const hasMore = useMemo(() => {
    return (data?.list.length || 0) < (data?.total || 0)
  }, [data?.list.length, data?.total])

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-slate-500">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载电子书书评...
        </div>
      </main>
    )
  }

  if (error) {
    return (
      <main className="space-y-6">
        <Card className="border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700">{error}</Card>
      </main>
    )
  }

  return (
    <main className="space-y-6">
      <section className="rounded-3xl border border-white/70 bg-white/90 p-6 shadow-soft backdrop-blur">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div>
            <p className="text-sm text-slate-500">电子书书评</p>
            <h2 className="mt-2 text-3xl font-semibold text-slate-950">{title}</h2>
            <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">
              这里承接已购电子书的评分和书评内容，先对齐 `dedao-gui` 的核心消费链路：看整体评分、浏览读者书评，再决定是否回到详情页继续阅读。
            </p>
          </div>
          <Button onClick={() => navigate("/purchased/ebooks")} variant="outline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            返回已购电子书
          </Button>
        </div>
        <div className="mt-5 flex flex-wrap gap-3 text-sm text-slate-500">
          <span className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-3 py-1.5">
            <MessageSquare className="h-4 w-4" />
            总评论 {data?.total || 0}
          </span>
          <span className="inline-flex items-center gap-2 rounded-full bg-amber-50 px-3 py-1.5 text-amber-700">
            <Star className="h-4 w-4" />
            平均评分 {Number(data?.ebook_score.average_score || 0).toFixed(1)}
          </span>
        </div>
      </section>

      {data?.list.length ? (
        <section className="columns-1 gap-5 md:columns-2 xl:columns-3">
          {data.list.map((item) => (
            <CommentCard item={item} key={item.note_id} />
          ))}
        </section>
      ) : (
        <Card className="p-10 text-center text-slate-500">
          <p className="text-lg font-medium text-slate-900">当前还没有可展示的电子书书评</p>
          <p className="mt-2 text-sm">后续如果需要，还可以继续补我的评分和更多排序能力。</p>
        </Card>
      )}

      {hasMore ? (
        <div className="flex justify-center">
          <Button disabled={loadingMore} onClick={() => setPage((current) => current + 1)} variant="outline">
            {loadingMore ? "加载中..." : "加载更多"}
          </Button>
        </div>
      ) : null}
    </main>
  )
}
