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
        "relative inline-flex h-10 w-10 items-center justify-center border border-border bg-transparent text-text-secondary transition-[background-color,border-color,color] duration-300 hover:bg-surface-soft hover:text-text-primary focus:outline-none focus:ring-2 focus:ring-ring/30 rounded-none",
        className,
      )}
      onClick={toggleTheme}
      title={label}
      type="button"
    >
      <Sun
        className={cn(
          "absolute size-icon-md transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
          isDark ? "rotate-0 scale-100 opacity-100" : "rotate-90 scale-75 opacity-0",
        )}
      />
      <Moon
        className={cn(
          "absolute size-icon-md transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]",
          isDark ? "-rotate-90 scale-75 opacity-0" : "rotate-0 scale-100 opacity-100",
        )}
      />
    </button>
  )
}
