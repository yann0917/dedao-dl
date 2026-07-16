import { BookOpen, Clock3, Loader2, LogOut, Radio, RefreshCcw, School, SwitchCamera } from "lucide-react"
import { InfoBlock, StatCard, StatusBadge } from "@/components/ui/Semantic"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import {
  semanticMetaTextClass,
  semanticSecondaryTextClass,
} from "@/lib/semanticStyles"
import { useUserCenter } from "@/hooks/useUserCenter"

function formatTimestamp(timestamp?: number) {
  if (!timestamp) {
    return "未获取"
  }

  return new Date(timestamp * 1000).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  })
}

export function UserCenterPage() {
  const { data, loading, switchingUID, error, reload, switchAccount, logout } = useUserCenter()

  if (loading) {
    return (
      <main className="flex min-h-[60vh] items-center justify-center">
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
          正在加载用户中心...
        </div>
      </main>
    )
  }

  const user = data?.user
  const ebookVip = data?.ebookVip
  const odobVip = data?.odobVip?.user

  return (
    <main className="space-y-6">

      {error ? (
        <Card className="border-danger bg-danger-soft">
          <div className="p-4 text-sm text-danger">{error}</div>
        </Card>
      ) : null}

      <section className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <Card className="p-6">
          <div className="flex flex-col gap-6 md:flex-row md:items-start">
            <div className="relative">
              <img
                alt={user?.nickname ?? "avatar"}
                className="h-24 w-24 rounded-3xl border-4 border-primary/10 object-cover"
                src={user?.avatar || "https://placehold.co/120x120/e2e8f0/334155?text=DD"}
              />
              {user?.is_teacher ? (
                <div className="absolute -right-2 -top-2 inline-flex items-center gap-1 rounded-full bg-primary px-3 py-1 text-xs font-medium text-white">
                  <School className="h-3.5 w-3.5" />
                  教师
                </div>
              ) : null}
            </div>

            <div className="flex-1 space-y-4">
              <div>
                <h3 className="text-3xl font-semibold text-text-primary">{user?.nickname ?? "未命名用户"}</h3>
                <p className={`mt-2 ${semanticMetaTextClass}`}>当前登录账号</p>
              </div>

              <div className="flex flex-wrap gap-3">
                {odobVip?.is_vip ? (
                  <div className="inline-flex items-center gap-3 rounded-2xl bg-warning-soft px-4 py-3 text-sm text-warning">
                    <Radio className="h-4 w-4" />
                    <div>
                      <div className="font-medium">听书会员</div>
                      <div className="text-xs text-warning/80">剩余 {odobVip.surplus_time ?? 0} 天</div>
                    </div>
                  </div>
                ) : null}
                {ebookVip?.is_vip ? (
                  <div className="inline-flex items-center gap-3 rounded-2xl bg-accent-soft px-4 py-3 text-sm text-accent">
                    <BookOpen className="h-4 w-4" />
                    <div>
                      <div className="font-medium">电子书会员</div>
                      <div className="text-xs text-accent/80">剩余 {ebookVip.surplus_time ?? 0} 天</div>
                    </div>
                  </div>
                ) : null}
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <StatCard hint="分钟" label="今日学习" value={Math.round((user?.today_study_time ?? 0) / 60)} />
                <StatCard hint="天" label="连续学习" value={user?.study_serial_days ?? 0} />
              </div>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-xl font-semibold text-text-primary">账号管理</h3>
              <p className={`mt-1 text-sm ${semanticSecondaryTextClass}`}>支持查看当前登录账号并切换活跃用户。</p>
            </div>
            <Button onClick={() => void reload()} variant="outline">
              <RefreshCcw className="mr-2 h-4 w-4" />
              刷新
            </Button>
          </div>

          <div className="mt-5 space-y-3">
            {data?.accounts.map((account) => (
              <div
                className="flex items-center justify-between rounded-2xl border border-border px-4 py-3"
                key={account.uidHazy}
              >
                <div className="flex min-w-0 items-center gap-3">
                  <img
                    alt={account.name}
                    className="h-12 w-12 rounded-2xl object-cover"
                    src={account.avatar || "https://placehold.co/80x80/e2e8f0/334155?text=DD"}
                  />
                  <div className="min-w-0">
                    <p className="truncate font-medium text-text-primary">{account.name}</p>
                    <p className={`truncate text-sm ${semanticMetaTextClass}`}>{account.active ? "当前使用中" : "可切换账号"}</p>
                  </div>
                </div>

                {account.active ? (
                  <StatusBadge variant="accent">当前账号</StatusBadge>
                ) : (
                  <Button
                    disabled={switchingUID === account.uidHazy}
                    onClick={() => void switchAccount(account.uidHazy)}
                    variant="outline"
                  >
                    <SwitchCamera className="mr-2 h-4 w-4" />
                    {switchingUID === account.uidHazy ? "切换中..." : "切换"}
                  </Button>
                )}
              </div>
            ))}
          </div>

          <div className="mt-5">
            <Button className="w-full" onClick={() => void logout()} variant="ghost">
              <LogOut className="mr-2 h-4 w-4" />
              退出登录
            </Button>
          </div>
        </Card>
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <Card className="p-6">
          <div className="flex items-center gap-3">
            <Radio className="h-5 w-5 text-primary" />
            <div>
              <h3 className="text-xl font-semibold text-text-primary">听书会员</h3>
            </div>
          </div>

          {odobVip?.is_vip ? (
            <div className="mt-5 space-y-4">
              <InfoBlock className="rounded-3xl p-5">
                <div className="flex items-center justify-between gap-3">
                  <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                    <Clock3 className="h-4 w-4" />
                    到期时间
                  </div>
                  <StatusBadge variant="danger">剩余 {odobVip.surplus_time} 天</StatusBadge>
                </div>
                <p className="mt-3 text-lg font-medium text-text-primary">{formatTimestamp(odobVip.expire_time)}</p>
                {odobVip.err_tips ? <p className="mt-3 text-sm text-warning">{odobVip.err_tips}</p> : null}
              </InfoBlock>

              <div className="grid gap-4 sm:grid-cols-2">
                <StatCard label="本周听书" value={odobVip.week_count} />
                <StatCard label="累计听书" value={odobVip.total_count} />
              </div>

              <div className="rounded-3xl bg-warning-soft p-5 text-warning">
                累计为你节省了 <span className="font-semibold">{odobVip.save_price}{odobVip.price_desc}</span>
              </div>
            </div>
          ) : (
            <div className="mt-5 rounded-3xl border border-dashed border-border p-8 text-sm text-text-muted">
              {data?.odobVipError || "当前账号未开通听书会员。"}
            </div>
          )}
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-3">
            <BookOpen className="h-5 w-5 text-primary" />
            <div>
              <h3 className="text-xl font-semibold text-text-primary">电子书会员</h3>
            </div>
          </div>

          {ebookVip?.is_vip ? (
            <div className="mt-5 space-y-4">
              <InfoBlock className="rounded-3xl p-5">
                <div className="flex items-center justify-between gap-3">
                  <div className="inline-flex items-center gap-2 text-sm text-text-muted">
                    <Clock3 className="h-4 w-4" />
                    到期时间
                  </div>
                  <StatusBadge variant="danger">剩余 {ebookVip.surplus_time} 天</StatusBadge>
                </div>
                <p className="mt-3 text-lg font-medium text-text-primary">{formatTimestamp(ebookVip.expire_time)}</p>
                {ebookVip.err_tips ? <p className="mt-3 text-sm text-warning">{ebookVip.err_tips}</p> : null}
              </InfoBlock>

              <div className="grid gap-4 sm:grid-cols-2">
                <StatCard label="本月读书" value={ebookVip.month_count} />
                <StatCard label="累计读书" value={ebookVip.total_count} />
              </div>

              <div className="rounded-3xl bg-accent-soft p-5 text-accent">
                累计为你节省了 <span className="font-semibold">{ebookVip.save_price}{ebookVip.price_desc}</span>
              </div>
            </div>
          ) : (
            <div className="mt-5 rounded-3xl border border-dashed border-border p-8 text-sm text-text-muted">
              {data?.ebookVipError || "当前账号未开通电子书会员。"}
            </div>
          )}
        </Card>
      </section>
    </main>
  )
}
