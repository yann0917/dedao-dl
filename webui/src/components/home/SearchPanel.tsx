import { Search } from "lucide-react"
import type { SearchSuggestResponse } from "@/api"
import { Card } from "@/components/ui/Card"

type SearchPanelProps = {
  query: string
  hotKeywords: string[]
  suggestions: SearchSuggestResponse | null
  onQueryChange: (value: string) => void
  onSelectCourse: (enid: string) => void
}

export function SearchPanel({
  query,
  hotKeywords,
  suggestions,
  onQueryChange,
  onSelectCourse,
}: SearchPanelProps) {
  return (
    <div className="space-y-4">
      <div className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3">
        <div className="flex items-center gap-3">
          <Search className="h-4 w-4 text-slate-500" />
          <input
            className="w-full border-0 bg-transparent text-sm outline-none placeholder:text-slate-400"
            onChange={(event) => onQueryChange(event.target.value)}
            placeholder="试试搜课程、作者或关键字"
            value={query}
          />
        </div>

        <div className="mt-3 flex flex-wrap gap-2">
          {hotKeywords.map((keyword) => (
            <button
              className="rounded-full border border-slate-200 bg-white px-3 py-1 text-xs text-slate-600 transition hover:border-primary hover:text-primary"
              key={keyword}
              onClick={() => onQueryChange(keyword)}
              type="button"
            >
              {keyword}
            </button>
          ))}
        </div>
      </div>

      <Card className="border-dashed shadow-none">
        <div className="space-y-3 p-5">
          <div>
            <h2 className="font-medium text-slate-900">搜索建议</h2>
            <p className="mt-1 text-sm text-slate-500">这里直接走 `/api/search/suggest`。</p>
          </div>

          {suggestions?.list.flatMap((group) => group.list).slice(0, 5).map((item) => (
            <button
              className="flex w-full items-center justify-between rounded-2xl border border-slate-200 px-4 py-3 text-left transition hover:border-primary hover:bg-primary/5"
              key={`${item.id}-${item.title}`}
              onClick={() => item.extra?.enid && onSelectCourse(item.extra.enid)}
              type="button"
            >
              <div>
                <p className="font-medium text-slate-900">{item.title}</p>
                <p className="mt-1 text-sm text-slate-500">{item.author || item.content || "来自搜索建议"}</p>
              </div>
              <span className="rounded-full bg-slate-100 px-3 py-1 text-xs text-slate-500">建议</span>
            </button>
          ))}

          {!suggestions?.list.length ? <p className="text-sm text-slate-500">输入关键字后，这里会显示建议。</p> : null}
        </div>
      </Card>
    </div>
  )
}
