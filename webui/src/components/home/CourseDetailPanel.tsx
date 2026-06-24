import { Loader2 } from "lucide-react"
import type { CourseInfoResponse } from "@/api"
import { InfoBlock, StatusBadge } from "@/components/ui/Semantic"
import { Card } from "@/components/ui/Card"
import { semanticMetaTextClass, semanticSecondaryTextClass } from "@/lib/semanticStyles"

type CourseDetailPanelProps = {
  detail: CourseInfoResponse | null
  loading: boolean
}

export function CourseDetailPanel({ detail, loading }: CourseDetailPanelProps) {
  const highlightItems = detail?.items.filter((item) => item.content?.trim()).slice(0, 4) ?? []

  return (
    <Card>
      <div className="space-y-4 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold text-text-primary">课程详情</h2>
            <p className={`mt-1 ${semanticMetaTextClass}`}>点击左侧列表或搜索建议查看详情。</p>
          </div>
          {loading ? <Loader2 className="h-4 w-4 animate-spin text-text-muted" /> : null}
        </div>

        {detail ? (
          <>
            <div className="grid gap-4 md:grid-cols-[180px_1fr]">
              <img
                alt={detail.class_info.name}
                className="h-44 w-full rounded-3xl object-cover"
                src={detail.class_info.square_img || "https://placehold.co/320x320/e2e8f0/334155?text=Course"}
              />
              <div className="space-y-3">
                <StatusBadge variant="accent">
                  {detail.class_info.price_desc || "课程详情"}
                </StatusBadge>
                <h3 className="text-2xl font-semibold text-text-primary">{detail.class_info.name}</h3>
                <p className={`text-sm leading-7 ${semanticSecondaryTextClass}`}>{detail.class_info.intro || detail.class_info.share_summary || "暂无简介"}</p>
                <div className="grid grid-cols-2 gap-3">
                  <InfoBlock>讲师：{detail.class_info.lecturer_name || "未知"}</InfoBlock>
                  <InfoBlock>学习人数：{detail.class_info.learn_user_count || 0}</InfoBlock>
                  <InfoBlock>文章数：{detail.class_info.current_article_count || 0}</InfoBlock>
                  <InfoBlock>头衔：{detail.class_info.lecturer_title || "-"}</InfoBlock>
                </div>
              </div>
            </div>

            <div className="space-y-3 border-t border-border pt-4">
              <h3 className="font-medium text-text-primary">课程亮点</h3>
              {highlightItems.map((item, index) => (
                <InfoBlock key={`${item.title}-${index}`}>
                  <p className="font-medium text-text-primary">{item.title}</p>
                  <p className="mt-2 whitespace-pre-wrap text-sm leading-6 text-text-secondary">{item.content}</p>
                </InfoBlock>
              ))}
            </div>
          </>
        ) : (
          <div className="flex min-h-80 items-center justify-center rounded-3xl border border-dashed border-border text-sm text-text-muted">
            先选择一门课程，这里会展示详情。
          </div>
        )}
      </div>
    </Card>
  )
}
