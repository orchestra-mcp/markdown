package parser

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HTMLSanitizer strips dangerous HTML elements and attributes using a
// DOM-based allowlist. Parses HTML into a node tree, walks every node,
// and removes anything not explicitly allowed.
type HTMLSanitizer struct{}

// NewSanitizer creates a new HTMLSanitizer.
func NewSanitizer() *HTMLSanitizer {
	return &HTMLSanitizer{}
}

// dropEntireSubtree lists tags whose entire subtree (including text
// children) is removed. These tags contain executable or opaque content.
var dropEntireSubtree = map[atom.Atom]bool{
	atom.Script:   true,
	atom.Style:    true,
	atom.Iframe:   true,
	atom.Object:   true,
	atom.Embed:    true,
	atom.Applet:   true,
	atom.Noscript: true,
}

// allowedTags is the set of tags that survive sanitization.
var allowedTags = map[atom.Atom]bool{
	atom.A: true, atom.Abbr: true, atom.B: true, atom.Blockquote: true,
	atom.Br: true, atom.Code: true, atom.Dd: true, atom.Del: true,
	atom.Details: true, atom.Div: true, atom.Dl: true, atom.Dt: true,
	atom.Em: true, atom.H1: true, atom.H2: true, atom.H3: true,
	atom.H4: true, atom.H5: true, atom.H6: true, atom.Hr: true,
	atom.I: true, atom.Img: true, atom.Ins: true, atom.Kbd: true,
	atom.Li: true, atom.Ol: true, atom.P: true, atom.Pre: true,
	atom.Q: true, atom.S: true, atom.Samp: true, atom.Small: true,
	atom.Span: true, atom.Strong: true, atom.Sub: true, atom.Summary: true,
	atom.Sup: true, atom.Table: true, atom.Tbody: true, atom.Td: true,
	atom.Tfoot: true, atom.Th: true, atom.Thead: true, atom.Tr: true,
	atom.U: true, atom.Ul: true, atom.Var: true,
}

// safeAttrs lists attributes kept on allowed tags.
var safeAttrs = map[string]bool{
	"href": true, "src": true, "alt": true, "title": true,
	"class": true, "id": true, "width": true, "height": true,
	"colspan": true, "rowspan": true, "align": true,
}

// Sanitize parses HTML into a DOM, removes disallowed elements and
// attributes, and renders the cleaned tree back to a string.
func (s *HTMLSanitizer) Sanitize(raw string) string {
	nodes, err := html.ParseFragment(
		strings.NewReader(raw),
		&html.Node{Type: html.ElementNode, DataAtom: atom.Body, Data: "body"},
	)
	if err != nil {
		return html.EscapeString(raw)
	}

	var buf bytes.Buffer
	for _, n := range nodes {
		renderClean(&buf, n)
	}
	return strings.TrimSpace(buf.String())
}

// renderClean handles a single top-level node: drops dangerous elements,
// unwraps disallowed-but-safe elements, and cleans allowed elements.
func renderClean(buf *bytes.Buffer, n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		if dropEntireSubtree[n.DataAtom] {
			return
		}
		if !allowedTags[n.DataAtom] {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderClean(buf, c)
			}
			return
		}
		cleanAttrs(n)
		walkAndClean(n)
		if err := html.Render(buf, n); err != nil {
			return
		}
	case html.TextNode:
		if err := html.Render(buf, n); err != nil {
			return
		}
	default:
		// drop comments, doctypes, etc.
	}
}

// walkAndClean recursively removes disallowed nodes and attributes.
func walkAndClean(n *html.Node) {
	var next *html.Node
	for c := n.FirstChild; c != nil; c = next {
		next = c.NextSibling

		switch c.Type {
		case html.ElementNode:
			if dropEntireSubtree[c.DataAtom] {
				n.RemoveChild(c)
				continue
			}
			if !allowedTags[c.DataAtom] {
				promoteChildren(n, c)
				continue
			}
			cleanAttrs(c)
			walkAndClean(c)
		case html.TextNode:
			// keep
		default:
			n.RemoveChild(c)
		}
	}
}

// promoteChildren moves all children of child into parent (before
// child's position) and removes child.
func promoteChildren(parent, child *html.Node) {
	for child.FirstChild != nil {
		moved := child.FirstChild
		child.RemoveChild(moved)
		parent.InsertBefore(moved, child)
	}
	parent.RemoveChild(child)
}

// cleanAttrs removes dangerous attributes from an element node.
func cleanAttrs(n *html.Node) {
	kept := make([]html.Attribute, 0, len(n.Attr))
	for _, attr := range n.Attr {
		key := strings.ToLower(attr.Key)
		if strings.HasPrefix(key, "on") {
			continue
		}
		if !safeAttrs[key] {
			continue
		}
		if key == "href" || key == "src" {
			val := strings.TrimSpace(strings.ToLower(attr.Val))
			if strings.HasPrefix(val, "javascript:") || strings.HasPrefix(val, "data:") {
				continue
			}
		}
		kept = append(kept, attr)
	}
	n.Attr = kept
}
