import { useEffect } from "react"
import { BadgeCheck, BookOpen, Search } from "lucide-react"
import { useNavigate } from "react-router-dom"
import { QrLoginCard } from "@/components/auth/QrLoginCard"
import { Card } from "@/components/ui/Card"
import { ThemeToggleButton } from "@/components/ui/ThemeToggleButton"
import { useAuth } from "@/providers/AuthProvider"

const loginHighlights = [
  { icon: BadgeCheck, title: "扫码登录", text: "打开页面后即可使用二维码登录当前账号。" },
  { icon: Search, title: "统一协议", text: "所有接口统一返回 code、msg、data。" },
  { icon: BookOpen, title: "学习工作台", text: "登录后可继续查看课程、电子书、听书和用户信息。" },
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
        <div className="flex flex-wrap items-center justify-between gap-3">
          <span className="w-fit rounded-full bg-primary/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-primary">
            dedao-dl Web
          </span>
          <ThemeToggleButton className="rounded-2xl" />
        </div>

        <div className="space-y-4">
          <h1 className="max-w-2xl text-4xl font-semibold tracking-tight text-slate-950 sm:text-5xl">
            用浏览器打开 dedao-dl，扫码后直接进入学习工作台。
          </h1>
          <p className="max-w-xl text-base leading-7 text-slate-600">
            登录成功后，可以继续查看课程、电子书、听书和用户中心等内容。
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
