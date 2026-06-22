import { ExternalLink, MessageSquare, PanelRightOpen } from "lucide-react"
import { type CourseListItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"

export function PurchasedEbookPage() {
  return (
    <PurchasedCollectionPage
      category="ebook"
      description="这里集中展示我已购买或已加入书架的电子书资产，详情页继续作为发现域和已购域共用的消费页。"
      emptyDescription="后续可以继续补充书评、下载和书架状态等交互。"
      emptyTitle="当前没有可展示的已购电子书"
      getPrimaryMeta={(item: CourseListItem) => item.author || item.lecturer_name || "电子书内容"}
      getSecondaryMeta={(item: CourseListItem) => (item.price ? `¥${item.price}` : "查看详情")}
      icon="ebook"
      itemLabel="电子书"
      loadingText="正在加载已购电子书..."
      renderActions={(item, helpers) =>
        item.is_group ? (
          <Button className="h-9 px-3" onClick={() => helpers.openItem(item)} variant="outline">
            <PanelRightOpen className="mr-2 h-4 w-4" />
            进入分组
          </Button>
        ) : (
          <>
            <Button className="h-9 px-3" onClick={() => helpers.openItem(item)} variant="outline">
              <PanelRightOpen className="mr-2 h-4 w-4" />
              查看详情
            </Button>
            {item.dd_url ? (
              <Button className="h-9 px-3" onClick={() => window.open(item.dd_url, "_blank", "noopener,noreferrer")} variant="ghost">
                <ExternalLink className="mr-2 h-4 w-4" />
                去得到打开
              </Button>
            ) : null}
            <Button className="h-9 px-3" disabled variant="ghost">
              <MessageSquare className="mr-2 h-4 w-4" />
              书评后续接入
            </Button>
          </>
        )
      }
      onOpenItem={(item, navigate) => {
        if (!item.enid) {
          return
        }

        navigate(`/ebooks/${encodeURIComponent(item.enid)}?from=purchased-ebook&parentTitle=${encodeURIComponent(item.title || item.name || "电子书详情")}`)
      }}
      title="已购电子书"
    />
  )
}
