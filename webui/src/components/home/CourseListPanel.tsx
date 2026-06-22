import { BookOpen, Loader2 } from "lucide-react"
import type { CourseListItem } from "@/api"
import { Card } from "@/components/ui/Card"

type CourseListPanelProps = {
  selectedCategory: string
  courses: CourseListItem[]
  loading: boolean
  onSelectCourse: (enid: string) => void
}

export function CourseListPanel({
  selectedCategory,
  courses,
  loading,
  onSelectCourse,
}: CourseListPanelProps) {
  return (
    <Card>
      <div className="space-y-4 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="flex items-center gap-2 text-xl font-semibold text-slate-950">
              <BookOpen className="h-5 w-5 text-primary" />
              课程列表
            </h2>
            <p className="mt-1 text-sm text-slate-500">{selectedCategory} 分类下的已购课程</p>
          </div>
          {loading ? <Loader2 className="h-4 w-4 animate-spin text-slate-400" /> : null}
        </div>

        <div className="space-y-3">
          {courses.map((course) => (
            <button
              className="grid w-full grid-cols-[84px_1fr] gap-4 rounded-3xl border border-slate-200 p-3 text-left transition hover:border-primary hover:bg-primary/5"
              key={course.enid || course.id}
              onClick={() => course.enid && onSelectCourse(course.enid)}
              type="button"
            >
              <img
                alt={course.title || course.name}
                className="h-20 w-20 rounded-2xl object-cover"
                src={course.cover || course.index_img || "https://placehold.co/160x160/e2e8f0/334155?text=DD"}
              />
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-3">
                  <p className="font-medium text-slate-900">{course.title || course.name}</p>
                  <span className="rounded-full bg-slate-100 px-3 py-1 text-xs text-slate-500">{course.price_desc || "已购"}</span>
                </div>
                <p className="line-clamp-2 text-sm leading-6 text-slate-500">{course.intro || course.subtitle || "暂无简介"}</p>
              </div>
            </button>
          ))}
        </div>
      </div>
    </Card>
  )
}
