import { BookOpen, Loader2, PlayCircle } from "lucide-react"
import { useNavigate } from "react-router-dom"
import {
  type HomeContentResponse,
  type HomeFreeResource,
  type HomeModule,
  type HomeNavigation,
} from "@/api"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { cn } from "@/lib/cn"

type SectionHeaderProps = {
  module?: HomeModule
  actions?: React.ReactNode
}

function SectionHeader({ module, actions }: SectionHeaderProps) {
  if (!module) {
    return null
  }

  return (
    <div className="mb-5 flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <h3 className="text-2xl font-semibold text-slate-950">{module.title}</h3>
        <p className="mt-2 text-sm leading-6 text-slate-500">{module.description}</p>
      </div>
      {actions}
    </div>
  )
}

type FreeResourceSectionProps = {
  module?: HomeModule
  resources: HomeFreeResource[]
  error?: string | null
  onOpenResource: (resource: HomeFreeResource) => void
}

export function FreeResourceSection({ module, resources, error, onOpenResource }: FreeResourceSectionProps) {
  if (!module) {
    return null
  }

  return (
    <section>
      <SectionHeader module={module} />
      {error ? (
        <Card className="border border-amber-200 bg-amber-50 p-4 text-sm text-amber-700">{error}</Card>
      ) : null}

      <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
        {resources.map((resource) => (
          <button className="text-left" key={resource.id} onClick={() => onOpenResource(resource)} type="button">
            <Card className="h-full overflow-hidden transition hover:-translate-y-1 hover:shadow-lg">
              <div className="relative aspect-video overflow-hidden">
                <img
                  alt={resource.name}
                  className="h-full w-full object-cover transition duration-500 hover:scale-105"
                  src={resource.logo}
                />
                <div className="absolute inset-0 flex items-center justify-center bg-slate-950/10">
                  <div className="rounded-full bg-white/90 p-3 text-primary shadow-md">
                    <PlayCircle className="h-6 w-6" />
                  </div>
                </div>
              </div>
              <div className="p-4">
                <h4 className="line-clamp-2 text-base font-semibold text-slate-950">{resource.name}</h4>
                <p className="mt-2 line-clamp-2 text-sm leading-6 text-slate-500">{resource.intro}</p>
                <p className="mt-4 text-xs text-slate-400">免费专区内容入口</p>
              </div>
            </Card>
          </button>
        ))}
      </div>
    </section>
  )
}

type LabeledShelfSectionProps = {
  module?: HomeModule
  labels: HomeNavigation[]
  content: HomeContentResponse | null
  selectedEnid: string
  loading: boolean
  error?: string | null
  onSelectLabel: (label: HomeNavigation) => void
  onOpenProduct: (productEnid: string) => void
  moreButtonText: string
  moreButtonDisabled?: boolean
  onClickMore?: () => void
  variant: "course" | "ebook"
}

export function LabeledShelfSection({
  module,
  labels,
  content,
  selectedEnid,
  loading,
  error,
  onSelectLabel,
  onOpenProduct,
  moreButtonText,
  moreButtonDisabled = false,
  onClickMore,
  variant,
}: LabeledShelfSectionProps) {
  const navigate = useNavigate()

  if (!module) {
    return null
  }

  const showCards = content?.product_list ?? []

  return (
    <section>
      <SectionHeader
        actions={
          labels.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {labels.map((label) => (
                <button
                  className={cn(
                    "rounded-full px-3 py-1.5 text-xs font-medium transition",
                    selectedEnid === label.enid
                      ? "bg-primary text-white"
                      : "bg-slate-100 text-slate-600 hover:bg-slate-200",
                  )}
                  key={label.enid}
                  onClick={() => onSelectLabel(label)}
                  type="button"
                >
                  {label.name}
                </button>
              ))}
            </div>
          ) : null
        }
        module={module}
      />

      {error ? (
        <Card className="mb-5 border border-amber-200 bg-amber-50 p-4 text-sm text-amber-700">{error}</Card>
      ) : null}

      {loading ? (
        <Card className="flex items-center justify-center p-8 text-slate-500">
          <Loader2 className="mr-2 h-5 w-5 animate-spin" />
          正在刷新{variant === "course" ? "课程" : "电子书"}内容...
        </Card>
      ) : null}

      <div className={cn("grid gap-5", variant === "course" ? "md:grid-cols-2 xl:grid-cols-4" : "md:grid-cols-3 xl:grid-cols-5")}>
        {showCards.map((product) => (
          <button
            className="text-left"
            key={`${product.product_enid}-${product.title}`}
            onClick={() => onOpenProduct(product.product_enid)}
            type="button"
          >
            <Card className="h-full overflow-hidden transition hover:-translate-y-1 hover:shadow-lg">
              <div className={cn("overflow-hidden bg-slate-100", variant === "course" ? "aspect-video" : "aspect-[3/4] p-3")}>
                <img
                  alt={product.title}
                  className={cn(
                    "h-full w-full transition duration-500 hover:scale-105",
                    variant === "course" ? "object-cover" : "rounded-2xl object-cover shadow-md",
                  )}
                  src={variant === "course" ? product.horizontal_image : product.index_image}
                />
              </div>
              <div className="p-4">
                <h4 className="line-clamp-2 text-base font-semibold text-slate-950">{product.title}</h4>
                <p className="mt-2 line-clamp-2 text-sm leading-6 text-slate-500">
                  {variant === "course" ? product.intro : product.author_list.join(" / ") || "电子书内容"}
                </p>
                <div className="mt-4 flex items-center justify-between text-xs text-slate-400">
                  {variant === "course" ? (
                    <span>{product.learn_user_count} 人加入</span>
                  ) : (
                    <span className="inline-flex items-center gap-1">
                      <BookOpen className="h-3.5 w-3.5" />
                      {product.score ? `评分 ${product.score}` : "暂无评分"}
                    </span>
                  )}
                  <span>{variant === "course" ? "进入工作区" : "详情后续接入"}</span>
                </div>
              </div>
            </Card>
          </button>
        ))}
      </div>

      <div className="mt-6 flex justify-center">
        <Button
          className="min-w-[240px]"
          disabled={moreButtonDisabled}
          onClick={() => {
            if (onClickMore) {
              onClickMore()
              return
            }

            navigate("/courses")
          }}
          variant="outline"
        >
          {moreButtonText}
        </Button>
      </div>
    </section>
  )
}
