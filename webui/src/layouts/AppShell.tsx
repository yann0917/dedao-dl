import { BookMarked, Compass, GraduationCap, Headphones, Home, Loader2, Menu, PanelLeftClose, PanelLeftOpen, Rows3, Search, Sparkles, Trophy, UserCircle2, X } from "lucide-react"
import type { ComponentType, ReactNode } from "react"
import { useEffect, useRef, useState } from "react"
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { ThemeToggleButton } from "@/components/ui/ThemeToggleButton"
import { cn } from "@/lib/cn"
import { useAuth } from "@/providers/AuthProvider"

const primaryNavItems = [
  { to: "/", label: "首页", icon: Home, end: true },
  { to: "/leaderboard", label: "得到榜单", icon: Trophy },
  { to: "/ai-channel", label: "AI 学习圈", icon: Sparkles },
]

const purchasedNavItems = [
  { to: "/purchased/manage", label: "已购管理", icon: Rows3 },
  { to: "/purchased/courses", label: "课程", icon: GraduationCap },
  { to: "/purchased/ebooks", label: "电子书", icon: BookMarked },
  { to: "/purchased/audios", label: "听书", icon: Headphones },
  { to: "/purchased/compass", label: "锦囊", icon: Compass },
]

const accountNavItems = [
  { to: "/user", label: "用户中心", icon: UserCircle2 },
]

type NavItem = { to: string; label: string; icon: ComponentType<{ className?: string }>; end?: boolean }

type DrawerSide = "left" | "right"

type FabPosition = {
  side: DrawerSide
  top: number
  left: number | null
}

const FAB_STORAGE_KEY = "dedao-dl-nav-fab-position"
const SIDEBAR_COLLAPSED_KEY = "dedao-dl-sidebar-collapsed"
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

function BrandMark({ collapsed = false }: { collapsed?: boolean }) {
  if (collapsed) {
    return (
      <span
        aria-hidden="true"
        className="flex h-9 w-9 -rotate-6 items-center justify-center bg-accent text-lg font-bold text-accent-foreground shadow-[inset_0_0_0_1.5px_hsl(0_0%_100%/0.35)]"
      >
        得
      </span>
    )
  }

  return (
    <div className="flex items-center gap-3">
      <span
        aria-hidden="true"
        className="flex h-9 w-9 -rotate-6 items-center justify-center bg-accent text-lg font-bold text-accent-foreground shadow-[inset_0_0_0_1.5px_hsl(0_0%_100%/0.35)]"
      >
        得
      </span>
      <div className="leading-tight">
        <p className="text-base font-semibold tracking-wide text-text-primary">dedao-dl</p>
        <p className="mt-0.5 font-mono text-[9px] uppercase tracking-[0.3em] text-text-muted">
          cli tool · web ui
        </p>
      </div>
    </div>
  )
}

function NavGroupLabel({ children }: { children: ReactNode }) {
  return (
    <p className="px-3 pb-1 pt-5 font-mono text-[10px] uppercase tracking-[0.28em] text-text-muted">
      {children}
    </p>
  )
}

