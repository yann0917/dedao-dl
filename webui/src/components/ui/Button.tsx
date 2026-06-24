import type { PropsWithChildren } from "react"
import { cn } from "@/lib/cn"

type ButtonProps = PropsWithChildren<{
  className?: string
  variant?: "default" | "outline" | "ghost" | "danger"
  onClick?: () => void
  disabled?: boolean
  type?: "button" | "submit"
}>

export function Button({
  className,
  variant = "default",
  onClick,
  disabled,
  type = "button",
  children,
}: ButtonProps) {
  const variantClass =
    variant === "outline"
      ? "border border-border bg-surface-panel text-text-secondary hover:bg-surface-soft"
      : variant === "danger"
        ? "bg-danger text-danger-foreground hover:bg-danger/90"
      : variant === "ghost"
        ? "bg-transparent text-text-secondary hover:bg-surface-soft"
        : "bg-accent text-accent-foreground hover:bg-accent/90"

  return (
    <button
      className={cn(
        "inline-flex h-10 items-center justify-center rounded-xl px-4 text-sm font-medium transition focus:outline-none focus:ring-2 focus:ring-ring/30 disabled:cursor-not-allowed disabled:opacity-50",
        variantClass,
        className,
      )}
      disabled={disabled}
      onClick={onClick}
      type={type}
    >
      {children}
    </button>
  )
}
