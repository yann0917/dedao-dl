import { useEffect, useRef, useState } from "react"
import {
  api,
  type CourseCategory,
  type CourseInfoResponse,
  type CourseListItem,
} from "@/api"

type UseCourseWorkspaceOptions = {
  initialCategory?: string | null
  initialEnid?: string | null
}

export function useCourseWorkspace(options: UseCourseWorkspaceOptions = {}) {
  const initialCategoryApplied = useRef(false)
  const [categories, setCategories] = useState<CourseCategory[]>([])
  const [selectedCategory, setSelectedCategory] = useState(options.initialCategory || "bauhinia")
  const [courses, setCourses] = useState<CourseListItem[]>([])
  const [detail, setDetail] = useState<CourseInfoResponse | null>(null)
  const [loadingCourses, setLoadingCourses] = useState(false)
  const [loadingDetail, setLoadingDetail] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    void bootstrap()
  }, [])

  useEffect(() => {
    if (options.initialCategory && !initialCategoryApplied.current) {
      setSelectedCategory(options.initialCategory)
      initialCategoryApplied.current = true
      return
    }

    if (!selectedCategory) {
      return
    }

    void loadCourses(selectedCategory, options.initialEnid || undefined)
  }, [options.initialCategory, options.initialEnid, selectedCategory])

  const bootstrap = async () => {
    setError(null)

    try {
      const categoryResp = await api.course.categories()
      const nextCategories = categoryResp.data.list
      setCategories(nextCategories)

      if (!options.initialCategory && nextCategories.length > 0) {
        setSelectedCategory(nextCategories[0].category)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "课程初始化失败")
    }
  }

  const loadCourses = async (category: string, preferredEnid?: string) => {
    setLoadingCourses(true)
    setError(null)

    try {
      const params = new URLSearchParams({ category, page: "1", limit: "18", order: "study" })
      const result = await api.course.list(params)
      setCourses(result.list)

      const targetEnid = preferredEnid || result.list[0]?.enid
      if (targetEnid) {
        await loadCourseDetail(targetEnid)
      } else {
        setDetail(null)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "课程列表加载失败")
    } finally {
      setLoadingCourses(false)
    }
  }

  const loadCourseDetail = async (enid: string) => {
    setLoadingDetail(true)
    setError(null)

    try {
      const result = await api.course.info(enid)
      setDetail(result)
    } catch (err) {
      setError(err instanceof Error ? err.message : "课程详情加载失败")
    } finally {
      setLoadingDetail(false)
    }
  }

  return {
    categories,
    selectedCategory,
    setSelectedCategory,
    courses,
    detail,
    loadingCourses,
    loadingDetail,
    error,
    bootstrap,
    loadCourseDetail,
  }
}
