package config

// MarkdownConfig holds configuration for the Markdown plugin.
type MarkdownConfig struct {
	Enabled             bool   `json:"enabled"`
	SanitizeHTML        bool   `json:"sanitize_html"`
	EnableMermaid       bool   `json:"enable_mermaid"`
	EnableMath          bool   `json:"enable_math"`
	EnableTableOfContents bool `json:"enable_table_of_contents"`
	MaxInputSize        int    `json:"max_input_size"`
	CodeTheme           string `json:"code_theme"`
}

// DefaultConfig returns the default markdown configuration.
func DefaultConfig() *MarkdownConfig {
	return &MarkdownConfig{
		Enabled:             true,
		SanitizeHTML:        true,
		EnableMermaid:       true,
		EnableMath:          true,
		EnableTableOfContents: true,
		MaxInputSize:        1048576, // 1MB
		CodeTheme:           "monokai",
	}
}
