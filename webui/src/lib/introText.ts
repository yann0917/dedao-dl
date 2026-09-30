// 得到接口的部分 intro 字段是富文本 JSON 串（如听书合集 spu.intro），
// 形如 [{"type":"paragraph","contents":[{"type":"text","text":{"content":"..."}}]}]。
// 这里把它安全地还原成纯文本；输入不是 JSON 时原样返回。
export function introToText(value: unknown): string {
  if (typeof value !== "string") {
    return ""
  }

  const text = value.trim()
  if (!text || (text[0] !== "[" && text[0] !== "{")) {
    return text
  }

  try {
    const parsed = JSON.parse(text) as unknown
    const parts: string[] = []

    const walk = (node: unknown) => {
      if (typeof node === "string") {
        return
      }

      if (Array.isArray(node)) {
        node.forEach(walk)
        return
      }

      if (node && typeof node === "object") {
        const record = node as Record<string, unknown>
        const textNode = record.text
        if (textNode && typeof textNode === "object" && typeof (textNode as Record<string, unknown>).content === "string") {
          parts.push((textNode as Record<string, unknown>).content as string)
          return
        }
        Object.values(record).forEach(walk)
      }
    }

    walk(parsed)
    const joined = parts.join("\n").trim()
    return joined || text
  } catch {
    return text
  }
}
