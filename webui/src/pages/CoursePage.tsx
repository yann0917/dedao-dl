import { PanelRightOpen } from "lucide-react"
import { api, type CourseListItem } from "@/api"
import { QuickDownloadButtons, type QuickDownloadOption } from "@/components/download/QuickDownloadButtons"
import { Button } from "@/components/ui/Button"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"

const courseQuickDownloadOptions: QuickDownloadOption[] = [
  { value: 1, label: "音频" },
  { value: 2, label: "PDF" },
  { value: 3, label: "MD" },
]

export function CoursePage() {
  return (
    <PurchasedCollectionPage
      category="bauhinia"
      emptyTitle="当前没有可展示的已购课程"
      getPrimaryMeta={(item: CourseListItem) => `已更 ${item.publish_num || 0}/${item.course_num || 0}`}
      getProgress={(item: CourseListItem) => item.progress || 0}
      getSecondaryMeta={(item: CourseListItem) => `${item.progress || 0}%`}
      icon="course"
      itemLabel="课程"
      loadingText="正在加载已购课程..."
      renderActions={(item, helpers) => (
        <>
          <Button className="h-9 px-3" onClick={() => helpers.openItem(item)} variant="ghost">
            <PanelRightOpen className="mr-2 h-4 w-4" />
            详情
          </Button>
          <QuickDownloadButtons
            onDownload={(downloadType) =>
              api.download.course({
                enid: item.enid,
                title: item.title || item.name || "课程下载",
                downloadType,
                isOrder: true,
              })
            }
            options={courseQuickDownloadOptions}
          />
        </>
      )}
      onOpenItem={(item, pageNavigate) => {
        if (!item.enid) {
          return
        }

        pageNavigate(
          `/courses/${encodeURIComponent(item.enid)}?from=purchased-course&parentTitle=${encodeURIComponent(item.title || item.name || "课程详情")}`,
        )
      }}
      title="已购课程"
    />
  )
}
