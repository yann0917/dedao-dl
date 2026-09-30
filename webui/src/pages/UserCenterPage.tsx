import type { ReactNode } from "react"
import { Loader2, LogOut, RefreshCcw, SwitchCamera } from "lucide-react"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"
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

function SectionHead({ title, actions }: { title: string; actions?: ReactNode }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-border pb-3">
      <div className="flex min-w-0 items-center gap-3">
        <span aria-hidden="true" className="h-4 w-[3px] shrink-0 bg-accent" />
        <h2 className="shrink-0 text-xl font-semibold tracking-wide text-text-primary">{title}</h2>
      </div>
      {actions}
    </div>
  )
}

function LedgerRow({
  label,
  value,
  accent = false,
}: {
  label: string
  value: ReactNode
  accent?: boolean
}) {
  return (
    <div className="flex items-baseline justify-between gap-4 py-3.5">
      <dt className="font-mono text-xs tracking-[0.12em] text-text-muted">{label}</dt>
      <dd
        className={cn(
          "text-right text-sm font-medium tabular-nums",
          accent ? "text-accent" : "text-text-primary",
        )}
      >
        {value}
      </dd>
    </div>
  )
}

function MemberTag({ children }: { children: ReactNode }) {
  return (
    <span className="border border-border-strong/60 px-2 py-1 font-mono text-[10px] tracking-[0.12em] text-text-secondary">
      {children}
    </span>
  )
}

function BigStat({
  label,
  value,
  hint,
}: {
  label: string
  value: ReactNode
  hint: string
}) {
  return (
    <div className="py-5">
      <dt className="font-mono text-[10px] uppercase tracking-[0.24em] text-text-muted">{label}</dt>
      <dd className="mt-2 text-3xl font-semibold tabular-nums text-text-primary">
        {value}
        <span className="ml-1.5 text-xs font-normal text-text-muted">{hint}</span>
      </dd>
    </div>
  )
}

