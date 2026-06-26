import { BookMarked, Compass, GraduationCap, Headphones, Home, Loader2, Menu, Rows3, Search, Trophy, UserCircle2, X } from "lucide-react"
import type { ComponentType } from "react"
import { useEffect, useRef, useState } from "react"
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { ThemeToggleButton } from "@/components/ui/ThemeToggleButton"
import { cn } from "@/lib/cn"
import { useAuth } from "@/providers/AuthProvider"

const primaryNavItems = [
  { to: "/", label: "首页", icon: Home, end: true },
  { to: "/leaderboard", label: "得到榜单", icon: Trophy },
]

const purchasedNavItems = [
  { to: "/purchased/manage", label: "已购管理", icon: Rows3 },
  { to: "/purchased/courses", label: "课程", icon: GraduationCap },
  { to: "/purchased/ebooks", label: "电子书", icon: BookMarked },
  { to: "/purchased/audios", label: "听书", icon: Headphones },
  { to: "/purchased/compass", label: "锦囊", icon: Compass },
]

const accountNavItems = [
  { to: "/user", label: "用户", icon: UserCircle2 },
]

type DrawerSide = "left" | "right"

type FabPosition = {
  side: DrawerSide
  top: number
  left: number | null
}

const FAB_STORAGE_KEY = "dedao-dl-nav-fab-position"
const FAB_SIZE = 56
const FAB_MARGIN = 24
const DRAG_THRESHOLD = 6

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max)
}

function getTopBounds() {
  const minTop = FAB_MARGIN
  const maxTop =
    typeof window === "undefined" ? minTop : Math.max(minTop, window.innerHeight - FAB_SIZE - FAB_MARGIN)

  return { minTop, maxTop }
}

function getLeftBounds() {
  const minLeft = FAB_MARGIN
  const maxLeft =
    typeof window === "undefined" ? minLeft : Math.max(minLeft, window.innerWidth - FAB_SIZE - FAB_MARGIN)

  return { minLeft, maxLeft }
}

function createDockedFabPosition(side: DrawerSide, top: number): FabPosition {
  return { side, top, left: null }
}

function getDefaultFabPosition(): FabPosition {
  if (typeof window === "undefined") {
    return createDockedFabPosition("left", FAB_MARGIN * 4)
  }

  const { minTop, maxTop } = getTopBounds()
  const defaultTop = clamp(window.innerHeight - FAB_SIZE - FAB_MARGIN * 2, minTop, maxTop)
  return createDockedFabPosition("left", defaultTop)
}

function readStoredFabPosition(): FabPosition {
  if (typeof window === "undefined") {
    return getDefaultFabPosition()
  }

  const fallback = getDefaultFabPosition()
  const rawValue = window.localStorage.getItem(FAB_STORAGE_KEY)
  if (!rawValue) {
    return fallback
  }

  try {
    const parsed = JSON.parse(rawValue) as Partial<{ side: DrawerSide; top: number }>
    const side = parsed.side === "right" ? "right" : "left"
    const { minTop, maxTop } = getTopBounds()
    const top = typeof parsed.top === "number" ? clamp(parsed.top, minTop, maxTop) : fallback.top
    return createDockedFabPosition(side, top)
  } catch {
    return fallback
  }
}

function persistFabPosition(position: FabPosition) {
  if (typeof window === "undefined") {
    return
  }

  window.localStorage.setItem(
    FAB_STORAGE_KEY,
    JSON.stringify({
      side: position.side,
      top: position.top,
    }),
  )
}

