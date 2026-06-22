import { useEffect } from "react"
import { Compass, LayoutGrid, Library, Sparkles } from "lucide-react"
import { useSearchParams } from "react-router-dom"
import { CategoryPanel } from "@/components/home/CategoryPanel"
import { CourseDetailPanel } from "@/components/home/CourseDetailPanel"
import { CourseListPanel } from "@/components/home/CourseListPanel"
import { Card } from "@/components/ui/Card"
import { useCourseWorkspace } from "@/hooks/useCourseWorkspace"

const sections = [
  {
    icon: LayoutGrid,
    title: "独立工作区",
    text: "分类、课程列表和详情已经从首页迁到这里，后续继续对齐 dedao-gui 的 Course 视图。",
  },
  {
    icon: Compass,
    title: "课程筛选",
    text: "后续可以继续补筛选器、分组浏览和分页加载，让课程页承担完整查询路径。",
  },
  {
    icon: Sparkles,
    title: "详情承载",
    text: "从首页搜索建议跳过来时，会直接尝试打开对应课程详情，方便搜索和浏览联动。",
  },
  {
    icon: Library,
    title: "执行入口",
    text: "下一步可以继续把 CLI 下载任务映射进来，让课程页成为查询与执行的统一入口。",
  },
]

export function CoursePage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const initialCategory = searchParams.get("category")
  const initialEnid = searchParams.get("enid")
  const workspace = useCourseWorkspace({
    initialCategory,
    initialEnid,
  })

  useEffect(() => {
    if (!workspace.selectedCategory) {
      return
    }

    if (searchParams.get("category") === workspace.selectedCategory) {
      return
    }

    const next = new URLSearchParams(searchParams)
    next.set("category", workspace.selectedCategory)
    setSearchParams(next, { replace: true })
  }, [searchParams, setSearchParams, workspace.selectedCategory])

  const handleSelectCategory = (category: string) => {
    const next = new URLSearchParams(searchParams)
    next.set("category", category)
    next.delete("enid")
    setSearchParams(next, { replace: true })
    workspace.setSelectedCategory(category)
  }

  const handleSelectCourse = (enid: string) => {
    const next = new URLSearchParams(searchParams)
    next.set("category", workspace.selectedCategory)
    next.set("enid", enid)
    setSearchParams(next, { replace: true })
    void workspace.loadCourseDetail(enid)
  }

  return (
    <main className="space-y-6">
      <section className="rounded-3xl border border-white/70 bg-white/90 p-6 shadow-soft backdrop-blur">
        <p className="text-sm text-slate-500">课程模块</p>
        <h2 className="mt-2 text-3xl font-semibold text-slate-950">课程工作区</h2>
        <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">
          这一页现在承载课程分类、课程列表和课程详情。首页只做总览和分发，课程查询路径收口到独立工作区里。
        </p>
      </section>

      {workspace.error ? (
        <Card className="border border-rose-200 bg-rose-50">
          <div className="p-4 text-sm text-rose-700">{workspace.error}</div>
        </Card>
      ) : null}

      <CategoryPanel
        categories={workspace.categories}
        onRefresh={() => void workspace.bootstrap()}
        onSelectCategory={handleSelectCategory}
        selectedCategory={workspace.selectedCategory}
      />

      <section className="grid gap-6 lg:grid-cols-4">
        {sections.map((section) => (
          <Card className="p-6" key={section.title}>
            <section.icon className="h-6 w-6 text-primary" />
            <h3 className="mt-4 text-lg font-semibold text-slate-950">{section.title}</h3>
            <p className="mt-2 text-sm leading-7 text-slate-600">{section.text}</p>
          </Card>
        ))}
      </section>

      <section className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
        <CourseListPanel
          courses={workspace.courses}
          loading={workspace.loadingCourses}
          onSelectCourse={handleSelectCourse}
          selectedCategory={workspace.selectedCategory}
        />
        <CourseDetailPanel detail={workspace.detail} loading={workspace.loadingDetail} />
      </section>
    </main>
  )
}
