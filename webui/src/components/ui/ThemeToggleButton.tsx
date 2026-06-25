import { Moon, Sun } from "lucide-react"
import { cn } from "@/lib/cn"
import { useTheme } from "@/providers/ThemeProvider"

type ThemeToggleButtonProps = {
  className?: string
}

export function ThemeToggleButton({ className }: ThemeToggleButtonProps) {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === "dark"
  const label = isDark ? "切换到亮色模式" : "切换到暗色模式"

  return (
    <button
      aria-label={label}
      className={cn(
        "group relative inline-flex h-12 w-12 items-center justify-center overflow-hidden rounded-2xl border border-border bg-surface-panel text-text-primary shadow-sm transition-[background-color,border-color,box-shadow,transform,color] duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] hover:-translate-y-0.5 hover:bg-surface-soft focus:outline-none focus:ring-2 focus:ring-ring/30",
        className,
      )}
      onClick={toggleTheme}
      title={label}
      type="button"
    >
      <span
        className={cn(
          "relative flex h-9 w-9 items-center justify-center rounded-xl bg-surface-soft transition-[background-color,transform] duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] group-hover:scale-100 group-hover:bg-surface-page",
          isDark ? "scale-100" : "scale-[0.98]",
        )}
      >
        <Sun
          className={cn(
            "absolute h-4 w-4 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
            isDark ? "scale-100 rotate-0 opacity-100 text-warning" : "scale-75 rotate-90 opacity-0 text-text-muted",
          )}
        />
        <Moon
          className={cn(
            "absolute h-4 w-4 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
            isDark ? "scale-75 -rotate-90 opacity-0 text-text-muted" : "scale-100 rotate-0 opacity-100 text-text-secondary",
          )}
        />
      </span>
      <span
        aria-hidden="true"
        className={cn(
          "absolute right-2 top-2 h-2.5 w-2.5 rounded-full ring-2 ring-surface-panel transition-[background-color,box-shadow] duration-300",
          isDark
            ? "bg-sky-400 shadow-[0_0_0_2px_rgba(56,189,248,0.18)]"
            : "bg-amber-500 shadow-[0_0_0_2px_rgba(245,158,11,0.18)]",
        )}
      />
    </button>
  )
}
