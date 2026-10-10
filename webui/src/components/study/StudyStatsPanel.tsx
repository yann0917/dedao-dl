import type { RecentNameCount, RecentStats } from "@/api"
import { StatCard } from "@/components/ui/Semantic"
import { semanticMetaTextClass, semanticPageSectionClass } from "@/lib/semanticStyles"

// 内容构成按类型名固定配色（accent 朱砂红 / success 绿 / warning 琥珀 / 中性灰），
// 顶层图例与各图表的堆叠分段共用这套颜色；未知类型落到中性灰
const TYPE_COLOR_CLASSES: Record<string, string> = {
  听书: "bg-accent",
  课程: "bg-success",
  电子书: "bg-warning",
  锦囊: "bg-text-muted/70",
}

const UNKNOWN_TYPE_COLOR = "bg-text-muted/50"

function typeColorClass(name: string) {
  return TYPE_COLOR_CLASSES[name] ?? UNKNOWN_TYPE_COLOR
}

function typeBreakdownText(types: RecentNameCount[]) {
  return types.map((item) => `${item.name} ${item.count}`).join("、")
}

function TypeComposition({ stats }: { stats: RecentStats }) {
  return (
    <div className="mt-6">
      <p className="text-sm font-medium text-text-primary">内容构成</p>
      <div className="mt-3 flex h-3 gap-0.5 overflow-hidden rounded-sm">
        {stats.types.map((item) => (
          <div
            className={`h-full ${typeColorClass(item.name)}`}
            key={item.name}
            style={{ flexGrow: item.count }}
            title={`${item.name} ${item.count} 条`}
          />
        ))}
      </div>
      <div className="mt-3 flex flex-wrap gap-x-5 gap-y-2 text-xs text-text-secondary">
        {stats.types.map((item) => (
          <span className="inline-flex items-center gap-2" key={item.name}>
            <span className={`size-2.5 rounded-sm ${typeColorClass(item.name)}`} />
            {item.name} {item.count}
          </span>
        ))}
      </div>
    </div>
  )
}

