import { useEffect } from "react"
import { BadgeCheck, BookOpen, Search } from "lucide-react"
import { useNavigate } from "react-router-dom"
import { QrLoginCard } from "@/components/auth/QrLoginCard"
import { Card } from "@/components/ui/Card"
import { useAuth } from "@/providers/AuthProvider"

const loginHighlights = [
  { icon: BadgeCheck, title: "登录主路径", text: "扫码登录不再是兜底，而是首屏入口。" },
  { icon: Search, title: "统一协议", text: "所有接口统一返回 code、msg、data。" },
  { icon: BookOpen, title: "后续可扩", text: "后面可以把现有 CLI 能力逐步映射成 API。" },
]

export function LoginPage() {
  const navigate = useNavigate()
  const { loggedIn, loading, completeLogin } = useAuth()

  useEffect(() => {
    if (!loading && loggedIn) {
      navigate("/", { replace: true })
    }
  }, [loading, loggedIn, navigate])

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-7xl items-center gap-10 px-6 py-12 lg:px-10">
      <section className="grid flex-1 gap-6">
        <span className="w-fit rounded-full bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-primary">
          dedao-dl Web
        </span>

        <div className="space-y-4">
          <h1 className="max-w-2xl text-4xl font-semibold tracking-tight text-slate-950 sm:text-5xl">
            用浏览器打开 dedao-dl，把扫码登录和查询工作台合成一条顺滑链路。
          </h1>
          <p className="max-w-xl text-base leading-7 text-slate-600">
            CLI 负责稳定能力边界，Web 负责交互壳。登录成功后，课程列表、搜索建议和课程详情都继续复用现有 services。
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-3">
          {loginHighlights.map((item) => (
            <Card key={item.title}>
              <div className="space-y-3 p-5">
                <item.icon className="h-5 w-5 text-primary" />
                <div>
                  <h2 className="font-medium text-slate-900">{item.title}</h2>
                  <p className="mt-1 text-sm leading-6 text-slate-600">{item.text}</p>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </section>

      <section className="w-full max-w-xl">
        <QrLoginCard
          onLoginSuccess={async (user) => {
            await completeLogin(user)
            navigate("/", { replace: true })
          }}
        />
      </section>
    </main>
  )
}
