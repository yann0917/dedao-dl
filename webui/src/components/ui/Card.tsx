import type { PropsWithChildren } from "react"
import { cn } from "@/lib/cn"

type CardProps = PropsWithChildren<{
  className?: string
}>

export function Card({ className, children }: CardProps) {
  return (
    <div
      className={cn(
        "rounded-3xl border border-border bg-surface-panel text-text-primary shadow-soft backdrop-blur",
        className,
      )}
    >
      {children}
    </div>
  )
}
