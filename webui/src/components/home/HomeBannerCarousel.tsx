import { useEffect, useState } from "react"
import { type HomeBanner } from "@/api"
import type { CarouselApi } from "@/components/ui/Carousel"
import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPagination, CarouselPrevious } from "@/components/ui/Carousel"
import { Card } from "@/components/ui/Card"

type HomeBannerCarouselProps = {
  banners: HomeBanner[]
}

export function HomeBannerCarousel({ banners }: HomeBannerCarouselProps) {
  const [api, setApi] = useState<CarouselApi>()

  useEffect(() => {
    if (!api || banners.length <= 1) {
      return
    }

    const timer = window.setInterval(() => {
      api.scrollNext()
    }, 5000)

    return () => window.clearInterval(timer)
  }, [api, banners.length])

  if (banners.length === 0) {
    return (
      <Card className="flex aspect-video w-full items-center justify-center p-8">
        <div className="text-center">
          <p className="text-lg font-medium text-text-primary">首页 Banner 暂未加载</p>
        </div>
      </Card>
    )
  }

  return (
    <Card className="relative aspect-video w-full self-start overflow-hidden">
      <Carousel className="h-full w-full" opts={{ align: "start", loop: banners.length > 1 }} setApi={setApi}>
        <CarouselContent className="h-full">
          {banners.map((banner) => (
            <CarouselItem className="h-full" key={banner.id}>
              <button
                className="group relative h-full w-full text-left"
                onClick={() => window.open(banner.url, "_blank", "noopener,noreferrer")}
                type="button"
              >
                <img
                  alt=""
                  aria-hidden="true"
                  className="absolute inset-0 h-full w-full scale-[1.03] object-cover opacity-48 blur-xl transition duration-500"
                  src={banner.img}
                />
                <div className="absolute inset-0 bg-secondary/18" />
                <div className="absolute inset-0">
                  <img
                    alt={banner.title || "banner"}
                    className="h-full w-full object-cover object-center transition duration-500"
                    src={banner.img}
                  />
                </div>
                <img
                  alt=""
                  aria-hidden="true"
                  className="absolute inset-0 h-full w-full object-cover opacity-[0.06]"
                  src={banner.img}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-secondary/34 via-secondary/14 to-secondary/8" />
              </button>
            </CarouselItem>
          ))}
        </CarouselContent>

        {banners.length > 1 ? (
          <>
            <CarouselPrevious />
            <CarouselNext />
            <CarouselPagination />
          </>
        ) : null}
      </Carousel>
    </Card>
  )
}
