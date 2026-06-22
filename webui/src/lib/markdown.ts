function escapeHtml(value: string) {
  return value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;")
}

function renderInline(value: string) {
  let html = escapeHtml(value)

  html = html.replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"([^"]+)")?\)/g, (_match: string, alt: string, src: string, title?: string) => {
    const titleAttr = title ? ` title="${escapeHtml(title)}"` : ""
    return `<img src="${escapeHtml(src)}" alt="${escapeHtml(alt)}"${titleAttr} />`
  })

  html = html.replace(/`([^`]+)`/g, "<code>$1</code>")
  html = html.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
  html = html.replace(/\*([^*]+)\*/g, "<em>$1</em>")

  return html
}

export function renderMarkdownToHtml(markdown: string) {
  const lines = markdown.replace(/\r\n/g, "\n").split("\n")
  const html: string[] = []
  let paragraph: string[] = []
  let listItems: Array<{ type: "ul" | "ol"; value: string }> = []
  let blockquote: string[] = []

  const flushParagraph = () => {
    if (paragraph.length === 0) {
      return
    }
    html.push(`<p>${renderInline(paragraph.join(" "))}</p>`)
    paragraph = []
  }

  const flushList = () => {
    if (listItems.length === 0) {
      return
    }
    const listType = listItems[0].type
    const tagName = listType === "ol" ? "ol" : "ul"
    html.push(`<${tagName}>${listItems.map((item) => `<li>${renderInline(item.value)}</li>`).join("")}</${tagName}>`)
    listItems = []
  }

  const flushBlockquote = () => {
    if (blockquote.length === 0) {
      return
    }
    html.push(`<blockquote>${blockquote.map((item) => `<p>${renderInline(item)}</p>`).join("")}</blockquote>`)
    blockquote = []
  }

  for (const rawLine of lines) {
    const line = rawLine.trim()

    if (line === "") {
      flushParagraph()
      flushList()
      flushBlockquote()
      continue
    }

    const heading = line.match(/^(#{1,6})\s+(.*)$/)
    if (heading) {
      flushParagraph()
      flushList()
      flushBlockquote()
      const level = heading[1].length
      html.push(`<h${level}>${renderInline(heading[2])}</h${level}>`)
      continue
    }

    if (line === "---") {
      flushParagraph()
      flushList()
      flushBlockquote()
      html.push("<hr />")
      continue
    }

    if (line.startsWith("> ")) {
      flushParagraph()
      flushList()
      blockquote.push(line.slice(2))
      continue
    }

    if (line.startsWith("* ")) {
      flushParagraph()
      flushBlockquote()
      listItems.push({ type: "ul", value: line.slice(2) })
      continue
    }

    const orderedList = line.match(/^\d+\.\s+(.*)$/)
    if (orderedList) {
      flushParagraph()
      flushBlockquote()
      listItems.push({ type: "ol", value: orderedList[1] })
      continue
    }

    flushList()
    flushBlockquote()
    paragraph.push(line)
  }

  flushParagraph()
  flushList()
  flushBlockquote()

  return html.join("\n")
}
