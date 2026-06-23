import { LogOut } from "lucide-react"
import type { UserInfo } from "@/api"
import { Button } from "@/components/ui/Button"

type UserSummaryCardProps = {
  user: UserInfo | null
  onLogout: () => void
}

export function UserSummaryCard({ user, onLogout }: UserSummaryCardProps) {
  return (
    <div className="rounded-3xl bg-slate-950 p-5 text-slate-50">
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-sm text-slate-300">当前用户</p>
          <h1 className="mt-2 text-2xl font-semibold">{user?.nickname ?? "未命名用户"}</h1>
        </div>
        <Button onClick={onLogout} variant="ghost">
          <LogOut className="mr-2 h-4 w-4" />
          退出
        </Button>
      </div>

      <div className="mt-6 flex items-center gap-4">
        <img
          alt={user?.nickname ?? "avatar"}
          className="h-16 w-16 rounded-2xl object-cover"
          src={user?.avatar || "https://placehold.co/128x128/e2e8f0/334155?text=DD"}
        />
        <div>
          <p className="text-sm text-slate-300">{user ? "当前账号已登录" : "等待登录状态"}</p>
          <p className="mt-1 text-sm text-slate-400">查询能力来自现有 services，UI 只是新壳。</p>
        </div>
      </div>

      <div className="mt-6 grid grid-cols-2 gap-4 text-sm">
        <div className="rounded-2xl bg-slate-900 p-4">
          <p className="text-slate-400">今日学习</p>
          <p className="mt-2 text-2xl font-semibold">{Math.round((user?.today_study_time ?? 0) / 60)}</p>
          <p className="text-slate-400">分钟</p>
        </div>
        <div className="rounded-2xl bg-slate-900 p-4">
          <p className="text-slate-400">连续学习</p>
          <p className="mt-2 text-2xl font-semibold">{user?.study_serial_days ?? 0}</p>
          <p className="text-slate-400">天</p>
        </div>
      </div>
    </div>
  )
}
