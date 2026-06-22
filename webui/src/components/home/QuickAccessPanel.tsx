import { ArrowRight, BookOpen, UserCircle2 } from "lucide-react"
import { Link } from "react-router-dom"
import { Card } from "@/components/ui/Card"

const quickLinks = [
  {
    to: "/purchased/courses",
    title: "已购课程",
    description: "进入已购资产域，继续查看课程列表、分组和文章消费链路。",
    icon: BookOpen,
  },
  {
    to: "/user",
    title: "用户中心",
    description: "后续把账号管理、设置和用户概览继续沉淀到独立页面。",
    icon: UserCircle2,
  },
]

export function QuickAccessPanel() {
  return (
    <section className="grid gap-6 lg:grid-cols-2">
      {quickLinks.map((item) => (
        <Card className="p-6" key={item.to}>
          <item.icon className="h-6 w-6 text-primary" />
          <h2 className="mt-4 text-xl font-semibold text-slate-950">{item.title}</h2>
          <p className="mt-2 text-sm leading-7 text-slate-600">{item.description}</p>
          <Link
            className="mt-4 inline-flex items-center gap-2 text-sm font-medium text-primary transition hover:text-primary/80"
            to={item.to}
          >
            进入页面
            <ArrowRight className="h-4 w-4" />
          </Link>
        </Card>
      ))}
    </section>
  )
}
