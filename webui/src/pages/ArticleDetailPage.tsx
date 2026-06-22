import { ArrowLeft, ChevronRight, Clock3, FileText, Headphones, Loader2, Play } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { api } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { renderMarkdownToHtml } from "@/lib/markdown"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

export function ArticleDetailPage() {
  const navigate = useNavigate()
  const { setQueue } = useAudioPlayer()
  const { aType = "2", enid = "" } = useParams()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [markdown, setMarkdown] = useState("")
  const [html, setHtml] = useState("")
  const [data, setData] = useState<Awaited<ReturnType<typeof api.article.detail>> | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      setError(null)
      try {
        const result = await api.article.detail(Number(aType) as 1 | 2, enid)
        if (!cancelled) {
          setData(result)
          setMarkdown(result.markdown)
          setHtml(renderMarkdownToHtml(result.markdown))
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "文章详情加载失败")
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
  }, [aType, enid])

  const articleInfo = data?.info?.article_info
  const classInfo = data?.info?.class_info
  const audio = data?.info?.audio

  const parentCrumb = useMemo(() => {
    const from = searchParams.get("from")
    const groupEnid = searchParams.get("groupEnid")
    const parentEnid = searchParams.get("parentEnid")
    const parentTitle = searchParams.get("parentTitle")
    const listFrom = searchParams.get("listFrom") || "course"

    if (from === "audio-group" && groupEnid) {
      return {
        label: parentTitle || "听书文稿",
        action: () => navigate(`/audio-groups/${encodeURIComponent(groupEnid)}/articles`),
      }
    }

    if (from === "audio" && parentEnid) {
      return {
        label: parentTitle || classInfo?.name || "听书详情",
        action: () => navigate(`/audios/${encodeURIComponent(parentEnid)}`),
      }
    }

    if (from === "course") {
      return {
        label: parentTitle || classInfo?.name || "课程文章",
        action: () => {
          if (parentEnid) {
            navigate(
              `/courses/${encodeURIComponent(parentEnid)}/articles?from=${encodeURIComponent(listFrom)}&parentTitle=${encodeURIComponent(
                parentTitle || classInfo?.name || "课程文章",
              )}`,
            )
            return
          }

          navigate("/purchased/courses")
        },
      }
    }

    return {
      label: parentTitle || classInfo?.name || "返回上级",
      action: () => navigate(-1),
    }
  }, [classInfo?.name, navigate, searchParams])

  const handlePlay = () => {
    if (!audio?.mp3_play_url || !articleInfo) {
      return
    }

    setQueue(
      [
        {
          id: enid,
          title: articleInfo.title,
          src: audio.mp3_play_url,
          poster: articleInfo.logo || audio.icon,
          subtitle: classInfo?.name || "文稿音频",
        },
      ],
      0,
    )
  }

  const handleExportMarkdown = () => {
    if (!articleInfo) {
      return
    }

    const blob = new Blob([markdown], { type: "text/markdown;charset=utf-8" })
    const url = URL.createObjectURL(blob)
    const link = document.createElement("a")
    link.href = url
    link.download = `${articleInfo.title || "article"}.md`
    link.click()
    URL.revokeObjectURL(url)
  }

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-slate-500">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载文稿详情...
        </div>
      </main>
    )
  }

  if (error || !articleInfo) {
    return (
      <main className="space-y-6">
        <Card className="border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700">
          {error || "未找到文稿详情"}
        </Card>
      </main>
    )
  }

  return (
    <main className="mx-auto max-w-5xl space-y-5">
      <div className="flex flex-wrap items-center gap-2 px-1 text-sm text-slate-500">
        <button
          className="inline-flex items-center rounded-lg px-2 py-1 transition hover:bg-white/70 hover:text-slate-700"
          onClick={parentCrumb.action}
          type="button"
        >
          <ArrowLeft className="mr-1 h-4 w-4" />
          {parentCrumb.label}
        </button>
        <ChevronRight className="h-4 w-4 text-slate-300" />
        <span className="font-medium text-slate-700">文稿</span>
      </div>

      <Card className="mx-auto max-w-4xl overflow-hidden">
        <div className="border-b border-slate-100 px-8 py-7">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="min-w-0">
              <p className="text-sm text-slate-500">文章详情</p>
              <h1 className="mt-2 text-3xl font-semibold leading-tight text-slate-950">{articleInfo.title}</h1>
              <p className="mt-4 max-w-3xl text-sm leading-7 text-slate-600">
                {articleInfo.summary || classInfo?.share_summary || "当前文稿已转换为 markdown 并在 Web 端直接渲染。"}
              </p>
            </div>

            <div className="flex flex-wrap gap-2">
              {audio?.mp3_play_url ? (
                <Button onClick={handlePlay} variant="outline">
                  <Play className="mr-2 h-4 w-4" />
                  播放音频
                </Button>
              ) : null}
              <Button onClick={handleExportMarkdown} variant="outline">
                <FileText className="mr-2 h-4 w-4" />
                导出 Markdown
              </Button>
            </div>
          </div>

          <div className="mt-6 flex flex-wrap gap-3 text-sm text-slate-500">
            <span className="rounded-full bg-slate-100 px-3 py-1.5">{classInfo?.name || "文稿"}</span>
            <span className="rounded-full bg-slate-100 px-3 py-1.5">{classInfo?.lecturer_name || "未知作者"}</span>
            <span className="inline-flex items-center rounded-full bg-slate-100 px-3 py-1.5">
              <Headphones className="mr-1.5 h-4 w-4" />
              {articleInfo.cur_learn_count || 0} 人学习
            </span>
            <span className="inline-flex items-center rounded-full bg-slate-100 px-3 py-1.5">
              <Clock3 className="mr-1.5 h-4 w-4" />
              {articleInfo.publish_time ? new Date(articleInfo.publish_time * 1000).toLocaleString("zh-CN") : "未知发布时间"}
            </span>
          </div>
        </div>

        <div className="bg-slate-50/40 px-8 py-10">
          <div className="mx-auto max-w-3xl">
            {(articleInfo.logo || classInfo?.square_img) ? (
              <img
                alt={articleInfo.title}
                className="mb-8 h-56 w-full rounded-3xl object-cover shadow-soft"
                src={articleInfo.logo || classInfo?.square_img}
              />
            ) : null}

            <article
              className={[
                "text-left text-[16px] leading-8 text-slate-700",
                "[&_h1]:mb-6 [&_h1]:mt-10 [&_h1]:text-3xl [&_h1]:font-semibold [&_h1]:text-slate-950",
                "[&_h2]:mb-4 [&_h2]:mt-10 [&_h2]:text-2xl [&_h2]:font-semibold [&_h2]:text-slate-950",
                "[&_h3]:mb-4 [&_h3]:mt-8 [&_h3]:text-xl [&_h3]:font-semibold [&_h3]:text-slate-900",
                "[&_p]:mb-5 [&_p]:leading-8 [&_p]:text-slate-700",
                "[&_blockquote]:my-6 [&_blockquote]:rounded-r-2xl [&_blockquote]:border-l-4 [&_blockquote]:border-primary [&_blockquote]:bg-white [&_blockquote]:px-5 [&_blockquote]:py-4 [&_blockquote]:text-slate-600",
                "[&_ul]:mb-5 [&_ul]:list-disc [&_ul]:space-y-2 [&_ul]:pl-6",
                "[&_ol]:mb-5 [&_ol]:list-decimal [&_ol]:space-y-2 [&_ol]:pl-6",
                "[&_li]:pl-1",
                "[&_img]:my-6 [&_img]:w-full [&_img]:rounded-2xl [&_img]:shadow-soft",
                "[&_hr]:my-8 [&_hr]:border-0 [&_hr]:border-t [&_hr]:border-slate-200",
                "[&_code]:rounded-md [&_code]:bg-slate-100 [&_code]:px-1.5 [&_code]:py-0.5 [&_code]:text-[0.95em]",
                "[&_h2>code]:bg-primary [&_h2>code]:text-white",
                "[&_em]:not-italic [&_em]:text-primary",
                "[&_strong]:font-semibold [&_strong]:text-slate-900",
              ].join(" ")}
              dangerouslySetInnerHTML={{ __html: html }}
            />
          </div>
        </div>
      </Card>
    </main>
  )
}
