import { ArrowLeft, Loader2, LibraryBig, MessageSquareText, Star, UserRound } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"
import { api, type CourseInfoResponse } from "@/api"
import { DownloadActionsPanel, type DownloadOption } from "@/components/download/DownloadActionsPanel"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  getSemanticStatusBadgeClass,
  semanticInfoBlockClass,
  semanticMetaTextClass,
  semanticPageSectionClass,
  semanticSecondaryTextClass,
} from "@/lib/semanticStyles"

function buildCourseArticlesPath(enid: string, parentTitle: string, from: string) {
  return `/courses/${encodeURIComponent(enid)}/articles?from=${encodeURIComponent(from)}&parentTitle=${encodeURIComponent(parentTitle)}`
}

function resolveCourseAccess(detail: CourseInfoResponse | null) {
  if (!detail) {
    return {
      canOpenArticles: false,
      actionLabel: "查看课程详情",
      badgeLabel: "未购",
      badgeVariant: "neutral" as const,
      statusText: "未购，先查看详情与介绍",
    }
  }

  if (detail.class_info.is_subscribe === 1) {
    return {
      canOpenArticles: true,
      actionLabel: "查看课程内容",
      badgeLabel: "已购",
      badgeVariant: "success" as const,
      statusText: "已购，可继续进入内容页",
    }
  }

  if (detail.class_info.is_in_vip || detail.class_info.is_vip) {
    return {
      canOpenArticles: true,
      actionLabel: "以会员身份查看内容",
      badgeLabel: "会员可看",
      badgeVariant: "warning" as const,
      statusText: "会员可看，可直接进入内容页",
    }
  }

  return {
    canOpenArticles: false,
    actionLabel: "查看课程详情",
    badgeLabel: "未购",
    badgeVariant: "neutral" as const,
    statusText: "未购，先查看详情与介绍",
  }
}

const courseDownloadOptions: DownloadOption[] = [
  { value: 1, label: "下载音频 MP3" },
  { value: 2, label: "下载课程 PDF" },
  { value: 3, label: "下载课程 Markdown" },
]

