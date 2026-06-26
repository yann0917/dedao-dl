import { RefreshCcw } from "lucide-react"
import type { CourseCategory } from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { semanticMetaTextClass } from "@/lib/semanticStyles"

type CategoryPanelProps = {
  categories: CourseCategory[]
  selectedCategory: string
  onRefresh: () => void
  onSelectCategory: (category: string) => void
}

export function CategoryPanel({
  categories,
  selectedCategory,
  onRefresh,
  onSelectCategory,
}: CategoryPanelProps) {
  return (
    <Card>
      <div className="space-y-4 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-text-primary">分类切换</h2>
            <p className={`mt-1 ${semanticMetaTextClass}`}>第一版先覆盖高频查询接口。</p>
          </div>
          <Button onClick={onRefresh} variant="outline">
            <RefreshCcw className="mr-2 h-4 w-4" />
            刷新
          </Button>
        </div>

        <div className="flex flex-wrap gap-2">
          {categories.map((category) => (
            <Button
              className="rounded-full"
              key={category.category}
              onClick={() => onSelectCategory(category.category)}
              variant={selectedCategory === category.category ? "default" : "outline"}
            >
              {category.name}
            </Button>
          ))}
        </div>
      </div>
    </Card>
  )
}