// 近 14 天每日条数：柱内按内容类型堆叠着色（与内容构成图例同色），看得出每天学了什么类型
function DailyChart({ stats }: { stats: RecentStats }) {
  const max = Math.max(...stats.daily.map((day) => day.count), 1)

  return (
    <div>
      <p className="text-sm font-medium text-text-primary">每日学习条数（近 {stats.daily.length} 天）</p>
      <div className="mt-3 flex h-36 items-end gap-1.5">
        {stats.daily.map((day, index) => (
          <div
            className="flex h-full min-w-0 flex-1 flex-col items-center justify-end gap-1"
            key={day.date}
            title={day.count > 0 ? `${day.date} · ${day.count} 条：${typeBreakdownText(day.types)}` : `${day.date} · 无记录`}
          >
            {day.count > 0 ? <span className="text-[10px] leading-none text-text-muted">{day.count}</span> : null}
            {day.count > 0 ? (
              <div
                className="flex w-full flex-col-reverse overflow-hidden rounded-sm"
                style={{ height: `${Math.max((day.count / max) * 100, 6)}%` }}
              >
                {day.types.map((segment) => (
                  <div
                    className={`w-full ${typeColorClass(segment.name)}`}
                    key={segment.name}
                    style={{ flexGrow: segment.count }}
                  />
                ))}
              </div>
            ) : (
              <div className="w-full rounded-sm bg-surface-soft" style={{ height: "2px" }} />
            )}
            <span className={`text-[10px] leading-none text-text-muted ${index % 2 === 1 ? "hidden sm:inline" : ""}`}>{day.date}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

// 时段分桶横向条：条内按内容类型堆叠着色（与内容构成图例同色），悬停可看类型细分与桶内高峰小时
const HOUR_BUCKETS = [
  { label: "凌晨", range: "00–06", from: 0, to: 6 },
  { label: "清晨", range: "06–09", from: 6, to: 9 },
  { label: "上午", range: "09–12", from: 9, to: 12 },
  { label: "午间", range: "12–14", from: 12, to: 14 },
  { label: "下午", range: "14–18", from: 14, to: 18 },
  { label: "晚间", range: "18–21", from: 18, to: 21 },
  { label: "深夜", range: "21–24", from: 21, to: 24 },
]

function HourChart({ stats }: { stats: RecentStats }) {
  const byHour = new Map(stats.hours.map((item) => [item.hour, item]))
  const peak = stats.hours.reduce((top, item) => (item.count > top.count ? item : top), { hour: 0, count: 0 })

  const buckets = HOUR_BUCKETS.map((bucket) => {
    const typeMap = new Map<string, number>()
    let count = 0
    let peakHour = -1
    let peakCount = 0
    for (let hour = bucket.from; hour < bucket.to; hour++) {
      const hourStat = byHour.get(hour)
      if (!hourStat) {
        continue
      }
      count += hourStat.count
      if (hourStat.count > peakCount) {
        peakCount = hourStat.count
        peakHour = hour
      }
      for (const segment of hourStat.types) {
        typeMap.set(segment.name, (typeMap.get(segment.name) ?? 0) + segment.count)
      }
    }
    // 按 stats.types 的全局顺序还原，保证各图分段顺序一致
    const types = stats.types
      .filter((item) => typeMap.has(item.name))
      .map((item) => ({ name: item.name, count: typeMap.get(item.name) ?? 0 }))
    return { ...bucket, count, peakHour, peakCount, types }
  })
  const max = Math.max(...buckets.map((bucket) => bucket.count), 1)

  return (
    <div>
      <div className="flex flex-wrap items-baseline justify-between gap-2">
        <p className="text-sm font-medium text-text-primary">学习时段分布</p>
        {peak.count > 0 ? (
          <span className="text-xs text-text-muted">
            单小时高峰 {String(peak.hour).padStart(2, "0")}:00（{peak.count} 条）
          </span>
        ) : null}
      </div>
      <div className="mt-3 space-y-2">
        {buckets.map((bucket) => (
          <div
            className="flex items-center gap-3"
            key={bucket.label}
            title={
              bucket.count > 0
                ? `${bucket.label} ${bucket.range} · ${bucket.count} 条：${typeBreakdownText(bucket.types)}`
                : `${bucket.label} ${bucket.range} · 无记录`
            }
          >
            <span className="w-24 shrink-0 text-xs text-text-secondary">
              {bucket.label}
              <span className="ml-1 text-text-muted">{bucket.range}</span>
            </span>
            <div className="h-4 flex-1 overflow-hidden rounded-sm bg-surface-soft">
              <div className="flex h-full" style={{ width: `${(bucket.count / max) * 100}%` }}>
                {bucket.types.map((segment) => (
                  <div
                    className={`h-full ${typeColorClass(segment.name)}`}
                    key={segment.name}
                    style={{ flexGrow: segment.count }}
                  />
                ))}
              </div>
            </div>
            <span className="w-8 shrink-0 text-right text-xs tabular-nums text-text-secondary">{bucket.count}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

export function StudyStatsPanel({ stats }: { stats: RecentStats }) {
  return (
    <section className={`${semanticPageSectionClass} p-6`} aria-label="近30天学习统计">
      <div className="flex flex-wrap items-baseline justify-between gap-2">
        <p className={semanticMetaTextClass}>近 {stats.window_days} 天统计</p>
        {stats.truncated ? <span className="text-xs text-text-muted">窗口内记录较多，仅统计最近 500 条</span> : null}
      </div>

      <div className="mt-4 grid grid-cols-2 gap-3 lg:grid-cols-4">
        <StatCard label="记录数" value={stats.total} />
        <StatCard
          label="学习天数"
          value={
            <>
              {stats.active_days}
              <span className="text-sm font-normal text-text-muted">/{stats.window_days}</span>
            </>
          }
        />
        <StatCard
          label="完成率"
          value={
            <>
              {stats.finish_rate}
              <span className="text-sm font-normal text-text-muted">%</span>
            </>
          }
        />
        <StatCard label="主力类型" value={stats.top_type || "—"} />
      </div>

      <TypeComposition stats={stats} />

      <div className="mt-6 grid gap-6 lg:grid-cols-2">
        <DailyChart stats={stats} />
        <HourChart stats={stats} />
      </div>
    </section>
  )
}
