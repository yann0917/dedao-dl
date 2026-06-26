import { useCallback, useEffect, useMemo, useState } from "react"
import { api, type HomeContentResponse, type HomeNavigation, type HomePortalResponse } from "@/api"

type PortalSectionKind = "course" | "ebook"

export function useHomePortal() {
  const [data, setData] = useState<HomePortalResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [switchingSection, setSwitchingSection] = useState<PortalSectionKind | null>(null)
  const [sectionError, setSectionError] = useState<string | null>(null)
  const [selectedCourseEnid, setSelectedCourseEnid] = useState("")
  const [selectedEbookEnid, setSelectedEbookEnid] = useState("")

  const applySectionSelection = useCallback((portal: HomePortalResponse) => {
    setSelectedCourseEnid(portal.courseContent?.current_enid || portal.courseLabels?.list[0]?.enid || "")
    setSelectedEbookEnid(portal.ebookContent?.current_enid || portal.ebookLabels?.list[0]?.enid || "")
    return portal
  }, [])

  const loadPortal = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const portal = await api.home.portal()
      setData(portal)
      applySectionSelection(portal)
    } catch (err) {
      setError(err instanceof Error ? err.message : "首页加载失败")
    } finally {
      setLoading(false)
    }
  }, [applySectionSelection])

  useEffect(() => {
    void loadPortal()
  }, [loadPortal])

  const updateSectionContent = useCallback((kind: PortalSectionKind, content: HomeContentResponse) => {
    setData((current) => {
      if (!current) {
        return current
      }

      if (kind === "course") {
        return {
          ...current,
          courseContent: content,
          courseContentError: undefined,
        }
      }

      return {
        ...current,
        ebookContent: content,
        ebookContentError: undefined,
      }
    })
  }, [])

  const selectLabel = useCallback(async (kind: PortalSectionKind, label: HomeNavigation) => {
    setSectionError(null)
    setSwitchingSection(kind)

    if (kind === "course") {
      setSelectedCourseEnid(label.enid)
    } else {
      setSelectedEbookEnid(label.enid)
    }

    try {
      const content = await api.home.labelContent(kind === "course" ? 4 : 2, label.enid)
      updateSectionContent(kind, content)
    } catch (err) {
      setSectionError(err instanceof Error ? err.message : "首页模块刷新失败")
    } finally {
      setSwitchingSection(null)
    }
  }, [updateSectionContent])

  const moduleMap = useMemo(
    () =>
      new Map(
        (data?.homeData.moduleList ?? [])
          .filter((item) => item.isShow === 3)
          .map((item) => [item.type, item] as const),
      ),
    [data?.homeData.moduleList],
  )

  return {
    loading,
    error,
    sectionError,
    switchingSection,
    homeData: data?.homeData,
    freeResources: data?.freeResources?.list ?? [],
    freeResourcesError: data?.freeResourcesError ?? null,
    ebookLabels: data?.ebookLabels?.list ?? [],
    ebookContent: data?.ebookContent ?? null,
    ebookLabelsError: data?.ebookLabelsError ?? null,
    ebookContentError: data?.ebookContentError ?? null,
    courseLabels: data?.courseLabels?.list ?? [],
    courseContent: data?.courseContent ?? null,
    courseLabelsError: data?.courseLabelsError ?? null,
    courseContentError: data?.courseContentError ?? null,
    selectedCourseEnid,
    selectedEbookEnid,
    moduleMap,
    reload: loadPortal,
    selectLabel,
  }
}
