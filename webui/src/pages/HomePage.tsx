import { useMemo } from "react"
import { useNavigate } from "react-router-dom"
import { type HomeCategory, type HomeFreeResource, type HomeNavigation } from "@/api"
import { HomeBannerCarousel } from "@/components/home/HomeBannerCarousel"
import { HomeCategoryMenu } from "@/components/home/HomeCategoryMenu"
import { FreeResourceSection, LabeledShelfSection } from "@/components/home/HomePortalSections"
import { HomePortalUserCard } from "@/components/home/HomePortalUserCard"
import { Card } from "@/components/ui/Card"
import { Skeleton } from "@/components/ui/Semantic"
import { Button } from "@/components/ui/Button"
import { buildCategoryQuery } from "@/hooks/useCategoryExplorer"
import { useHomePortal } from "@/hooks/useHomePortal"
import { semanticMetaTextClass, semanticPageSectionClass } from "@/lib/semanticStyles"

function resolveCategoryProductType(navType: number) {
  if (navType === 2) {
    return "2"
  }

  if (navType === 4) {
    return "66"
  }

  return "0"
}

export function HomePage() {
  const navigate = useNavigate()
  const portal = useHomePortal()

  const handleSelectCourse = (enid: string) => {
    // 精选课程为 AI 学习圈时，点击任意卡片跳转 AI 学习圈
    if (portal.selectedCourseEnid.startsWith("aiSphereGroupType:")) {
      navigate("/ai-channel")
      return
    }

    navigate(`/courses/${encodeURIComponent(enid)}?from=home`)
  }

  const handleSelectEbook = (enid: string) => {
    navigate(`/ebooks/${encodeURIComponent(enid)}?from=home`)
  }

  const handleNavigateCategory = (category: HomeCategory, labelEnid: string) => {
    navigate(
      `/category?${buildCategoryQuery({
        id: category.id,
        name: category.name,
        navType: category.navType,
        enid: category.enid,
        labelId: labelEnid,
        productType: resolveCategoryProductType(category.navType),
      })}`,
    )
  }

  const handleOpenFreeResource = (resource: HomeFreeResource) => {
    if (!resource.enid) {
      return
    }

    navigate(
      `/courses/${encodeURIComponent(resource.enid)}?from=home&parentTitle=${encodeURIComponent(resource.name || "课程详情")}`,
    )
  }

  const handleSelectShelfLabel = (kind: "course" | "ebook", label: HomeNavigation) => {
    void portal.selectLabel(kind, label)
  }

  const handleOpenCategoryFromShelf = (kind: "course" | "ebook") => {
    const selectedLabel = (kind === "course" ? portal.courseLabels : portal.ebookLabels).find((item) =>
      kind === "course" ? item.enid === portal.selectedCourseEnid : item.enid === portal.selectedEbookEnid,
    )

    if (!selectedLabel) {
      return
    }

    // 精选课程的 enid 带 aiSphereGroupType: 前缀时跳转 AI 学习圈
    if (kind === "course" && selectedLabel.enid.startsWith("aiSphereGroupType:")) {
      navigate("/ai-channel")
      return
    }

    navigate(
      `/category?${buildCategoryQuery({
        id: selectedLabel.id,
        name: selectedLabel.name,
        navType: selectedLabel.nav_type || (kind === "course" ? 4 : 2),
        enid: selectedLabel.enid,
        labelId: "",
        productType: kind === "course" ? "66" : "2",
      })}`,
    )
  }

  // 当前选中的课程 label 名称，用于「查看更多」按钮文案
  const selectedCourseLabelName = useMemo(
    () => portal.courseLabels.find((item) => item.enid === portal.selectedCourseEnid)?.name,
    [portal.courseLabels, portal.selectedCourseEnid],
  )

  if (portal.loading) {
    return (
      <main aria-busy="true" className="space-y-10">
        <section className="grid gap-6 xl:grid-cols-[260px_minmax(0,1fr)_280px]">
          <Card className="h-[407px] space-y-2 p-4">
            {Array.from({ length: 9 }).map((_, index) => (
              <Skeleton className="h-7 w-full" key={index} />
            ))}
          </Card>
          <Card className="aspect-video w-full">
            <Skeleton className="h-full w-full" />
          </Card>
          <Card className="h-[407px] p-5">
            <div className="flex items-center gap-4">
              <Skeleton className="h-14 w-14" />
              <div className="flex-1 space-y-2">
                <Skeleton className="h-4 w-20" />
                <Skeleton className="h-3 w-14" />
              </div>
            </div>
            <div className="mt-5 grid grid-cols-2">
              <Skeleton className="h-16" />
              <Skeleton className="h-16" />
            </div>
            <Skeleton className="mt-5 h-10 w-full" />
          </Card>
        </section>

        <section className="space-y-5">
          <div className="border-b border-border pb-3">
            <Skeleton className="h-5 w-40" />
          </div>
          <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
            {Array.from({ length: 4 }).map((_, index) => (
              <Card key={index}>
                <Skeleton className="aspect-16/10 w-full" />
                <div className="space-y-2 p-4">
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-3 w-full" />
                </div>
              </Card>
            ))}
          </div>
        </section>
      </main>
    )
  }

  return (
    <main className="space-y-8">
      {portal.error ? (
        <Card className="border-danger bg-danger-soft">
          <div className="flex flex-col gap-3 p-4 text-sm text-danger md:flex-row md:items-center md:justify-between">
            <span>{portal.error}</span>
            <Button onClick={() => void portal.reload()} variant="outline">
              重试
            </Button>
          </div>
        </Card>
      ) : null}

      {portal.sectionError ? (
        <Card className="border-warning bg-warning-soft">
          <div className="p-4 text-sm text-warning">{portal.sectionError}</div>
        </Card>
      ) : null}

      <section className="grid gap-6 xl:grid-cols-[260px_minmax(0,1fr)_280px]">
        <HomeCategoryMenu
          categories={portal.homeData?.categoryList ?? []}
          onNavigateCategory={handleNavigateCategory}
        />
        <HomeBannerCarousel banners={portal.homeData?.banner ?? []} />
        <HomePortalUserCard />
      </section>

      <FreeResourceSection
        error={portal.freeResourcesError}
        module={portal.moduleMap.get("free_class")}
        onOpenResource={handleOpenFreeResource}
        resources={portal.freeResources}
      />

      <LabeledShelfSection
        content={portal.courseContent}
        error={portal.courseLabelsError || portal.courseContentError}
        labels={portal.courseLabels}
        loading={portal.switchingSection === "course"}
        module={portal.moduleMap.get("class")}
        moreButtonText={`查看更多${selectedCourseLabelName || "精选课程"}`}
        onClickMore={() => handleOpenCategoryFromShelf("course")}
        onOpenProduct={handleSelectCourse}
        onSelectLabel={(label) => handleSelectShelfLabel("course", label)}
        selectedEnid={portal.selectedCourseEnid}
        variant="course"
      />

      <LabeledShelfSection
        content={portal.ebookContent}
        error={portal.ebookLabelsError || portal.ebookContentError}
        labels={portal.ebookLabels}
        loading={portal.switchingSection === "ebook"}
        module={portal.moduleMap.get("ebook")}
        moreButtonText="查看更多电子书"
        onClickMore={() => handleOpenCategoryFromShelf("ebook")}
        onOpenProduct={handleSelectEbook}
        onSelectLabel={(label) => handleSelectShelfLabel("ebook", label)}
        selectedEnid={portal.selectedEbookEnid}
        variant="ebook"
      />
    </main>
  )
}
