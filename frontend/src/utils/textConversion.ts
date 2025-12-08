import { sify, tify } from 'chinese-conv'

type TargetLang = 'zh-Hans' | 'zh-Hant'

// Restrict conversion to CJK ideographs to avoid touching punctuation/quotes.
function isCJK(char: string): boolean {
  const code = char.codePointAt(0)
  if (code === undefined) return false
  // Basic + Extension A
  if ((code >= 0x4e00 && code <= 0x9fff) || (code >= 0x3400 && code <= 0x4dbf)) {
    return true
  }
  // Extension B..F (rare but safe)
  if (code >= 0x20000 && code <= 0x2ebef) return true
  return false
}

export function convertText(text: string, target: TargetLang): string {
  try {
    const convertChar = target === 'zh-Hans' ? sify : tify
    let out = ''
    for (const ch of text) {
      out += isCJK(ch) ? convertChar(ch) : ch
    }
    return out
  } catch (err) {
    console.error('Text conversion failed', err)
    return text
  }
}

export function convertLines(lines: string[], target: TargetLang): string[] {
  return lines.map((line) => convertText(line, target))
}