export function AppShell() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, loading, logout } = useAuth()
  const [isNavOpen, setIsNavOpen] = useState(false)
  const [fabPosition, setFabPosition] = useState<FabPosition>(() => readStoredFabPosition())
  const [isDragging, setIsDragging] = useState(false)
  const fabPositionRef = useRef<FabPosition>(fabPosition)
  const suppressClickRef = useRef(false)
  const dragStateRef = useRef({
    active: false,
    pointerId: -1,
    startX: 0,
    startY: 0,
    startTop: 0,
    startLeft: 0,
    dragged: false,
  })

  const updateFabPosition = (updater: FabPosition | ((current: FabPosition) => FabPosition)) => {
    setFabPosition((current) => {
      const next =
        typeof updater === "function" ? (updater as (current: FabPosition) => FabPosition)(current) : updater
      fabPositionRef.current = next
      return next
    })
  }

  const handleLogout = async () => {
    await logout()
    navigate("/login", { replace: true })
  }

  useEffect(() => {
    const handleResize = () => {
      updateFabPosition((current) => {
        const { minTop, maxTop } = getTopBounds()
        const nextTop = clamp(current.top, minTop, maxTop)

        if (current.left === null) {
          if (nextTop !== current.top) {
            const next = createDockedFabPosition(current.side, nextTop)
            persistFabPosition(next)
            return next
          }

          return current
        }

        const { minLeft, maxLeft } = getLeftBounds()
        const nextLeft = clamp(current.left, minLeft, maxLeft)

        if (nextTop === current.top && nextLeft === current.left) {
          return current
        }

        return {
          ...current,
          top: nextTop,
          left: nextLeft,
        }
      })
    }

    window.addEventListener("resize", handleResize)
    return () => window.removeEventListener("resize", handleResize)
  }, [])

  useEffect(() => {
    const handlePointerMove = (event: PointerEvent) => {
      const dragState = dragStateRef.current
      if (!dragState.active || dragState.pointerId !== event.pointerId) {
        return
      }

      const deltaX = event.clientX - dragState.startX
      const deltaY = event.clientY - dragState.startY
      const movement = Math.hypot(deltaX, deltaY)
      if (!dragState.dragged && movement < DRAG_THRESHOLD) {
        return
      }

      if (!dragState.dragged) {
        dragState.dragged = true
        setIsDragging(true)
        setIsNavOpen(false)
      }

      const { minTop, maxTop } = getTopBounds()
      const { minLeft, maxLeft } = getLeftBounds()
      const nextTop = clamp(dragState.startTop + deltaY, minTop, maxTop)
      const nextLeft = clamp(dragState.startLeft + deltaX, minLeft, maxLeft)
      const nextSide: DrawerSide = nextLeft + FAB_SIZE / 2 < window.innerWidth / 2 ? "left" : "right"

      updateFabPosition({
        side: nextSide,
        top: nextTop,
        left: nextLeft,
      })
    }

    const handlePointerEnd = (event: PointerEvent) => {
      const dragState = dragStateRef.current
      if (!dragState.active || dragState.pointerId !== event.pointerId) {
        return
      }

      dragState.active = false
      dragState.pointerId = -1

      if (!dragState.dragged) {
        return
      }

      dragState.dragged = false
      setIsDragging(false)
      suppressClickRef.current = true
      window.setTimeout(() => {
        suppressClickRef.current = false
      }, 0)

      const current = fabPositionRef.current
      const resolvedLeft =
        current.left === null
          ? current.side === "left"
            ? FAB_MARGIN
            : Math.max(FAB_MARGIN, window.innerWidth - FAB_SIZE - FAB_MARGIN)
          : current.left
      const nextSide: DrawerSide = resolvedLeft + FAB_SIZE / 2 < window.innerWidth / 2 ? "left" : "right"
      const { minTop, maxTop } = getTopBounds()
      const dockedPosition = createDockedFabPosition(nextSide, clamp(current.top, minTop, maxTop))

      updateFabPosition(dockedPosition)
      persistFabPosition(dockedPosition)
    }

    window.addEventListener("pointermove", handlePointerMove)
    window.addEventListener("pointerup", handlePointerEnd)
    window.addEventListener("pointercancel", handlePointerEnd)

    return () => {
      window.removeEventListener("pointermove", handlePointerMove)
      window.removeEventListener("pointerup", handlePointerEnd)
      window.removeEventListener("pointercancel", handlePointerEnd)
    }
  }, [])

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
      onClick={() => setIsNavOpen(false)}
      to={item.to}
    >
      <item.icon className="h-4 w-4" />
      {item.label}
    </Link>
  )

  const handleFabPointerDown = (event: React.PointerEvent<HTMLButtonElement>) => {
    if (event.button !== 0) {
      return
    }

    const target = event.currentTarget
    const rect = target.getBoundingClientRect()
    const resolvedLeft = fabPositionRef.current.left ?? rect.left

    dragStateRef.current = {
      active: true,
      pointerId: event.pointerId,
      startX: event.clientX,
      startY: event.clientY,
      startTop: rect.top,
      startLeft: resolvedLeft,
      dragged: false,
    }

    target.setPointerCapture(event.pointerId)
  }

  const handleFabClick = () => {
    if (suppressClickRef.current) {
      return
    }

    setIsNavOpen((current) => !current)
  }

  const fabStyle =
    fabPosition.left === null
      ? fabPosition.side === "left"
        ? { top: `${fabPosition.top}px`, left: `${FAB_MARGIN}px` }
        : { top: `${fabPosition.top}px`, right: `${FAB_MARGIN}px` }
      : { top: `${fabPosition.top}px`, left: `${fabPosition.left}px` }

  return (
    <div className="min-h-screen bg-surface-page text-text-primary">
      <div
        className={cn(
          "fixed inset-0 z-40 bg-black/30 backdrop-blur-[2px] transition-all duration-500 ease-[cubic-bezier(0.22,1,0.36,1)]",
          isNavOpen ? "pointer-events-auto opacity-100" : "pointer-events-none opacity-0 backdrop-blur-0",
        )}
        onClick={() => setIsNavOpen(false)}
      />

      <aside
        className={cn(
          "fixed inset-y-4 z-50 flex w-[280px] max-w-[calc(100vw-2rem)] flex-col rounded-3xl border border-border bg-surface-inverse p-5 text-text-inverse shadow-soft transition-[transform,opacity,filter,scale] duration-500 ease-[cubic-bezier(0.22,1,0.36,1)] will-change-transform will-change-[opacity,filter]",
          fabPosition.side === "left" ? "left-4" : "right-4",
          isNavOpen
            ? "translate-x-0 scale-100 opacity-100 blur-0"
            : fabPosition.side === "left"
              ? "-translate-x-[calc(100%+1.5rem)] scale-[0.97] opacity-0 blur-[2px]"
              : "translate-x-[calc(100%+1.5rem)] scale-[0.97] opacity-0 blur-[2px]",
        )}
      >
        <div className="flex items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-accent/15 text-accent">
              <GraduationCap className="h-6 w-6" />
            </div>
            <div>
              <p className="text-lg font-semibold">dedao-dl Web</p>
            </div>
          </div>

          <button
            aria-label="关闭导航"
            className="flex h-10 w-10 items-center justify-center rounded-2xl text-text-inverse/70 transition hover:bg-white/10 hover:text-text-inverse"
            onClick={() => setIsNavOpen(false)}
            type="button"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="mt-8 flex-1 space-y-6 overflow-y-auto pr-1">
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

      <button
        aria-expanded={isNavOpen}
        aria-label={isNavOpen ? "关闭导航" : "打开导航"}
        className={cn(
          "group fixed z-[60] flex h-14 w-14 touch-none items-center justify-center overflow-hidden rounded-full border border-border ring-1 ring-black/5",
          isNavOpen
            ? "border-transparent bg-accent text-accent-foreground shadow-[0_18px_45px_rgba(251,146,60,0.32)] ring-transparent"
            : "bg-surface-panel text-text-primary shadow-soft",
          isDragging
            ? ""
            : "transition-[left,right,top,background-color,border-color,box-shadow,transform,color] duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] hover:-translate-y-0.5 hover:scale-[1.02] hover:bg-surface-soft",
        )}
        onClick={handleFabClick}
        onPointerDown={handleFabPointerDown}
        style={fabStyle}
        type="button"
      >
        <span
          aria-hidden="true"
          className={cn(
            "pointer-events-none absolute inset-0 rounded-full bg-accent/20 opacity-0",
            !isNavOpen && !isDragging ? "nav-fab-breathe" : "",
          )}
        />
        <span
          className={cn(
            "relative flex h-11 w-11 items-center justify-center rounded-full transition-transform duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
            isNavOpen
              ? "scale-100 bg-white/14"
              : "scale-[0.98] bg-surface-soft text-text-primary group-hover:scale-100 group-hover:bg-surface-page",
          )}
        >
          <Menu
            className={cn(
              "absolute h-5 w-5 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
              isNavOpen ? "scale-75 -rotate-90 opacity-0" : "scale-100 rotate-0 opacity-100",
            )}
          />
          <X
            className={cn(
              "absolute h-5 w-5 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
              isNavOpen ? "scale-100 rotate-0 opacity-100" : "scale-75 rotate-90 opacity-0",
            )}
          />
        </span>
      </button>

      <div className="mx-auto flex min-h-screen max-w-[1600px] flex-col gap-6 px-4 py-4 lg:px-6">
        <div className="flex min-h-0 flex-1 flex-col gap-6 pb-28">
          <header className="rounded-3xl border border-border bg-surface-panel px-5 py-4 shadow-soft backdrop-blur">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-center">
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-3 rounded-2xl border border-border bg-surface-soft px-4 py-3 text-sm text-text-muted xl:max-w-3xl">
                  <Search className="h-4 w-4" />
                  <span>全局搜索建设中</span>
                </div>
              </div>

              <div className="flex items-center justify-end gap-3">
                <ThemeToggleButton className="rounded-2xl" />

                <div className="flex items-center gap-3 rounded-2xl border border-border bg-surface-panel px-4 py-3 text-text-primary shadow-sm">
                  {loading ? (
                    <Loader2 className="h-4 w-4 animate-spin text-text-muted" />
                  ) : (
                    <img
                      alt={user?.nickname ?? "avatar"}
                      className="h-10 w-10 rounded-2xl object-cover"
                      src={user?.avatar || "https://placehold.co/80x80/e2e8f0/334155?text=DD"}
                    />
                  )}
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium">{user?.nickname ?? "未登录"}</p>
                    <p className="truncate text-xs text-text-muted">{user ? "当前账号" : "等待用户信息"}</p>
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
