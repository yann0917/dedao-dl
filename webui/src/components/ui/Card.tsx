import type { PropsWithChildren } from "react"
import { cn } from "@/lib/cn"

type CardProps = PropsWithChildren<{
  className?: string
}>

export function Card({ className, children }: CardProps) {
  return <div className={cn("rounded-3xl border border-white/70 bg-white/90 shadow-soft backdrop-blur", className)}>{children}</div>
}
