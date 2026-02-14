# Orchestra Markdown Plugin

Pure Go markdown-to-HTML renderer with sanitization, table-of-contents extraction, code block extraction, and syntax highlighting. Built on goldmark with GFM extensions.

## Features

- **Goldmark rendering** — GFM tables, strikethrough, autolinks, task lists, typographer
- **Syntax highlighting** — Chroma-based code highlighting with configurable themes
- **HTML sanitization** — DOM-based allowlist sanitizer (strips scripts, iframes, event handlers)
- **TOC extraction** — structured heading tree with levels and anchors
- **Code block extraction** — fenced blocks with language detection and line counts
- **Input size limits** — configurable maximum input size (default 1MB)

## Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `Enabled` | true | Plugin on/off |
| `SanitizeHTML` | true | Enable HTML sanitization |
| `EnableMermaid` | true | Mermaid diagram support |
| `EnableMath` | true | Math notation support |
| `EnableTOC` | true | Table of contents extraction |
| `MaxInputSize` | 1048576 | Max input bytes (1MB) |
| `CodeTheme` | `monokai` | Syntax highlighting theme |

## MCP Tools

| Tool | Description |
|------|-------------|
| `render_markdown` | Render markdown to HTML |
| `extract_toc` | Extract heading tree |
| `extract_code_blocks` | Extract fenced code blocks |

## REST API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/markdown/render` | Render markdown to HTML |
| `POST` | `/markdown/toc` | Extract table of contents |
| `POST` | `/markdown/code-blocks` | Extract code blocks |

## Package Structure

```
plugins/markdown/
├── config/markdown.go           # MarkdownConfig
├── providers/
│   ├── plugin.go                # MarkdownPlugin (activate, services, tools)
│   ├── routes.go                # REST endpoints
│   └── tools.go                 # 3 MCP tool definitions
├── src/
│   ├── parser/
│   │   ├── parser.go            # MarkdownParser (goldmark + highlighting)
│   │   └── sanitize.go          # HTMLSanitizer (DOM-based allowlist)
│   ├── service/service.go       # MarkdownService (render, TOC, code blocks)
│   └── types/types.go           # RenderRequest, RenderResult, TOCEntry, CodeBlock
├── tests/parser_test.go         # 18 tests (rendering, sanitization, extraction)
└── go.mod
```