export function CourseDetailPage() {
  const navigate = useNavigate()
  const { enid = "" } = useParams()
  const [searchParams] = useSearchParams()
  const [detail, setDetail] = useState<CourseInfoResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const load = async () => {
      setLoading(true)
      setError(null)

      try {
        const result = await api.course.info(enid)
        if (!cancelled) {
          setDetail(result)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "课程详情加载失败")
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    if (enid) {
      void load()
    } else {
      setLoading(false)
      setError("缺少课程 enid")
    }

    return () => {
      cancelled = true
    }
  }, [enid])

  const title = detail?.class_info.name || searchParams.get("parentTitle") || "课程详情"
  const from = searchParams.get("from") || "course-detail"
  const highlights = useMemo(
    () =>
      detail?.items
        .filter((item) => item.content?.trim())
        .slice(0, 6) ?? [],
    [detail],
  )
  const commentItems = detail?.class_comment_info?.comment_list?.slice(0, 3) ?? []
  const averageScore = Number(detail?.class_comment_info?.average_score || "0")
  const access = useMemo(() => resolveCourseAccess(detail), [detail])
  const updateStatusText = useMemo(() => {
    if (!detail) {
      return ""
    }

    const phaseNum = detail.class_info.phase_num || 0
    const priceDesc = detail.class_info.price_desc || "讲"
    const currentArticleCount = detail.class_info.current_article_count || 0

    if (detail.class_info.is_finished === 1) {
      return `共${phaseNum}${priceDesc}`
    }

    return `共${phaseNum}${priceDesc}，已更新${currentArticleCount}${priceDesc}`
  }, [detail])

  const backAction = () => {
    if (from === "algo" || from === "home") {
      navigate(-1)
      return
    }

    navigate("/purchased/courses")
  }

  const backLabel =
    from === "algo" ? "返回分类结果" : from === "home" ? "返回首页" : "返回已购课程"

  const openArticles = () => {
    navigate(buildCourseArticlesPath(enid, title, from))
  }

  const handleCourseDownload = (downloadType: number) =>
    api.download.course({
      enid,
      title,
      downloadType,
      isOrder: true,
    })

  if (loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载课程详情...
        </div>
      </main>
    )
  }

  if (error || !detail?.class_info) {
    return (
      <main className="space-y-6">
        <Card className="border-danger bg-danger-soft p-6 text-sm text-danger">
          {error || "未找到课程详情"}
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
              alt={title}
              className="h-28 w-28 rounded-3xl object-cover shadow-soft"
              src={detail.class_info.square_img || detail.class_info.index_img || "https://placehold.co/240x240/e2e8f0/334155?text=Course"}
            />
            <div className="min-w-0">
              <p className={semanticMetaTextClass}>课程详情</p>
              <h2 className="mt-2 text-3xl font-semibold text-text-primary">{title}</h2>
              <p className={`mt-3 max-w-3xl text-sm leading-7 ${semanticSecondaryTextClass}`}>
                {detail.class_info.intro || detail.class_info.share_summary || "暂无课程介绍"}
              </p>
              <div className="mt-4 flex flex-wrap gap-3 text-sm text-text-muted">
                <span className={getSemanticStatusBadgeClass(access.badgeVariant, "px-3 py-1.5 text-sm")}>{access.badgeLabel}</span>
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {detail.class_info.lecturer_name_and_title || detail.class_info.lecturer_name || "未知讲师"}
                </span>
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {detail.class_info.current_article_count || 0} 篇内容
                </span>
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {detail.class_info.learn_user_count || 0} 人学习
                </span>
                <span className={getSemanticStatusBadgeClass("neutral", "px-3 py-1.5 text-sm")}>
                  {detail.class_info.price_desc || "课程详情"}
                </span>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
            {access.canOpenArticles ? (
              <Button onClick={openArticles}>
                <LibraryBig className="mr-2 h-4 w-4" />
                {access.actionLabel}
              </Button>
            ) : null}
          </div>
        </div>
      </section>

      <section className="grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)]">
        <Card className="p-6">
          <div className="space-y-5">
            <div>
              <p className={semanticMetaTextClass}>讲师信息</p>
              <div className="mt-4 flex items-center gap-4">
                <img
                  alt={detail.class_info.lecturer_name || "讲师"}
                  className="h-16 w-16 rounded-full object-cover"
                  src={detail.class_info.lecturer_avatar || detail.class_info.square_img || "https://placehold.co/128x128/e2e8f0/334155?text=Tutor"}
                />
                <div>
                  <h3 className="text-lg font-semibold text-text-primary">{detail.class_info.lecturer_name || "未知讲师"}</h3>
                  <p className={semanticMetaTextClass}>{detail.class_info.lecturer_title || "暂无头衔信息"}</p>
                </div>
              </div>
              <p className={`mt-4 text-sm leading-7 ${semanticSecondaryTextClass}`}>
                {detail.class_info.lecturer_intro || detail.class_info.share_summary || "暂无讲师介绍"}
              </p>
            </div>

            <div className="grid gap-3">
              <div className={semanticInfoBlockClass}>课程状态：{access.statusText}</div>
              <div className={semanticInfoBlockClass}>{updateStatusText || `课程内容：${detail.class_info.current_article_count || 0} 篇`}</div>
              <div className={semanticInfoBlockClass}>学习人数：{detail.class_info.learn_user_count || 0}</div>
              <div className={semanticInfoBlockClass}>
                收藏状态：{detail.class_info.collection?.is_collected ? "已收藏" : "未收藏"}
                {detail.class_info.collection?.collection_count ? ` · ${detail.class_info.collection.collection_count} 人收藏` : ""}
              </div>
              {detail.user_type ? <div className={semanticInfoBlockClass}>用户类型：{detail.user_type}</div> : null}
            </div>
          </div>
        </Card>

        <div className="space-y-6">
          <DownloadActionsPanel
            description="将当前课程内容直接下载到服务端本地 output 目录，并按章节自动整理文件。"
            disabled={!access.canOpenArticles}
            disabledReason="当前账号暂无课程访问权限，暂不支持直接下载。"
            onDownload={handleCourseDownload}
            options={courseDownloadOptions}
            title="课程下载"
          />

          <Card className="p-6">
            <div className="flex flex-wrap items-center gap-3">
              <span className={getSemanticStatusBadgeClass("warning", "inline-flex items-center gap-2 px-3 py-1.5 text-sm")}>
                <Star className="h-4 w-4" />
                {averageScore > 0 ? `${averageScore.toFixed(1)} 分` : "暂无评分"}
              </span>
              <span className={getSemanticStatusBadgeClass("neutral", "inline-flex items-center gap-2 px-3 py-1.5 text-sm")}>
                <MessageSquareText className="h-4 w-4" />
                {detail.class_comment_info?.count || detail.class_reviews_count || 0} 条评价
              </span>
              <span className={getSemanticStatusBadgeClass("neutral", "inline-flex items-center gap-2 px-3 py-1.5 text-sm")}>
                <UserRound className="h-4 w-4" />
                {detail.class_info.learn_user_count || 0} 人加入学习
              </span>
            </div>

            {commentItems.length > 0 ? (
              <div className="mt-4 grid gap-4 md:grid-cols-3">
                {commentItems.map((item) => (
                  <div className={semanticInfoBlockClass} key={item.id}>
                    <div className="flex items-center gap-3">
                      <img
                        alt={item.nickname || "学员"}
                        className="h-9 w-9 rounded-full object-cover"
                        src={item.avatar_s || "https://placehold.co/72x72/e2e8f0/334155?text=U"}
                      />
                      <div>
                        <p className="text-sm font-medium text-text-primary">{item.nickname || "匿名学员"}</p>
                        <p className="text-xs text-text-muted">{item.score ? `${item.score} 分评价` : "学员评价"}</p>
                      </div>
                    </div>
                    {item.title ? <p className="mt-3 text-sm font-medium text-text-primary">{item.title}</p> : null}
                    <p className="mt-2 text-sm leading-6 text-text-secondary">{item.no_style_content || "暂无评价内容"}</p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="mt-4 text-sm text-text-muted">当前暂无更多学员评价。</p>
            )}
          </Card>

          <Card className="p-6">
            <h3 className="text-xl font-semibold text-text-primary">课程亮点</h3>
            {detail.class_info.highlight ? (
              <div className={`${semanticInfoBlockClass} mt-4 whitespace-pre-wrap leading-7`}>
                {detail.class_info.highlight}
              </div>
            ) : null}
            {highlights.length > 0 ? (
              <div className="mt-4 space-y-4">
                {highlights.map((item, index) => (
                  <div className={semanticInfoBlockClass} key={`${item.title}-${index}`}>
                    <p className="font-medium text-text-primary">{item.title || `亮点 ${index + 1}`}</p>
                    <p className="mt-2 whitespace-pre-wrap text-sm leading-6 text-text-secondary">{item.content}</p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="mt-4 text-sm text-text-muted">当前暂无更多课程亮点。</p>
            )}
          </Card>

          {detail.class_info.outline_img ? (
            <Card className="p-6">
              <h3 className="text-xl font-semibold text-text-primary">课程大纲</h3>
              <div className="mt-4 overflow-hidden rounded-2xl border border-border bg-surface-soft">
                <img
                  alt={`${title} 课程大纲`}
                  className="w-full object-contain"
                  src={detail.class_info.outline_img}
                />
              </div>
            </Card>
          ) : null}
        </div>
      </section>
    </main>
  )
}
