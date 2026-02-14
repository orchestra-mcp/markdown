package service

import (
	"fmt"

	"github.com/orchestra-mcp/markdown/src/parser"
	"github.com/orchestra-mcp/markdown/src/types"
)

// MarkdownService provides the full markdown rendering pipeline.
type MarkdownService struct {
	parser       *parser.MarkdownParser
	sanitizer    *parser.HTMLSanitizer
	maxInputSize int
	sanitize     bool
}

// New creates a MarkdownService with the given parser, sanitizer, and limits.
func New(p *parser.MarkdownParser, sanitize bool, maxInputSize int) *MarkdownService {
	return &MarkdownService{
		parser:       p,
		sanitizer:    parser.NewSanitizer(),
		maxInputSize: maxInputSize,
		sanitize:     sanitize,
	}
}

// Render executes the full pipeline: validate, parse, sanitize, extract.
func (s *MarkdownService) Render(req types.RenderRequest) (*types.RenderResult, error) {
	if len(req.Content) == 0 {
		return &types.RenderResult{HTML: ""}, nil
	}

	if s.maxInputSize > 0 && len(req.Content) > s.maxInputSize {
		return nil, fmt.Errorf("input exceeds maximum size of %d bytes", s.maxInputSize)
	}

	result, err := s.parser.Render([]byte(req.Content))
	if err != nil {
		return nil, fmt.Errorf("render failed: %w", err)
	}

	if s.sanitize {
		result.HTML = s.sanitizer.Sanitize(result.HTML)
	}

	return result, nil
}

// RenderString is a convenience method: string in, HTML string out.
func (s *MarkdownService) RenderString(content string) (string, error) {
	result, err := s.Render(types.RenderRequest{Content: content})
	if err != nil {
		return "", err
	}
	return result.HTML, nil
}

// ExtractTOC returns the table of contents for the given markdown.
func (s *MarkdownService) ExtractTOC(content string) ([]types.TOCEntry, error) {
	if len(content) == 0 {
		return nil, nil
	}
	if s.maxInputSize > 0 && len(content) > s.maxInputSize {
		return nil, fmt.Errorf("input exceeds maximum size of %d bytes", s.maxInputSize)
	}
	return s.parser.ExtractTOC([]byte(content)), nil
}

// ExtractCodeBlocks returns all fenced code blocks from the given markdown.
func (s *MarkdownService) ExtractCodeBlocks(content string) ([]types.CodeBlock, error) {
	if len(content) == 0 {
		return nil, nil
	}
	if s.maxInputSize > 0 && len(content) > s.maxInputSize {
		return nil, fmt.Errorf("input exceeds maximum size of %d bytes", s.maxInputSize)
	}
	return s.parser.ExtractCodeBlocks([]byte(content)), nil
}
