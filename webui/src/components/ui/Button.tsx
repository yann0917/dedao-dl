import type { PropsWithChildren } from "react"
import { cn } from "@/lib/cn"

type ButtonProps = PropsWithChildren<{
  className?: string
  variant?: "default" | "outline" | "ghost"
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
      ? "border border-slate-200 bg-white text-slate-700 hover:bg-slate-50"
      : variant === "ghost"
        ? "bg-transparent text-slate-600 hover:bg-slate-100"
        : "bg-primary text-white hover:bg-primary/90"

  return (
    <button
      className={cn(
        "inline-flex h-10 items-center justify-center rounded-xl px-4 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-50",
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
