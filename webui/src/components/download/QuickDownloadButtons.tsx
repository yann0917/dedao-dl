import { ChevronDown, Loader2 } from "lucide-react"
import { useState } from "react"
import type { DownloadSessionResponse } from "@/api"
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
} from "@/components/ui/Menubar"
import { useDownloadProgress } from "@/providers/DownloadProgressProvider"

export type QuickDownloadOption = {
  value: number
  label: string
}

type QuickDownloadButtonsProps = {
  options: QuickDownloadOption[]
  onDownload: (downloadType: number) => Promise<DownloadSessionResponse>
  disabled?: boolean
}

export function QuickDownloadButtons({ options, onDownload, disabled = false }: QuickDownloadButtonsProps) {
  const { beginSession } = useDownloadProgress()
  const [pendingType, setPendingType] = useState<number | null>(null)

  const handleDownload = async (downloadType: number) => {
    setPendingType(downloadType)
    try {
      const result = await onDownload(downloadType)
      beginSession(result)
    } finally {
      setPendingType(null)
    }
  }

  return (
    <Menubar className="border-0 bg-transparent p-0">
      <MenubarMenu>
        <MenubarTrigger className="h-9 gap-2 px-3" disabled={disabled || pendingType !== null}>
          {pendingType !== null ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
          <span>{pendingType !== null ? "下载中..." : "下载"}</span>
          <ChevronDown className="h-4 w-4" />
        </MenubarTrigger>
        <MenubarContent align="end" className="min-w-[10rem]">
          {options.map((option) => (
            <MenubarItem key={option.value} onClick={() => void handleDownload(option.value)}>
              {option.label}
            </MenubarItem>
          ))}
        </MenubarContent>
      </MenubarMenu>
    </Menubar>
  )
}
