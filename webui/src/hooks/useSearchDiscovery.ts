import { useEffect, useMemo, useState } from "react"
import {
  api,
  type SearchHotResponse,
  type SearchSuggestResponse,
} from "@/api"

export function useSearchDiscovery() {
  const [hot, setHot] = useState<SearchHotResponse | null>(null)
  const [suggestions, setSuggestions] = useState<SearchSuggestResponse | null>(null)
  const [query, setQuery] = useState("")
  const [error, setError] = useState<string | null>(null)

  const hotKeywords = useMemo(
    () =>
      hot?.hot_tab_list
        .flatMap((item) => item.list.slice(0, 3).map((entry) => entry.searchKey || entry.title))
        .slice(0, 8) ?? [],
    [hot],
  )

  useEffect(() => {
    api.search
      .hot()
      .then(setHot)
      .catch((err: Error) => setError(err.message))
  }, [])

  useEffect(() => {
    if (!query.trim()) {
      setSuggestions(null)
      return
    }

    const timer = window.setTimeout(() => {
      api.search
        .suggest(query)
        .then(setSuggestions)
        .catch((err: Error) => setError(err.message))
    }, 300)

    return () => window.clearTimeout(timer)
  }, [query])

  return {
    hotKeywords,
    suggestions,
    query,
    setQuery,
    error,
  }
}
