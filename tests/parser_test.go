package tests

import (
	"strings"
	"testing"

	"github.com/orchestra-mcp/markdown/src/parser"
	"github.com/orchestra-mcp/markdown/src/service"
	"github.com/orchestra-mcp/markdown/src/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ──────────────────────────────────────────────────────

func newParser(enableTOC bool) *parser.MarkdownParser {
	return parser.New(types.RenderOptions{
		EnableTOC: enableTOC,
		CodeTheme: "monokai",
	})
}

func newService() *service.MarkdownService {
	p := newParser(true)
	return service.New(p, true, 1048576)
}

// ── Basic Rendering ──────────────────────────────────────────────

func TestRenderBasicMarkdown(t *testing.T) {
	svc := newService()
	html, err := svc.RenderString("# Hello\n\nWorld\n\n- one\n- two\n")
	require.NoError(t, err)

	assert.Contains(t, html, "<h1")
	assert.Contains(t, html, "Hello")
	assert.Contains(t, html, "<p>World</p>")
	assert.Contains(t, html, "<li>one</li>")
	assert.Contains(t, html, "<li>two</li>")
}

func TestRenderCodeBlock(t *testing.T) {
	md := "```go\nfmt.Println(\"hello\")\n```\n"
	svc := newService()
	html, err := svc.RenderString(md)
	require.NoError(t, err)

	assert.Contains(t, html, "Println")
	assert.Contains(t, html, "<pre")
}

func TestRenderTable(t *testing.T) {
	md := "| A | B |\n|---|---|\n| 1 | 2 |\n"
	svc := newService()
	html, err := svc.RenderString(md)
	require.NoError(t, err)

	assert.Contains(t, html, "<table>")
	assert.Contains(t, html, "<td>1</td>")
	assert.Contains(t, html, "<td>2</td>")
}

func TestRenderLinks(t *testing.T) {
	md := "[Google](https://google.com)\n"
	svc := newService()
	html, err := svc.RenderString(md)
	require.NoError(t, err)

	assert.Contains(t, html, `<a href="https://google.com"`)
	assert.Contains(t, html, "Google")
}

func TestRenderInlineCode(t *testing.T) {
	md := "Use `fmt.Println` to print.\n"
	svc := newService()
	html, err := svc.RenderString(md)
	require.NoError(t, err)

	assert.Contains(t, html, "<code>fmt.Println</code>")
}

func TestRenderBoldItalic(t *testing.T) {
	md := "**bold** and *italic* text\n"
	svc := newService()
	html, err := svc.RenderString(md)
	require.NoError(t, err)

	assert.Contains(t, html, "<strong>bold</strong>")
	assert.Contains(t, html, "<em>italic</em>")
}

// ── TOC Extraction ───────────────────────────────────────────────

func TestExtractTOC(t *testing.T) {
	md := "# Title\n## Section A\n### Sub A1\n## Section B\n"
	svc := newService()
	toc, err := svc.ExtractTOC(md)
	require.NoError(t, err)

	require.Len(t, toc, 4)
	assert.Equal(t, 1, toc[0].Level)
	assert.Equal(t, "Title", toc[0].Text)
	assert.Equal(t, "title", toc[0].ID)
	assert.Equal(t, 2, toc[1].Level)
	assert.Equal(t, "Section A", toc[1].Text)
	assert.Equal(t, 3, toc[2].Level)
}

// ── Code Block Extraction ────────────────────────────────────────

func TestExtractCodeBlocks(t *testing.T) {
	md := "```python\nprint('hi')\n```\n\ntext\n\n```js\nconsole.log('hi')\n```\n"
	svc := newService()
	blocks, err := svc.ExtractCodeBlocks(md)
	require.NoError(t, err)

	require.Len(t, blocks, 2)
	assert.Equal(t, "python", blocks[0].Language)
	assert.Contains(t, blocks[0].Code, "print('hi')")
	assert.Equal(t, "js", blocks[1].Language)
	assert.Contains(t, blocks[1].Code, "console.log")
}

// ── Frontmatter ──────────────────────────────────────────────────

func TestFrontmatter(t *testing.T) {
	md := "---\ntitle: My Doc\nauthor: Alice\n---\n# Hello\n"
	p := newParser(false)
	meta, body := p.ExtractFrontmatter([]byte(md))

	assert.Equal(t, "My Doc", meta["title"])
	assert.Equal(t, "Alice", meta["author"])
	assert.Contains(t, string(body), "# Hello")
}

func TestFrontmatterAbsent(t *testing.T) {
	md := "# No frontmatter\n"
	p := newParser(false)
	meta, body := p.ExtractFrontmatter([]byte(md))

	assert.Nil(t, meta)
	assert.Equal(t, md, string(body))
}

// ── Sanitization ─────────────────────────────────────────────────

func TestSanitizeHTML(t *testing.T) {
	s := parser.NewSanitizer()

	input := `<p>Hello</p><script>alert('xss')</script><p>World</p>`
	result := s.Sanitize(input)

	assert.NotContains(t, result, "script")
	assert.NotContains(t, result, "alert")
	assert.Contains(t, result, "<p>Hello</p>")
	assert.Contains(t, result, "<p>World</p>")
}

func TestSanitizeEventHandlers(t *testing.T) {
	s := parser.NewSanitizer()

	input := `<div onclick="alert('xss')">Click</div>`
	result := s.Sanitize(input)

	assert.NotContains(t, result, "onclick")
	assert.Contains(t, result, "Click")
}

func TestSanitizeIframe(t *testing.T) {
	s := parser.NewSanitizer()

	input := `<p>Before</p><iframe src="evil.com"></iframe><p>After</p>`
	result := s.Sanitize(input)

	assert.NotContains(t, result, "iframe")
	assert.Contains(t, result, "<p>Before</p>")
	assert.Contains(t, result, "<p>After</p>")
}

// ── Size Limit ───────────────────────────────────────────────────

func TestMaxInputSize(t *testing.T) {
	svc := service.New(newParser(false), false, 100)

	big := strings.Repeat("x", 101)
	_, err := svc.RenderString(big)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum size")
}

func TestMaxInputSizeAllowed(t *testing.T) {
	svc := service.New(newParser(false), false, 100)

	small := strings.Repeat("x", 50)
	html, err := svc.RenderString(small)
	require.NoError(t, err)
	assert.NotEmpty(t, html)
}

// ── Empty Input ──────────────────────────────────────────────────

func TestEmptyInput(t *testing.T) {
	svc := newService()
	html, err := svc.RenderString("")
	require.NoError(t, err)
	assert.Equal(t, "", html)
}

func TestEmptyTOC(t *testing.T) {
	svc := newService()
	toc, err := svc.ExtractTOC("")
	require.NoError(t, err)
	assert.Nil(t, toc)
}

func TestEmptyCodeBlocks(t *testing.T) {
	svc := newService()
	blocks, err := svc.ExtractCodeBlocks("")
	require.NoError(t, err)
	assert.Nil(t, blocks)
}