export function AppShell() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, loading, logout } = useAuth()
  const [isNavOpen, setIsNavOpen] = useState(false)
  const [fabPosition, setFabPosition] = useState<FabPosition>(() => readStoredFabPosition())
  const [isDragging, setIsDragging] = useState(false)
  const [sidebarCollapsed, setSidebarCollapsed] = useState<boolean>(() => {
    if (typeof window === "undefined") {
      return false
    }

    return window.localStorage.getItem(SIDEBAR_COLLAPSED_KEY) === "1"
  })
  const fabPositionRef = useRef<FabPosition>(fabPosition)
  const suppressClickRef = useRef(false)
  const drawerCloseRef = useRef<HTMLButtonElement | null>(null)
  const fabRef = useRef<HTMLButtonElement | null>(null)
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

  const toggleSidebarCollapsed = () => {
    setSidebarCollapsed((current) => {
      const next = !current
      window.localStorage.setItem(SIDEBAR_COLLAPSED_KEY, next ? "1" : "0")
      return next
    })
  }

  // 抽屉打开时：焦点移入抽屉、Esc 关闭并把焦点还给悬浮球。
  useEffect(() => {
    if (!isNavOpen) {
      return
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsNavOpen(false)
        fabRef.current?.focus()
      }
    }

    window.addEventListener("keydown", handleKeyDown)
    drawerCloseRef.current?.focus()
    return () => window.removeEventListener("keydown", handleKeyDown)
  }, [isNavOpen])

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

  const renderNavItem = (item: NavItem, opts?: { collapsed?: boolean; onNavigate?: () => void }) => {
    const active = isActiveNav(item.to, item.end)
    return (
      <Link
        aria-current={active ? "page" : undefined}
        className={cn(
          "flex items-center rounded-md text-sm transition",
          opts?.collapsed ? "h-10 justify-center px-0" : "gap-3 px-3 py-2.5",
          active
            ? "bg-accent-soft font-medium text-accent"
            : "text-text-secondary hover:bg-surface-soft hover:text-text-primary",
        )}
        key={item.to}
        onClick={() => {
          opts?.onNavigate?.()
        }}
        title={opts?.collapsed ? item.label : undefined}
        to={item.to}
      >
        <item.icon className="size-icon-md shrink-0" />
        {opts?.collapsed ? null : <span className="truncate">{item.label}</span>}
      </Link>
    )
  }

  const renderNav = (opts?: { collapsed?: boolean; onNavigate?: () => void }) => (
    <nav className="flex-1 overflow-y-auto">
      <div className="space-y-1">
        {primaryNavItems.map((item) => renderNavItem(item, opts))}
      </div>
      {opts?.collapsed ? (
        <div aria-hidden="true" className="mx-auto my-3 h-px w-6 bg-border" />
      ) : (
        <NavGroupLabel>已购</NavGroupLabel>
      )}
      <div className="space-y-1">
        {purchasedNavItems.map((item) => renderNavItem(item, opts))}
      </div>
      {opts?.collapsed ? (
        <div aria-hidden="true" className="mx-auto my-3 h-px w-6 bg-border" />
      ) : (
        <NavGroupLabel>账号</NavGroupLabel>
      )}
      <div className="space-y-1">{accountNavItems.map((item) => renderNavItem(item, opts))}</div>
    </nav>
  )

  const userBlock = (
    <div className="flex items-center gap-3">
      {loading ? (
        <div className="flex h-9 w-9 items-center justify-center border border-border">
          <Loader2 className="size-icon-md animate-spin text-text-muted" />
        </div>
      ) : (
        <img
          alt={user?.nickname ?? "avatar"}
          className="h-9 w-9 rounded-md border border-border object-cover"
          src={user?.avatar || "https://placehold.co/72x72/e2e8f0/334155?text=DD"}
        />
      )}
      <div className="min-w-0">
        <p className="truncate text-sm font-medium text-text-primary">{user?.nickname ?? "未登录"}</p>
        <p className="truncate font-mono text-[10px] uppercase tracking-[0.2em] text-text-muted">
          {user ? "当前账号" : "等待用户信息"}
        </p>
      </div>
    </div>
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
          "fixed inset-0 z-40 bg-black/30 transition-opacity duration-300 xl:hidden",
          isNavOpen ? "pointer-events-auto opacity-100" : "pointer-events-none opacity-0",
        )}
        onClick={() => setIsNavOpen(false)}
      />

      {/* 桌面端固定侧边栏 */}
      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-40 hidden flex-col border-r border-border bg-surface-page py-5 transition-[width] duration-200 xl:flex",
          sidebarCollapsed ? "w-16 px-2" : "w-60 px-3",
        )}
      >
        <div className={cn("flex", sidebarCollapsed ? "justify-center" : "px-2")}>
          <BrandMark collapsed={sidebarCollapsed} />
        </div>
        <div className="mt-6 flex min-h-0 flex-1 flex-col">
          {renderNav({ collapsed: sidebarCollapsed })}
        </div>
        <div className="border-t border-border pt-3">
          <button
            className={cn(
              "flex w-full items-center justify-center gap-2 rounded-md py-2 text-sm text-text-muted transition hover:bg-surface-soft hover:text-text-primary",
              sidebarCollapsed ? "px-0" : "px-3",
            )}
            onClick={toggleSidebarCollapsed}
            title={sidebarCollapsed ? "展开侧栏" : "收起侧栏"}
            type="button"
          >
            {sidebarCollapsed ? (
              <PanelLeftOpen className="size-icon-md" />
            ) : (
              <>
                <PanelLeftClose className="size-icon-md" />
                <span>收起侧栏</span>
              </>
            )}
          </button>
        </div>
      </aside>

      {/* 移动端抽屉 */}
      <aside
        className={cn(
          "fixed inset-y-3 z-50 flex w-[280px] max-w-[calc(100vw-2rem)] flex-col border border-border-strong/60 bg-surface-page px-3 py-5 shadow-soft transition-transform duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] will-change-transform xl:hidden",
          fabPosition.side === "left" ? "left-3" : "right-3",
          isNavOpen
            ? "translate-x-0"
            : fabPosition.side === "left"
              ? "-translate-x-[calc(100%+1.5rem)]"
              : "translate-x-[calc(100%+1.5rem)]",
        )}
      >
        <div className="flex items-center justify-between gap-3 px-1">
          <BrandMark />
          <button
            aria-label="关闭导航"
            className="flex h-9 w-9 items-center justify-center text-text-muted transition hover:bg-surface-soft hover:text-text-primary"
            onClick={() => setIsNavOpen(false)}
            ref={drawerCloseRef}
            type="button"
          >
            <X className="size-icon-lg" />
          </button>
        </div>
        <div className="mt-6 flex min-h-0 flex-1 flex-col">
          {renderNav({ onNavigate: () => setIsNavOpen(false) })}
        </div>
        <div className="border-t border-border pt-4">
          {userBlock}
          <div className="mt-3 flex items-center justify-between">
            <ThemeToggleButton className="h-9 w-9" />
            <Button className="h-9 px-3 text-xs" onClick={() => void handleLogout()} variant="outline">
              退出
            </Button>
          </div>
        </div>
      </aside>

      {/* 移动端导航悬浮球 */}
      <button
        aria-expanded={isNavOpen}
        aria-label={isNavOpen ? "关闭导航" : "打开导航"}
        className={cn(
          "fixed z-[60] hidden h-14 w-14 touch-none items-center justify-center rounded-full border sm:flex xl:hidden",
          isNavOpen
            ? "border-transparent bg-accent text-accent-foreground"
            : "border-border-strong/60 bg-surface-panel text-text-primary shadow-soft",
          isDragging ? "" : "transition-transform duration-300",
        )}
        onClick={handleFabClick}
        onPointerDown={handleFabPointerDown}
        ref={fabRef}
        style={fabStyle}
        type="button"
      >
        <Menu
          className={cn(
            "absolute size-icon-lg transition-all duration-300",
            isNavOpen ? "-rotate-90 scale-75 opacity-0" : "rotate-0 scale-100 opacity-100",
          )}
        />
        <X
          className={cn(
            "absolute size-icon-lg transition-all duration-300",
            isNavOpen ? "rotate-0 scale-100 opacity-100" : "rotate-90 scale-75 opacity-0",
          )}
        />
      </button>

      <div
        className={cn(
          "min-h-screen transition-[padding] duration-200",
          sidebarCollapsed ? "xl:pl-16" : "xl:pl-60",
        )}
      >
        <header className="sticky top-0 z-30 border-b border-border bg-surface-page/90 backdrop-blur">
          <div className="mx-auto flex max-w-[1500px] items-center gap-4 px-4 py-3 lg:px-6">
            <div className="xl:hidden">
              <BrandMark />
            </div>

            <div className="flex min-w-0 flex-1 items-center gap-2.5 border-b border-border pb-1.5 text-sm text-text-muted xl:max-w-md">
              <Search className="size-icon-md shrink-0" />
              <span className="truncate">全局搜索建设中</span>
            </div>

            <div className="ml-auto flex shrink-0 items-center gap-3">
              <div className="hidden sm:block">{userBlock}</div>
              <ThemeToggleButton className="rounded-none" />
              <Button className="h-9 px-3 text-xs" onClick={() => void handleLogout()} variant="outline">
                退出
              </Button>
            </div>
          </div>
        </header>

        <main className="mx-auto max-w-[1500px] px-4 pb-28 pt-6 lg:px-6 xl:pb-10">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
