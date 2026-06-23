import type { PropsWithChildren, ReactNode } from "react"
import { cn } from "@/lib/cn"
import {
  type SemanticStatusVariant,
  getSemanticStatusBadgeClass,
  semanticInfoBlockClass,
  semanticMetaTextClass,
  semanticStatCardClass,
} from "@/lib/semanticStyles"

type StatusBadgeProps = PropsWithChildren<{
  variant?: SemanticStatusVariant
  className?: string
  icon?: ReactNode
}>

export function StatusBadge({ variant = "neutral", className, icon, children }: StatusBadgeProps) {
  return (
    <span className={cn(getSemanticStatusBadgeClass(variant), icon ? "inline-flex items-center gap-2" : "", className)}>
      {icon}
      {children}
    </span>
  )
}

type InfoBlockProps = PropsWithChildren<{
  className?: string
}>

export function InfoBlock({ className, children }: InfoBlockProps) {
  return <div className={cn(semanticInfoBlockClass, className)}>{children}</div>
}

type StatCardProps = {
  className?: string
  icon?: ReactNode
  label: ReactNode
  value: ReactNode
  hint?: ReactNode
  valueClassName?: string
}

export function StatCard({ className, icon, label, value, hint, valueClassName }: StatCardProps) {
  return (
    <div className={cn(semanticStatCardClass, className)}>
      <div className={cn("inline-flex items-center gap-2", semanticMetaTextClass)}>
        {icon}
        {label}
      </div>
      <p className={cn("mt-3 text-2xl font-semibold text-text-primary", valueClassName)}>{value}</p>
      {hint ? <p className={cn("mt-1", semanticMetaTextClass)}>{hint}</p> : null}
    </div>
  )
}
