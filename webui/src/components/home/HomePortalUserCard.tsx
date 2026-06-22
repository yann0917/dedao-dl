import { BookOpen, Clock3, UserCircle2 } from "lucide-react"
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { useAuth } from "@/providers/AuthProvider"

export function HomePortalUserCard() {
  const navigate = useNavigate()
  const { user } = useAuth()

  return (
    <Card className="flex h-full min-h-[360px] flex-col justify-between p-6">
      <div>
        <p className="text-sm text-slate-500">学习状态</p>
        <div className="mt-5 flex flex-col items-center text-center">
          <img
            alt={user?.nickname ?? "avatar"}
            className="h-20 w-20 rounded-3xl border-4 border-primary/10 object-cover"
            src={user?.avatar || "https://placehold.co/120x120/e2e8f0/334155?text=DD"}
          />
          <h3 className="mt-4 text-xl font-semibold text-slate-950">{user?.nickname ?? "得到用户"}</h3>
          <p className="mt-1 text-sm text-slate-500">{user?.uid_hazy ?? "当前账号"}</p>
        </div>
      </div>

      <div className="mt-6 grid gap-3">
        <div className="grid grid-cols-2 gap-3">
          <div className="rounded-2xl bg-slate-50 p-4">
            <div className="inline-flex items-center gap-2 text-sm text-slate-500">
              <Clock3 className="h-4 w-4" />
              今日学习
            </div>
            <p className="mt-3 text-2xl font-semibold text-slate-950">{Math.round((user?.today_study_time ?? 0) / 60)}</p>
            <p className="text-xs text-slate-500">分钟</p>
          </div>
          <div className="rounded-2xl bg-slate-50 p-4">
            <div className="inline-flex items-center gap-2 text-sm text-slate-500">
              <BookOpen className="h-4 w-4" />
              连续学习
            </div>
            <p className="mt-3 text-2xl font-semibold text-slate-950">{user?.study_serial_days ?? 0}</p>
            <p className="text-xs text-slate-500">天</p>
          </div>
        </div>

        <Button className="w-full" onClick={() => navigate("/user")}>
          <UserCircle2 className="mr-2 h-4 w-4" />
          进入用户中心
        </Button>
      </div>
    </Card>
  )
}
