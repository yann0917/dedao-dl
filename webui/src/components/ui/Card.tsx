import type { PropsWithChildren } from "react"
import { cn } from "@/lib/cn"

type CardProps = PropsWithChildren<{
  className?: string
}>

export function Card({ className, children }: CardProps) {
  return (
    <div
      className={cn(
        "border border-border-strong/60 bg-surface-panel text-text-primary",
        className,
      )}
    >
      {children}
    </div>
  )
}
