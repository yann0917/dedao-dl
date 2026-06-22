import { ArrowLeft, Clock3, FileText, Headphones, Loader2, Play, Rows4, SortAsc, SortDesc } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { api, type CourseArticleItem, type CourseArticleListResponse, type CourseInfoResponse } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

const articlePageSize = 30

function formatPublishTime(value?: number) {
  if (!value) {
    return "未知发布时间"
  }

  return new Date(value * 1000).toLocaleDateString("zh-CN")
}

function buildArticleTrack(item: CourseArticleItem, courseTitle: string) {
  if (!item.audio?.mp3_play_url) {
    return null
  }

  return {
    id: item.enid,
    title: item.title,
    src: item.audio.mp3_play_url,
    poster: item.logo || item.audio.icon,
    subtitle: courseTitle,
  }
}

export function CourseArticleListPage() {
  const navigate = useNavigate()
  const { setQueue } = useAudioPlayer()
  const { enid = "" } = useParams()
  const [searchParams] = useSearchParams()
  const [reverse, setReverse] = useState(false)
  const [courseInfo, setCourseInfo] = useState<CourseInfoResponse | null>(null)
  const [articles, setArticles] = useState<CourseArticleListResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [finished, setFinished] = useState(false)
  const [maxId, setMaxId] = useState(0)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const loadCourseInfo = async () => {
      try {
        const nextCourseInfo = await api.course.info(enid)
        if (!cancelled) {
          setCourseInfo(nextCourseInfo)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "课程信息加载失败")
        }
      }
    }

    if (enid) {
      void loadCourseInfo()
    }

    return () => {
      cancelled = true
    }
  }, [enid])

  useEffect(() => {
    let cancelled = false

    const loadFirstPage = async () => {
      setLoading(true)
      setError(null)
      setFinished(false)
      setMaxId(0)

      try {
        const nextArticles = await api.course.articles(enid, {
          count: articlePageSize,
          maxId: 0,
          reverse,
        })

        if (!cancelled) {
          setArticles(nextArticles)
          const nextList = nextArticles.article_list ?? []
          setFinished(nextList.length < articlePageSize)
          setMaxId(nextList.length > 0 ? nextList[nextList.length - 1].id : 0)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "课程文章列表加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    if (enid) {
      void loadFirstPage()
    }

    return () => {
      cancelled = true
    }
  }, [enid, reverse])

  const courseTitle = courseInfo?.class_info.name || searchParams.get("parentTitle") || "课程文章"
  const articleItems = articles?.article_list ?? []

  const playerTracks = useMemo(
    () =>
      articleItems
        .map((item) => buildArticleTrack(item, courseTitle))
        .filter((item): item is NonNullable<ReturnType<typeof buildArticleTrack>> => item !== null),
    [articleItems, courseTitle],
  )

  const backAction = () => {
    const from = searchParams.get("from")
    if (from === "algo" || from === "home") {
      navigate(-1)
      return
    }

    navigate("/purchased/courses")
  }

  const backLabel =
    searchParams.get("from") === "algo" ? "返回分类结果" : searchParams.get("from") === "home" ? "返回首页" : "返回已购课程"

  const openArticleDetail = (item: CourseArticleItem) => {
    const listFrom = searchParams.get("from") || "course"
    navigate(
      `/articles/1/${encodeURIComponent(item.enid)}?from=course&listFrom=${encodeURIComponent(listFrom)}&parentEnid=${encodeURIComponent(enid)}&parentTitle=${encodeURIComponent(courseTitle)}`,
    )
  }

  const playArticle = (item: CourseArticleItem) => {
    const target = buildArticleTrack(item, courseTitle)
    if (!target) {
      return
    }

    const startIndex = playerTracks.findIndex((track) => track.id === target.id)
    setQueue(playerTracks, startIndex >= 0 ? startIndex : 0)
  }

  const loadMore = async () => {
    if (loadingMore || finished || !enid) {
      return
    }

    setLoadingMore(true)
    setError(null)

    try {
      const nextArticles = await api.course.articles(enid, {
        count: articlePageSize,
        maxId,
        reverse,
      })
      const nextList = nextArticles.article_list ?? []

      setArticles((current) => ({
        ...(current ?? nextArticles),
        ...nextArticles,
        article_list: [...(current?.article_list ?? []), ...nextList],
      }))
      setFinished(nextList.length < articlePageSize)
      if (nextList.length > 0) {
        setMaxId(nextList[nextList.length - 1].id)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "课程文章列表加载失败")
    } finally {
      setLoadingMore(false)
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-slate-500">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载课程文章列表...
        </div>
      </main>
    )
  }

  if (error || !courseInfo?.class_info) {
    return (
      <main className="space-y-6">
        <Card className="border border-rose-200 bg-rose-50 p-6 text-sm text-rose-700">
          {error || "未找到课程文章列表"}
        </Card>
      </main>
    )
  }

  return (
    <main className="space-y-6">
      <div className="flex flex-wrap items-center gap-2 px-1 text-sm text-slate-500">
        <button
          className="inline-flex items-center rounded-lg px-2 py-1 transition hover:bg-white/70 hover:text-slate-700"
          onClick={backAction}
          type="button"
        >
          <ArrowLeft className="mr-1 h-4 w-4" />
          {backLabel}
        </button>
      </div>

      <section className="rounded-3xl border border-white/70 bg-white/90 p-6 shadow-soft backdrop-blur">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
          <div className="flex gap-5">
            <img
              alt={courseTitle}
              className="h-28 w-28 rounded-3xl object-cover shadow-soft"
              src={courseInfo.class_info.square_img || courseInfo.class_info.index_img || "https://placehold.co/240x240/e2e8f0/334155?text=Course"}
            />
            <div className="min-w-0">
              <p className="text-sm text-slate-500">课程文章列表</p>
              <h2 className="mt-2 text-3xl font-semibold text-slate-950">{courseTitle}</h2>
              <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">
                {courseInfo.class_info.intro || courseInfo.class_info.share_summary || "这里承接课程内容消费链路，可继续进入单篇文章详情。"}
              </p>
              <div className="mt-4 flex flex-wrap gap-3 text-sm text-slate-500">
                <span className="rounded-full bg-slate-100 px-3 py-1.5">{courseInfo.class_info.lecturer_name || "未知讲师"}</span>
                <span className="rounded-full bg-slate-100 px-3 py-1.5">{courseInfo.class_info.current_article_count || articleItems.length} 篇内容</span>
                <span className="rounded-full bg-slate-100 px-3 py-1.5">{courseInfo.class_info.learn_user_count || 0} 人学习</span>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
            <Button disabled={playerTracks.length === 0} onClick={() => setQueue(playerTracks, 0)} variant="outline">
              <Play className="mr-2 h-4 w-4" />
              播放全部
            </Button>
            <Button onClick={() => setReverse((current) => !current)} variant="outline">
              {reverse ? <SortAsc className="mr-2 h-4 w-4" /> : <SortDesc className="mr-2 h-4 w-4" />}
              {reverse ? "切回正序" : "切换倒序"}
            </Button>
          </div>
        </div>
      </section>

      <section className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        {articleItems.map((item) => (
          <Card className="h-full p-5" key={item.enid}>
            <div className="flex h-full flex-col gap-4">
              <img
                alt={item.title}
                className="aspect-[16/9] w-full rounded-3xl object-cover"
                src={item.logo || courseInfo.class_info.square_img || "https://placehold.co/640x360/e2e8f0/334155?text=Article"}
              />

              <div className="flex flex-wrap gap-2">
                {item.is_read ? <span className="rounded-full bg-emerald-100 px-3 py-1 text-xs font-medium text-emerald-700">已读</span> : null}
                {item.audio?.mp3_play_url ? (
                  <span className="rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-700">含音频</span>
                ) : null}
                {item.video_status === 1 ? <span className="rounded-full bg-violet-100 px-3 py-1 text-xs font-medium text-violet-700">含视频</span> : null}
              </div>

              <div className="min-w-0 flex-1">
                <h3 className="line-clamp-2 text-xl font-semibold text-slate-950">{item.title}</h3>
                <p className="mt-3 line-clamp-3 text-sm leading-7 text-slate-600">{item.summary || "暂无摘要"}</p>
              </div>

              <div className="flex flex-wrap gap-3 text-sm text-slate-500">
                <span className="inline-flex items-center gap-1">
                  <Headphones className="h-4 w-4" />
                  {item.cur_learn_count || 0} 人学习
                </span>
                <span className="inline-flex items-center gap-1">
                  <Clock3 className="h-4 w-4" />
                  {formatPublishTime(item.publish_time)}
                </span>
              </div>

              <div className="flex gap-3">
                {item.audio?.mp3_play_url ? (
                  <Button onClick={() => playArticle(item)} variant="outline">
                    <Play className="mr-2 h-4 w-4" />
                    播放
                  </Button>
                ) : null}
                <Button onClick={() => openArticleDetail(item)}>
                  <FileText className="mr-2 h-4 w-4" />
                  查看文章
                </Button>
              </div>
            </div>
          </Card>
        ))}
      </section>

      {articleItems.length === 0 ? (
        <Card className="p-10 text-center text-slate-500">
          <Rows4 className="mx-auto h-8 w-8" />
          <p className="mt-3">当前课程还没有可展示的文章内容。</p>
        </Card>
      ) : null}

      {articleItems.length > 0 && !finished ? (
        <div className="flex justify-center">
          <Button disabled={loadingMore} onClick={() => void loadMore()} variant="outline">
            {loadingMore ? "加载中..." : "加载更多"}
          </Button>
        </div>
      ) : null}
    </main>
  )
}
