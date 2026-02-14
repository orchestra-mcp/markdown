package providers

import (
	"fmt"

	"github.com/orchestra-mcp/framework/app/plugins"
	"github.com/orchestra-mcp/markdown/src/types"
)

// McpTools returns MCP tool definitions contributed by the Markdown plugin.
func (p *MarkdownPlugin) McpTools() []plugins.McpToolDefinition {
	return []plugins.McpToolDefinition{
		{
			Name:        "render_markdown",
			Description: "Render markdown content to HTML",
			InputSchema: map[string]any{
				"content": map[string]any{"type": "string", "description": "Markdown content to render"},
				"format":  map[string]any{"type": "string", "description": "Output format: html, text, ast"},
			},
			Handler: p.toolRenderMarkdown,
		},
		{
			Name:        "extract_toc",
			Description: "Extract table of contents from markdown",
			InputSchema: map[string]any{
				"content": map[string]any{"type": "string", "description": "Markdown content"},
			},
			Handler: p.toolExtractTOC,
		},
		{
			Name:        "extract_code_blocks",
			Description: "Extract fenced code blocks from markdown",
			InputSchema: map[string]any{
				"content": map[string]any{"type": "string", "description": "Markdown content"},
			},
			Handler: p.toolExtractCodeBlocks,
		},
	}
}

func (p *MarkdownPlugin) toolRenderMarkdown(input map[string]any) (any, error) {
	content, _ := input["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	format, _ := input["format"].(string)
	req := types.RenderRequest{Content: content, Format: format}

	result, err := p.svc.Render(req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *MarkdownPlugin) toolExtractTOC(input map[string]any) (any, error) {
	content, _ := input["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	toc, err := p.svc.ExtractTOC(content)
	if err != nil {
		return nil, err
	}

	return map[string]any{"toc": toc}, nil
}

func (p *MarkdownPlugin) toolExtractCodeBlocks(input map[string]any) (any, error) {
	content, _ := input["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	blocks, err := p.svc.ExtractCodeBlocks(content)
	if err != nil {
		return nil, err
	}

	return map[string]any{"code_blocks": blocks}, nil
}
