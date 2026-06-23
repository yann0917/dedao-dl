import { ChevronLeft, ChevronRight } from "lucide-react"
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react"
import useEmblaCarousel, { type UseEmblaCarouselType } from "embla-carousel-react"
import { cn } from "@/lib/cn"

export type CarouselApi = UseEmblaCarouselType[1]

type CarouselContextValue = {
  carouselRef: UseEmblaCarouselType[0]
  api: CarouselApi
  canScrollPrev: boolean
  canScrollNext: boolean
  selectedIndex: number
  scrollSnaps: number[]
  scrollPrev: () => void
  scrollNext: () => void
  scrollTo: (index: number) => void
}

type CarouselProps = {
  className?: string
  children: React.ReactNode
  opts?: Parameters<typeof useEmblaCarousel>[0]
  setApi?: (api: CarouselApi) => void
}

type CarouselContentProps = {
  className?: string
  children: React.ReactNode
}

type CarouselItemProps = {
  className?: string
  children: React.ReactNode
}

type CarouselArrowProps = {
  className?: string
}

const CarouselContext = createContext<CarouselContextValue | null>(null)

function useCarousel() {
  const context = useContext(CarouselContext)
  if (!context) {
    throw new Error("Carousel components must be used within Carousel")
  }

  return context
}

export function Carousel({ className, children, opts, setApi }: CarouselProps) {
  const [carouselRef, api] = useEmblaCarousel(opts)
  const [canScrollPrev, setCanScrollPrev] = useState(false)
  const [canScrollNext, setCanScrollNext] = useState(false)
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [scrollSnaps, setScrollSnaps] = useState<number[]>([])

  const updateState = useCallback((carouselApi: CarouselApi) => {
    if (!carouselApi) {
      return
    }

    setCanScrollPrev(carouselApi.canScrollPrev())
    setCanScrollNext(carouselApi.canScrollNext())
    setSelectedIndex(carouselApi.selectedScrollSnap())
    setScrollSnaps(carouselApi.scrollSnapList())
  }, [])

  useEffect(() => {
    if (!api) {
      return
    }

    setApi?.(api)
    updateState(api)

    const handleSelect = () => updateState(api)
    api.on("select", handleSelect)
    api.on("reInit", handleSelect)

    return () => {
      api.off("select", handleSelect)
      api.off("reInit", handleSelect)
    }
  }, [api, setApi, updateState])

  const value = useMemo<CarouselContextValue>(
    () => ({
      carouselRef,
      api,
      canScrollPrev,
      canScrollNext,
      selectedIndex,
      scrollSnaps,
      scrollPrev: () => api?.scrollPrev(),
      scrollNext: () => api?.scrollNext(),
      scrollTo: (index: number) => api?.scrollTo(index),
    }),
    [api, canScrollNext, canScrollPrev, carouselRef, scrollSnaps, selectedIndex],
  )

  return <CarouselContext.Provider value={value}><div className={cn("relative", className)}>{children}</div></CarouselContext.Provider>
}

export function CarouselContent({ className, children }: CarouselContentProps) {
  const { carouselRef } = useCarousel()

  return (
    <div className="h-full overflow-hidden" ref={carouselRef}>
      <div className={cn("flex h-full", className)}>{children}</div>
    </div>
  )
}

export function CarouselItem({ className, children }: CarouselItemProps) {
  return <div className={cn("h-full min-w-0 shrink-0 grow-0 basis-full", className)}>{children}</div>
}

export function CarouselPrevious({ className }: CarouselArrowProps) {
  const { canScrollPrev, scrollPrev } = useCarousel()

  return (
    <button
      aria-label="上一张"
      className={cn(
        "absolute left-4 top-1/2 z-20 inline-flex h-11 w-11 -translate-y-1/2 items-center justify-center rounded-full border border-white/20 bg-black/30 text-white shadow-lg backdrop-blur transition hover:bg-black/45 disabled:pointer-events-none disabled:opacity-35",
        className,
      )}
      disabled={!canScrollPrev}
      onClick={scrollPrev}
      type="button"
    >
      <ChevronLeft className="h-5 w-5" />
    </button>
  )
}

export function CarouselNext({ className }: CarouselArrowProps) {
  const { canScrollNext, scrollNext } = useCarousel()

  return (
    <button
      aria-label="下一张"
      className={cn(
        "absolute right-4 top-1/2 z-20 inline-flex h-11 w-11 -translate-y-1/2 items-center justify-center rounded-full border border-white/20 bg-black/30 text-white shadow-lg backdrop-blur transition hover:bg-black/45 disabled:pointer-events-none disabled:opacity-35",
        className,
      )}
      disabled={!canScrollNext}
      onClick={scrollNext}
      type="button"
    >
      <ChevronRight className="h-5 w-5" />
    </button>
  )
}

export function CarouselPagination({ className }: { className?: string }) {
  const { scrollSnaps, selectedIndex, scrollTo } = useCarousel()

  if (scrollSnaps.length <= 1) {
    return null
  }

  return (
    <div className={cn("absolute bottom-4 right-4 z-20 flex items-center gap-2 rounded-full bg-black/28 px-3 py-2 backdrop-blur-md", className)}>
      {scrollSnaps.map((_, index) => (
        <button
          aria-label={`切换到第 ${index + 1} 张`}
          className={cn(
            "h-2 w-2 rounded-full transition-all duration-300",
            index === selectedIndex ? "w-5 bg-white" : "bg-white/45 hover:bg-white/75",
          )}
          key={index}
          onClick={() => scrollTo(index)}
          type="button"
        />
      ))}
    </div>
  )
}
