import { useEffect, useState } from "react"
import { ArrowUpRight } from "lucide-react"
import { type HomeBanner } from "@/api"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"
import { getSemanticStatusBadgeClass } from "@/lib/semanticStyles"

type HomeBannerCarouselProps = {
  banners: HomeBanner[]
}

export function HomeBannerCarousel({ banners }: HomeBannerCarouselProps) {
  const [activeIndex, setActiveIndex] = useState(0)

  useEffect(() => {
    if (banners.length <= 1) {
      return
    }

    const timer = window.setInterval(() => {
      setActiveIndex((current) => (current + 1) % banners.length)
    }, 5000)

    return () => window.clearInterval(timer)
  }, [banners.length])

  useEffect(() => {
    if (activeIndex < banners.length) {
      return
    }

    setActiveIndex(0)
  }, [activeIndex, banners.length])

  if (banners.length === 0) {
    return (
      <Card className="flex h-full min-h-[360px] items-center justify-center p-8">
        <div className="text-center">
          <p className="text-lg font-medium text-text-primary">首页 Banner 暂未加载</p>
        </div>
      </Card>
    )
  }

  const currentBanner = banners[activeIndex]

  return (
    <Card className="relative h-full min-h-[360px] overflow-hidden">
      <button
        className="group relative h-full w-full text-left"
        onClick={() => window.open(currentBanner.url, "_blank", "noopener,noreferrer")}
        type="button"
      >
        <img
          alt={currentBanner.title || "banner"}
          className="h-full w-full object-cover transition duration-500 group-hover:scale-[1.02]"
          src={currentBanner.img}
        />
        <div className="absolute inset-0 bg-gradient-to-t from-secondary/85 via-secondary/15 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 p-6 text-white">
          <div className={cn(getSemanticStatusBadgeClass("accent"), "inline-flex items-center gap-2 bg-white/15 text-white backdrop-blur")}>
            首页推荐
            <ArrowUpRight className="h-3.5 w-3.5" />
          </div>
          <h3 className="mt-3 text-2xl font-semibold">{currentBanner.title || "得到首页推荐内容"}</h3>
        </div>
      </button>

      {banners.length > 1 ? (
        <div className="absolute bottom-5 right-5 flex items-center gap-2 rounded-full bg-secondary/60 px-3 py-2 backdrop-blur">
          {banners.map((banner, index) => (
            <button
              aria-label={`切换到第 ${index + 1} 个 Banner`}
              className={cn(
                "h-2.5 w-2.5 rounded-full transition",
                index === activeIndex ? "bg-white" : "bg-white/40 hover:bg-white/70",
              )}
              key={banner.id}
              onClick={() => setActiveIndex(index)}
              type="button"
            />
          ))}
        </div>
      ) : null}
    </Card>
  )
}
