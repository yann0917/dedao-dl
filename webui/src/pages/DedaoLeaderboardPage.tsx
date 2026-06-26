import { BookOpen, Loader2, RefreshCcw, Sparkles, Trophy } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { api, type RankBaseInfoResponse, type RankBoard, type RankListItem, type RankListResponse } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  getSemanticChipClass,
  getSemanticStatusBadgeClass,
  semanticMetaTextClass,
  semanticPageSectionClass,
} from "@/lib/semanticStyles"

function normalizeAssetUrl(value?: string) {
  return (value || "").replace(/`/g, "").trim()
}

function normalizeColor(value?: string) {
  const color = (value || "").replace("#", "").trim()
  return color ? `#${color}` : "#B88771"
}

function formatCount(value?: number) {
  if (!value) {
    return "0"
  }
  return new Intl.NumberFormat("zh-CN").format(value)
}

function resolveRankItemImage(item: RankListItem) {
  return normalizeAssetUrl(item.index_img || item.logo || item.square_img)
}

function resolveRankTypeBadge(board: RankBoard | null) {
  const resourceType = board?.resource_type ?? 0
  if (resourceType === 66 || resourceType === 106) {
    return { label: "课程", className: getSemanticStatusBadgeClass("accent") }
  }
  if (resourceType === 13) {
    return { label: "听书", className: getSemanticStatusBadgeClass("success") }
  }
  if (resourceType === 2 || resourceType === 910) {
    return { label: "图书", className: getSemanticStatusBadgeClass("warning") }
  }
  return { label: "榜单", className: getSemanticStatusBadgeClass("neutral") }
}

function RankBoardHero({ board }: { board: RankBoard | null }) {
  const badge = resolveRankTypeBadge(board)
  const backgroundColor = normalizeColor(board?.background_color)

  return (
    <div
      className="overflow-hidden rounded-3xl border border-border px-6 py-6 text-white shadow-soft"
      style={{
        background: `linear-gradient(135deg, ${backgroundColor} 0%, #1f2937 100%)`,
      }}
    >
      <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
        <div className="min-w-0 space-y-3">
          <span className={badge.className}>{badge.label}</span>
          <h2 className="text-2xl font-semibold">{board?.title || "得到榜单"}</h2>
          {board?.rn_right_title ? (
            <div className="flex flex-wrap gap-3 text-sm text-white/85">
              <span className="rounded-full bg-white/12 px-3 py-1.5">{board.rn_right_title}</span>
            </div>
          ) : null}
        </div>

        <div className="grid min-w-[14rem] grid-cols-2 gap-3">
          <div className="rounded-2xl bg-white/10 p-4 backdrop-blur">
            <p className="text-xs text-white/70">榜单条目</p>
            <p className="mt-2 text-2xl font-semibold">{board?.count || board?.list.length || 0}</p>
          </div>
          <div className="rounded-2xl bg-white/10 p-4 backdrop-blur">
            <p className="text-xs text-white/70">当前类型</p>
            <p className="mt-2 text-2xl font-semibold">{badge.label}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

function RankItemCardView({ item, index }: { item: RankListItem; index: number }) {
  const cover = resolveRankItemImage(item)
  const hotText = item.hot_intro?.number ? `${formatCount(item.hot_intro.number)}${item.hot_intro.intro || ""}` : ""
  const priceText = item.cost_intro?.price ? `￥${item.cost_intro.price}` : ""
  const showBought = item.authority_intro?.is_buy

  return (
    <Card className="h-full overflow-visible p-5 transition hover:-translate-y-1 hover:shadow-lg">
      <div className="flex h-full flex-col gap-3">
        <div className="grid gap-4 sm:grid-cols-[120px_minmax(0,1fr)]">
          <div className="relative overflow-visible">
            <div className="overflow-hidden rounded-2xl bg-surface-soft">
              {cover ? (
                <img alt={item.title || "榜单封面"} className="h-full w-full object-cover" src={cover} />
              ) : (
                <div className="flex h-full min-h-[9rem] items-center justify-center bg-surface-soft text-text-muted">
                  <BookOpen className="h-8 w-8" />
                </div>
              )}
            </div>

            <div className="pointer-events-none absolute -left-3 -top-3 z-10">
              <span className="inline-flex items-center rounded-full bg-black/45 px-2.5 py-1 text-xs font-semibold text-white backdrop-blur ring-1 ring-white/20 shadow-soft">
                #{index + 1}
              </span>
            </div>
          </div>

          <div className="min-w-0 space-y-3">
            <div className="space-y-1">
              <div className="flex flex-wrap items-center gap-2">
                <h3 className="line-clamp-2 text-base font-semibold text-text-primary">{item.title || "未命名内容"}</h3>
                {showBought ? <span className={getSemanticStatusBadgeClass("success")}>已购</span> : null}
                {item.is_today ? <span className={getSemanticStatusBadgeClass("warning")}>今日</span> : null}
              </div>
              <p className={`text-sm ${semanticMetaTextClass}`}>{item.author || item.type_name || "得到内容"}</p>
            </div>

            <p className="line-clamp-3 text-sm leading-7 text-text-secondary">
              {item.recommend_title || item.recommend_intro || item.intro || "暂无简介"}
            </p>

            <div className="flex flex-wrap gap-2">
              {item.metrics ? <span className={getSemanticStatusBadgeClass("neutral")}>{item.metrics}</span> : null}
              {hotText ? <span className={getSemanticStatusBadgeClass("accent")}>{hotText}</span> : null}
              {priceText ? <span className={getSemanticStatusBadgeClass("warning")}>{priceText}</span> : null}
              {item.douban_info?.use_douban && item.douban_info.score ? (
                <span className={getSemanticStatusBadgeClass("neutral")}>豆瓣 {item.douban_info.score}</span>
              ) : null}
            </div>
          </div>
        </div>
      </div>
    </Card>
  )
}

export function DedaoLeaderboardPage() {
  const [baseInfo, setBaseInfo] = useState<RankBaseInfoResponse | null>(null)
  const [rankData, setRankData] = useState<RankListResponse | null>(null)
  const [selectedPType, setSelectedPType] = useState<number | null>(null)
  const [selectedRankType, setSelectedRankType] = useState<number | null>(null)
  const [baseLoading, setBaseLoading] = useState(true)
  const [listLoading, setListLoading] = useState(false)
  const [baseError, setBaseError] = useState("")
  const [listError, setListError] = useState("")

  const currentNav = useMemo(
    () => baseInfo?.nav_list.find((item) => item.ptype === selectedPType) ?? baseInfo?.nav_list[0] ?? null,
    [baseInfo, selectedPType],
  )

  const currentRank = useMemo(
    () => currentNav?.rank_type_list.find((item) => item.rank_id === selectedRankType) ?? currentNav?.rank_type_list[0] ?? null,
    [currentNav, selectedRankType],
  )

  const currentBoard = rankData?.list[0] ?? null

  const loadBaseInfo = async () => {
    setBaseLoading(true)
    setBaseError("")

    try {
      const response = await api.rank.baseInfo()
      setBaseInfo(response)

      const firstNav = response.nav_list[0]
      const firstRank = firstNav?.rank_type_list[0]
      setSelectedPType(firstNav?.ptype ?? null)
      setSelectedRankType(firstRank?.rank_id ?? null)
    } catch (error) {
      setBaseError(error instanceof Error ? error.message : "榜单分类加载失败")
    } finally {
      setBaseLoading(false)
    }
  }

  const loadRankList = async (rankType: number) => {
    setListLoading(true)
    setListError("")

    try {
      const response = await api.rank.list(rankType)
      setRankData(response)
    } catch (error) {
      setListError(error instanceof Error ? error.message : "榜单内容加载失败")
      setRankData(null)
    } finally {
      setListLoading(false)
    }
  }

  useEffect(() => {
    void loadBaseInfo()
  }, [])

  useEffect(() => {
    if (!selectedRankType) {
      return
    }
    void loadRankList(selectedRankType)
  }, [selectedRankType])

  const handleSelectNav = (ptype: number) => {
    if (!baseInfo) {
      return
    }

    const nextNav = baseInfo.nav_list.find((item) => item.ptype === ptype)
    setSelectedPType(ptype)
    setSelectedRankType(nextNav?.rank_type_list[0]?.rank_id ?? null)
  }

  const handleSelectRank = (rankType: number) => {
    setSelectedRankType(rankType)
  }

  return (
    <main className="space-y-6">
      {baseError ? (
        <Card className="border-danger bg-danger-soft">
          <div className="flex flex-col gap-3 p-4 text-sm text-danger md:flex-row md:items-center md:justify-between">
            <span>{baseError}</span>
            <Button onClick={() => void loadBaseInfo()} variant="outline">
              重试
            </Button>
          </div>
        </Card>
      ) : null}

      <section className={`${semanticPageSectionClass} space-y-6 p-6`}>
        <header className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-accent/12 text-accent">
              <Trophy className="h-5 w-5" />
            </div>
            <h1 className="text-lg font-semibold leading-none text-text-primary">得到榜单</h1>
          </div>

          <Button className="gap-2 self-start md:self-auto" onClick={() => void loadBaseInfo()} variant="outline">
            <RefreshCcw className="h-4 w-4" />
            刷新榜单
          </Button>
        </header>

        {baseLoading ? (
          <div className="flex min-h-[18rem] items-center justify-center gap-3 text-sm text-text-muted">
            <Loader2 className="h-5 w-5 animate-spin" />
            正在加载榜单分类...
          </div>
        ) : (
          <div className="space-y-5">
            <div className="grid gap-3 lg:grid-cols-[120px_1fr] lg:items-start">
              <p className={`font-medium ${semanticMetaTextClass}`}>内容类型</p>
              <div className="flex flex-wrap gap-2">
                {(baseInfo?.nav_list ?? []).map((item) => (
                  <button
                    className={getSemanticChipClass((currentNav?.ptype ?? selectedPType) === item.ptype, "strong")}
                    key={item.ptype}
                    onClick={() => handleSelectNav(item.ptype)}
                    type="button"
                  >
                    {item.name}
                  </button>
                ))}
              </div>
            </div>

            <div className="grid gap-3 lg:grid-cols-[120px_1fr] lg:items-start">
              <p className={`font-medium ${semanticMetaTextClass}`}>榜单分类</p>
              <div className="flex flex-wrap gap-2">
                {(currentNav?.rank_type_list ?? []).map((item) => (
                  <button
                    className={getSemanticChipClass((currentRank?.rank_id ?? selectedRankType) === item.rank_id)}
                    key={item.rank_id}
                    onClick={() => handleSelectRank(item.rank_id)}
                    type="button"
                  >
                    {item.title}
                  </button>
                ))}
              </div>
            </div>

            {currentRank?.sub_title ? (
              <div className="flex items-start gap-3 rounded-2xl bg-surface-soft px-4 py-3 text-sm text-text-secondary">
                <Sparkles className="mt-0.5 h-4 w-4 shrink-0 text-accent" />
                <span>{currentRank.sub_title}</span>
              </div>
            ) : null}

            <RankBoardHero board={currentBoard} />
          </div>
        )}
      </section>

      {listError ? (
        <Card className="border-danger bg-danger-soft">
          <div className="flex flex-col gap-3 p-4 text-sm text-danger md:flex-row md:items-center md:justify-between">
            <span>{listError}</span>
            <Button disabled={!selectedRankType} onClick={() => selectedRankType && void loadRankList(selectedRankType)} variant="outline">
              重试
            </Button>
          </div>
        </Card>
      ) : null}

      <section className={`${semanticPageSectionClass} p-6`}>
        <div className="mb-4">
          <h2 className="text-lg font-semibold text-text-primary">榜单内容</h2>
        </div>

        {listLoading ? (
          <div className="flex min-h-[18rem] items-center justify-center gap-3 text-sm text-text-muted">
            <Loader2 className="h-5 w-5 animate-spin" />
            正在加载榜单内容...
          </div>
        ) : currentBoard?.list?.length ? (
          <div className="grid gap-4 xl:grid-cols-2">
            {currentBoard.list.map((item, index) => (
              <RankItemCardView item={item} index={index} key={`${currentBoard.rank_type}-${item.product_id}-${index}`} />
            ))}
          </div>
        ) : (
          <div className="flex min-h-[12rem] items-center justify-center rounded-3xl bg-surface-soft text-sm text-text-muted">
            暂无榜单数据
          </div>
        )}
      </section>

    </main>
  )
}
