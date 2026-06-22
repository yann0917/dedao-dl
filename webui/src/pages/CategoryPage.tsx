import { BookOpen, ChevronRight, Loader2 } from "lucide-react"
import { useMemo } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { type AlgoOption, type AlgoProductItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  buildCategoryQuery,
  useCategoryExplorer,
} from "@/hooks/useCategoryExplorer"
import { cn } from "@/lib/cn"

function ProductTypeBadge({ item }: { item: AlgoProductItem }) {
  if (item.item_type === 66) {
    return <span className="rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary">课程</span>
  }

  if (item.item_type === 2) {
    return <span className="rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-700">电子书</span>
  }

  if (item.item_type === 13) {
    return <span className="rounded-full bg-orange-100 px-3 py-1 text-xs font-medium text-orange-700">听书</span>
  }

  return <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-500">内容</span>
}

export function CategoryPage() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()

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
      navigate(`/courses?enid=${encodeURIComponent(item.id_out)}`)
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

  return (
    <main className="space-y-6">
      <section className="rounded-3xl border border-white/70 bg-white/90 p-6 shadow-soft backdrop-blur">
        <p className="text-sm text-slate-500">分类结果页</p>
        <h2 className="mt-2 text-3xl font-semibold text-slate-950">{init.classfcName}</h2>
        <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">
          这里承接首页分类点击后的筛选与结果流。当前优先对齐 `Algo.vue` 的主链路，先覆盖筛选、排序、分页和内容跳转。
        </p>
      </section>

      {explorer.error ? (
        <Card className="border border-rose-200 bg-rose-50">
          <div className="p-4 text-sm text-rose-700">{explorer.error}</div>
        </Card>
      ) : null}

      <Card className="p-6">
        <div className="space-y-5">
          <div className="grid gap-3 lg:grid-cols-[120px_1fr] lg:items-start">
            <p className="text-sm font-medium text-slate-500">内容类型</p>
            <div className="flex flex-wrap gap-2">
              {explorer.productTypeOptions.map((option) => (
                <button
                  className={cn(
                    "rounded-full px-3 py-1.5 text-sm font-medium transition",
                    explorer.selectedProductType === option.value
                      ? "bg-primary text-white"
                      : "bg-slate-100 text-slate-600 hover:bg-slate-200",
                  )}
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
            <p className="text-sm font-medium text-slate-500">内容分类</p>
            <div className="flex flex-wrap gap-2">
              {explorer.navigationOptions.map((option) => (
                <button
                  className={cn(
                    "rounded-full px-3 py-1.5 text-sm font-medium transition",
                    explorer.selectedNavigationId === option.value
                      ? "bg-primary text-white"
                      : "bg-slate-100 text-slate-600 hover:bg-slate-200",
                  )}
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
              <p className="text-sm font-medium text-slate-500">子标签</p>
              <div className="flex flex-wrap gap-2">
                {explorer.subOptions.map((option) => (
                  <button
                    className={cn(
                      "rounded-full px-3 py-1.5 text-sm font-medium transition",
                      explorer.selectedLabelId === option.value
                        ? "bg-slate-950 text-white"
                        : "bg-slate-100 text-slate-600 hover:bg-slate-200",
                    )}
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
            <p className="text-sm text-slate-500">结果概览</p>
            <h3 className="mt-2 text-xl font-semibold text-slate-950">已为你找到 {explorer.total} 个内容</h3>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <span className="text-sm text-slate-500">排序</span>
            {explorer.sortOptions.map((option) => (
              <button
                className={cn(
                  "rounded-full px-3 py-1.5 text-sm font-medium transition",
                  explorer.selectedSortStrategy === option.value
                    ? "bg-primary text-white"
                    : "bg-slate-100 text-slate-600 hover:bg-slate-200",
                )}
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
        <Card className="flex items-center justify-center p-8 text-slate-500">
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

          return (
            <button
              className="text-left"
              key={`${item.id}-${item.id_out}`}
              onClick={() => handleOpenProduct(item)}
              type="button"
            >
              <Card className="h-full overflow-hidden transition hover:-translate-y-1 hover:shadow-lg">
                <div className="aspect-[16/10] overflow-hidden bg-slate-100">
                  <img
                    alt={title}
                    className="h-full w-full object-cover transition duration-500 hover:scale-105"
                    src={item.index_img || item.horizontal_image || "https://placehold.co/640x400/e2e8f0/334155?text=DD"}
                  />
                </div>
                <div className="space-y-4 p-5">
                  <div className="flex items-center justify-between gap-3">
                    <ProductTypeBadge item={item} />
                    {item.learn_user_count > 0 ? <span className="text-xs text-slate-400">{item.learn_user_count} 人学习</span> : null}
                  </div>

                  <div>
                    <h4 className="line-clamp-2 text-lg font-semibold text-slate-950">{title}</h4>
                    <p className="mt-2 line-clamp-3 text-sm leading-6 text-slate-500">{summary}</p>
                  </div>

                  <div className="flex items-center justify-between text-sm text-slate-500">
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
              </Card>
            </button>
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
        <Card className="p-10 text-center text-slate-500">
          <p className="text-lg font-medium text-slate-900">当前筛选下没有找到内容</p>
          <p className="mt-2 text-sm">可以尝试切换内容类型、分类或排序条件。</p>
        </Card>
      ) : null}
    </main>
  )
}
