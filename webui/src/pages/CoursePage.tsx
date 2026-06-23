import { type CourseListItem } from "@/api"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"

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
      onOpenItem={(item, pageNavigate) => {
        if (!item.enid) {
          return
        }

        pageNavigate(
          `/courses/${encodeURIComponent(item.enid)}/articles?from=purchased-course&parentTitle=${encodeURIComponent(item.title || item.name || "课程内容")}`,
        )
      }}
      title="已购课程"
    />
  )
}
