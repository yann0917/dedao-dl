import { ChevronDown } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { type HomeCategory } from "@/api"
import { cn } from "@/lib/cn"
import { semanticPageSectionClass } from "@/lib/semanticStyles"

type HomeCategoryMenuProps = {
  categories: HomeCategory[]
  onNavigateCategory: (category: HomeCategory, labelEnid: string) => void
}

// 面板内直接展示的分类数量，其余收进「全部分类」悬浮索引。
const VISIBLE_COUNT = 9
const CLOSE_DELAY_MS = 140

export function HomeCategoryMenu({
  categories,
  onNavigateCategory,
}: HomeCategoryMenuProps) {
  const [expanded, setExpanded] = useState(false)
  const closeTimerRef = useRef<number | null>(null)

  function clearCloseTimer() {
    if (closeTimerRef.current !== null) {
      window.clearTimeout(closeTimerRef.current)
      closeTimerRef.current = null
    }
  }

  function openPanel() {
    clearCloseTimer()
    setExpanded(true)
  }

  function scheduleClosePanel() {
    clearCloseTimer()
    closeTimerRef.current = window.setTimeout(() => {
      setExpanded(false)
      closeTimerRef.current = null
    }, CLOSE_DELAY_MS)
  }

  useEffect(() => {
    return () => {
      clearCloseTimer()
    }
  }, [])

  const visibleCategories = categories.slice(0, VISIBLE_COUNT)
  const hiddenCount = categories.length - visibleCategories.length

  return (
    <div className={cn(semanticPageSectionClass, "flex h-full flex-col py-2")}>
      <nav>
        {visibleCategories.map((category) => (
          <button
            className="flex w-full items-center justify-between gap-2 px-4 py-2 text-sm text-text-secondary transition hover:bg-surface-soft hover:text-text-primary"
            key={category.enid}
            onClick={() => onNavigateCategory(category, "")}
            type="button"
          >
            <span className="truncate">{category.name}</span>
            {category.labelList.length > 0 ? (
              <span className="font-mono text-[10px] tabular-nums text-text-muted">
                {category.labelList.length}
              </span>
            ) : null}
          </button>
        ))}
      </nav>

      {hiddenCount > 0 ? (
        <div
          className="relative mt-auto px-2 pb-1 pt-1"
          onPointerEnter={openPanel}
          onPointerLeave={scheduleClosePanel}
        >
          <button
            aria-expanded={expanded}
            className={cn(
              "flex w-full items-center justify-center gap-1.5 border border-dashed px-3 py-2 text-xs transition",
              expanded
                ? "border-accent/60 text-accent"
                : "border-border text-text-muted hover:border-accent/60 hover:text-accent",
            )}
            onClick={() => setExpanded((current) => !current)}
            onFocus={openPanel}
            type="button"
          >
            全部分类（{categories.length}）
            <ChevronDown
              className={cn("size-icon-sm transition-transform", expanded ? "rotate-180" : "")}
            />
          </button>

          {expanded ? (
            <div className="absolute left-2 top-full z-50 mt-2 w-[560px] max-w-[calc(100vw-2rem)] border border-border-strong/60 bg-surface-panel p-5 shadow-soft">
              <div className="max-h-[420px] space-y-4 overflow-y-auto">
                {categories.map((category) => (
                  <div key={category.enid}>
                    <button
                      className="text-sm font-medium text-text-primary transition hover:text-accent"
                      onClick={() => onNavigateCategory(category, "")}
                      type="button"
                    >
                      {category.name}
                    </button>
                    {category.labelList.length > 0 ? (
                      <div className="mt-1.5 flex flex-wrap gap-x-4 gap-y-1">
                        {category.labelList.map((label) => (
                          <button
                            className="text-xs text-text-muted transition hover:text-accent"
                            key={label.enid}
                            onClick={() => onNavigateCategory(category, label.enid)}
                            type="button"
                          >
                            {label.name}
                          </button>
                        ))}
                      </div>
                    ) : null}
                  </div>
                ))}
              </div>
            </div>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}
