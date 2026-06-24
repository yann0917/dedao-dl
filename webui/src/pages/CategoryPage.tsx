import { BookOpen, ChevronRight, Loader2 } from "lucide-react"
import { useMemo, useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { toast } from "sonner"
import { api, type AlgoOption, type AlgoProductItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
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

  const handleSelectProductType = (option: AlgoOption) => {
    const nextName = option.name || init.classfcName
    syncQuery({
      name: nextName,
      navType: 0,
      enid: "",
      labelId: "",
      productType: option.value,
    })

    void explorer.applyParams(
      {
        classfc_name: nextName,
        nav_type: 0,
        navigation_id: "",
        label_id: "",
        product_types: option.value,
      },
      true,
    )
  }

  const handleSelectNavigation = (option: AlgoOption) => {
    syncQuery({
      name: option.name,
      navType: 0,
      enid: option.value,
      labelId: "",
      productType: explorer.selectedProductType,
    })

    void explorer.applyParams(
      {
        classfc_name: option.name,
        navigation_id: option.value,
        label_id: "",
      },
      true,
    )
  }

  const handleSelectSubOption = (option: AlgoOption) => {
    syncQuery({
      name: init.classfcName,
      navType: init.navType,
      enid: explorer.selectedNavigationId,
      labelId: option.value,
      productType: explorer.selectedProductType,
    })

    void explorer.applyParams({
      label_id: option.value,
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

      <Card className="p-6">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p className={semanticMetaTextClass}>结果概览</p>
            <h3 className="mt-2 text-xl font-semibold text-text-primary">已为你找到 {explorer.total} 个内容</h3>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <span className={semanticMetaTextClass}>排序</span>
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
      </Card>

      {(explorer.loadingFilter || explorer.loadingProducts) && explorer.products.length === 0 ? (
        <Card className="flex items-center justify-center p-8 text-text-muted">
          <Loader2 className="mr-2 h-5 w-5 animate-spin" />
          正在加载分类内容...
        </Card>
      ) : null}

      <section className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        {explorer.products.map((item) => {
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

          return (
            <Card className="h-full overflow-hidden transition hover:-translate-y-1 hover:shadow-lg" key={`${item.id}-${item.id_out}`}>
              <button
                className="w-full text-left"
                onClick={() => handleOpenProduct(item)}
                type="button"
              >
                <div className="aspect-[16/10] overflow-hidden bg-surface-soft">
                  <img
                    alt={title}
                    className="h-full w-full object-cover transition duration-500 hover:scale-105"
                    src={item.index_img || item.horizontal_image || "https://placehold.co/640x400/e2e8f0/334155?text=DD"}
                  />
                </div>
                <div className="space-y-4 p-5">
                  <div className="flex items-center justify-between gap-3">
                    <ProductTypeBadge item={item} />
                    {item.learn_user_count > 0 ? <span className="text-xs text-text-muted/80">{item.learn_user_count} 人学习</span> : null}
                  </div>

                  <div>
                    <h4 className="line-clamp-2 text-lg font-semibold text-text-primary">{title}</h4>
                    <p className="mt-2 line-clamp-3 text-sm leading-6 text-text-secondary">{summary}</p>
                  </div>

                  <div className="flex items-center justify-between text-sm text-text-muted">
                    <span className="inline-flex items-center gap-1">
                      <BookOpen className="h-4 w-4" />
                      {item.score ? `评分 ${item.score}` : item.price_desc || "内容详情"}
                    </span>
                    <span className="inline-flex items-center gap-1 text-primary">
                      {canOpen ? "打开" : "待接入"}
                      <ChevronRight className="h-4 w-4" />
                    </span>
                  </div>
                </div>
              </button>

              {(isAudioShelfItem || isEbookShelfItem) ? (
                <div className="flex flex-wrap items-center gap-2 border-t border-border px-5 pb-5 pt-4">
                  {isAudioShelfItem ? (
                    isAudioInShelf ? (
                      <span className={getSemanticStatusBadgeClass("neutral")}>已加入书架</span>
                    ) : (
                      <Button disabled={isAddingAudio} onClick={() => void handleAddAudioToShelf(item)}>
                        {isAddingAudio ? "处理中..." : "加入书架"}
                      </Button>
                    )
                  ) : null}

                  {isEbookShelfItem ? (
                    isEbookInShelf ? (
                      <Button
                        disabled={isRemovingEbook}
                        onClick={() => void handleRemoveEbookFromShelf(item)}
                        variant="danger"
                      >
                        {isRemovingEbook ? "处理中..." : "移出书架"}
                      </Button>
                    ) : (
                      <Button disabled={isAddingEbook} onClick={() => void handleAddEbookToShelf(item)}>
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
            {explorer.loadingMore ? "加载中..." : "加载更多"}
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
