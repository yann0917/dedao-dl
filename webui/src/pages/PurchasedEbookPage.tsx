import { MessageSquare, PanelRightOpen } from "lucide-react"
import { type CourseListItem } from "@/api"
import { Button } from "@/components/ui/Button"
import { PurchasedCollectionPage } from "@/components/purchased/PurchasedCollectionPage"

export function PurchasedEbookPage() {
  return (
    <PurchasedCollectionPage
      category="ebook"
      coverContainerClassName="aspect-[3/4] p-3"
      coverImageClassName="object-contain"
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
            <Button
              className="h-9 px-3"
              onClick={() =>
                helpers.navigate(
                  `/ebooks/${encodeURIComponent(item.enid)}/comments?title=${encodeURIComponent(item.title || item.name || "电子书书评")}`,
                )
              }
              variant="ghost"
            >
              <MessageSquare className="mr-2 h-4 w-4" />
              查看书评
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
