import { Loader2 } from "lucide-react"
import { useNavigate } from "react-router-dom"
import { type HomeCategory, type HomeFreeResource, type HomeNavigation } from "@/api"
import { HomeBannerCarousel } from "@/components/home/HomeBannerCarousel"
import { HomeCategoryMenu } from "@/components/home/HomeCategoryMenu"
import { FreeResourceSection, LabeledShelfSection } from "@/components/home/HomePortalSections"
import { HomePortalUserCard } from "@/components/home/HomePortalUserCard"
import { Card } from "@/components/ui/Card"
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

  if (portal.loading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center">
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载首页门户...
        </div>
      </main>
    )
  }

  return (
    <main className="space-y-8">
      <section className={`${semanticPageSectionClass} p-6`}>
        <p className={semanticMetaTextClass}>首页门户</p>
        <h2 className="mt-2 text-3xl font-semibold text-text-primary">内容发现主入口</h2>
      </section>

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
        moreButtonText="查看更多精选课程"
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
