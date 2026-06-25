import { type CourseListItem } from "@/api"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"

export function PurchasedCompassPage() {
  return (
    <PurchasedCollectionPage
      category="compass"
      emptyTitle="当前没有可展示的已购锦囊"
      getPrimaryMeta={(item: CourseListItem) => item.author || "锦囊内容"}
      getSecondaryMeta={() => ""}
      icon="compass"
      isItemInteractive={() => false}
      itemLabel="锦囊"
      loadingText="正在加载已购锦囊..."
      onOpenItem={() => {}}
      title="已购锦囊"
    />
  )
}
