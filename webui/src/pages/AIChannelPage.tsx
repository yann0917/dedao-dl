import { BookOpen, Clock, Crown, Loader2, Play, RefreshCcw, Sparkles, Users } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { useNavigate } from "react-router-dom"
import {
  api,
  type ChannelHomepageCategory,
  type ChannelInfo,
  type ChannelItem,
  type ChannelTopicCategory,
  type ChannelVipInfo,
} from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"
import {
  getSemanticChipClass,
  semanticMetaTextClass,
  semanticPageSectionClass,
} from "@/lib/semanticStyles"

// AI 学习圈固定频道 ID
const AI_CHANNEL_ID = 1000

// 本地存储 key：记住上次选中的分类（一级 + 二级）
const CHANNEL_STORAGE_KEY = `ai_channel_state_${AI_CHANNEL_ID}`

function normalizeAssetUrl(value?: string) {
  return (value || "").replace(/`/g, "").trim()
}

function formatCount(value?: number) {
  if (!value) {
    return "0"
  }
  return new Intl.NumberFormat("zh-CN").format(value)
}

function formatDuration(seconds?: number) {
  if (!seconds) {
    return ""
  }
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return hours > 0 ? `${hours}时${minutes}分` : `${minutes}分`
}

function getDifficultyLabel(level?: number) {
  const map: Record<number, string> = { 1: "入门", 2: "进阶", 3: "高阶" }
  return level ? map[level] || "" : ""
}

function getDifficultyClass(level?: number) {
  const map: Record<number, string> = {
    1: "bg-emerald-600/90",
    2: "bg-amber-500/90",
    3: "bg-rose-600/90",
  }
  return level ? map[level] || "bg-black/55" : "bg-black/55"
}

function getProductTypeName(productType?: number) {
  const map: Record<number, string> = { 65: "视频", 66: "课程" }
  return productType ? map[productType] || `类型${productType}` : ""
}

export function AIChannelPage() {
  const navigate = useNavigate()
  const [channelInfo, setChannelInfo] = useState<ChannelInfo | null>(null)
  const [vipInfo, setVipInfo] = useState<ChannelVipInfo | null>(null)
  const [categories, setCategories] = useState<ChannelHomepageCategory[]>([])
  const [activeCategoryId, setActiveCategoryId] = useState<number | null>(null)
  const [activeSubcategoryId, setActiveSubcategoryId] = useState<number | null>(null)
  const [currentItems, setCurrentItems] = useState<ChannelItem[]>([])
  const [loading, setLoading] = useState(true)
  const [listLoading, setListLoading] = useState(false)
  const [error, setError] = useState("")

  const currentCategory = useMemo(
    () => categories.find((cat) => cat.category_id === activeCategoryId) ?? null,
    [categories, activeCategoryId],
  )

  const currentSubcategories = currentCategory?.list ?? []

  const activeSubcategory = useMemo(
    () => currentSubcategories.find((sub) => sub.id === activeSubcategoryId) ?? null,
    [currentSubcategories, activeSubcategoryId],
  )

  // 拉取某个二级分类下的完整内容列表
  const loadTopicItems = async (subcategory: ChannelTopicCategory) => {
    setListLoading(true)

    try {
      const detail = await api.channel.topicDetail(subcategory.id)
      setCurrentItems(detail?.items ?? subcategory.items ?? [])
    } catch {
      // 详情接口失败时回退到首页内嵌的 items
      setCurrentItems(subcategory.items ?? [])
    } finally {
      setListLoading(false)
    }
  }

  const saveState = (categoryId: number | null, subcategoryId: number | null) => {
    try {
      window.localStorage.setItem(CHANNEL_STORAGE_KEY, JSON.stringify({ categoryId, subcategoryId }))
    } catch {
      // 忽略存储失败
    }
  }

  const restoreState = (): { categoryId: number; subcategoryId: number | null } | null => {
    try {
      const raw = window.localStorage.getItem(CHANNEL_STORAGE_KEY)
      if (!raw) {
        return null
      }
      const parsed = JSON.parse(raw) as { categoryId?: number; subcategoryId?: number | null }
      if (typeof parsed.categoryId === "number") {
        return { categoryId: parsed.categoryId, subcategoryId: parsed.subcategoryId ?? null }
      }
    } catch {
      // 忽略解析失败
    }
    return null
  }

  const loadAll = async () => {
    setLoading(true)
    setError("")

    try {
      const [info, homepage, vip] = await Promise.all([
        api.channel.info(AI_CHANNEL_ID),
        api.channel.homepage(AI_CHANNEL_ID),
        api.channel.vip(AI_CHANNEL_ID).catch(() => null),
      ])

      setChannelInfo(info)
      setCategories(homepage)
      setVipInfo(vip)

      if (homepage.length === 0) {
        return
      }

      // 恢复上次选中的分类；不存在或无效时回退到第一个一级分类 / 第一个二级分类
      const saved = restoreState()
      const savedCategory = saved ? homepage.find((cat) => cat.category_id === saved.categoryId) : undefined
      const targetCategory = savedCategory ?? homepage[0]
      const targetSubcategory =
        savedCategory?.list.find((sub) => sub.id === saved?.subcategoryId) ?? targetCategory.list[0] ?? undefined

      setActiveCategoryId(targetCategory.category_id)

      if (targetSubcategory) {
        setActiveSubcategoryId(targetSubcategory.id)
        void loadTopicItems(targetSubcategory)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "AI 学习圈加载失败")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadAll()
  }, [])

  // 切换一级分类时重置二级分类选中状态并记住
  const handleCategoryChange = (categoryId: number) => {
    setActiveCategoryId(categoryId)
    setActiveSubcategoryId(null)
    setCurrentItems([])
    saveState(categoryId, null)
  }

  // 点击二级分类，拉取该主题下的完整内容列表并记住
  const handleSubcategoryClick = async (subcategory: ChannelTopicCategory) => {
    setActiveSubcategoryId(subcategory.id)
    saveState(activeCategoryId, subcategory.id)
    await loadTopicItems(subcategory)
  }

  const handleItemClick = (item: ChannelItem) => {
    const enId = String(item.en_id ?? "").trim()

    // 视频 → 文章详情
    if (item.product_type === 65 && enId) {
      navigate(
        `/articles/1/${encodeURIComponent(enId)}?from=aiChannel&parentTitle=${encodeURIComponent(
          channelInfo?.title ?? "AI 学习圈",
        )}`,
      )
      return
    }

    // 课程 → 课程详情
    if (item.class_en_id) {
      navigate(`/courses/${encodeURIComponent(item.class_en_id)}`)
    }
  }

  const subscriberCount = channelInfo?.statistics?.total_subscribers ?? 0
  const vipSurplusDays = vipInfo?.surplus_days ?? 0

  return (
    <main className="space-y-6">
      {error ? (
        <Card className="border-danger bg-danger-soft">
          <div className="flex flex-col gap-3 p-4 text-sm text-danger md:flex-row md:items-center md:justify-between">
            <span>{error}</span>
            <Button onClick={() => void loadAll()} variant="outline">
              重试
            </Button>
          </div>
        </Card>
      ) : null}

      {loading ? (
        <div className="flex min-h-[18rem] items-center justify-center gap-3 text-sm text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载 AI 学习圈...
        </div>
      ) : (
        <>
          {/* 频道头部 */}
          {channelInfo ? (
            <section className={`${semanticPageSectionClass} overflow-hidden`}>
              <div className="flex flex-col gap-6 px-6 py-6 lg:flex-row lg:items-center lg:justify-between">
                <div className="flex min-w-0 items-center gap-4">
                  {channelInfo.logo ? (
                    <img
                      alt={channelInfo.title}
                      className="h-16 w-16 shrink-0 rounded-2xl border border-border object-cover"
                      src={normalizeAssetUrl(channelInfo.logo)}
                    />
                  ) : (
                    <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-surface-soft text-text-muted">
                      <Sparkles className="h-7 w-7" />
                    </div>
                  )}

                  <div className="min-w-0 space-y-2">
                    <div className="flex flex-wrap items-center gap-2">
                      <h1 className="text-2xl font-semibold text-text-primary">{channelInfo.title}</h1>
                      {channelInfo.is_vip ? (
                        <span className="rounded-full bg-accent-soft px-3 py-1 text-xs font-medium text-accent">VIP</span>
                      ) : null}
                    </div>
                    {channelInfo.description ? (
                      <p className="line-clamp-2 max-w-2xl text-sm text-text-secondary">{channelInfo.description}</p>
                    ) : null}
                  </div>
                </div>

                <div className="grid shrink-0 grid-cols-2 gap-3">
                  <div className="rounded-2xl bg-surface-soft p-4">
                    <div className="flex items-center gap-2 text-xs text-text-muted">
                      <Users className="h-3.5 w-3.5" />
                      已加入
                    </div>
                    <p className="mt-2 text-2xl font-semibold text-text-primary">{formatCount(subscriberCount)}</p>
                  </div>
                  <div className="rounded-2xl bg-surface-soft p-4">
                    <div className="flex items-center gap-2 text-xs text-text-muted">
                      <Crown className="h-3.5 w-3.5" />
                      VIP 状态
                    </div>
                    <p className="mt-2 text-2xl font-semibold text-text-primary">
                      {vipInfo?.is_vip ? `${vipSurplusDays}天` : vipInfo?.is_ever_subscribed ? "已过期" : "未开通"}
                    </p>
                  </div>
                </div>
              </div>
            </section>
          ) : null}

          {/* 一级分类 Tab */}
          {categories.length > 0 ? (
            <section className={`${semanticPageSectionClass} p-6`}>
              <div className="flex flex-wrap gap-2">
                {categories.map((cat) => (
                  <button
                    className={getSemanticChipClass(activeCategoryId === cat.category_id, "strong")}
                    key={cat.category_id}
                    onClick={() => handleCategoryChange(cat.category_id)}
                    type="button"
                  >
                    {cat.category_name}
                  </button>
                ))}
              </div>

              {/* 二级分类 + 内容列表 */}
              <div className="mt-6 flex gap-6">
                {/* 二级分类侧栏 */}
                {currentSubcategories.length > 0 ? (
                  <aside className="w-64 shrink-0 space-y-1">
                    {currentSubcategories.map((sub) => (
                      <button
                        className={cn(
                          "flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-left transition",
                          activeSubcategoryId === sub.id
                            ? "bg-accent-soft text-accent"
                            : "text-text-secondary hover:bg-surface-soft hover:text-text-primary",
                        )}
                        key={sub.id}
                        onClick={() => void handleSubcategoryClick(sub)}
                        type="button"
                      >
                        <div className="min-w-0 flex-1">
                          <p className="truncate text-sm font-medium">{sub.title}</p>
                          {sub.intro ? (
                            <p className="mt-0.5 line-clamp-1 text-xs text-text-muted">{sub.intro}</p>
                          ) : null}
                        </div>
                        {sub.length ? (
                          <span className="shrink-0 rounded-full bg-surface-soft px-2 py-0.5 text-xs text-text-muted">
                            {sub.length}
                          </span>
                        ) : null}
                      </button>
                    ))}
                  </aside>
                ) : null}

                {/* 内容列表 */}
                <div className="min-w-0 flex-1">
                  {activeSubcategory ? (
                    <div className="mb-4 flex items-center justify-between">
                      <h2 className="text-lg font-semibold text-text-primary">{activeSubcategory.title}</h2>
                      {!listLoading ? (
                        <span className={`${semanticMetaTextClass}`}>{currentItems.length} 项内容</span>
                      ) : null}
                    </div>
                  ) : null}

                  {listLoading ? (
                    <div className="flex min-h-[12rem] items-center justify-center gap-3 text-sm text-text-muted">
                      <Loader2 className="h-5 w-5 animate-spin" />
                      正在加载内容...
                    </div>
                  ) : currentItems.length > 0 ? (
                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                      {currentItems.map((item, index) => {
                        const cover = normalizeAssetUrl(item.cover || item.logo)
                        return (
                          <Card
                            className="h-full overflow-hidden p-0 transition hover:-translate-y-1 hover:shadow-lg"
                            key={`${item.product_id}-${index}`}
                          >
                            <button
                              className="flex h-full w-full flex-col text-left"
                              onClick={() => handleItemClick(item)}
                              type="button"
                            >
                              <div className="relative aspect-video overflow-hidden bg-surface-soft">
                                {cover ? (
                                  <img alt={item.title} className="h-full w-full object-cover" src={cover} />
                                ) : (
                                  <div className="flex h-full items-center justify-center text-text-muted">
                                    <BookOpen className="h-8 w-8" />
                                  </div>
                                )}

                                {getDifficultyLabel(item.difficulty_level) ? (
                                  <span
                                    className={cn(
                                      "absolute left-2 top-2 rounded-md px-2 py-0.5 text-xs text-white backdrop-blur",
                                      getDifficultyClass(item.difficulty_level),
                                    )}
                                  >
                                    {getDifficultyLabel(item.difficulty_level)}
                                  </span>
                                ) : null}
                                <span className="absolute right-2 top-2 rounded-md bg-black/55 px-2 py-0.5 text-xs text-white backdrop-blur">
                                  {getProductTypeName(item.product_type)}
                                </span>
                              </div>

                              <div className="flex flex-1 flex-col gap-2 p-4">
                                <h3 className="line-clamp-2 text-sm font-semibold text-text-primary">{item.title}</h3>
                                {item.summary ? (
                                  <p className="line-clamp-2 text-xs leading-6 text-text-secondary">{item.summary}</p>
                                ) : null}

                                <div className="mt-auto flex items-center gap-3 pt-2 text-xs text-text-muted">
                                  {formatDuration(item.duration) ? (
                                    <span className="flex items-center gap-1">
                                      <Clock className="h-3.5 w-3.5" />
                                      {formatDuration(item.duration)}
                                    </span>
                                  ) : null}
                                  {item.learn_count ? (
                                    <span className="flex items-center gap-1">
                                      <Play className="h-3.5 w-3.5" />
                                      {formatCount(item.learn_count)}
                                    </span>
                                  ) : null}
                                </div>
                              </div>
                            </button>
                          </Card>
                        )
                      })}
                    </div>
                  ) : activeSubcategory ? (
                    <div className="flex min-h-[12rem] items-center justify-center rounded-3xl bg-surface-soft text-sm text-text-muted">
                      该分类下暂无内容
                    </div>
                  ) : (
                    <div className="flex min-h-[12rem] items-center justify-center rounded-3xl bg-surface-soft text-sm text-text-muted">
                      请选择左侧分类查看内容
                    </div>
                  )}
                </div>
              </div>
            </section>
          ) : (
            <section className={`${semanticPageSectionClass} flex min-h-[18rem] items-center justify-center p-6`}>
              <div className="flex items-center gap-3 text-sm text-text-muted">
                <RefreshCcw className="h-4 w-4" />
                暂无分类数据
              </div>
            </section>
          )}
        </>
      )}
    </main>
  )
}
