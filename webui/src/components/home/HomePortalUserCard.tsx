import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { useAuth } from "@/providers/AuthProvider"

export function HomePortalUserCard() {
  const navigate = useNavigate()
  const { user } = useAuth()

  return (
    <Card className="flex h-full flex-col p-5">
      <div className="flex min-w-0 items-center gap-4">
        <img
          alt={user?.nickname ?? "avatar"}
          className="h-14 w-14 shrink-0 rounded-md border border-border object-cover"
          src={user?.avatar || "https://placehold.co/112x112/e2e8f0/334155?text=DD"}
        />
        <div className="min-w-0">
          <p className="truncate text-base font-semibold text-text-primary">
            {user?.nickname ?? "得到用户"}
          </p>
          <p className="mt-1 font-mono text-[10px] uppercase tracking-[0.2em] text-text-muted">
            {user ? "当前账号" : "欢迎登录"}
          </p>
        </div>
      </div>

      <dl className="mt-5 grid grid-cols-2 divide-x divide-border border-y border-border">
        <div className="py-4 pr-4">
          <dt className="font-mono text-[10px] uppercase tracking-[0.24em] text-text-muted">今日学习</dt>
          <dd className="mt-1.5 text-2xl font-semibold tabular-nums text-text-primary">
            {Math.round((user?.today_study_time ?? 0) / 60)}
            <span className="ml-1 text-xs font-normal text-text-muted">分钟</span>
          </dd>
        </div>
        <div className="py-4 pl-4">
          <dt className="font-mono text-[10px] uppercase tracking-[0.24em] text-text-muted">连续学习</dt>
          <dd className="mt-1.5 text-2xl font-semibold tabular-nums text-text-primary">
            {user?.study_serial_days ?? 0}
            <span className="ml-1 text-xs font-normal text-text-muted">天</span>
          </dd>
        </div>
      </dl>

      <Button className="mt-auto w-full" onClick={() => navigate("/user")}>
        进入用户中心
      </Button>
    </Card>
  )
}
