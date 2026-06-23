import { BookMarked, Compass, GraduationCap, Headphones, Home, Loader2, Search, UserCircle2 } from "lucide-react"
import type { ComponentType } from "react"
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { ThemeToggleButton } from "@/components/ui/ThemeToggleButton"
import { cn } from "@/lib/cn"
import { useAuth } from "@/providers/AuthProvider"

const primaryNavItems = [
  { to: "/", label: "首页", icon: Home, end: true },
]

const purchasedNavItems = [
  { to: "/purchased/courses", label: "课程", icon: GraduationCap },
  { to: "/purchased/ebooks", label: "电子书", icon: BookMarked },
  { to: "/purchased/audios", label: "听书", icon: Headphones },
  { to: "/purchased/compass", label: "锦囊", icon: Compass },
]

const accountNavItems = [
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

  const isActiveNav = (to: string, end?: boolean) => {
    if (end) {
      return location.pathname === to
    }

    if (to === "/purchased/courses") {
      return location.pathname.startsWith("/purchased/courses") || location.pathname.startsWith("/courses/")
    }

    if (to === "/purchased/ebooks") {
      return location.pathname.startsWith("/purchased/ebooks") || location.pathname.startsWith("/ebooks/")
    }

    if (to === "/purchased/audios") {
      return (
        location.pathname.startsWith("/purchased/audios") ||
        location.pathname.startsWith("/audios/") ||
        location.pathname.startsWith("/audio-groups/")
      )
    }

    return location.pathname.startsWith(to)
  }

  const renderNavLink = (item: { to: string; label: string; icon: ComponentType<{ className?: string }>; end?: boolean }) => (
    <Link
      className={cn(
        "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-medium transition",
        isActiveNav(item.to, item.end)
          ? "bg-accent-soft text-accent"
          : "text-text-inverse/78 hover:bg-white/10 hover:text-text-inverse",
      )}
      key={item.to}
      to={item.to}
    >
      <item.icon className="h-4 w-4" />
      {item.label}
    </Link>
  )

  return (
    <div className="min-h-screen bg-surface-page text-text-primary">
      <div className="mx-auto grid min-h-screen max-w-[1600px] grid-cols-1 gap-6 px-4 py-4 lg:grid-cols-[260px_1fr] lg:px-6">
        <aside className="rounded-3xl border border-border bg-surface-inverse p-5 text-text-inverse shadow-soft">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-accent/15 text-accent">
              <GraduationCap className="h-6 w-6" />
            </div>
            <div>
              <p className="text-lg font-semibold">dedao-dl Web</p>
              <p className="text-xs text-text-inverse/68">统一 Web 工作台</p>
            </div>
          </div>

          <nav className="mt-8 space-y-6">
            <div className="space-y-2">
              {primaryNavItems.map((item) => renderNavLink(item))}
            </div>

            <div>
              <p className="px-4 text-xs font-medium uppercase tracking-[0.18em] text-text-inverse/46">已购</p>
              <div className="mt-2 space-y-2">
                {purchasedNavItems.map((item) => renderNavLink(item))}
              </div>
            </div>

            <div className="space-y-2">
              {accountNavItems.map((item) => renderNavLink(item))}
            </div>
          </nav>

        </aside>

        <div className="flex min-h-0 flex-col gap-6 pb-28">
          <header className="rounded-3xl border border-border bg-surface-panel px-5 py-4 shadow-soft backdrop-blur">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
              <div>
                <p className="text-sm text-text-muted">工作台</p>
                <h1 className="text-2xl font-semibold text-text-primary">登录后工作台</h1>
              </div>

              <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
                <div className="flex items-center gap-3 rounded-2xl border border-border bg-surface-soft px-4 py-3 text-sm text-text-muted">
                  <Search className="h-4 w-4" />
                  <span>全局搜索建设中</span>
                </div>

                <ThemeToggleButton className="rounded-2xl" />

                <div className="flex items-center gap-3 rounded-2xl bg-surface-inverse px-4 py-3 text-text-inverse">
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
                    <p className="truncate text-xs text-text-inverse/66">{user ? "当前账号" : "等待用户信息"}</p>
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
