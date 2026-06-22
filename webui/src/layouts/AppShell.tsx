import { Home, LayoutGrid, Loader2, Search, UserCircle2 } from "lucide-react"
import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { cn } from "@/lib/cn"
import { useAuth } from "@/providers/AuthProvider"

const navItems = [
  { to: "/", label: "首页", icon: Home, end: true },
  { to: "/courses", label: "课程", icon: LayoutGrid },
  { to: "/user", label: "用户", icon: UserCircle2 },
]

export function AppShell() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, loading, logout } = useAuth()

  const handleLogout = async () => {
    await logout()
    navigate("/login", { replace: true })
  }

  return (
    <div className="min-h-screen">
      <div className="mx-auto grid min-h-screen max-w-[1600px] grid-cols-1 gap-6 px-4 py-4 lg:grid-cols-[260px_1fr] lg:px-6">
        <aside className="rounded-3xl border border-white/70 bg-slate-950 p-5 text-slate-50 shadow-soft">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/20 text-primary">
              <LayoutGrid className="h-6 w-6" />
            </div>
            <div>
              <p className="text-lg font-semibold">dedao-dl Web</p>
              <p className="text-xs text-slate-400">统一 Web 工作台</p>
            </div>
          </div>

          <nav className="mt-8 space-y-2">
            {navItems.map((item) => (
              <NavLink
                className={({ isActive }) =>
                  cn(
                    "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-medium transition",
                    isActive
                      ? "bg-white text-slate-950"
                      : "text-slate-300 hover:bg-slate-900 hover:text-white",
                  )
                }
                end={item.end}
                key={item.to}
                to={item.to}
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </NavLink>
            ))}
          </nav>

          <div className="mt-8 rounded-3xl bg-slate-900 p-4">
            <p className="text-xs uppercase tracking-[0.18em] text-slate-500">当前页面</p>
            <p className="mt-2 text-sm text-slate-200">{location.pathname}</p>
            <p className="mt-3 text-sm leading-6 text-slate-400">这一层只承载登录后页面的共用导航、顶部栏和内容容器。</p>
          </div>
        </aside>

        <div className="flex min-h-0 flex-col gap-6 pb-28">
          <header className="rounded-3xl border border-white/70 bg-white/90 px-5 py-4 shadow-soft backdrop-blur">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
              <div>
                <p className="text-sm text-slate-500">应用壳层</p>
                <h1 className="text-2xl font-semibold text-slate-950">登录后工作台</h1>
              </div>

              <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
                <div className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-500">
                  <Search className="h-4 w-4" />
                  <span>全局搜索入口预留，后续可接课程/电子书/文章统一检索</span>
                </div>

                <div className="flex items-center gap-3 rounded-2xl bg-slate-950 px-4 py-3 text-slate-50">
                  {loading ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <img
                      alt={user?.nickname ?? "avatar"}
                      className="h-10 w-10 rounded-2xl object-cover"
                      src={user?.avatar || "https://placehold.co/80x80/e2e8f0/334155?text=DD"}
                    />
                  )}
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium">{user?.nickname ?? "未登录"}</p>
                    <p className="truncate text-xs text-slate-400">{user?.uid_hazy ?? "等待用户信息"}</p>
                  </div>
                  <Button className="rounded-xl" onClick={() => void handleLogout()} variant="outline">
                    退出
                  </Button>
                </div>
              </div>
            </div>
          </header>

          <div className="min-h-0 flex-1">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  )
}
