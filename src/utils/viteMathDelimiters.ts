import type { Plugin } from "vite";

// Rewrites the LaTeX delimiters \(...\) and \[...\] to $$...$$ in raw
// Markdown/MDX *before* remark parses it. remark/CommonMark strips the
// backslash from \( and \{ (they are ASCII-punctuation escapes), so this has
// to run on the source string, not on the mdast. Block $$…$$ is captured raw
// by remark-math, so \\ and \{ survive intact inside it.
//
// shortcut: only \(…\) / \[…\] are handled and inline code spans / fenced code
// blocks are skipped; upgrade to a full tokenizer if a formula is ever split
// across inline nodes (emphasis, links, …).
export function viteMathDelimiters(): Plugin {
  return {
    name: "vite-math-delimiters",
    enforce: "pre",
    transform(code, id) {
      if (!/\.mdx?$/.test(id)) return null;
      return { code: rewriteDelimiters(code), map: null };
    },
  };
}

function rewriteDelimiters(src: string): string {
  let inFence = false;
  let fenceChar = "";
  let fenceLen = 0;

  return src
    .split("\n")
    .map(line => {
      const fence = line.match(/^\s*(`{3,}|~{3,})/);
      if (fence) {
        const marker = fence[1];
        if (!inFence) {
          inFence = true;
          fenceChar = marker[0];
          fenceLen = marker.length;
        } else if (marker[0] === fenceChar && marker.length >= fenceLen) {
          inFence = false;
        }
        return line;
      }
      return inFence ? line : replaceOutsideInlineCode(line);
    })
    .join("\n");
}

function replaceOutsideInlineCode(line: string): string {
  // Odd-indexed parts are inline `code` spans; leave them untouched.
  return line
    .split(/(`+[^`]*`+)/g)
    .map((part, i) =>
      i % 2 === 1
        ? part
        : part
            .replace(/\\\(/g, "$$$$")
            .replace(/\\\)/g, "$$$$")
            .replace(/\\\[/g, "$$$$")
            .replace(/\\\]/g, "$$$$")
    )
    .join("");
}
