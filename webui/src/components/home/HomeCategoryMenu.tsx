import { useEffect, useRef, useState } from "react"
import { type HomeCategory } from "@/api"
import { Card } from "@/components/ui/Card"
import {
  Menubar,
  MenubarMenu,
  MenubarTrigger,
} from "@/components/ui/Menubar"

type HomeCategoryMenuProps = {
  categories: HomeCategory[]
  onNavigateCategory: (category: HomeCategory, labelEnid: string) => void
}

export function HomeCategoryMenu({
  categories,
  onNavigateCategory,
}: HomeCategoryMenuProps) {
  const [hoveredCategoryEnid, setHoveredCategoryEnid] = useState<string | null>(null)
  const closeTimerRef = useRef<number | null>(null)

  function clearCloseTimer() {
    if (closeTimerRef.current !== null) {
      window.clearTimeout(closeTimerRef.current)
      closeTimerRef.current = null
    }
  }

  function openCategoryMenu(categoryEnid: string) {
    clearCloseTimer()
    setHoveredCategoryEnid(categoryEnid)
  }

  function scheduleCloseCategoryMenu(categoryEnid: string) {
    clearCloseTimer()
    closeTimerRef.current = window.setTimeout(() => {
      setHoveredCategoryEnid((current) => (current === categoryEnid ? null : current))
      closeTimerRef.current = null
    }, 120)
  }

  useEffect(() => {
    return () => {
      clearCloseTimer()
    }
  }, [])

  return (
    <Card className="h-full p-4">
      <div>
        <p className="text-sm font-medium text-slate-900">内容分类</p>
        <p className="mt-1 text-xs leading-5 text-slate-500">点击一级分类可直接进入结果页，悬浮后可继续选择二级标签。</p>
      </div>

      <div className="mt-4">
        <Menubar className="w-full justify-start bg-transparent p-0">
          {categories.map((category) => {
            const hasLabels = category.labelList.length > 0
            const isHovered = hoveredCategoryEnid === category.enid

            return (
              <MenubarMenu key={category.enid}>
                {/* Radix menubar 的顶层菜单不提供这里需要的受控开合，悬浮层改用局部状态管理。 */}
                <div
                  className="relative"
                  onBlurCapture={(event) => {
                    if (!event.currentTarget.contains(event.relatedTarget as Node | null)) {
                      scheduleCloseCategoryMenu(category.enid)
                    }
                  }}
                  onFocusCapture={() => {
                    if (hasLabels) {
                      openCategoryMenu(category.enid)
                    }
                  }}
                  onPointerEnter={() => {
                    if (hasLabels) {
                      openCategoryMenu(category.enid)
                    }
                  }}
                  onPointerLeave={() => {
                    if (hasLabels) {
                      scheduleCloseCategoryMenu(category.enid)
                    }
                  }}
                >
                  <MenubarTrigger onClick={() => onNavigateCategory(category, "")}>
                    {category.name}
                  </MenubarTrigger>

                  {hasLabels && isHovered ? (
                    <div className="absolute left-0 top-full z-50 min-w-[14rem] pt-2">
                      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white p-2 text-slate-950 shadow-xl">
                      <p className="px-3 py-2 text-xs font-semibold uppercase tracking-wide text-slate-400">
                        选择标签
                      </p>
                      <div className="space-y-1">
                        {category.labelList.map((label) => (
                          <button
                            key={label.enid}
                            className="flex w-full items-center rounded-xl px-3 py-2 text-left text-sm text-slate-700 transition hover:bg-slate-100 focus:bg-slate-100 focus:outline-none"
                            onClick={() => onNavigateCategory(category, label.enid)}
                            type="button"
                          >
                            <span className="truncate">{label.name}</span>
                          </button>
                        ))}
                      </div>
                      </div>
                    </div>
                  ) : null}
                </div>
              </MenubarMenu>
            )
          })}
        </Menubar>
      </div>
    </Card>
  )
}
