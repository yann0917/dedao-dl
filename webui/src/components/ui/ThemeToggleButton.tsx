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
        "relative inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-border bg-surface-panel text-text-secondary transition hover:bg-surface-soft focus:outline-none focus:ring-2 focus:ring-ring/30",
        className,
      )}
      onClick={toggleTheme}
      title={label}
      type="button"
    >
      {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
      <span
        aria-hidden="true"
        className={cn(
          "absolute right-2 top-2 h-2.5 w-2.5 rounded-full ring-2 ring-surface-panel",
          isDark ? "bg-warning" : "bg-accent",
        )}
      />
    </button>
  )
}
