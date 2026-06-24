import { ArrowLeft, Clock3, FileText, Headphones, Loader2, Play, Rows4, SortAsc, SortDesc } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { api, type CourseArticleItem, type CourseArticleListResponse, type CourseInfoResponse } from "@/api"
import { QuickDownloadButtons, type QuickDownloadOption } from "@/components/download/QuickDownloadButtons"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { getSemanticStatusBadgeClass, semanticMetaTextClass, semanticPageSectionClass } from "@/lib/semanticStyles"
import { useAudioPlayer } from "@/providers/AudioPlayerProvider"

const articlePageSize = 30
const chapterQuickDownloadOptions: QuickDownloadOption[] = [
  { value: 1, label: "音频" },
  { value: 2, label: "PDF" },
  { value: 3, label: "MD" },
]

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

    openCourseDetail()
  }

  const backLabel =
    searchParams.get("from") === "algo" ? "返回分类结果" : searchParams.get("from") === "home" ? "返回首页" : "返回课程详情"

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

  const openCourseDetail = () => {
    navigate(
      `/courses/${encodeURIComponent(enid)}?from=${encodeURIComponent(searchParams.get("from") || "purchased-course")}&parentTitle=${encodeURIComponent(courseTitle)}`,
    )
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
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载课程文章列表...
        </div>
      </main>
    )
  }

  if (error || !courseInfo?.class_info) {
    return (
      <main className="space-y-6">
        <Card className="border-danger bg-danger-soft p-6 text-sm text-danger">
          {error || "未找到课程文章列表"}
        </Card>
      </main>
    )
  }

  return (
    <main className="space-y-6">
      <div className="flex flex-wrap items-center gap-2 px-1 text-sm text-text-muted">
        <button
          className="inline-flex items-center rounded-lg px-2 py-1 transition hover:bg-surface-panel hover:text-text-secondary"
          onClick={backAction}
          type="button"
        >
          <ArrowLeft className="mr-1 h-4 w-4" />
          {backLabel}
        </button>
      </div>

      <section className={`${semanticPageSectionClass} p-6`}>
        <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
          <div className="flex gap-5">
            <img
              alt={courseTitle}
              className="h-28 w-28 rounded-3xl object-cover shadow-soft"
              src={courseInfo.class_info.square_img || courseInfo.class_info.index_img || "https://placehold.co/240x240/e2e8f0/334155?text=Course"}
            />
            <div className="min-w-0">
              <p className={semanticMetaTextClass}>课程文章列表</p>
              <h2 className="mt-2 text-3xl font-semibold text-text-primary">{courseTitle}</h2>
              <p className="mt-3 max-w-3xl text-sm leading-7 text-text-secondary">
                {courseInfo.class_info.intro || courseInfo.class_info.share_summary || "暂无课程介绍"}
              </p>
              <div className="mt-4 flex flex-wrap gap-3 text-sm text-text-muted">
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {courseInfo.class_info.lecturer_name || "未知讲师"}
                </span>
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {courseInfo.class_info.current_article_count || articleItems.length} 篇内容
                </span>
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {courseInfo.class_info.learn_user_count || 0} 人学习
                </span>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
            <Button disabled={playerTracks.length === 0} onClick={() => setQueue(playerTracks, 0)} variant="outline">
              <Play className="mr-2 h-4 w-4" />
              播放全部
            </Button>
            <Button onClick={openCourseDetail} variant="outline">
              <FileText className="mr-2 h-4 w-4" />
              课程详情
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
                {item.is_read ? <span className={getSemanticStatusBadgeClass("success")}>已读</span> : null}
                {item.audio?.mp3_play_url ? (
                  <span className={getSemanticStatusBadgeClass("accent")}>含音频</span>
                ) : null}
                {item.video_status === 1 ? <span className={getSemanticStatusBadgeClass("warning")}>含视频</span> : null}
              </div>

              <div className="min-w-0 flex-1">
                <h3 className="line-clamp-2 text-xl font-semibold text-text-primary">{item.title}</h3>
                <p className="mt-3 line-clamp-3 text-sm leading-7 text-text-secondary">{item.summary || "暂无摘要"}</p>
              </div>

              <div className="flex flex-wrap gap-3 text-sm text-text-muted">
                <span className="inline-flex items-center gap-1">
                  <Headphones className="h-4 w-4" />
                  {item.cur_learn_count || 0} 人学习
                </span>
                <span className="inline-flex items-center gap-1">
                  <Clock3 className="h-4 w-4" />
                  {formatPublishTime(item.publish_time)}
                </span>
              </div>

              <div className="flex flex-wrap gap-3">
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
                <QuickDownloadButtons
                  onDownload={(downloadType) =>
                    api.download.course({
                      enid,
                      title: `${courseTitle} - ${item.title}`,
                      articleId: item.id,
                      downloadType,
                      isOrder: true,
                    })
                  }
                  options={chapterQuickDownloadOptions}
                />
              </div>
            </div>
          </Card>
        ))}
      </section>

      {articleItems.length === 0 ? (
        <Card className="p-10 text-center text-text-muted">
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
