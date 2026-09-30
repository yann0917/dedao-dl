import { useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { QrLoginCard } from "@/components/auth/QrLoginCard"
import { ThemeToggleButton } from "@/components/ui/ThemeToggleButton"
import { useAuth } from "@/providers/AuthProvider"

const tocEntries = [
  { no: "01", title: "扫码登录", desc: "得到 App 或微信扫一扫，登录你的账号。" },
  { no: "02", title: "浏览检索", desc: "已购的课程、电子书、听书，支持分类和搜索。" },
  { no: "03", title: "勾选下载", desc: "导出 PDF、Markdown 或音频，保存到指定目录。" },
]

export function LoginPage() {
  const navigate = useNavigate()
  const { loggedIn, loading, completeLogin } = useAuth()

  useEffect(() => {
    if (!loading && loggedIn) {
      navigate("/", { replace: true })
    }
  }, [loading, loggedIn, navigate])

  return (
    <main className="relative mx-auto flex min-h-screen w-full max-w-[1240px] flex-col px-6 pb-6 pt-7 sm:px-12 sm:pb-8 sm:pt-9 lg:px-16">
      <div aria-hidden="true" className="ed-frame" />
      <span
        aria-hidden="true"
        className="pointer-events-none absolute -right-6 bottom-0 hidden select-none overflow-hidden font-display text-[24rem] font-bold leading-none text-text-primary/[0.035] xl:block"
      >
        得
      </span>

      <header
        className="ed-rise relative z-10 flex items-center justify-between"
        style={{ animationDelay: "40ms" }}
      >
        <div className="flex items-center gap-4">
          <span className="ed-stamp ed-seal flex h-11 w-11 items-center justify-center bg-accent font-display text-[22px] font-bold text-accent-foreground">
            得
          </span>
          <div className="leading-tight">
            <p className="font-display text-lg font-semibold tracking-[0.08em] text-text-primary">
              dedao-dl
            </p>
            <p className="mt-0.5 font-mono text-[10px] uppercase tracking-[0.32em] text-text-muted">
              cli tool · web ui
            </p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <span className="hidden font-mono text-[11px] tracking-[0.18em] text-text-muted sm:block">
            127.0.0.1:17878
          </span>
          <ThemeToggleButton className="rounded-none" />
        </div>
      </header>

      <div className="relative z-10 grid flex-1 items-center gap-14 py-8 lg:grid-cols-[1.1fr_0.9fr] lg:gap-20 lg:py-6">
        <span
          aria-hidden="true"
          className="ed-vertical ed-rise absolute -left-12 top-1/2 hidden -translate-y-1/2 font-display text-sm text-text-muted/70 xl:block"
          style={{ animationDelay: "500ms" }}
        >
          知识就在得到
        </span>

        <section>
          <p
            className="ed-rise flex items-center gap-3 font-mono text-[11px] uppercase tracking-[0.34em] text-accent"
            style={{ animationDelay: "120ms" }}
          >
            <span aria-hidden="true" className="h-px w-9 bg-accent" />
            得到内容下载工具
          </p>

          <h1
            className="ed-rise mt-7 max-w-xl font-display text-[clamp(34px,4.2vw,54px)] font-bold leading-[1.32] tracking-[0.02em] text-text-primary"
            style={{ animationDelay: "200ms" }}
          >
            扫码登录得到账号
            <br />
            把<span className="ed-emphasis">已购内容</span>下载到本地
          </h1>

          <p
            className="ed-rise mt-7 max-w-md text-[15px] leading-7 text-text-secondary"
            style={{ animationDelay: "300ms" }}
          >
            dedao-dl 命令行工具的网页版：浏览、搜索、勾选，然后下载，不用再记命令。
          </p>

          <ol
            className="ed-rise mt-11 max-w-lg"
            style={{ animationDelay: "400ms" }}
          >
            {tocEntries.map((item) => (
              <li
                key={item.no}
                className="grid grid-cols-[2.5rem_6rem_1fr] items-baseline gap-2 border-t border-border py-3.5 last:border-b"
              >
                <span className="font-mono text-xs text-accent">{item.no}</span>
                <span className="text-sm font-medium text-text-primary">{item.title}</span>
                <span className="text-sm leading-6 text-text-muted">{item.desc}</span>
              </li>
            ))}
          </ol>
        </section>

        <section className="ed-rise" style={{ animationDelay: "340ms" }}>
          <QrLoginCard
            onLoginSuccess={async (user) => {
              await completeLogin(user)
              navigate("/", { replace: true })
            }}
          />
        </section>
      </div>

      <footer className="relative z-10 flex items-center justify-between border-t border-border pt-4 font-mono text-[11px] tracking-[0.14em] text-text-muted">
        <span>dedao-dl · web</span>
        <span className="hidden sm:block">下载内容仅供个人学习</span>
      </footer>
    </main>
  )
}
