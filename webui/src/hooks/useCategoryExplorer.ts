import { useCallback, useEffect, useMemo, useRef, useState } from "react"
import {
  api,
  type AlgoFilterRequest,
  type AlgoFilterResponse,
  type AlgoOption,
  type AlgoProductItem,
} from "@/api"

type CategoryExplorerInit = {
  classfcName: string
  labelId: string
  navType: number
  navigationId: string
  productTypes: string
}

const PAGE_SIZE = 18

function buildBaseParams(init: CategoryExplorerInit): AlgoFilterRequest {
  return {
    classfc_name: init.classfcName || "全部",
    label_id: init.labelId || "",
    nav_type: init.navType || 0,
    navigation_id: init.navigationId || "",
    page: 0,
    page_size: PAGE_SIZE,
    product_types: init.productTypes || "66",
    request_id: "",
    sort_strategy: "HOT",
  }
}

export function useCategoryExplorer(init: CategoryExplorerInit) {
  const baseParams = useMemo(
    () => buildBaseParams(init),
    [init.classfcName, init.labelId, init.navType, init.navigationId, init.productTypes],
  )

  const [params, setParams] = useState<AlgoFilterRequest>(baseParams)
  const [filter, setFilter] = useState<AlgoFilterResponse | null>(null)
  const [products, setProducts] = useState<AlgoProductItem[]>([])
  const [total, setTotal] = useState(0)
  const [isMore, setIsMore] = useState(0)
  const [loadingFilter, setLoadingFilter] = useState(true)
  const [loadingProducts, setLoadingProducts] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // 请求代际号：每次“替换式”加载（切换筛选、首次进入）递增；迟到的旧响应整体丢弃，
  // 这是防止切换内容类型时新旧数据交错（卡片重复/不更新）的关键。
  const requestGenRef = useRef(0)
  // loadMore 同步锁：观察器可能在 React 状态落地前连续触发，挡掉同拍内的重复追加。
  const loadMoreLockRef = useRef(false)
  // params / isMore / loadingMore 的同步镜像，保证懒加载回调读到的永远是最新值，
  // 不受闭包时序影响。
  const paramsRef = useRef(params)
  const isMoreRef = useRef(isMore)
  const loadingMoreRef = useRef(loadingMore)

  const commitParams = useCallback((next: AlgoFilterRequest) => {
    paramsRef.current = next
    setParams(next)
  }, [])

  useEffect(() => {
    commitParams(baseParams)
  }, [baseParams, commitParams])

  const loadFilter = useCallback(
    async (nextParams: AlgoFilterRequest) => {
      requestGenRef.current += 1
      const gen = requestGenRef.current
      setLoadingFilter(true)
      try {
        const result = await api.algo.filter(nextParams)
        if (gen !== requestGenRef.current) {
          return
        }
        setFilter(result)
        setTotal(result.total)
      } finally {
        setLoadingFilter(false)
      }
    },
    [],
  )

  const loadProducts = useCallback(async (nextParams: AlgoFilterRequest, append: boolean) => {
    if (append) {
      if (loadingMoreRef.current) {
        return
      }
      loadingMoreRef.current = true
      setLoadingMore(true)
    } else {
      requestGenRef.current += 1
      setLoadingProducts(true)
    }
    const gen = requestGenRef.current

    try {
      const result = await api.algo.products(nextParams)
      if (gen !== requestGenRef.current) {
        return
      }
      setProducts((current) => {
        if (!append) {
          return result.product_list
        }
        // 兜底去重：algo 接口的 id 字段恒为 0，只能按 id_out 识别条目；
        // 后端翻页存在 1 条重叠、同页内也可能出现重复 id_out，统一在此过滤。
        const seen = new Set(
          current.map((item) => `${item.id_out}-${item.product_id}-${item.name}`),
        )
        return [
          ...current,
          ...result.product_list.filter(
            (item) => !seen.has(`${item.id_out}-${item.product_id}-${item.name}`),
          ),
        ]
      })
      setTotal(result.total)
      isMoreRef.current = result.is_more
      setIsMore(result.is_more)
      commitParams({
        ...paramsRef.current,
        page: nextParams.page,
        request_id: result.request_id || paramsRef.current.request_id,
      })
    } finally {
      if (append) {
        loadingMoreRef.current = false
        setLoadingMore(false)
      } else {
        setLoadingProducts(false)
      }
    }
  }, [commitParams])

  useEffect(() => {
    let cancelled = false

    const bootstrap = async () => {
      setError(null)
      setProducts([])

      try {
        await loadFilter(baseParams)
        if (!cancelled) {
          await loadProducts(baseParams, false)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "分类页加载失败")
        }
      }
    }

    void bootstrap()

    return () => {
      cancelled = true
    }
  }, [baseParams, loadFilter, loadProducts])

  const applyParams = useCallback(
    async (patch: Partial<AlgoFilterRequest>, reloadFilter = false) => {
      const nextParams: AlgoFilterRequest = {
        ...paramsRef.current,
        ...patch,
        page: 0,
        request_id: "",
        page_size: PAGE_SIZE,
      }

      commitParams(nextParams)
      setError(null)
      // 立即清空旧列表，让骨架屏接管，也断掉与旧筛选结果的关联。
      setProducts([])

      try {
        if (reloadFilter) {
          await loadFilter(nextParams)
        }
        await loadProducts(nextParams, false)
      } catch (err) {
        setError(err instanceof Error ? err.message : "分类页刷新失败")
      }
    },
    [commitParams, loadFilter, loadProducts],
  )

  const loadMore = useCallback(async () => {
    if (loadMoreLockRef.current || loadingMoreRef.current || isMoreRef.current !== 1) {
      return
    }

    loadMoreLockRef.current = true
    const nextParams: AlgoFilterRequest = {
      ...paramsRef.current,
      page: paramsRef.current.page + 1,
    }
    commitParams(nextParams)

    try {
      await loadProducts(nextParams, true)
    } catch (err) {
      setError(err instanceof Error ? err.message : "加载更多失败")
    } finally {
      loadMoreLockRef.current = false
    }
  }, [commitParams, loadProducts])

  const productTypeOptions = filter?.filter.product_types.options ?? []
  const navigationOptions = filter?.filter.navigations.options ?? []
  const sortOptions = filter?.filter.sort_strategy.options ?? []
  const activeNavigation = navigationOptions.find((item) => item.value === params.navigation_id)
  const subOptions = activeNavigation?.sub_options ?? []

  return {
    params,
    total,
    products,
    isMore,
    loadingFilter,
    loadingProducts,
    loadingMore,
    error,
    productTypeOptions,
    navigationOptions,
    subOptions,
    sortOptions,
    applyParams,
    loadMore,
    selectedProductType: params.product_types,
    selectedNavigationId: params.navigation_id,
    selectedLabelId: params.label_id,
    selectedSortStrategy: params.sort_strategy,
  }
}

export function buildCategoryQuery(next: {
  id?: string | number
  name: string
  navType: number
  enid: string
  labelId?: string
  productType?: string
}) {
  const query = new URLSearchParams()
  if (next.id !== undefined && next.id !== null && String(next.id) !== "") {
    query.set("id", String(next.id))
  }
  query.set("name", next.name || "全部")
  query.set("nav_type", String(next.navType || 0))
  query.set("enid", next.enid || "")
  query.set("label_id", next.labelId || "")
  query.set("product_type", next.productType || (next.navType === 2 ? "2" : "66"))
  return query.toString()
}

export function resolveCategoryOptionName(options: AlgoOption[], value: string, fallback: string) {
  return options.find((item) => item.value === value)?.name || fallback
}
