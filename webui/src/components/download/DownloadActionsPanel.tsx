import { useState } from "react"
import { toast } from "sonner"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import type { DownloadSessionResponse } from "@/api"
import { useDownloadProgress } from "@/providers/DownloadProgressProvider"

export type DownloadOption = {
  value: number
  label: string
}

type DownloadActionsPanelProps = {
  title: string
  description?: string
  options: DownloadOption[]
  disabled?: boolean
  disabledReason?: string
  onDownload: (downloadType: number) => Promise<DownloadSessionResponse>
}

export function DownloadActionsPanel({
  title,
  description = "下载过程在服务端执行，时间较长时请耐心等待。",
  options,
  disabled = false,
  disabledReason,
  onDownload,
}: DownloadActionsPanelProps) {
  const { beginSession } = useDownloadProgress()
  const [pendingType, setPendingType] = useState<number | null>(null)

  const handleDownload = async (downloadType: number) => {
    setPendingType(downloadType)

    try {
      const result = await onDownload(downloadType)
      beginSession(result)
      toast.success("下载已开始", {
        description: `正在连接进度流。输出目录：${result.outputDir}`,
      })
    } catch (error) {
      toast.error("下载失败", {
        description: error instanceof Error ? error.message : "请稍后重试",
      })
    } finally {
      setPendingType(null)
    }
  }

  return (
    <Card className="p-6">
      <h3 className="text-xl font-semibold text-text-primary">{title}</h3>
      <p className="mt-3 text-sm leading-7 text-text-secondary">{description}</p>
      <div className="mt-5 flex flex-wrap gap-3">
        {options.map((option) => (
          <Button
            disabled={disabled || pendingType !== null}
            key={option.value}
            onClick={() => void handleDownload(option.value)}
            variant={option.value === 1 ? "default" : "outline"}
          >
            {pendingType === option.value ? "下载中..." : option.label}
          </Button>
        ))}
      </div>
      {disabled && disabledReason ? (
        <div className="mt-5 rounded-2xl border border-warning bg-warning-soft p-4 text-sm text-warning">
          {disabledReason}
        </div>
      ) : null}
    </Card>
  )
}
