package types

// RenderRequest represents a request to render markdown content.
type RenderRequest struct {
	Content string        `json:"content"`
	Format  string        `json:"format"` // "html", "text", "ast"
	Options RenderOptions `json:"options"`
}

// RenderOptions configures how markdown is rendered.
type RenderOptions struct {
	SanitizeHTML  bool   `json:"sanitize_html"`
	EnableMermaid bool   `json:"enable_mermaid"`
	EnableMath    bool   `json:"enable_math"`
	EnableTOC     bool   `json:"enable_toc"`
	CodeTheme     string `json:"code_theme"`
}

// RenderResult holds the output of a markdown render operation.
type RenderResult struct {
	HTML       string            `json:"html"`
	TOC        []TOCEntry        `json:"toc,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CodeBlocks []CodeBlock       `json:"code_blocks,omitempty"`
}

// TOCEntry represents one heading in the table of contents.
type TOCEntry struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
	ID    string `json:"id"`
}

// CodeBlock represents a fenced code block extracted from markdown.
type CodeBlock struct {
	Language  string `json:"language"`
	Code      string `json:"code"`
	LineCount int    `json:"line_count"`
}
