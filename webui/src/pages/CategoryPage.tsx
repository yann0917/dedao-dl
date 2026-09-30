import { BookOpen, ChevronRight, Loader2 } from "lucide-react"
import { useMemo, useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { toast } from "sonner"
import { api, type AlgoOption, type AlgoProductItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { Skeleton } from "@/components/ui/Semantic"
import {
  buildCategoryQuery,
  useCategoryExplorer,
} from "@/hooks/useCategoryExplorer"
import { cn } from "@/lib/cn"
import { getSemanticChipClass, getSemanticStatusBadgeClass, semanticMetaTextClass } from "@/lib/semanticStyles"

function ProductTypeBadge({ item }: { item: AlgoProductItem }) {
  if (item.item_type === 66) {
    return <span className={getSemanticStatusBadgeClass("accent")}>课程</span>
  }

  if (item.item_type === 2) {
    return <span className={getSemanticStatusBadgeClass("warning")}>电子书</span>
  }

  if (item.item_type === 13) {
    return <span className={getSemanticStatusBadgeClass("success")}>听书</span>
  }

  return <span className={getSemanticStatusBadgeClass("neutral")}>内容</span>
}

export function CategoryPage() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const [actionLoadingKey, setActionLoadingKey] = useState<string | null>(null)

  const init = useMemo(
    () => ({
      classfcName: searchParams.get("name") || "全部",
      labelId: searchParams.get("label_id") || searchParams.get("labelId") || "",
      navType: Number(searchParams.get("nav_type") || searchParams.get("navType") || 0),
      navigationId: searchParams.get("enid") || "",
      productTypes: searchParams.get("product_type") || searchParams.get("productType") || "66",
    }),
    [searchParams],
  )

  const explorer = useCategoryExplorer(init)

  // 「加载更多」为手动触发：懒加载的自动追加分页在此页不可靠，改为显式按钮。
  const { loadMore } = explorer

  const syncQuery = (patch: {
    id?: string | number
    name: string
    navType: number
    enid: string
    labelId?: string
    productType?: string
  }) => {
    setSearchParams(buildCategoryQuery(patch), { replace: true })
  }

  // 类型/分类/子标签的切换都只改 URL：baseParams 随之变化并触发 hook 里的 bootstrap
  // 统一加载，避免手动 applyParams 与 bootstrap 双管线竞态（曾导致分类 chips 与
  // 列表类型不同步）。排序不走 URL，仍由 applyParams 处理。
  const handleSelectProductType = (option: AlgoOption) => {
    syncQuery({
      name: init.classfcName,
      navType: 0,
      enid: "",
      labelId: "",
      productType: option.value,
    })
  }

  const handleSelectNavigation = (option: AlgoOption) => {
    syncQuery({
      name: option.name,
      navType: 0,
      enid: option.value,
      labelId: "",
      productType: explorer.selectedProductType,
    })
  }

  const handleSelectSubOption = (option: AlgoOption) => {
    syncQuery({
      name: init.classfcName,
      navType: init.navType,
      enid: explorer.selectedNavigationId,
      labelId: option.value,
      productType: explorer.selectedProductType,
    })
  }

  const handleSelectSort = (option: AlgoOption) => {
    void explorer.applyParams({
      sort_strategy: option.value,
    })
  }

  const handleOpenProduct = (item: AlgoProductItem) => {
    if (item.item_type === 66 && item.id_out) {
      navigate(
        `/courses/${encodeURIComponent(item.id_out)}?from=algo&parentTitle=${encodeURIComponent(item.name || item.intro || "课程详情")}`,
      )
      return
    }

    if (item.item_type === 2 && item.id_out) {
      navigate(`/ebooks/${encodeURIComponent(item.id_out)}`)
      return
    }

    if (item.item_type === 13 && item.product_type === 13 && item.id_out) {
      navigate(`/audios/${encodeURIComponent(item.id_out)}`)
      return
    }

    if (item.item_type === 13 && item.product_type === 1013 && item.id_out) {
      navigate(`/audio-groups/${encodeURIComponent(item.id_out)}`)
      return
    }

    if (item.dd_url) {
      window.open(item.dd_url, "_blank", "noopener,noreferrer")
    }
  }

  const handleAddAudioToShelf = async (item: AlgoProductItem) => {
    if (!item.id_out) {
      return
    }

    const actionKey = `audio-add-${item.id_out}`
    setActionLoadingKey(actionKey)
    try {
      await api.audio.addToShelf([item.id_out])
      await explorer.applyParams({})
      toast.success("已加入书架", {
        description: item.name || "当前听书已加入书架",
      })
    } catch (err) {
      toast.error("加入书架失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setActionLoadingKey(null)
    }
  }

  const handleAddEbookToShelf = async (item: AlgoProductItem) => {
    if (!item.id_out) {
      return
    }

    const actionKey = `ebook-add-${item.id_out}`
    setActionLoadingKey(actionKey)
    try {
      await api.ebook.addToShelf([item.id_out])
      await explorer.applyParams({})
      toast.success("已加入书架", {
        description: item.name || "当前电子书已加入书架",
      })
    } catch (err) {
      toast.error("加入书架失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setActionLoadingKey(null)
    }
  }

  const handleRemoveEbookFromShelf = async (item: AlgoProductItem) => {
    if (!item.id_out) {
      return
    }

    if (!window.confirm("确定要将这本电子书移出书架吗？")) {
      return
    }

    const actionKey = `ebook-remove-${item.id_out}`
    setActionLoadingKey(actionKey)
    try {
      await api.ebook.removeFromShelf([item.id_out])
      await explorer.applyParams({})
      toast.success("已移出书架", {
        description: item.name || "当前电子书已从书架移除",
      })
    } catch (err) {
      toast.error("移出书架失败", {
        description: err instanceof Error ? err.message : "请稍后重试",
      })
    } finally {
      setActionLoadingKey(null)
    }
  }

  return (
    <main className="space-y-6">
      {explorer.error ? (
        <Card className="border-danger bg-danger-soft">
          <div className="p-4 text-sm text-danger">{explorer.error}</div>
        </Card>
      ) : null}

      <Card className="p-6">
        <div className="space-y-5">
          <div className="grid gap-3 lg:grid-cols-[120px_1fr] lg:items-start">
            <p className={cn("font-medium", semanticMetaTextClass)}>内容类型</p>
            <div className="flex flex-wrap gap-2">
              {explorer.productTypeOptions.map((option) => (
                <button
                  className={getSemanticChipClass(explorer.selectedProductType === option.value)}
                  key={option.value}
                  onClick={() => handleSelectProductType(option)}
                  type="button"
                >
                  {option.name}
                </button>
              ))}
            </div>
          </div>

          <div className="grid gap-3 lg:grid-cols-[120px_1fr] lg:items-start">
            <p className={cn("font-medium", semanticMetaTextClass)}>内容分类</p>
            <div className="flex flex-wrap gap-2">
              {explorer.navigationOptions.map((option) => (
                <button
                  className={getSemanticChipClass(explorer.selectedNavigationId === option.value)}
                  key={option.value}
                  onClick={() => handleSelectNavigation(option)}
                  type="button"
                >
                  {option.name}
                </button>
              ))}
            </div>
          </div>

          {explorer.subOptions.length > 0 ? (
            <div className="grid gap-3 lg:grid-cols-[120px_1fr] lg:items-start">
              <p className={cn("font-medium", semanticMetaTextClass)}>子标签</p>
              <div className="flex flex-wrap gap-2">
                {explorer.subOptions.map((option) => (
                  <button
                    className={getSemanticChipClass(explorer.selectedLabelId === option.value, "strong")}
                    key={option.value}
                    onClick={() => handleSelectSubOption(option)}
                    type="button"
                  >
                    {option.name}
                  </button>
                ))}
              </div>
            </div>
          ) : null}
        </div>
      </Card>

      <div className="flex flex-wrap items-center justify-between gap-x-6 gap-y-3 border-b border-border pb-3">
        <p className="font-mono text-xs tracking-[0.12em] text-text-muted">
          共 <span className="font-semibold text-text-primary">{explorer.total}</span> 个内容
        </p>

        <div className="flex flex-wrap items-center gap-2">
          <span className="font-mono text-[10px] uppercase tracking-[0.24em] text-text-muted">排序</span>
          {explorer.sortOptions.map((option) => (
            <button
              className={getSemanticChipClass(explorer.selectedSortStrategy === option.value)}
              key={option.value}
              onClick={() => handleSelectSort(option)}
              type="button"
            >
              {option.name}
            </button>
          ))}
        </div>
      </div>

      {/* 筛选骨架只在筛选数据本身未就绪时显示，避免与真实筛选卡同屏重复。 */}
      {explorer.loadingFilter && explorer.products.length === 0 ? (
        <Card className="space-y-5 p-6">
          <div className="flex gap-4">
            <Skeleton className="h-6 w-16" />
            <div className="flex flex-1 gap-2">
              <Skeleton className="h-8 w-16" />
              <Skeleton className="h-8 w-16" />
              <Skeleton className="h-8 w-16" />
            </div>
          </div>
          <div className="flex gap-4">
            <Skeleton className="h-6 w-16" />
            <div className="flex flex-1 gap-2">
              <Skeleton className="h-8 w-20" />
              <Skeleton className="h-8 w-20" />
              <Skeleton className="h-8 w-14" />
              <Skeleton className="h-8 w-14" />
            </div>
          </div>
        </Card>
      ) : null}

      {(explorer.loadingFilter || explorer.loadingProducts) && explorer.products.length === 0 ? (
        <div aria-busy="true" className="space-y-6">
          <div className="flex items-center justify-between border-b border-border pb-3">
            <Skeleton className="h-4 w-28" />
            <Skeleton className="h-6 w-40" />
          </div>

          <section className="grid gap-4 md:grid-cols-2">
            {Array.from({ length: 6 }).map((_, index) => (
              <Card className="flex gap-4 p-4" key={index}>
                <Skeleton className="aspect-video w-40 shrink-0 self-center" />
                <div className="flex-1 space-y-3 py-1">
                  <Skeleton className="h-5 w-16" />
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-3 w-full" />
                  <Skeleton className="h-3 w-1/2" />
                </div>
              </Card>
            ))}
          </section>
        </div>
      ) : null}

      <section className="grid gap-4 md:grid-cols-2">
        {explorer.products.map((item, index) => {
          const title = item.name || item.intro || "得到内容"
          const summary = item.intro || item.lecturer_name_and_title || item.author_list.join(" / ") || "暂无简介"
          const canOpen =
            (item.item_type === 66 && !!item.id_out) ||
            (item.item_type === 2 && !!item.id_out) ||
            (item.item_type === 13 && !!item.id_out) ||
            !!item.dd_url
          const isAudioShelfItem = item.item_type === 13 && item.product_type === 13 && !!item.id_out
          const isEbookShelfItem = item.item_type === 2 && !!item.id_out
          const isAudioInShelf = Boolean(item.in_bookrack)
          const isEbookInShelf = Boolean(item.is_on_bookshelf)
          const isAddingAudio = actionLoadingKey === `audio-add-${item.id_out}`
          const isAddingEbook = actionLoadingKey === `ebook-add-${item.id_out}`
          const isRemovingEbook = actionLoadingKey === `ebook-remove-${item.id_out}`
          // 电子书/听书封面为竖版 3:4，课程封面为横版 16:9，各按原始比例展示避免裁切。
          const isVerticalCover = item.item_type === 2 || item.item_type === 13
          const coverSrc =
            (isVerticalCover ? item.index_img : item.horizontal_image || item.index_img) ||
            "https://placehold.co/480x640/e2e8f0/334155?text=DD"

          return (
            <Card
              className="h-full overflow-hidden transition hover:shadow-soft"
              key={`${item.id_out}-${item.product_id}-${item.name}-${index}`}
            >
              <button
                className="flex w-full text-left"
                onClick={() => handleOpenProduct(item)}
                type="button"
              >
                <div
                  className={cn(
                    "relative ml-4 shrink-0 self-center overflow-hidden bg-surface-soft",
                    isVerticalCover ? "aspect-[3/4] w-24 sm:w-28" : "aspect-video w-40 sm:w-44",
                  )}
                >
                  <img
                    alt={title}
                    className="absolute inset-0 h-full w-full object-cover"
                    src={coverSrc}
                  />
                </div>

                <div className="flex min-w-0 flex-1 flex-col gap-2 p-4 sm:pr-5">
                  <div className="flex items-center justify-between gap-2">
                    <ProductTypeBadge item={item} />
                    {item.learn_user_count > 0 ? (
                      <span className="font-mono text-[10px] tabular-nums text-text-muted">
                        {item.learn_user_count} 人学习
                      </span>
                    ) : null}
                  </div>

                  <h4 className="line-clamp-2 text-base font-semibold leading-snug text-text-primary">{title}</h4>
                  <p className="line-clamp-2 text-sm leading-6 text-text-secondary">{summary}</p>

                  <div className="mt-auto flex items-center justify-between pt-1 text-xs text-text-muted">
                    <span className="inline-flex items-center gap-1">
                      <BookOpen className="size-icon-sm" />
                      {item.score ? `评分 ${item.score}` : item.price_desc || "内容详情"}
                    </span>
                    <span className="inline-flex items-center gap-0.5 text-accent">
                      {canOpen ? "打开" : "待接入"}
                      <ChevronRight className="size-icon-sm" />
                    </span>
                  </div>
                </div>
              </button>

              {(isAudioShelfItem || isEbookShelfItem) ? (
                <div className="flex flex-wrap items-center gap-2 border-t border-border px-4 py-3">
                  {isAudioShelfItem ? (
                    isAudioInShelf ? (
                      <span className={getSemanticStatusBadgeClass("neutral")}>已加入书架</span>
                    ) : (
                      <Button className="h-9 px-3 text-xs" disabled={isAddingAudio} onClick={() => void handleAddAudioToShelf(item)}>
                        {isAddingAudio ? "处理中..." : "加入书架"}
                      </Button>
                    )
                  ) : null}

                  {isEbookShelfItem ? (
                    isEbookInShelf ? (
                      <Button
                        className="h-9 px-3 text-xs"
                        disabled={isRemovingEbook}
                        onClick={() => void handleRemoveEbookFromShelf(item)}
                        variant="danger"
                      >
                        {isRemovingEbook ? "处理中..." : "移出书架"}
                      </Button>
                    ) : (
                      <Button className="h-9 px-3 text-xs" disabled={isAddingEbook} onClick={() => void handleAddEbookToShelf(item)}>
                        {isAddingEbook ? "处理中..." : "加入书架"}
                      </Button>
                    )
                  ) : null}
                </div>
              ) : null}
            </Card>
          )
        })}
      </section>

      {explorer.isMore === 1 ? (
        <div className="flex justify-center">
          <Button disabled={explorer.loadingMore} onClick={() => void explorer.loadMore()} variant="outline">
            {explorer.loadingMore ? (
              <>
                <Loader2 className="mr-2 size-icon-sm animate-spin" />
                加载中...
              </>
            ) : (
              "加载更多"
            )}
          </Button>
        </div>
      ) : null}

      {!explorer.loadingProducts && explorer.products.length === 0 ? (
        <Card className="p-10 text-center text-text-muted">
          <p className="text-lg font-medium text-text-primary">当前筛选下没有找到内容</p>
          <p className="mt-2 text-sm">可以尝试切换内容类型、分类或排序条件。</p>
        </Card>
      ) : null}
    </main>
  )
}
