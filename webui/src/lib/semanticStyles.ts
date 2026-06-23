import { cn } from "@/lib/cn"

export const semanticPageSectionClass =
  "rounded-3xl border border-border bg-surface-panel shadow-soft backdrop-blur"

export const semanticSubtlePanelClass = "rounded-2xl bg-surface-soft"

export const semanticMetaTextClass = "text-sm text-text-muted"

export const semanticSecondaryTextClass = "text-text-secondary"

export const semanticInfoBlockClass = `${semanticSubtlePanelClass} p-4 text-sm text-text-secondary`

export const semanticStatCardClass = `${semanticSubtlePanelClass} p-4`

export const semanticStatusBadgeBaseClass = "rounded-full px-3 py-1 text-xs font-medium"

const semanticStatusBadgeVariantClass = {
  accent: "bg-accent-soft text-accent",
  neutral: "bg-surface-soft text-text-secondary",
  success: "bg-success-soft text-success",
  warning: "bg-warning-soft text-warning",
  danger: "bg-danger-soft text-danger",
} as const

export type SemanticStatusVariant = keyof typeof semanticStatusBadgeVariantClass

export function getSemanticStatusBadgeClass(
  variant: SemanticStatusVariant = "neutral",
  className?: string,
) {
  return cn(semanticStatusBadgeBaseClass, semanticStatusBadgeVariantClass[variant], className)
}

export function getSemanticChipClass(selected: boolean, emphasis: "primary" | "strong" = "primary") {
  if (selected) {
    return emphasis === "strong"
      ? "rounded-full bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground transition"
      : "rounded-full bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground transition"
  }

  return "rounded-full bg-surface-soft px-3 py-1.5 text-sm font-medium text-text-secondary transition hover:bg-accent-soft hover:text-accent"
}