export function UserCenterPage() {
  const { data, loading, switchingUID, error, reload, switchAccount, logout } = useUserCenter()

  if (loading) {
    return (
      <main className="flex min-h-[60vh] items-center justify-center">
        <div className="flex items-center gap-3 text-text-muted">
          <Loader2 className="size-icon-lg animate-spin" />
          正在加载用户中心...
        </div>
      </main>
    )
  }

  const user = data?.user
  const ebookVip = data?.ebookVip
  const odobVip = data?.odobVip?.user

  return (
    <main className="space-y-12">
      {error ? (
        <Card className="border-danger bg-danger-soft">
          <div className="p-4 text-sm text-danger">{error}</div>
        </Card>
      ) : null}

      {/* 档案头 */}
      <section>
        <div className="flex flex-col gap-6 md:flex-row md:items-start md:justify-between">
          <div className="flex min-w-0 items-center gap-6">
            <img
              alt={user?.nickname ?? "avatar"}
              className="h-24 w-24 shrink-0 rounded-md border border-border object-cover"
              src={user?.avatar || "https://placehold.co/120x120/e2e8f0/334155?text=DD"}
            />
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2.5">
                <h1 className="text-3xl font-bold tracking-wide text-text-primary">
                  {user?.nickname ?? "未命名用户"}
                </h1>
                {user?.is_teacher ? (
                  <span className="bg-accent px-1.5 py-0.5 text-[10px] font-medium text-accent-foreground">
                    教师
                  </span>
                ) : null}
              </div>
              <p className="mt-2 truncate font-mono text-[11px] uppercase tracking-[0.2em] text-text-muted">
                当前登录账号{user?.uid_hazy ? ` · ${user.uid_hazy}` : ""}
              </p>
              <div className="mt-3 flex flex-wrap gap-2">
                {odobVip?.is_vip ? <MemberTag>听书会员 · 余 {odobVip.surplus_time} 天</MemberTag> : null}
                {ebookVip?.is_vip ? <MemberTag>电子书会员 · 余 {ebookVip.surplus_time} 天</MemberTag> : null}
              </div>
            </div>
          </div>

          <Button className="h-9 shrink-0 px-4 text-xs" onClick={() => void reload()} variant="outline">
            <RefreshCcw className="mr-2 size-icon-sm" />
            刷新资料
          </Button>
        </div>

        <dl className="mt-8 grid grid-cols-2 divide-x divide-border border-y border-border">
          <div className="pr-6">
            <BigStat label="今日学习" value={Math.round((user?.today_study_time ?? 0) / 60)} hint="分钟" />
          </div>
          <div className="pl-6">
            <BigStat label="连续学习" value={user?.study_serial_days ?? 0} hint="天" />
          </div>
        </dl>
      </section>

      {/* 会员台账 */}
      <section className="grid gap-x-10 gap-y-12 lg:grid-cols-2">
        <div>
          <SectionHead title="听书会员" />
          {odobVip?.is_vip ? (
            <dl className="divide-y divide-border">
              <LedgerRow label="到期时间" value={formatTimestamp(odobVip.expire_time)} />
              <LedgerRow label="剩余天数" value={`${odobVip.surplus_time} 天`} />
              <LedgerRow label="本周听书" value={odobVip.week_count} />
              <LedgerRow label="累计听书" value={odobVip.total_count} />
              <LedgerRow accent label="累计节省" value={`${odobVip.save_price}${odobVip.price_desc}`} />
            </dl>
          ) : (
            <p className="pt-4 font-mono text-xs tracking-[0.12em] text-text-muted">
              {data?.odobVipError || "当前账号未开通听书会员。"}
            </p>
          )}
          {odobVip?.err_tips ? <p className="pt-3 text-sm text-warning">{odobVip.err_tips}</p> : null}
        </div>

        <div>
          <SectionHead title="电子书会员" />
          {ebookVip?.is_vip ? (
            <dl className="divide-y divide-border">
              <LedgerRow label="到期时间" value={formatTimestamp(ebookVip.expire_time)} />
              <LedgerRow label="剩余天数" value={`${ebookVip.surplus_time} 天`} />
              <LedgerRow label="本月读书" value={ebookVip.month_count} />
              <LedgerRow label="累计读书" value={ebookVip.total_count} />
              <LedgerRow accent label="累计节省" value={`${ebookVip.save_price}${ebookVip.price_desc}`} />
            </dl>
          ) : (
            <p className="pt-4 font-mono text-xs tracking-[0.12em] text-text-muted">
              {data?.ebookVipError || "当前账号未开通电子书会员。"}
            </p>
          )}
          {ebookVip?.err_tips ? <p className="pt-3 text-sm text-warning">{ebookVip.err_tips}</p> : null}
        </div>
      </section>

      {/* 账号管理 */}
      <section>
        <SectionHead title="账号管理" />
        <ul className="divide-y divide-border">
          {data?.accounts.map((account) => (
            <li className="flex items-center gap-4 py-4" key={account.uidHazy}>
              <img
                alt={account.name}
                className="h-10 w-10 shrink-0 rounded-md border border-border object-cover"
                src={account.avatar || "https://placehold.co/80x80/e2e8f0/334155?text=DD"}
              />
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-text-primary">{account.name}</p>
                <p className="mt-0.5 truncate font-mono text-[10px] uppercase tracking-[0.14em] text-text-muted">
                  {account.uidHazy}
                </p>
              </div>
              {account.active ? (
                <span className="inline-flex shrink-0 items-center gap-1.5 font-mono text-xs text-accent">
                  <span aria-hidden="true" className="h-1.5 w-1.5 rounded-full bg-accent" />
                  使用中
                </span>
              ) : (
                <Button
                  className="h-8 shrink-0 px-3 text-xs"
                  disabled={switchingUID === account.uidHazy}
                  onClick={() => void switchAccount(account.uidHazy)}
                  variant="outline"
                >
                  <SwitchCamera className="mr-1.5 size-icon-sm" />
                  {switchingUID === account.uidHazy ? "切换中..." : "切换"}
                </Button>
              )}
            </li>
          ))}
        </ul>

        <button
          className="mt-6 inline-flex items-center gap-1.5 text-sm text-danger underline decoration-danger/30 underline-offset-4 transition hover:decoration-danger"
          onClick={() => void logout()}
          type="button"
        >
          <LogOut className="size-icon-sm" />
          退出登录
        </button>
      </section>
    </main>
  )
}
