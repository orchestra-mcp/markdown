package providers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/orchestra-mcp/markdown/src/types"
)

// RegisterRoutes registers REST API endpoints at /markdown.
func (p *MarkdownPlugin) RegisterRoutes(group fiber.Router) {
	g := group.Group("/markdown")

	g.Post("/render", p.handleRender)
	g.Post("/toc", p.handleTOC)
	g.Post("/code-blocks", p.handleCodeBlocks)
}

func (p *MarkdownPlugin) handleRender(c fiber.Ctx) error {
	var body struct {
		Content string               `json:"content"`
		Options *types.RenderOptions `json:"options,omitempty"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_body", "message": err.Error(),
		})
	}

	req := types.RenderRequest{Content: body.Content}
	if body.Options != nil {
		req.Options = *body.Options
	}

	result, err := p.svc.Render(req)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "render_failed", "message": err.Error(),
		})
	}

	return c.JSON(result)
}

func (p *MarkdownPlugin) handleTOC(c fiber.Ctx) error {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_body", "message": err.Error(),
		})
	}

	toc, err := p.svc.ExtractTOC(body.Content)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "extract_failed", "message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"toc": toc})
}

func (p *MarkdownPlugin) handleCodeBlocks(c fiber.Ctx) error {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_body", "message": err.Error(),
		})
	}

	blocks, err := p.svc.ExtractCodeBlocks(body.Content)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "extract_failed", "message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"code_blocks": blocks})
}
