import { Loader2 } from "lucide-react"
import type { CourseInfoResponse } from "@/api"
import { Card } from "@/components/ui/Card"

type CourseDetailPanelProps = {
  detail: CourseInfoResponse | null
  loading: boolean
}

export function CourseDetailPanel({ detail, loading }: CourseDetailPanelProps) {
  return (
    <Card>
      <div className="space-y-4 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold text-slate-950">课程详情</h2>
            <p className="mt-1 text-sm text-slate-500">点击左侧列表或搜索建议查看详情。</p>
          </div>
          {loading ? <Loader2 className="h-4 w-4 animate-spin text-slate-400" /> : null}
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
                <span className="inline-flex rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary">
                  {detail.class_info.price_desc || "课程详情"}
                </span>
                <h3 className="text-2xl font-semibold text-slate-950">{detail.class_info.name}</h3>
                <p className="text-sm leading-7 text-slate-600">{detail.class_info.intro || detail.class_info.share_summary || "暂无简介"}</p>
                <div className="grid grid-cols-2 gap-3 text-sm text-slate-600">
                  <div className="rounded-2xl bg-slate-50 p-3">讲师：{detail.class_info.lecturer_name || "未知"}</div>
                  <div className="rounded-2xl bg-slate-50 p-3">学习人数：{detail.class_info.learn_user_count || 0}</div>
                  <div className="rounded-2xl bg-slate-50 p-3">文章数：{detail.class_info.current_article_count || 0}</div>
                  <div className="rounded-2xl bg-slate-50 p-3">头衔：{detail.class_info.lecturer_title || "-"}</div>
                </div>
              </div>
            </div>

            <div className="space-y-3 border-t border-slate-200 pt-4">
              <h3 className="font-medium text-slate-900">课程亮点</h3>
              {detail.items.slice(0, 4).map((item, index) => (
                <div className="rounded-2xl bg-slate-50 p-4" key={`${item.title}-${index}`}>
                  <p className="font-medium text-slate-900">{item.title}</p>
                  <p className="mt-2 text-sm leading-6 text-slate-600">{item.content || "暂无内容"}</p>
                </div>
              ))}
            </div>
          </>
        ) : (
          <div className="flex min-h-80 items-center justify-center rounded-3xl border border-dashed border-slate-200 text-sm text-slate-500">
            先选择一门课程，这里会展示详情。
          </div>
        )}
      </div>
    </Card>
  )
}
