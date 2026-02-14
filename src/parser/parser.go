package parser

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/orchestra-mcp/markdown/src/types"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownParser renders markdown to HTML using goldmark.
type MarkdownParser struct {
	md   goldmark.Markdown
	opts types.RenderOptions
}

// New creates a MarkdownParser with the given options.
func New(opts types.RenderOptions) *MarkdownParser {
	theme := opts.CodeTheme
	if theme == "" {
		theme = "monokai"
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	return &MarkdownParser{md: md, opts: opts}
}

// Render converts markdown bytes to an HTML string.
func (p *MarkdownParser) Render(input []byte) (*types.RenderResult, error) {
	var buf bytes.Buffer
	if err := p.md.Convert(input, &buf); err != nil {
		return nil, err
	}

	result := &types.RenderResult{
		HTML:       buf.String(),
		CodeBlocks: p.ExtractCodeBlocks(input),
	}

	if p.opts.EnableTOC {
		result.TOC = p.ExtractTOC(input)
	}

	meta, _ := p.ExtractFrontmatter(input)
	if len(meta) > 0 {
		result.Metadata = meta
	}

	return result, nil
}

// headingRe matches ATX-style headings (# Heading).
var headingRe = regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)

// ExtractTOC extracts ATX-style headings from markdown source.
func (p *MarkdownParser) ExtractTOC(input []byte) []types.TOCEntry {
	matches := headingRe.FindAllSubmatch(input, -1)
	entries := make([]types.TOCEntry, 0, len(matches))

	for _, m := range matches {
		text := strings.TrimSpace(string(m[2]))
		id := slugify(text)
		entries = append(entries, types.TOCEntry{
			Level: len(m[1]),
			Text:  text,
			ID:    id,
		})
	}

	return entries
}

// codeBlockRe matches fenced code blocks with optional language.
var codeBlockRe = regexp.MustCompile("(?ms)^```(\\w*)\\n(.*?)^```")

// ExtractCodeBlocks extracts fenced code blocks from markdown source.
func (p *MarkdownParser) ExtractCodeBlocks(input []byte) []types.CodeBlock {
	matches := codeBlockRe.FindAllSubmatch(input, -1)
	blocks := make([]types.CodeBlock, 0, len(matches))

	for _, m := range matches {
		code := string(m[2])
		blocks = append(blocks, types.CodeBlock{
			Language:  string(m[1]),
			Code:      code,
			LineCount: strings.Count(code, "\n") + 1,
		})
	}

	return blocks
}

// ExtractFrontmatter extracts YAML frontmatter delimited by --- lines.
// Returns the parsed key-value pairs and the remaining body.
func (p *MarkdownParser) ExtractFrontmatter(input []byte) (map[string]string, []byte) {
	s := string(input)
	if !strings.HasPrefix(s, "---\n") {
		return nil, input
	}

	end := strings.Index(s[4:], "\n---")
	if end < 0 {
		return nil, input
	}

	block := s[4 : 4+end]
	body := []byte(s[4+end+4:])
	meta := make(map[string]string)

	for _, line := range strings.Split(block, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if key != "" {
				meta[key] = val
			}
		}
	}

	return meta, body
}

// slugify converts a heading text to a URL-safe ID.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
