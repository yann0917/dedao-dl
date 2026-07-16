import { useCallback, useEffect, useMemo, useState } from "react"
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

  useEffect(() => {
    setParams(baseParams)
  }, [baseParams])

  const loadFilter = useCallback(async (nextParams: AlgoFilterRequest) => {
    setLoadingFilter(true)
    try {
      const result = await api.algo.filter(nextParams)
      setFilter(result)
      setTotal(result.total)
    } finally {
      setLoadingFilter(false)
    }
  }, [])

  const loadProducts = useCallback(async (nextParams: AlgoFilterRequest, append: boolean) => {
    if (append) {
      setLoadingMore(true)
    } else {
      setLoadingProducts(true)
    }

    try {
      const result = await api.algo.products(nextParams)
      setProducts((current) => (append ? [...current, ...result.product_list] : result.product_list))
      setTotal(result.total)
      setIsMore(result.is_more)
      setParams((current) => ({
        ...current,
        request_id: result.request_id || current.request_id,
      }))
    } finally {
      if (append) {
        setLoadingMore(false)
      } else {
        setLoadingProducts(false)
      }
    }
  }, [])

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

  const applyParams = useCallback(async (patch: Partial<AlgoFilterRequest>, reloadFilter = false) => {
    const nextParams: AlgoFilterRequest = {
      ...params,
      ...patch,
      page: 0,
      request_id: "",
      page_size: PAGE_SIZE,
    }

    setParams(nextParams)
    setError(null)

    try {
      if (reloadFilter) {
        await loadFilter(nextParams)
      }
      await loadProducts(nextParams, false)
    } catch (err) {
      setError(err instanceof Error ? err.message : "分类页刷新失败")
    }
  }, [loadFilter, loadProducts, params])

  const loadMore = useCallback(async () => {
    if (loadingMore || isMore !== 1) {
      return
    }

    const nextParams: AlgoFilterRequest = {
      ...params,
      page: params.page + 1,
    }

    setParams(nextParams)

    try {
      await loadProducts(nextParams, true)
    } catch (err) {
      setError(err instanceof Error ? err.message : "加载更多失败")
    }
  }, [isMore, loadingMore, loadProducts, params])

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
