import { type CourseListItem } from "@/api"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"

export function PurchasedCompassPage() {
  return (
    <PurchasedCollectionPage
      category="compass"
      description="这里集中展示我已拥有的锦囊内容。锦囊当前先保留资产列表与外链打开能力，后续再补站内详情或下载链路。"
      emptyDescription="如果后面确认锦囊也要做站内详情页，可以在这条资产链路上继续接。"
      emptyTitle="当前没有可展示的已购锦囊"
      getPrimaryMeta={(item: CourseListItem) => item.author || "锦囊内容"}
      getSecondaryMeta={(item: CourseListItem) => (item.dd_url ? "去得到打开" : "详情后续接入")}
      icon="compass"
      itemLabel="锦囊"
      loadingText="正在加载已购锦囊..."
      onOpenItem={(item) => {
        if (!item.dd_url) {
          return
        }

        window.open(item.dd_url, "_blank", "noopener,noreferrer")
      }}
      title="已购锦囊"
    />
  )
}
