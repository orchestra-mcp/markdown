package providers

import (
	"github.com/orchestra-mcp/framework/app/plugins"
	"github.com/orchestra-mcp/markdown/config"
	"github.com/orchestra-mcp/markdown/src/parser"
	"github.com/orchestra-mcp/markdown/src/service"
	"github.com/orchestra-mcp/markdown/src/types"
)

// MarkdownPlugin implements the Orchestra plugin interface for markdown parsing.
type MarkdownPlugin struct {
	active bool
	ctx    *plugins.PluginContext
	cfg    *config.MarkdownConfig
	svc    *service.MarkdownService
}

// NewMarkdownPlugin creates a new Markdown plugin instance.
func NewMarkdownPlugin() *MarkdownPlugin { return &MarkdownPlugin{} }

func (p *MarkdownPlugin) ID() string             { return "orchestra/markdown" }
func (p *MarkdownPlugin) Name() string           { return "Markdown Parser" }
func (p *MarkdownPlugin) Version() string        { return "0.1.0" }
func (p *MarkdownPlugin) Dependencies() []string { return nil }
func (p *MarkdownPlugin) IsActive() bool         { return p.active }
func (p *MarkdownPlugin) FeatureFlag() string    { return "markdown" }
func (p *MarkdownPlugin) ConfigKey() string      { return "markdown" }

func (p *MarkdownPlugin) DefaultConfig() map[string]any {
	return map[string]any{
		"sanitize_html":          true,
		"enable_mermaid":         true,
		"enable_math":            true,
		"enable_table_of_contents": true,
		"max_input_size":         1048576,
		"code_theme":             "monokai",
	}
}

// Activate initializes the parser and service with config values.
func (p *MarkdownPlugin) Activate(ctx *plugins.PluginContext) error {
	p.ctx = ctx
	p.cfg = config.DefaultConfig()

	if v, ok := ctx.GetConfig("code_theme"); ok {
		if s, ok := v.(string); ok && s != "" {
			p.cfg.CodeTheme = s
		}
	}

	opts := types.RenderOptions{
		SanitizeHTML:  p.cfg.SanitizeHTML,
		EnableMermaid: p.cfg.EnableMermaid,
		EnableMath:    p.cfg.EnableMath,
		EnableTOC:     p.cfg.EnableTableOfContents,
		CodeTheme:     p.cfg.CodeTheme,
	}

	mdParser := parser.New(opts)
	p.svc = service.New(mdParser, p.cfg.SanitizeHTML, p.cfg.MaxInputSize)

	p.active = true
	ctx.Logger.Info().Str("plugin", p.ID()).Msg("markdown plugin activated")
	return nil
}

// Deactivate shuts down the markdown plugin.
func (p *MarkdownPlugin) Deactivate() error {
	p.active = false
	return nil
}

// Service returns the underlying MarkdownService.
func (p *MarkdownPlugin) Service() *service.MarkdownService {
	return p.svc
}

// Services returns ServiceDefinitions for the DI registry.
func (p *MarkdownPlugin) Services() []plugins.ServiceDefinition {
	return []plugins.ServiceDefinition{
		{
			ID:      "markdown",
			Factory: func() any { return p.svc },
		},
	}
}

// Compile-time interface assertions.
var (
	_ plugins.Plugin         = (*MarkdownPlugin)(nil)
	_ plugins.HasConfig      = (*MarkdownPlugin)(nil)
	_ plugins.HasFeatureFlag = (*MarkdownPlugin)(nil)
	_ plugins.HasServices    = (*MarkdownPlugin)(nil)
	_ plugins.HasRoutes      = (*MarkdownPlugin)(nil)
	_ plugins.HasMcpTools    = (*MarkdownPlugin)(nil)
)
